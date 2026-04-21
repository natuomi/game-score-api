package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"game-score-api/internal/model"
	"game-score-api/internal/repository"

	"github.com/redis/go-redis/v9"
)

const rankingCacheKey = "rankings:top"

// ScoreService はスコア管理とランキングのビジネスロジックを担当する
type ScoreService struct {
	scoreRepo *repository.ScoreRepository
	redis     *redis.Client
}

func NewScoreService(scoreRepo *repository.ScoreRepository, rdb *redis.Client) *ScoreService {
	return &ScoreService{scoreRepo: scoreRepo, redis: rdb}
}

// PostScore はスコアを登録し、Redisキャッシュを無効化する
func (s *ScoreService) PostScore(ctx context.Context, userID string, req model.PostScoreRequest) (*model.Score, error) {
	gameMode := req.GameMode
	if gameMode == "" {
		gameMode = "classic"
	}

	score, err := s.scoreRepo.Create(ctx, userID, req.Score, gameMode)
	if err != nil {
		return nil, err
	}

	// 新スコア登録でランキングが変わる可能性があるためキャッシュ削除
	s.redis.Del(ctx, rankingCacheKey)

	return score, nil
}

// GetMyScores は認証ユーザー自身のスコア履歴を返す
func (s *ScoreService) GetMyScores(ctx context.Context, userID string) ([]model.Score, error) {
	return s.scoreRepo.FindByUserID(ctx, userID)
}

// GetRankings はランキングを返す
// Redisにキャッシュがある場合はキャッシュから、なければDBから取得してキャッシュに保存する
func (s *ScoreService) GetRankings(ctx context.Context) (*model.RankingResponse, error) {
	limit := getRankingLimit()

	// ① Redisキャッシュを確認
	cached, err := s.redis.Get(ctx, rankingCacheKey).Result()
	if err == nil {
		// キャッシュヒット → JSONをデシリアライズして返す
		var rankings []model.RankingEntry
		if jsonErr := json.Unmarshal([]byte(cached), &rankings); jsonErr == nil {
			total, _ := s.scoreRepo.CountAll(ctx)
			return &model.RankingResponse{
				Rankings: rankings,
				Cached:   true,
				Total:    total,
			}, nil
		}
	}

	// ② キャッシュミス → DBから取得
	rankings, err := s.scoreRepo.GetRankings(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get rankings from db: %w", err)
	}

	// ③ Redisにキャッシュ保存（TTLはenv変数で設定可能）
	ttl := getRankingCacheTTL()
	if data, jsonErr := json.Marshal(rankings); jsonErr == nil {
		s.redis.Set(ctx, rankingCacheKey, data, ttl)
	}

	total, _ := s.scoreRepo.CountAll(ctx)
	return &model.RankingResponse{
		Rankings: rankings,
		Cached:   false,
		Total:    total,
	}, nil
}

// getRankingLimit は環境変数からランキング表示件数を取得する（デフォルト100）
func getRankingLimit() int {
	if v := os.Getenv("RANKING_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 100
}

// getRankingCacheTTL は環境変数からキャッシュTTLを取得する（デフォルト60秒）
func getRankingCacheTTL() time.Duration {
	if v := os.Getenv("RANKING_CACHE_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return 60 * time.Second
}
