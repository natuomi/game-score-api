package handler

import (
	"net/http"

	"game-score-api/internal/model"
	"game-score-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler はユーザー認証のHTTPハンドラーを担当する
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register はユーザー登録を処理する
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		// メールアドレス重複はUNIQUE制約エラーとして返ってくる
		c.JSON(http.StatusConflict, gin.H{"error": "email or name already in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login はログインを処理し、JWTトークンを返す
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
