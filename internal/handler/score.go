package handler

import (
	"net/http"

	"game-score-api/internal/model"
	"game-score-api/internal/service"

	"github.com/gin-gonic/gin"
)

// ScoreHandler はスコア操作のHTTPハンドラーを担当する
type ScoreHandler struct {
	scoreService *service.ScoreService
}

func NewScoreHandler(scoreService *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{scoreService: scoreService}
}

// PostScore はスコアを登録する（要認証）
// POST /api/v1/scores
func (h *ScoreHandler) PostScore(c *gin.Context) {
	var req model.PostScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JWTミドルウェアがセットしたuser_idを取得
	userID := c.GetString("user_id")

	score, err := h.scoreService.PostScore(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, score)
}

// GetMyScores は自分のスコア履歴を返す（要認証）
// GET /api/v1/scores/me
func (h *ScoreHandler) GetMyScores(c *gin.Context) {
	userID := c.GetString("user_id")

	scores, err := h.scoreService.GetMyScores(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"scores": scores})
}

// GetRankings はランキングを返す（認証不要・Redisキャッシュ対応）
// GET /api/v1/rankings
func (h *ScoreHandler) GetRankings(c *gin.Context) {
	resp, err := h.scoreService.GetRankings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
