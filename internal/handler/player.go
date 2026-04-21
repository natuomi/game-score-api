package handler

import (
	"net/http"

	"game-score-api/internal/repository"

	"github.com/gin-gonic/gin"
)

// PlayerHandler はプレイヤー情報のHTTPハンドラーを担当する
type PlayerHandler struct {
	userRepo *repository.UserRepository
}

func NewPlayerHandler(userRepo *repository.UserRepository) *PlayerHandler {
	return &PlayerHandler{userRepo: userRepo}
}

// GetPlayers はプレイヤー一覧を返す（認証不要）
// GET /api/v1/players
func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	users, err := h.userRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": users})
}

// GetPlayer はプレイヤー詳細を返す（認証不要）
// GET /api/v1/players/:id
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
