package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"game-score-api/internal/model"

	"github.com/redis/go-redis/v9"
)

// ── スコアリポジトリ モック ────────────────────────────────────────────

// mockScoreRepo は ScoreRepositoryInterface の手書きモック実装
type mockScoreRepo struct {
	CreateFn      func(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error)
	FindByUserIDFn func(ctx context.Context, userID string) ([]model.Score, error)
	GetRankingsFn  func(ctx context.Context, limit int) ([]model.RankingEntry, error)
	CountAllFn     func(ctx context.Context) (int, error)

	CreateCalled      int
	GetRankingsCalled int
	CountAllCalled    int
}

func (m *mockScoreRepo) Create(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error) {
	m.CreateCalled++
	if m.CreateFn != nil {
		return m.CreateFn(ctx, userID, score, gameMode)
	}
	return nil, errors.New("CreateFn not set")
}

func (m *mockScoreRepo) FindByUserID(ctx context.Context, userID string) ([]model.Score, error) {
	if m.FindByUserIDFn != nil {
		return m.FindByUserIDFn(ctx, userID)
	}
	return nil, errors.New("FindByUserIDFn not set")
}

func (m *mockScoreRepo) GetRankings(ctx context.Context, limit int) ([]model.RankingEntry, error) {
	m.GetRankingsCalled++
	if m.GetRankingsFn != nil {
		return m.GetRankingsFn(ctx, limit)
	}
	return nil, errors.New("GetRankingsFn not set")
}

func (m *mockScoreRepo) CountAll(ctx context.Context) (int, error) {
	m.CountAllCalled++
	if m.CountAllFn != nil {
		return m.CountAllFn(ctx)
	}
	return 0, nil
}

// ── Redis モック ───────────────────────────────────────────────────────

// mockRedis は RedisClient インターフェースの手書きモック実装
// key→value のインメモリストアとして動作する
type mockRedis struct {
	store     map[string]string
	DelCalled []string // Del で削除されたキーを記録
	SetCalled int      // Set 呼び出し回数
	GetCalled int      // Get 呼び出し回数
}

func newMockRedis() *mockRedis {
	return &mockRedis{store: make(map[string]string)}
}

func (m *mockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	m.GetCalled++
	cmd := redis.NewStringCmd(ctx)
	val, ok := m.store[key]
	if !ok {
		cmd.SetErr(redis.Nil)
	} else {
		cmd.SetVal(val)
	}
	return cmd
}

func (m *mockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.SetCalled++
	cmd := redis.NewStatusCmd(ctx)
	switch v := value.(type) {
	case string:
		m.store[key] = v
	case []byte:
		m.store[key] = string(v)
	}
	cmd.SetVal("OK")
	return cmd
}

func (m *mockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	m.DelCalled = append(m.DelCalled, keys...)
	cmd := redis.NewIntCmd(ctx)
	deleted := int64(0)
	for _, k := range keys {
		if _, ok := m.store[k]; ok {
			delete(m.store, k)
			deleted++
		}
	}
	cmd.SetVal(deleted)
	return cmd
}

// ── テストヘルパー ────────────────────────────────────────────────────

// newTestScore はテスト用のスコアを生成する
func newTestScore(id, userID string, score int, gameMode string) *model.Score {
	return &model.Score{
		ID:        id,
		UserID:    userID,
		Score:     score,
		GameMode:  gameMode,
		CreatedAt: time.Now(),
	}
}

// newTestRankingEntries はテスト用のランキングエントリを生成する
func newTestRankingEntries() []model.RankingEntry {
	return []model.RankingEntry{
		{Rank: 1, PlayerName: "Alice", Score: 1000, GameMode: "classic", CreatedAt: time.Now()},
		{Rank: 2, PlayerName: "Bob", Score: 800, GameMode: "classic", CreatedAt: time.Now()},
	}
}

// ── PostScore テスト ──────────────────────────────────────────────────

// TestPostScore_Success はPostScore成功ケースを検証する
// - repo.Create が呼ばれること
// - Redisキャッシュが無効化（Del）されること
// - GameMode が空の場合は "classic" になること
func TestPostScore_Success(t *testing.T) {
	t.Parallel()

	expectedScore := newTestScore("score-1", "uid-1", 500, "classic")

	repo := &mockScoreRepo{
		CreateFn: func(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error) {
			if userID != "uid-1" {
				t.Errorf("expected userID=uid-1, got %s", userID)
			}
			if score != 500 {
				t.Errorf("expected score=500, got %d", score)
			}
			// GameMode空の場合に "classic" が補完されることを確認
			if gameMode != "classic" {
				t.Errorf("expected gameMode=classic, got %s", gameMode)
			}
			return expectedScore, nil
		},
	}
	rdb := newMockRedis()

	svc := NewScoreService(repo, rdb)
	req := model.PostScoreRequest{
		Score:    500,
		GameMode: "", // 空→ "classic" に補完される
	}

	score, err := svc.PostScore(context.Background(), "uid-1", req)
	if err != nil {
		t.Fatalf("PostScore() returned unexpected error: %v", err)
	}
	if score.ID != expectedScore.ID {
		t.Errorf("expected score.ID=%s, got %s", expectedScore.ID, score.ID)
	}
	if repo.CreateCalled != 1 {
		t.Errorf("expected repo.Create called 1 time, got %d", repo.CreateCalled)
	}

	// Redisキャッシュが無効化されていることを確認
	if len(rdb.DelCalled) == 0 {
		t.Error("expected Redis Del to be called for cache invalidation")
	}
	if rdb.DelCalled[0] != rankingCacheKey {
		t.Errorf("expected Del key=%s, got %s", rankingCacheKey, rdb.DelCalled[0])
	}
}

