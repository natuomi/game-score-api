package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient は Redis クライアントを作成する
// 使い方:
//
//	rdb, err := database.NewRedisClient()
//	defer rdb.Close()
func NewRedisClient() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s",
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"), // 本番では必ず設定
		DB:       0,
	})

	// 接続確認
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
