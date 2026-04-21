package service

import (
	"context"
	"errors"
	"fmt"

	"game-score-api/internal/model"
	"game-score-api/internal/repository"
	"game-score-api/pkg/auth"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService はユーザー認証のビジネスロジックを担当する
type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register はユーザー登録を行う
// パスワードはbcryptでハッシュ化してDBに保存する
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	// パスワードをbcryptでハッシュ化（コスト12は安全性と速度のバランス）
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, req.Name, req.Email, string(hashed))
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}

// Login はメール+パスワードを検証してJWTトークンを返す
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	// メールアドレスでユーザー検索
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// pgx.ErrNoRows でもエラーメッセージを統一する（ユーザー名列挙攻撃対策）
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// JWTトークン生成
	token, err := auth.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
