// Package infrastructure ユーザーインフラストラクチャ層テスト
package infrastructure

import (
	"database/sql"
	"testing"
	"time"

	"shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
)

func TestUserModel_ToDomain(t *testing.T) {
	now := time.Now()

	t.Run("全フィールドが設定されている場合", func(t *testing.T) {
		orgID := uuid.New()
		loginAt := now.Add(-time.Hour)

		model := &UserModel{
			ID:             uuid.New(),
			OrganizationID: uuid.NullUUID{UUID: orgID, Valid: true},
			Email:          "test@example.com",
			PasswordHash:   "hashed_password",
			FirstName:      "太郎",
			LastName:       "田中",
			Role:           "admin",
			IsActive:       true,
			LastLoginAt:    sql.NullTime{Time: loginAt, Valid: true},
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		user := model.ToDomain()

		if user.ID != sharedDomain.ID(model.ID) {
			t.Error("ID should match")
		}
		if user.OrganizationID == nil {
			t.Fatal("OrganizationID should not be nil")
		}
		if *user.OrganizationID != sharedDomain.ID(orgID) {
			t.Error("OrganizationID should match")
		}
		if user.Email != model.Email {
			t.Errorf("expected email %s but got %s", model.Email, user.Email)
		}
		if user.PasswordHash != model.PasswordHash {
			t.Error("PasswordHash should match")
		}
		if user.FirstName != model.FirstName {
			t.Errorf("expected first name %s but got %s", model.FirstName, user.FirstName)
		}
		if user.LastName != model.LastName {
			t.Errorf("expected last name %s but got %s", model.LastName, user.LastName)
		}
		if user.Role != domain.RoleAdmin {
			t.Errorf("expected role admin but got %s", user.Role)
		}
		if !user.IsActive {
			t.Error("IsActive should be true")
		}
		if user.LastLoginAt == nil {
			t.Fatal("LastLoginAt should not be nil")
		}
		if !user.LastLoginAt.Equal(loginAt) {
			t.Error("LastLoginAt should match")
		}
	})

	t.Run("OrganizationIDがnullの場合", func(t *testing.T) {
		model := &UserModel{
			ID:             uuid.New(),
			OrganizationID: uuid.NullUUID{Valid: false},
			Email:          "test@example.com",
			PasswordHash:   "hashed_password",
			FirstName:      "太郎",
			LastName:       "田中",
			Role:           "user",
			IsActive:       true,
			LastLoginAt:    sql.NullTime{Valid: false},
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		user := model.ToDomain()

		if user.OrganizationID != nil {
			t.Error("OrganizationID should be nil")
		}
		if user.LastLoginAt != nil {
			t.Error("LastLoginAt should be nil")
		}
	})
}

func TestUserModelFromDomain(t *testing.T) {
	now := time.Now()

	t.Run("全フィールドが設定されている場合", func(t *testing.T) {
		orgID := sharedDomain.NewID()
		loginAt := now.Add(-time.Hour)

		user := &domain.User{
			ID:             sharedDomain.NewID(),
			OrganizationID: &orgID,
			Email:          "test@example.com",
			PasswordHash:   "hashed_password",
			FirstName:      "太郎",
			LastName:       "田中",
			Role:           domain.RoleAdmin,
			IsActive:       true,
			LastLoginAt:    &loginAt,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		model := UserModelFromDomain(user)

		if model.ID != uuid.UUID(user.ID) {
			t.Error("ID should match")
		}
		if !model.OrganizationID.Valid {
			t.Error("OrganizationID should be valid")
		}
		if model.OrganizationID.UUID != uuid.UUID(orgID) {
			t.Error("OrganizationID should match")
		}
		if model.Email != user.Email {
			t.Errorf("expected email %s but got %s", user.Email, model.Email)
		}
		if model.Role != "admin" {
			t.Errorf("expected role admin but got %s", model.Role)
		}
		if !model.LastLoginAt.Valid {
			t.Error("LastLoginAt should be valid")
		}
	})

	t.Run("OrganizationIDがnilの場合", func(t *testing.T) {
		user := &domain.User{
			ID:           sharedDomain.NewID(),
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
			FirstName:    "太郎",
			LastName:     "田中",
			Role:         domain.RoleUser,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		model := UserModelFromDomain(user)

		if model.OrganizationID.Valid {
			t.Error("OrganizationID should not be valid")
		}
		if model.LastLoginAt.Valid {
			t.Error("LastLoginAt should not be valid")
		}
	})
}

func TestRefreshTokenModel_ToDomain(t *testing.T) {
	now := time.Now()

	t.Run("正常な変換", func(t *testing.T) {
		model := &RefreshTokenModel{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			TokenHash: "test_token_hash",
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		}

		token := model.ToDomain()

		if token.ID != sharedDomain.ID(model.ID) {
			t.Error("ID should match")
		}
		if token.UserID != sharedDomain.ID(model.UserID) {
			t.Error("UserID should match")
		}
		if token.TokenHash != model.TokenHash {
			t.Errorf("expected token hash %s but got %s", model.TokenHash, token.TokenHash)
		}
		if !token.ExpiresAt.Equal(model.ExpiresAt) {
			t.Error("ExpiresAt should match")
		}
		if !token.CreatedAt.Equal(model.CreatedAt) {
			t.Error("CreatedAt should match")
		}
	})
}

func TestRefreshTokenModelFromDomain(t *testing.T) {
	now := time.Now()

	t.Run("正常な変換", func(t *testing.T) {
		token := &domain.RefreshToken{
			ID:        sharedDomain.NewID(),
			UserID:    sharedDomain.NewID(),
			TokenHash: "test_token_hash",
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		}

		model := RefreshTokenModelFromDomain(token)

		if model.ID != uuid.UUID(token.ID) {
			t.Error("ID should match")
		}
		if model.UserID != uuid.UUID(token.UserID) {
			t.Error("UserID should match")
		}
		if model.TokenHash != token.TokenHash {
			t.Errorf("expected token hash %s but got %s", token.TokenHash, model.TokenHash)
		}
		if !model.ExpiresAt.Equal(token.ExpiresAt) {
			t.Error("ExpiresAt should match")
		}
		if !model.CreatedAt.Equal(token.CreatedAt) {
			t.Error("CreatedAt should match")
		}
	})
}

// 境界値テスト

func TestUserModel_ToDomain_BoundaryValues(t *testing.T) {
	t.Run("空文字列のフィールド", func(t *testing.T) {
		model := &UserModel{
			ID:           uuid.New(),
			Email:        "",
			PasswordHash: "",
			FirstName:    "",
			LastName:     "",
			Role:         "",
			IsActive:     false,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
		}

		user := model.ToDomain()

		if user.Email != "" {
			t.Error("Email should be empty")
		}
		if user.FirstName != "" {
			t.Error("FirstName should be empty")
		}
		if user.IsActive {
			t.Error("IsActive should be false")
		}
	})

	t.Run("ゼロ値の時刻", func(t *testing.T) {
		model := &UserModel{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Role:      "user",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		user := model.ToDomain()

		if !user.CreatedAt.IsZero() {
			t.Error("CreatedAt should be zero time")
		}
		if !user.UpdatedAt.IsZero() {
			t.Error("UpdatedAt should be zero time")
		}
	})
}
