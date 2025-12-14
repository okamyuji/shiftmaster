// Package domain ユーザードメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// User ユーザーエンティティ
type User struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID *domain.ID
	// Email メールアドレス
	Email string
	// PasswordHash パスワードハッシュ
	PasswordHash string
	// FirstName 名
	FirstName string
	// LastName 姓
	LastName string
	// Role ロール
	Role UserRole
	// IsActive 有効フラグ
	IsActive bool
	// LastLoginAt 最終ログイン日時
	LastLoginAt *time.Time
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// FullName フルネーム
func (u *User) FullName() string {
	return u.LastName + " " + u.FirstName
}

// IsSuperAdmin スーパー管理者判定（全テナント管理可能）
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsAdmin テナント管理者以上判定
func (u *User) IsAdmin() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin
}

// IsManager マネージャー以上判定
func (u *User) IsManager() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin || u.Role == RoleManager
}

// CanManageUsers ユーザー管理権限判定（テナント管理者以上）
func (u *User) CanManageUsers() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin
}

// CanManageSchedules 勤務表管理権限判定（マネージャー以上）
func (u *User) CanManageSchedules() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin || u.Role == RoleManager
}

// CanViewReports レポート閲覧権限判定（マネージャー以上）
func (u *User) CanViewReports() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin || u.Role == RoleManager
}

// CanAccessAllTenants 全テナントアクセス権限判定
func (u *User) CanAccessAllTenants() bool {
	return u.Role == RoleSuperAdmin
}

// UserRole ユーザーロール
type UserRole string

const (
	// RoleSuperAdmin スーパー管理者（全テナント管理可能）
	RoleSuperAdmin UserRole = "super_admin"
	// RoleAdmin テナント管理者（自テナントのみ管理）
	RoleAdmin UserRole = "admin"
	// RoleManager マネージャー（勤務表・スタッフ管理）
	RoleManager UserRole = "manager"
	// RoleUser 一般ユーザー（自分の操作のみ）
	RoleUser UserRole = "user"
)

// String 文字列変換
func (r UserRole) String() string {
	return string(r)
}

// Label 表示ラベル
func (r UserRole) Label() string {
	switch r {
	case RoleSuperAdmin:
		return "全体管理者"
	case RoleAdmin:
		return "テナント管理者"
	case RoleManager:
		return "マネージャー"
	case RoleUser:
		return "一般ユーザー"
	default:
		return "不明"
	}
}

// IsValid ロール有効性チェック
func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleManager, RoleUser:
		return true
	default:
		return false
	}
}

// RefreshToken リフレッシュトークンエンティティ
type RefreshToken struct {
	// ID 一意識別子
	ID domain.ID
	// UserID ユーザーID
	UserID domain.ID
	// TokenHash トークンハッシュ
	TokenHash string
	// ExpiresAt 有効期限
	ExpiresAt time.Time
	// CreatedAt 作成日時
	CreatedAt time.Time
}

// IsExpired 期限切れ判定
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
