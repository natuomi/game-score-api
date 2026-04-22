package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"game-score-api/internal/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// ── モック定義 ────────────────────────────────────────────────────────

// mockUserRepo は UserRepositoryInterface の手書きモック実装
type mockUserRepo struct {
	// CreateFn をセットすると Create 呼び出し時にその関数が実行される
	CreateFn      func(ctx context.Context, name, email, hashedPassword string) (*model.User, error)
	FindByEmailFn func(ctx context.Context, email string) (*model.User, error)
	FindByIDFn    func(ctx context.Context, id string) (*model.User, error)
	FindAllFn     func(ctx context.Context) ([]model.User, error)

	// 呼び出し確認用カウンタ
	CreateCalled      int
	FindByEmailCalled int
}

func (m *mockUserRepo) Create(ctx context.Context, name, email, hashedPassword string) (*model.User, error) {
	m.CreateCalled++
	if m.CreateFn != nil {
		return m.CreateFn(ctx, name, email, hashedPassword)
	}
	return nil, errors.New("CreateFn not set")
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	m.FindByEmailCalled++
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, errors.New("FindByEmailFn not set")
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, errors.New("FindByIDFn not set")
}

func (m *mockUserRepo) FindAll(ctx context.Context) ([]model.User, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, errors.New("FindAllFn not set")
}

// ── テストヘルパー ────────────────────────────────────────────────────

// newTestUser はテスト用の User を生成する
func newTestUser(id, name, email, hashedPwd string) *model.User {
	return &model.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  hashedPwd,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ── Register テスト ───────────────────────────────────────────────────

// TestRegister_Success はRegister成功ケースを検証する
// - bcryptハッシュ化されたパスワードでrepo.Createが呼ばれること
// - 返却されたユーザーがrepoの返却値と一致すること
func TestRegister_Success(t *testing.T) {
	t.Parallel()

	expectedUser := newTestUser("uid-1", "Alice", "alice@example.com", "")

	repo := &mockUserRepo{
		CreateFn: func(ctx context.Context, name, email, hashedPassword string) (*model.User, error) {
			// bcryptハッシュが正しく生成されているか検証
			if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("password123")); err != nil {
				t.Errorf("bcrypt hash mismatch: %v", err)
			}
			// 元のパスワードが渡されていないことを確認
			if hashedPassword == "password123" {
				t.Error("password should be hashed, not plaintext")
			}
			if name != "Alice" {
				t.Errorf("expected name=Alice, got %s", name)
			}
			if email != "alice@example.com" {
				t.Errorf("expected email=alice@example.com, got %s", email)
			}
			return expectedUser, nil
		},
	}

	svc := NewAuthService(repo)
	req := model.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	user, err := svc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Register() returned unexpected error: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("expected user.ID=%s, got %s", expectedUser.ID, user.ID)
	}
	if repo.CreateCalled != 1 {
		t.Errorf("expected repo.Create called 1 time, got %d", repo.CreateCalled)
	}
}

// TestRegister_EmailAlreadyExists はメールアドレスが既に存在する場合の失敗ケースを検証する
func TestRegister_EmailAlreadyExists(t *testing.T) {
	t.Parallel()

	repo := &mockUserRepo{
		CreateFn: func(ctx context.Context, name, email, hashedPassword string) (*model.User, error) {
			// PostgreSQL の unique violation を模倣
			return nil, errors.New("create user: ERROR: duplicate key value violates unique constraint")
		},
	}

	svc := NewAuthService(repo)
	req := model.RegisterRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "password123",
	}

	_, err := svc.Register(context.Background(), req)
	if err == nil {
		t.Fatal("Register() should return error when email already exists")
	}
	if repo.CreateCalled != 1 {
		t.Errorf("expected repo.Create called 1 time, got %d", repo.CreateCalled)
	}
}

// ── Login テスト ──────────────────────────────────────────────────────

// TestLogin_Success はLogin成功ケースを検証する
// - JWT が生成されて返ること
// - 返却されたユーザー情報が正しいこと
// Note: t.Setenv と t.Parallel は併用不可のため逐次実行
func TestLogin_Success(t *testing.T) {
	// テスト用JWTのためenv設定（t.Setenvはt.Parallel()と併用不可）
	t.Setenv("JWT_SECRET", "test-secret-for-unit-test")

	const plainPwd = "mypassword"
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash: %v", err)
	}

	storedUser := newTestUser("uid-99", "Carol", "carol@example.com", string(hashed))

	repo := &mockUserRepo{
		FindByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			if email != "carol@example.com" {
				return nil, pgx.ErrNoRows
			}
			return storedUser, nil
		},
	}

	svc := NewAuthService(repo)
	req := model.LoginRequest{
		Email:    "carol@example.com",
		Password: plainPwd,
	}

	resp, err := svc.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("Login() returned unexpected error: %v", err)
	}
	if resp.Token == "" {
		t.Error("Login() should return a non-empty JWT token")
	}
	if resp.User.ID != storedUser.ID {
		t.Errorf("expected user.ID=%s, got %s", storedUser.ID, resp.User.ID)
	}
	if repo.FindByEmailCalled != 1 {
		t.Errorf("expected repo.FindByEmail called 1 time, got %d", repo.FindByEmailCalled)
	}
}

// TestLogin_PasswordMismatch はパスワード不一致の失敗ケースを検証する
// Note: t.Setenv と t.Parallel は併用不可のため逐次実行
func TestLogin_PasswordMismatch(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-unit-test")

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	storedUser := newTestUser("uid-1", "Dave", "dave@example.com", string(hashed))

	repo := &mockUserRepo{
		FindByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return storedUser, nil
		},
	}

	svc := NewAuthService(repo)
	req := model.LoginRequest{
		Email:    "dave@example.com",
		Password: "wrongpassword",
	}

	_, err := svc.Login(context.Background(), req)
	if err == nil {
		t.Fatal("Login() should return error on password mismatch")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}

// TestLogin_UserNotFound はユーザーが存在しない場合の失敗ケースを検証する
// ユーザー名列挙攻撃対策として「invalid credentials」で統一されることを確認する
// Note: t.Setenv と t.Parallel は併用不可のため逐次実行
func TestLogin_UserNotFound(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-unit-test")

	repo := &mockUserRepo{
		FindByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			// pgx.ErrNoRows = ユーザーが存在しない
			return nil, pgx.ErrNoRows
		},
	}

	svc := NewAuthService(repo)
	req := model.LoginRequest{
		Email:    "nobody@example.com",
		Password: "anypassword",
	}

	_, err := svc.Login(context.Background(), req)
	if err == nil {
		t.Fatal("Login() should return error when user not found")
	}
	// ユーザー名列挙攻撃対策：パスワード不一致と同じエラーメッセージになること
	if err.Error() != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%s'", err.Error())
	}
}