// TestPostScore_WithGameMode はGameModeが明示指定されている場合を検証する
func TestPostScore_WithGameMode(t *testing.T) {
	t.Parallel()

	repo := &mockScoreRepo{
		CreateFn: func(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error) {
			if gameMode != "ranked" {
				t.Errorf("expected gameMode=ranked, got %s", gameMode)
			}
			return newTestScore("score-2", userID, score, gameMode), nil
		},
	}
	rdb := newMockRedis()

	svc := NewScoreService(repo, rdb)
	req := model.PostScoreRequest{Score: 750, GameMode: "ranked"}

	_, err := svc.PostScore(context.Background(), "uid-2", req)
	if err != nil {
		t.Fatalf("PostScore() returned unexpected error: %v", err)
	}
}

// ── GetRankings テスト ────────────────────────────────────────────────

// TestGetRankings_CacheHit はRedisキャッシュヒットケースを検証する
// - Redisから返ってきたデータを使うこと
// - DBのGetRankingsが呼ばれないこと
// - Cached=true が返ること
func TestGetRankings_CacheHit(t *testing.T) {
	t.Parallel()

	entries := newTestRankingEntries()
	cached, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	rdb := newMockRedis()
	// キャッシュに事前データをセット
	rdb.store[rankingCacheKey] = string(cached)

	repo := &mockScoreRepo{
		CountAllFn: func(ctx context.Context) (int, error) {
			return 2, nil
		},
	}

	svc := NewScoreService(repo, rdb)

	resp, err := svc.GetRankings(context.Background())
	if err != nil {
		t.Fatalf("GetRankings() returned unexpected error: %v", err)
	}
	if !resp.Cached {
		t.Error("expected Cached=true on cache hit")
	}
	if len(resp.Rankings) != 2 {
		t.Errorf("expected 2 rankings, got %d", len(resp.Rankings))
	}
	if resp.Rankings[0].PlayerName != "Alice" {
		t.Errorf("expected first player=Alice, got %s", resp.Rankings[0].PlayerName)
	}
	// DBのGetRankingsが呼ばれていないことを確認
	if repo.GetRankingsCalled != 0 {
		t.Errorf("expected repo.GetRankings NOT called, but called %d times", repo.GetRankingsCalled)
	}
	// Redisのみアクセスされていることを確認
	if rdb.GetCalled != 1 {
		t.Errorf("expected Redis Get called 1 time, got %d", rdb.GetCalled)
	}
}

// TestGetRankings_CacheMiss はRedisキャッシュミスケースを検証する
// - Redisにデータがない場合はDBから取得すること
// - 取得したデータをRedisにキャッシュすること
// - Cached=false が返ること
func TestGetRankings_CacheMiss(t *testing.T) {
	t.Parallel()

	entries := newTestRankingEntries()

	repo := &mockScoreRepo{
		GetRankingsFn: func(ctx context.Context, limit int) ([]model.RankingEntry, error) {
			return entries, nil
		},
		CountAllFn: func(ctx context.Context) (int, error) {
			return 2, nil
		},
	}

	// Redisは空（キャッシュミス状態）
	rdb := newMockRedis()

	svc := NewScoreService(repo, rdb)

	resp, err := svc.GetRankings(context.Background())
	if err != nil {
		t.Fatalf("GetRankings() returned unexpected error: %v", err)
	}
	if resp.Cached {
		t.Error("expected Cached=false on cache miss")
	}
	if len(resp.Rankings) != 2 {
		t.Errorf("expected 2 rankings, got %d", len(resp.Rankings))
	}
	// DBのGetRankingsが呼ばれていることを確認
	if repo.GetRankingsCalled != 1 {
		t.Errorf("expected repo.GetRankings called 1 time, got %d", repo.GetRankingsCalled)
	}
	// Redisにキャッシュ保存されていることを確認
	if rdb.SetCalled != 1 {
		t.Errorf("expected Redis Set called 1 time for caching, got %d", rdb.SetCalled)
	}
	// キャッシュに実際にデータが保存されていることを確認
	if _, ok := rdb.store[rankingCacheKey]; !ok {
		t.Error("expected ranking data to be stored in Redis cache")
	}
	// キャッシュのデータを検証
	var cachedRankings []model.RankingEntry
	if err := json.Unmarshal([]byte(rdb.store[rankingCacheKey]), &cachedRankings); err != nil {
		t.Fatalf("failed to unmarshal cached rankings: %v", err)
	}
	if len(cachedRankings) != 2 {
		t.Errorf("expected 2 cached rankings, got %d", len(cachedRankings))
	}
}

// TestGetRankings_DBError はDB取得エラー時の挙動を検証する
func TestGetRankings_DBError(t *testing.T) {
	t.Parallel()

	repo := &mockScoreRepo{
		GetRankingsFn: func(ctx context.Context, limit int) ([]model.RankingEntry, error) {
			return nil, errors.New("db connection timeout")
		},
	}

	rdb := newMockRedis()
	svc := NewScoreService(repo, rdb)

	_, err := svc.GetRankings(context.Background())
	if err == nil {
		t.Fatal("GetRankings() should return error on DB failure")
	}
}
