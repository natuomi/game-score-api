package main

import (
	"log"
	"os"

	"game-score-api/internal/handler"
	"game-score-api/internal/repository"
	"game-score-api/internal/service"
	"game-score-api/pkg/database"
	"game-score-api/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env 読み込み（本番環境では環境変数を直接使用）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// ── DB接続 ───────────────────────────────────────
	db, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	// ── 依存関係の組み立て（Repository → Service → Handler）───
	userRepo  := repository.NewUserRepository(db)
	scoreRepo := repository.NewScoreRepository(db)

	authSvc  := service.NewAuthService(userRepo)
	scoreSvc := service.NewScoreService(scoreRepo, rdb)

	authH   := handler.NewAuthHandler(authSvc)
	scoreH  := handler.NewScoreHandler(scoreSvc)
	playerH := handler.NewPlayerHandler(userRepo)

	// ── ルーティング ──────────────────────────────────
	r := gin.Default()

	// CORS設定（Vercelフロントエンドからのアクセスを許可）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://game-score-frontend-zeta.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// ヘルスチェック（ECS の ALB ヘルスチェック用）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "1.0.0"})
	})

	v1 := r.Group("/api/v1")
	{
		// 認証不要
		auth := v1.Group("/auth")
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)

		v1.GET("/rankings", scoreH.GetRankings)
		v1.GET("/players", playerH.GetPlayers)
		v1.GET("/players/:id", playerH.GetPlayer)

		// 認証必要（JWTミドルウェア）
		secured := v1.Group("/")
		secured.Use(middleware.JWTAuth())
		{
			secured.POST("/scores", scoreH.PostScore)
			secured.GET("/scores/me", scoreH.GetMyScores)
		}
	}

	// ── サーバー起動 ──────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
