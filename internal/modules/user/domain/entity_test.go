// Package domain ユーザードメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ============================================
// User関連テスト
// ============================================

func TestUser_FullName(t *testing.T) {
	tests := []struct {
		name      string
		lastName  string
		firstName string
		expected  string
	}{
		{
			name:      "正常系_日本語名前",
			lastName:  "山田",
			firstName: "太郎",
			expected:  "山田 太郎",
		},
		{
			name:      "正常系_英語名前",
			lastName:  "Smith",
			firstName: "John",
			expected:  "Smith John",
		},
		{
			name:      "姓のみ",
			lastName:  "山田",
			firstName: "",
			expected:  "山田 ",
		},
		{
			name:      "名のみ",
			lastName:  "",
			firstName: "太郎",
			expected:  " 太郎",
		},
		{
			name:      "空の名前",
			lastName:  "",
			firstName: "",
			expected:  " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				LastName:  tt.lastName,
				FirstName: tt.firstName,
			}

			result := user.FullName()
			if result != tt.expected {
				t.Errorf("FullName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_IsSuperAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール",
			role:     RoleAdmin,
			expected: false,
		},
		{
			name:     "managerロール",
			role:     RoleManager,
			expected: false,
		},
		{
			name:     "userロール",
			role:     RoleUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.IsSuperAdmin()
			if result != tt.expected {
				t.Errorf("IsSuperAdmin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール",
			role:     RoleManager,
			expected: false,
		},
		{
			name:     "userロール",
			role:     RoleUser,
			expected: false,
		},
		{
			name:     "空のロール",
			role:     UserRole(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.IsAdmin()
			if result != tt.expected {
				t.Errorf("IsAdmin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_IsManager(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "userロール",
			role:     RoleUser,
			expected: false,
		},
		{
			name:     "不明なロール",
			role:     UserRole("unknown"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.IsManager()
			if result != tt.expected {
				t.Errorf("IsManager() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_CanManageUsers(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール_許可",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール_許可",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール_不許可",
			role:     RoleManager,
			expected: false,
		},
		{
			name:     "userロール_不許可",
			role:     RoleUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.CanManageUsers()
			if result != tt.expected {
				t.Errorf("CanManageUsers() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_CanManageSchedules(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール_許可",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール_許可",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール_許可",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "userロール_不許可",
			role:     RoleUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.CanManageSchedules()
			if result != tt.expected {
				t.Errorf("CanManageSchedules() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_CanViewReports(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール_許可",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール_許可",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール_許可",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "userロール_不許可",
			role:     RoleUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.CanViewReports()
			if result != tt.expected {
				t.Errorf("CanViewReports() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_CanAccessAllTenants(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール_許可",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール_不許可",
			role:     RoleAdmin,
			expected: false,
		},
		{
			name:     "managerロール_不許可",
			role:     RoleManager,
			expected: false,
		},
		{
			name:     "userロール_不許可",
			role:     RoleUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			result := user.CanAccessAllTenants()
			if result != tt.expected {
				t.Errorf("CanAccessAllTenants() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_Structure(t *testing.T) {
	t.Run("User構造体の完全な初期化", func(t *testing.T) {
		userID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()
		loginTime := now.Add(-1 * time.Hour)

		user := &User{
			ID:             userID,
			OrganizationID: &orgID,
			Email:          "test@example.com",
			PasswordHash:   "hashed_password",
			FirstName:      "太郎",
			LastName:       "山田",
			Role:           RoleAdmin,
			IsActive:       true,
			LastLoginAt:    &loginTime,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if user.ID != userID {
			t.Errorf("ID = %v, want %v", user.ID, userID)
		}
		if user.OrganizationID == nil || *user.OrganizationID != orgID {
			t.Error("OrganizationID not set correctly")
		}
		if user.Email != "test@example.com" {
			t.Errorf("Email = %v, want %v", user.Email, "test@example.com")
		}
		if user.PasswordHash != "hashed_password" {
			t.Errorf("PasswordHash = %v, want %v", user.PasswordHash, "hashed_password")
		}
		if user.FirstName != "太郎" {
			t.Errorf("FirstName = %v, want %v", user.FirstName, "太郎")
		}
		if user.LastName != "山田" {
			t.Errorf("LastName = %v, want %v", user.LastName, "山田")
		}
		if user.Role != RoleAdmin {
			t.Errorf("Role = %v, want %v", user.Role, RoleAdmin)
		}
		if !user.IsActive {
			t.Error("IsActive should be true")
		}
		if user.LastLoginAt == nil {
			t.Error("LastLoginAt should not be nil")
		}
		if user.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", user.CreatedAt, now)
		}
		if user.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", user.UpdatedAt, now)
		}
	})
}

// ============================================
// UserRole関連テスト
// ============================================

func TestUserRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected string
	}{
		{
			name:     "adminロール",
			role:     RoleAdmin,
			expected: "admin",
		},
		{
			name:     "managerロール",
			role:     RoleManager,
			expected: "manager",
		},
		{
			name:     "userロール",
			role:     RoleUser,
			expected: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUserRole_Label(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected string
	}{
		{
			name:     "super_adminロール",
			role:     RoleSuperAdmin,
			expected: "全体管理者",
		},
		{
			name:     "adminロール",
			role:     RoleAdmin,
			expected: "テナント管理者",
		},
		{
			name:     "managerロール",
			role:     RoleManager,
			expected: "マネージャー",
		},
		{
			name:     "userロール",
			role:     RoleUser,
			expected: "一般ユーザー",
		},
		{
			name:     "不明なロール",
			role:     UserRole("unknown"),
			expected: "不明",
		},
		{
			name:     "空のロール",
			role:     UserRole(""),
			expected: "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "super_adminロール_有効",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "adminロール_有効",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "managerロール_有効",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "userロール_有効",
			role:     RoleUser,
			expected: true,
		},
		{
			name:     "不明なロール_無効",
			role:     UserRole("unknown"),
			expected: false,
		},
		{
			name:     "空のロール_無効",
			role:     UserRole(""),
			expected: false,
		},
		{
			name:     "大文字のロール_無効",
			role:     UserRole("ADMIN"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// RefreshToken関連テスト
// ============================================

func TestRefreshToken_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "期限切れ_過去",
			expiresAt: time.Now().Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "有効_未来",
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
			token := &RefreshToken{
				ExpiresAt: tt.expiresAt,
			}

			result := token.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRefreshToken_Structure(t *testing.T) {
	t.Run("RefreshToken構造体の初期化", func(t *testing.T) {
		tokenID := sharedDomain.NewID()
		userID := sharedDomain.NewID()
		now := time.Now()
		expiresAt := now.Add(7 * 24 * time.Hour)

		token := &RefreshToken{
			ID:        tokenID,
			UserID:    userID,
			TokenHash: "hashed_token_value",
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		if token.ID != tokenID {
			t.Errorf("ID = %v, want %v", token.ID, tokenID)
		}
		if token.UserID != userID {
			t.Errorf("UserID = %v, want %v", token.UserID, userID)
		}
		if token.TokenHash != "hashed_token_value" {
			t.Errorf("TokenHash = %v, want %v", token.TokenHash, "hashed_token_value")
		}
		if token.ExpiresAt != expiresAt {
			t.Errorf("ExpiresAt = %v, want %v", token.ExpiresAt, expiresAt)
		}
		if token.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", token.CreatedAt, now)
		}
	})
}

// ============================================
// 定数テスト
// ============================================

func TestRoleConstants(t *testing.T) {
	t.Run("ロール定数の値", func(t *testing.T) {
		if RoleAdmin != "admin" {
			t.Errorf("RoleAdmin = %v, want %v", RoleAdmin, "admin")
		}
		if RoleManager != "manager" {
			t.Errorf("RoleManager = %v, want %v", RoleManager, "manager")
		}
		if RoleUser != "user" {
			t.Errorf("RoleUser = %v, want %v", RoleUser, "user")
		}
	})

	t.Run("ロール定数の一意性", func(t *testing.T) {
		roles := []UserRole{RoleAdmin, RoleManager, RoleUser}
		seen := make(map[UserRole]bool)

		for _, role := range roles {
			if seen[role] {
				t.Errorf("Duplicate role constant: %v", role)
			}
			seen[role] = true
		}
	})
}
