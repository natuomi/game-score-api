package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool は PostgreSQL 接続プールを作成する
// 使い方:
//
//	pool, err := database.NewPostgresPool()
//	defer pool.Close()
func NewPostgresPool() (*pgxpool.Pool, error) {
	dsn := buildDSN()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 接続確認
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}

// buildDSN は環境変数から DSN 文字列を組み立てる
func buildDSN() string {
	// DATABASE_URL が直接設定されている場合はそのまま使う（AWS RDS など）
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// 個別環境変数からDSNを構築（ローカル開発用）
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "game_score")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
