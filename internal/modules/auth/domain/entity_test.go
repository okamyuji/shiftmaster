// Package domain 認証ドメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
)

// ============================================
// TokenPair関連テスト
// ============================================

func TestTokenPair(t *testing.T) {
	t.Run("TokenPair構造体の初期化", func(t *testing.T) {
		now := time.Now()
		accessExpires := now.Add(15 * time.Minute)
		refreshExpires := now.Add(7 * 24 * time.Hour)

		tp := TokenPair{
			AccessToken:           "access_token_string",
			RefreshToken:          "refresh_token_string",
			AccessTokenExpiresAt:  accessExpires,
			RefreshTokenExpiresAt: refreshExpires,
		}

		if tp.AccessToken != "access_token_string" {
			t.Errorf("AccessToken = %v, want %v", tp.AccessToken, "access_token_string")
		}
		if tp.RefreshToken != "refresh_token_string" {
			t.Errorf("RefreshToken = %v, want %v", tp.RefreshToken, "refresh_token_string")
		}
		if !tp.AccessTokenExpiresAt.Equal(accessExpires) {
			t.Errorf("AccessTokenExpiresAt = %v, want %v", tp.AccessTokenExpiresAt, accessExpires)
		}
		if !tp.RefreshTokenExpiresAt.Equal(refreshExpires) {
			t.Errorf("RefreshTokenExpiresAt = %v, want %v", tp.RefreshTokenExpiresAt, refreshExpires)
		}
	})

	t.Run("ゼロ値のTokenPair", func(t *testing.T) {
		var tp TokenPair

		if tp.AccessToken != "" {
			t.Error("AccessToken zero value should be empty string")
		}
		if tp.RefreshToken != "" {
			t.Error("RefreshToken zero value should be empty string")
		}
		if !tp.AccessTokenExpiresAt.IsZero() {
			t.Error("AccessTokenExpiresAt zero value should be zero time")
		}
		if !tp.RefreshTokenExpiresAt.IsZero() {
			t.Error("RefreshTokenExpiresAt zero value should be zero time")
		}
	})
}

// ============================================
// Claims関連テスト
// ============================================

