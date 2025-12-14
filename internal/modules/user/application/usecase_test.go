// Package application ユーザーユースケーステスト
package application

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"golang.org/x/crypto/bcrypt"
)

// モックリポジトリ

type mockUserRepository struct {
	users        map[sharedDomain.ID]*domain.User
	usersByEmail map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:        make(map[sharedDomain.ID]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) FindByID(_ context.Context, id sharedDomain.ID) (*domain.User, error) {
	return m.users[id], nil
}

func (m *mockUserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	return m.usersByEmail[email], nil
}

func (m *mockUserRepository) FindByRole(_ context.Context, role domain.UserRole) ([]domain.User, error) {
	var result []domain.User
	for _, user := range m.users {
		if user.Role == role {
			result = append(result, *user)
		}
	}
	return result, nil
}

func (m *mockUserRepository) FindAll(_ context.Context) ([]domain.User, error) {
	result := make([]domain.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, *user)
	}
	return result, nil
}

func (m *mockUserRepository) Save(_ context.Context, user *domain.User) error {
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

func (m *mockUserRepository) FindByOrganizationID(_ context.Context, _ sharedDomain.ID) ([]domain.User, error) {
	return nil, nil
}

// モックリフレッシュトークンリポジトリ

type mockRefreshTokenRepository struct {
	tokens map[sharedDomain.ID]*domain.RefreshToken
}

func newMockRefreshTokenRepository() *mockRefreshTokenRepository {
	return &mockRefreshTokenRepository{
		tokens: make(map[sharedDomain.ID]*domain.RefreshToken),
	}
}

func (m *mockRefreshTokenRepository) FindByID(_ context.Context, id sharedDomain.ID) (*domain.RefreshToken, error) {
	return m.tokens[id], nil
}

func (m *mockRefreshTokenRepository) FindByTokenHash(_ context.Context, _ string) (*domain.RefreshToken, error) {
	return nil, nil
}

func (m *mockRefreshTokenRepository) Save(_ context.Context, token *domain.RefreshToken) error {
	m.tokens[token.ID] = token
	return nil
}

func (m *mockRefreshTokenRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	delete(m.tokens, id)
	return nil
}

func (m *mockRefreshTokenRepository) DeleteByUserID(_ context.Context, userID sharedDomain.ID) error {
	for id, token := range m.tokens {
		if token.UserID == userID {
			delete(m.tokens, id)
		}
	}
	return nil
}

func (m *mockRefreshTokenRepository) DeleteExpired(_ context.Context) error {
	for id, token := range m.tokens {
		if token.IsExpired() {
			delete(m.tokens, id)
		}
	}
	return nil
}

func (m *mockRefreshTokenRepository) FindByUserID(_ context.Context, userID sharedDomain.ID) ([]domain.RefreshToken, error) {
	var result []domain.RefreshToken
	for _, token := range m.tokens {
		if token.UserID == userID {
			result = append(result, *token)
		}
	}
	return result, nil
}

// テスト

func TestNewUserUseCase(t *testing.T) {
	userRepo := newMockUserRepository()
	tokenRepo := newMockRefreshTokenRepository()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := NewUserUseCase(userRepo, tokenRepo, logger)

	if useCase == nil {
		t.Fatal("NewUserUseCase returned nil")
	}
	if useCase.userRepo == nil {
		t.Error("userRepo should not be nil")
	}
	if useCase.tokenRepo == nil {
		t.Error("tokenRepo should not be nil")
	}
	if useCase.logger == nil {
		t.Error("logger should not be nil")
	}
}

func TestUserUseCase_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name      string
		input     *CreateUserInput
		setupUser func(*mockUserRepository)
		wantErr   bool
		errCode   string
	}{
		{
			name: "正常なユーザー作成",
			input: &CreateUserInput{
				Email:     "new@example.com",
				Password:  "password123",
				FirstName: "新規",
				LastName:  "ユーザー",
				Role:      "admin",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   false,
		},
		{
			name: "メールアドレス重複",
			input: &CreateUserInput{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "重複",
				LastName:  "ユーザー",
				Role:      "user",
			},
			setupUser: func(repo *mockUserRepository) {
				user := &domain.User{
					ID:    sharedDomain.NewID(),
					Email: "existing@example.com",
				}
				_ = repo.Save(context.Background(), user)
			},
			wantErr: true,
			errCode: sharedDomain.ErrCodeConflict,
		},
		{
			name: "メールアドレス未入力",
			input: &CreateUserInput{
				Email:     "",
				Password:  "password123",
				FirstName: "テスト",
				LastName:  "ユーザー",
				Role:      "user",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   true,
			errCode:   sharedDomain.ErrCodeValidation,
		},
		{
			name: "パスワード未入力",
			input: &CreateUserInput{
				Email:     "test@example.com",
				Password:  "",
				FirstName: "テスト",
				LastName:  "ユーザー",
				Role:      "user",
			},
			setupUser: func(_ *mockUserRepository) {},
			wantErr:   true,
			errCode:   sharedDomain.ErrCodeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			tokenRepo := newMockRefreshTokenRepository()
			tt.setupUser(userRepo)

			useCase := NewUserUseCase(userRepo, tokenRepo, logger)
			output, err := useCase.Create(context.Background(), tt.input)

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
					if output.Email != tt.input.Email {
						t.Errorf("expected email %s but got %s", tt.input.Email, output.Email)
					}
				}
			}
		})
	}
}

