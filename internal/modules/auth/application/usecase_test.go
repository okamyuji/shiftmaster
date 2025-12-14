// Package application 認証ユースケーステスト
package application

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	authDomain "shiftmaster/internal/modules/auth/domain"
	userDomain "shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"golang.org/x/crypto/bcrypt"
)

// モックリポジトリ

type mockUserRepository struct {
	users        map[sharedDomain.ID]*userDomain.User
	usersByEmail map[string]*userDomain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:        make(map[sharedDomain.ID]*userDomain.User),
		usersByEmail: make(map[string]*userDomain.User),
	}
}

func (m *mockUserRepository) FindByID(_ context.Context, id sharedDomain.ID) (*userDomain.User, error) {
	return m.users[id], nil
}

func (m *mockUserRepository) FindByEmail(_ context.Context, email string) (*userDomain.User, error) {
	return m.usersByEmail[email], nil
}

func (m *mockUserRepository) FindByRole(_ context.Context, _ userDomain.UserRole) ([]userDomain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindAll(_ context.Context) ([]userDomain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Save(_ context.Context, user *userDomain.User) error {
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	if user, ok := m.users[id]; ok {
		delete(m.usersByEmail, user.Email)
		delete(m.users, id)
	}
	return nil
}

func (m *mockUserRepository) UpdateLastLogin(_ context.Context, _ sharedDomain.ID) error {
	return nil
}

func (m *mockUserRepository) FindByOrganizationID(_ context.Context, _ sharedDomain.ID) ([]userDomain.User, error) {
	return nil, nil
}

// モックリフレッシュトークンリポジトリ

type mockRefreshTokenRepository struct {
	tokens       map[sharedDomain.ID]*userDomain.RefreshToken
	tokensByHash map[string]*userDomain.RefreshToken
}

func newMockRefreshTokenRepository() *mockRefreshTokenRepository {
	return &mockRefreshTokenRepository{
		tokens:       make(map[sharedDomain.ID]*userDomain.RefreshToken),
		tokensByHash: make(map[string]*userDomain.RefreshToken),
	}
}

func (m *mockRefreshTokenRepository) FindByID(_ context.Context, id sharedDomain.ID) (*userDomain.RefreshToken, error) {
	return m.tokens[id], nil
}

func (m *mockRefreshTokenRepository) FindByTokenHash(_ context.Context, tokenHash string) (*userDomain.RefreshToken, error) {
	return m.tokensByHash[tokenHash], nil
}

func (m *mockRefreshTokenRepository) Save(_ context.Context, token *userDomain.RefreshToken) error {
	m.tokens[token.ID] = token
	m.tokensByHash[token.TokenHash] = token
	return nil
}

func (m *mockRefreshTokenRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	if token, ok := m.tokens[id]; ok {
		delete(m.tokensByHash, token.TokenHash)
		delete(m.tokens, id)
	}
	return nil
}

func (m *mockRefreshTokenRepository) DeleteByUserID(_ context.Context, userID sharedDomain.ID) error {
	for id, token := range m.tokens {
		if token.UserID == userID {
			delete(m.tokensByHash, token.TokenHash)
			delete(m.tokens, id)
		}
	}
	return nil
}

func (m *mockRefreshTokenRepository) DeleteExpired(_ context.Context) error {
	for id, token := range m.tokens {
		if token.IsExpired() {
			delete(m.tokensByHash, token.TokenHash)
			delete(m.tokens, id)
		}
	}
	return nil
}

func (m *mockRefreshTokenRepository) FindByUserID(_ context.Context, userID sharedDomain.ID) ([]userDomain.RefreshToken, error) {
	var result []userDomain.RefreshToken
	for _, token := range m.tokens {
		if token.UserID == userID {
			result = append(result, *token)
		}
	}
	return result, nil
}

// モックトークンサービス

type mockTokenService struct {
	accessTokenExpiresAt  time.Time
	refreshTokenExpiresAt time.Time
}

func newMockTokenService() *mockTokenService {
	return &mockTokenService{
		accessTokenExpiresAt:  time.Now().Add(15 * time.Minute),
		refreshTokenExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
}

func (m *mockTokenService) GenerateTokenPair(_ *authDomain.Claims) (*authDomain.TokenPair, error) {
	return &authDomain.TokenPair{
		AccessToken:           "mock-access-token",
		RefreshToken:          "mock-refresh-token",
		AccessTokenExpiresAt:  m.accessTokenExpiresAt,
		RefreshTokenExpiresAt: m.refreshTokenExpiresAt,
	}, nil
}

func (m *mockTokenService) ValidateAccessToken(token string) (*authDomain.Claims, error) {
	if token == "valid-access-token" {
		return &authDomain.Claims{
			UserID:   sharedDomain.NewID(),
			Email:    "test@example.com",
			Role:     "admin",
			IssuedAt: time.Now(),
		}, nil
	}
	return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "無効なトークン")
}

func (m *mockTokenService) ValidateRefreshToken(token string) (*authDomain.Claims, error) {
	if token == "valid-refresh-token" {
		return &authDomain.Claims{
			UserID:   sharedDomain.NewID(),
			Email:    "test@example.com",
			Role:     "admin",
			IssuedAt: time.Now(),
		}, nil
	}
	return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "無効なトークン")
}

func (m *mockTokenService) HashToken(token string) string {
	return "hashed-" + token
}

// テスト

func TestNewAuthUseCase(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockRefreshTokenRepository()
	tokenService := newMockTokenService()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)

	if useCase == nil {
		t.Fatal("NewAuthUseCase returned nil")
	}
	if useCase.userRepo == nil {
		t.Error("userRepo should not be nil")
	}
	if useCase.tokenRepo == nil {
		t.Error("tokenRepo should not be nil")
	}
	if useCase.tokenService == nil {
		t.Error("tokenService should not be nil")
	}
	if useCase.logger == nil {
		t.Error("logger should not be nil")
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name      string
		input     *LoginInput
		setupUser func(*mockUserRepository)
		wantErr   bool
		errCode   string
	}{
		{
			name: "正常なログイン",
			input: &LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupUser: func(repo *mockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
				user := &userDomain.User{
					ID:           sharedDomain.NewID(),
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "テスト",
					LastName:     "ユーザー",
					Role:         userDomain.RoleAdmin,
					IsActive:     true,
				}
				_ = repo.Save(context.Background(), user)
			},
			wantErr: false,
		},
		{
			name: "メールアドレス未入力",
			input: &LoginInput{
				Email:    "",
				Password: "password123",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   true,
			errCode:   sharedDomain.ErrCodeValidation,
		},
		{
			name: "パスワード未入力",
			input: &LoginInput{
				Email:    "test@example.com",
				Password: "",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   true,
			errCode:   sharedDomain.ErrCodeValidation,
		},
		{
			name: "ユーザー未検出",
			input: &LoginInput{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   true,
			errCode:   sharedDomain.ErrCodeUnauthorized,
		},
		{
			name: "無効ユーザー",
			input: &LoginInput{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			setupUser: func(repo *mockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
				user := &userDomain.User{
					ID:           sharedDomain.NewID(),
					Email:        "inactive@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "無効",
					LastName:     "ユーザー",
					Role:         userDomain.RoleUser,
					IsActive:     false, // 無効
				}
				_ = repo.Save(context.Background(), user)
			},
			wantErr: true,
			errCode: sharedDomain.ErrCodeUnauthorized,
		},
		{
			name: "パスワード不一致",
			input: &LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupUser: func(repo *mockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
				user := &userDomain.User{
					ID:           sharedDomain.NewID(),
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "テスト",
					LastName:     "ユーザー",
					Role:         userDomain.RoleAdmin,
					IsActive:     true,
				}
				_ = repo.Save(context.Background(), user)
			},
			wantErr: true,
			errCode: sharedDomain.ErrCodeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			tokenRepo := newMockRefreshTokenRepository()
			tokenService := newMockTokenService()
			tt.setupUser(userRepo)

			useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
			output, err := useCase.Login(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if domainErr, ok := err.(*sharedDomain.DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("expected error code %s but got %s", tt.errCode, domainErr.Code)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if output == nil {
					t.Error("expected output but got nil")
				} else {
					if output.AccessToken == "" {
						t.Error("access token should not be empty")
					}
					if output.RefreshToken == "" {
						t.Error("refresh token should not be empty")
					}
					if output.TokenType != "Bearer" {
						t.Errorf("expected token type Bearer but got %s", output.TokenType)
					}
				}
			}
		})
	}
}

func TestAuthUseCase_Logout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なログアウト", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()
		tokenService := newMockTokenService()

		// トークンを保存
		token := &userDomain.RefreshToken{
			ID:        sharedDomain.NewID(),
			UserID:    sharedDomain.NewID(),
			TokenHash: "hashed-test-token",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		_ = tokenRepo.Save(context.Background(), token)

		useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
		err := useCase.Logout(context.Background(), "test-token")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// トークンが削除されているか確認
		found, _ := tokenRepo.FindByTokenHash(context.Background(), "hashed-test-token")
		if found != nil {
			t.Error("token should be deleted")
		}
	})

	t.Run("空のトークンでログアウト", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()
		tokenService := newMockTokenService()

		useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
		err := useCase.Logout(context.Background(), "")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("存在しないトークンでログアウト", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()
		tokenService := newMockTokenService()

		useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
		err := useCase.Logout(context.Background(), "nonexistent-token")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("有効なトークン_アクティブユーザー", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()
		tokenService := &mockTokenServiceWithUserID{}

		userID := sharedDomain.NewID()
		tokenService.userID = userID

		user := &userDomain.User{
			ID:       userID,
			Email:    "test@example.com",
			IsActive: true,
		}
		_ = userRepo.Save(context.Background(), user)

		useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
		claims, err := useCase.ValidateToken(context.Background(), "valid-access-token")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if claims == nil {
			t.Error("expected claims but got nil")
		}
	})

	t.Run("無効なトークン", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()
		tokenService := newMockTokenService()

		useCase := NewAuthUseCase(userRepo, tokenRepo, tokenService, logger)
		_, err := useCase.ValidateToken(context.Background(), "invalid-access-token")

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

// mockTokenServiceWithUserID ユーザーIDを設定可能なモック
type mockTokenServiceWithUserID struct {
	userID sharedDomain.ID
}

func (m *mockTokenServiceWithUserID) GenerateTokenPair(_ *authDomain.Claims) (*authDomain.TokenPair, error) {
	return &authDomain.TokenPair{
		AccessToken:           "mock-access-token",
		RefreshToken:          "mock-refresh-token",
		AccessTokenExpiresAt:  time.Now().Add(15 * time.Minute),
		RefreshTokenExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

func (m *mockTokenServiceWithUserID) ValidateAccessToken(_ string) (*authDomain.Claims, error) {
	return &authDomain.Claims{
		UserID:   m.userID,
		Email:    "test@example.com",
		Role:     "admin",
		IssuedAt: time.Now(),
	}, nil
}

func (m *mockTokenServiceWithUserID) ValidateRefreshToken(_ string) (*authDomain.Claims, error) {
	return &authDomain.Claims{
		UserID:   m.userID,
		Email:    "test@example.com",
		Role:     "admin",
		IssuedAt: time.Now(),
	}, nil
}

func (m *mockTokenServiceWithUserID) HashToken(token string) string {
	return "hashed-" + token
}
