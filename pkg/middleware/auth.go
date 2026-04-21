package middleware

import (
	"net/http"
	"strings"

	"game-score-api/pkg/auth"

	"github.com/gin-gonic/gin"
)

// JWTAuth は Authorization: Bearer <token> を検証するミドルウェア
// 使い方:
//
//	authorized := r.Group("/")
//	authorized.Use(middleware.JWTAuth())
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// "Bearer <token>" の形式チェック
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// ハンドラーから参照できるようにコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)
		c.Next()
	}
}