func TestUserUseCase_GetByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("存在するユーザー取得", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		userID := sharedDomain.NewID()
		user := &domain.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "テスト",
			LastName:  "ユーザー",
		}
		_ = userRepo.Save(context.Background(), user)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		output, err := useCase.GetByID(context.Background(), userID.String())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.ID != userID.String() {
			t.Errorf("expected ID %s but got %s", userID.String(), output.ID)
		}
	})

	t.Run("存在しないユーザー取得", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		_, err := useCase.GetByID(context.Background(), sharedDomain.NewID().String())

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})

	t.Run("不正なID形式", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		_, err := useCase.GetByID(context.Background(), "invalid-id")

		if err == nil {
			t.Error("expected error but got nil")
		}
		if domainErr, ok := err.(*sharedDomain.DomainError); ok {
			if domainErr.Code != sharedDomain.ErrCodeValidation {
				t.Errorf("expected ErrCodeValidation but got %s", domainErr.Code)
			}
		}
	})
}

func TestUserUseCase_List(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("ユーザー一覧取得_空", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		output, err := useCase.List(context.Background())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.Total != 0 {
			t.Errorf("expected total 0 but got %d", output.Total)
		}
	})

	t.Run("ユーザー一覧取得_複数", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		// ユーザーを追加
		for i := 0; i < 3; i++ {
			user := &domain.User{
				ID:    sharedDomain.NewID(),
				Email: "user" + string(rune('1'+i)) + "@example.com",
			}
			_ = userRepo.Save(context.Background(), user)
		}

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		output, err := useCase.List(context.Background())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.Total != 3 {
			t.Errorf("expected total 3 but got %d", output.Total)
		}
	})
}

func TestUserUseCase_Update(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なユーザー更新", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		userID := sharedDomain.NewID()
		user := &domain.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "テスト",
			LastName:  "ユーザー",
			Role:      domain.RoleUser,
			IsActive:  true,
		}
		_ = userRepo.Save(context.Background(), user)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		output, err := useCase.Update(context.Background(), &UpdateUserInput{
			ID:        userID.String(),
			Email:     "updated@example.com",
			FirstName: "更新",
			LastName:  "ユーザー",
			Role:      "manager",
			IsActive:  true,
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.Email != "updated@example.com" {
			t.Errorf("expected email updated@example.com but got %s", output.Email)
		}
	})

	t.Run("存在しないユーザー更新", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		_, err := useCase.Update(context.Background(), &UpdateUserInput{
			ID:        sharedDomain.NewID().String(),
			Email:     "test@example.com",
			FirstName: "テスト",
			LastName:  "ユーザー",
			Role:      "user",
			IsActive:  true,
		})

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})
}

