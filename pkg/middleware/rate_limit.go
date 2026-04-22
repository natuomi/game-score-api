package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipBucket はIPアドレスごとのリクエスト記録
type ipBucket struct {
	count    int
	resetAt  time.Time
	mu       sync.Mutex
}

var (
	buckets   = sync.Map{}
	limit     = 10              // 最大リクエスト数
	windowDur = time.Minute     // 計測ウィンドウ（1分）
)

// RateLimit は指定ウィンドウ内のリクエスト数を制限するミドルウェア
// 同一IPから1分間に10回を超えるとアクセスを拒否する（ブルートフォース対策）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		val, _ := buckets.LoadOrStore(ip, &ipBucket{resetAt: time.Now().Add(windowDur)})
		bucket := val.(*ipBucket)

		bucket.mu.Lock()
		defer bucket.mu.Unlock()

		// ウィンドウが過ぎていればリセット
		if time.Now().After(bucket.resetAt) {
			bucket.count = 0
			bucket.resetAt = time.Now().Add(windowDur)
		}

		bucket.count++
		if bucket.count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
