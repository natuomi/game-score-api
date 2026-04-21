package repository

import (
	"context"
	"fmt"

	"game-score-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ScoreRepository はスコアのDB操作を担当する
type ScoreRepository struct {
	db *pgxpool.Pool
}

func NewScoreRepository(db *pgxpool.Pool) *ScoreRepository {
	return &ScoreRepository{db: db}
}

// Create はスコアを登録する
func (r *ScoreRepository) Create(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error) {
	var s model.Score
	err := r.db.QueryRow(ctx,
		`INSERT INTO scores (user_id, score, game_mode)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, score, game_mode, created_at`,
		userID, score, gameMode,
	).Scan(&s.ID, &s.UserID, &s.Score, &s.GameMode, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create score: %w", err)
	}
	return &s, nil
}

// FindByUserID はユーザーIDでスコア一覧を返す（新しい順）
func (r *ScoreRepository) FindByUserID(ctx context.Context, userID string) ([]model.Score, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, score, game_mode, created_at
		 FROM scores WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("find scores by user: %w", err)
	}
	defer rows.Close()

	var scores []model.Score
	for rows.Next() {
		var s model.Score
		if err := rows.Scan(&s.ID, &s.UserID, &s.Score, &s.GameMode, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan score: %w", err)
		}
		scores = append(scores, s)
	}
	return scores, nil
}

// GetRankings はスコア降順でランキングを返す（users JOINでプレイヤー名を含む）
func (r *ScoreRepository) GetRankings(ctx context.Context, limit int) ([]model.RankingEntry, error) {
	rows, err := r.db.Query(ctx,
		`SELECT u.name, s.score, s.game_mode, s.created_at
		 FROM scores s
		 JOIN users u ON s.user_id = u.id
		 ORDER BY s.score DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get rankings: %w", err)
	}
	defer rows.Close()

	var rankings []model.RankingEntry
	rank := 1
	for rows.Next() {
		var e model.RankingEntry
		if err := rows.Scan(&e.PlayerName, &e.Score, &e.GameMode, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ranking: %w", err)
		}
		e.Rank = rank
		rankings = append(rankings, e)
		rank++
	}
	return rankings, nil
}

// CountAll はスコアの総件数を返す
func (r *ScoreRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM scores`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count scores: %w", err)
	}
	return count, nil
}
