package repository

import (
	"context"

	"game-score-api/internal/model"
)

// UserRepositoryInterface はユーザーのDB操作を抽象化するインターフェース
// テスト時にモック実装を差し込めるようにする
type UserRepositoryInterface interface {
	Create(ctx context.Context, name, email, hashedPassword string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
}

// ScoreRepositoryInterface はスコアのDB操作を抽象化するインターフェース
type ScoreRepositoryInterface interface {
	Create(ctx context.Context, userID string, score int, gameMode string) (*model.Score, error)
	FindByUserID(ctx context.Context, userID string) ([]model.Score, error)
	GetRankings(ctx context.Context, limit int) ([]model.RankingEntry, error)
	CountAll(ctx context.Context) (int, error)
}

// コンパイル時に UserRepository が UserRepositoryInterface を満たすことを保証する
var _ UserRepositoryInterface = (*UserRepository)(nil)

// コンパイル時に ScoreRepository が ScoreRepositoryInterface を満たすことを保証する
var _ ScoreRepositoryInterface = (*ScoreRepository)(nil)