func TestClaims_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "期限切れ_過去の時刻",
			expiresAt: time.Now().Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "有効_未来の時刻",
			expiresAt: time.Now().Add(1 * time.Hour),
			expected:  false,
		},
		{
			name:      "境界値_1秒前",
			expiresAt: time.Now().Add(-1 * time.Second),
			expected:  true,
		},
		{
			name:      "境界値_1秒後",
			expiresAt: time.Now().Add(1 * time.Second),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &Claims{
				ExpiresAt: tt.expiresAt,
			}

			result := claims.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClaims_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{
			name:     "adminロール",
			role:     "admin",
			expected: true,
		},
		{
			name:     "managerロール",
			role:     "manager",
			expected: false,
		},
		{
			name:     "userロール",
			role:     "user",
			expected: false,
		},
		{
			name:     "空のロール",
			role:     "",
			expected: false,
		},
		{
			name:     "大文字ADMIN",
			role:     "ADMIN",
			expected: false, // 大文字小文字は区別される
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &Claims{
				Role: tt.role,
			}

			result := claims.IsAdmin()
			if result != tt.expected {
				t.Errorf("IsAdmin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClaims_IsManager(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{
			name:     "adminロール",
			role:     "admin",
			expected: true,
		},
		{
			name:     "managerロール",
			role:     "manager",
			expected: true,
		},
		{
			name:     "userロール",
			role:     "user",
			expected: false,
		},
		{
			name:     "空のロール",
			role:     "",
			expected: false,
		},
		{
			name:     "大文字MANAGER",
			role:     "MANAGER",
			expected: false, // 大文字小文字は区別される
		},
		{
			name:     "不明なロール",
			role:     "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &Claims{
				Role: tt.role,
			}

			result := claims.IsManager()
			if result != tt.expected {
				t.Errorf("IsManager() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClaims_Structure(t *testing.T) {
	t.Run("Claims構造体の完全な初期化", func(t *testing.T) {
		userID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		issuedAt := time.Now()
		expiresAt := issuedAt.Add(15 * time.Minute)

		claims := &Claims{
			UserID:         userID,
			Email:          "test@example.com",
			Role:           "admin",
			OrganizationID: &orgID,
			IssuedAt:       issuedAt,
			ExpiresAt:      expiresAt,
		}

		if claims.UserID != userID {
			t.Errorf("UserID = %v, want %v", claims.UserID, userID)
		}
		if claims.Email != "test@example.com" {
			t.Errorf("Email = %v, want %v", claims.Email, "test@example.com")
		}
		if claims.Role != "admin" {
			t.Errorf("Role = %v, want %v", claims.Role, "admin")
		}
		if claims.OrganizationID == nil || *claims.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", claims.OrganizationID, orgID)
		}
		if claims.IssuedAt != issuedAt {
			t.Errorf("IssuedAt = %v, want %v", claims.IssuedAt, issuedAt)
		}
		if claims.ExpiresAt != expiresAt {
			t.Errorf("ExpiresAt = %v, want %v", claims.ExpiresAt, expiresAt)
		}
	})

	t.Run("OrganizationIDがnilの場合", func(t *testing.T) {
		claims := &Claims{
			UserID:         sharedDomain.NewID(),
			Email:          "test@example.com",
			Role:           "user",
			OrganizationID: nil,
		}

		if claims.UserID == sharedDomain.ID(uuid.Nil) {
			t.Error("UserID should not be nil")
		}
		if claims.Email != "test@example.com" {
			t.Errorf("Email = %v, want %v", claims.Email, "test@example.com")
		}
		if claims.Role != "user" {
			t.Errorf("Role = %v, want %v", claims.Role, "user")
		}
		if claims.OrganizationID != nil {
			t.Error("OrganizationID should be nil")
		}
		// IssuedAtとExpiresAtは明示的に設定していないためゼロ値
	})
}

// ============================================
// LoginRequest関連テスト
// ============================================

func TestLoginRequest(t *testing.T) {
	t.Run("正常なリクエスト", func(t *testing.T) {
		req := LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePassword123!",
		}

		if req.Email != "user@example.com" {
			t.Errorf("Email = %v, want %v", req.Email, "user@example.com")
		}
		if req.Password != "SecurePassword123!" {
			t.Errorf("Password = %v, want %v", req.Password, "SecurePassword123!")
		}
	})

	t.Run("空のリクエスト", func(t *testing.T) {
		req := LoginRequest{}

		if req.Email != "" {
			t.Error("Email should be empty")
		}
		if req.Password != "" {
			t.Error("Password should be empty")
		}
	})
}

// ============================================
// LoginResult関連テスト
// ============================================

func TestLoginResult(t *testing.T) {
	t.Run("正常なログイン結果", func(t *testing.T) {
		userID := sharedDomain.NewID()
		now := time.Now()
		tokenPair := TokenPair{
			AccessToken:           "access_token",
			RefreshToken:          "refresh_token",
			AccessTokenExpiresAt:  now.Add(15 * time.Minute),
			RefreshTokenExpiresAt: now.Add(7 * 24 * time.Hour),
		}

		result := LoginResult{
			UserID:    userID,
			Email:     "admin@example.com",
			FullName:  "山田 太郎",
			Role:      "admin",
			TokenPair: tokenPair,
		}

		if result.UserID != userID {
			t.Errorf("UserID = %v, want %v", result.UserID, userID)
		}
		if result.Email != "admin@example.com" {
			t.Errorf("Email = %v, want %v", result.Email, "admin@example.com")
		}
		if result.FullName != "山田 太郎" {
			t.Errorf("FullName = %v, want %v", result.FullName, "山田 太郎")
		}
		if result.Role != "admin" {
			t.Errorf("Role = %v, want %v", result.Role, "admin")
		}
		if result.TokenPair.AccessToken != "access_token" {
			t.Errorf("TokenPair.AccessToken = %v, want %v", result.TokenPair.AccessToken, "access_token")
		}
	})
}