func TestUserUseCase_ChangePassword(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なパスワード変更", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		userID := sharedDomain.NewID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.MinCost)
		user := &domain.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			IsActive:     true,
		}
		_ = userRepo.Save(context.Background(), user)

		// リフレッシュトークンを追加
		token := &domain.RefreshToken{
			ID:        sharedDomain.NewID(),
			UserID:    userID,
			TokenHash: "test-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		_ = tokenRepo.Save(context.Background(), token)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		err := useCase.ChangePassword(context.Background(), &ChangePasswordInput{
			UserID:          userID.String(),
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword",
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// 新しいパスワードで検証
		updatedUser, _ := userRepo.FindByID(context.Background(), userID)
		if bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("newpassword")) != nil {
			t.Error("password should be updated")
		}
	})

	t.Run("現在のパスワード不一致", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		userID := sharedDomain.NewID()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.MinCost)
		user := &domain.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			IsActive:     true,
		}
		_ = userRepo.Save(context.Background(), user)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		err := useCase.ChangePassword(context.Background(), &ChangePasswordInput{
			UserID:          userID.String(),
			CurrentPassword: "wrongpassword",
			NewPassword:     "newpassword",
		})

		if err == nil {
			t.Error("expected error but got nil")
		}
		if domainErr, ok := err.(*sharedDomain.DomainError); ok {
			if domainErr.Code != sharedDomain.ErrCodeUnauthorized {
				t.Errorf("expected ErrCodeUnauthorized but got %s", domainErr.Code)
			}
		}
	})
}

func TestUserUseCase_Delete(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なユーザー削除", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		// 管理者を2人作成
		adminID1 := sharedDomain.NewID()
		admin1 := &domain.User{
			ID:       adminID1,
			Email:    "admin1@example.com",
			Role:     domain.RoleAdmin,
			IsActive: true,
		}
		_ = userRepo.Save(context.Background(), admin1)

		adminID2 := sharedDomain.NewID()
		admin2 := &domain.User{
			ID:       adminID2,
			Email:    "admin2@example.com",
			Role:     domain.RoleAdmin,
			IsActive: true,
		}
		_ = userRepo.Save(context.Background(), admin2)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		err := useCase.Delete(context.Background(), adminID1.String())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// 削除確認
		deleted, _ := userRepo.FindByID(context.Background(), adminID1)
		if deleted != nil {
			t.Error("user should be deleted")
		}
	})

	t.Run("最後の管理者は削除不可", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		// 管理者を1人だけ作成
		adminID := sharedDomain.NewID()
		admin := &domain.User{
			ID:       adminID,
			Email:    "admin@example.com",
			Role:     domain.RoleAdmin,
			IsActive: true,
		}
		_ = userRepo.Save(context.Background(), admin)

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		err := useCase.Delete(context.Background(), adminID.String())

		if err == nil {
			t.Error("expected error but got nil")
		}
		if domainErr, ok := err.(*sharedDomain.DomainError); ok {
			if domainErr.Code != sharedDomain.ErrCodeValidation {
				t.Errorf("expected ErrCodeValidation but got %s", domainErr.Code)
			}
		}
	})

	t.Run("存在しないユーザー削除", func(t *testing.T) {
		userRepo := newMockUserRepository()
		tokenRepo := newMockRefreshTokenRepository()

		useCase := NewUserUseCase(userRepo, tokenRepo, logger)
		err := useCase.Delete(context.Background(), sharedDomain.NewID().String())

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})
}
