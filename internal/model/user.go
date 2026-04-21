package model

import "time"

// User はプレイヤー情報を表す
type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Password  string    `json:"-" db:"password"` // レスポンスには含めない
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest はユーザー登録リクエスト
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse はログインレスポンス（JWTトークンを含む）
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
