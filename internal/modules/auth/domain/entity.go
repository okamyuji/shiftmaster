// Package domain 認証ドメイン層
package domain

import (
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// TokenPair トークンペア
type TokenPair struct {
	// AccessToken アクセストークン
	AccessToken string
	// RefreshToken リフレッシュトークン
	RefreshToken string
	// AccessTokenExpiresAt アクセストークン有効期限
	AccessTokenExpiresAt time.Time
	// RefreshTokenExpiresAt リフレッシュトークン有効期限
	RefreshTokenExpiresAt time.Time
}

// Claims JWT クレーム
type Claims struct {
	// UserID ユーザーID
	UserID sharedDomain.ID
	// Email メールアドレス
	Email string
	// Role ロール
	Role string
	// OrganizationID 組織ID
	OrganizationID *sharedDomain.ID
	// IssuedAt 発行日時
	IssuedAt time.Time
	// ExpiresAt 有効期限
	ExpiresAt time.Time
}

// IsExpired 期限切れ判定
func (c *Claims) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// IsSuperAdmin スーパー管理者判定（全テナント管理可能）
func (c *Claims) IsSuperAdmin() bool {
	return c.Role == "super_admin"
}

// IsAdmin テナント管理者以上判定
func (c *Claims) IsAdmin() bool {
	return c.Role == "super_admin" || c.Role == "admin"
}

// IsManager マネージャー以上判定
func (c *Claims) IsManager() bool {
	return c.Role == "super_admin" || c.Role == "admin" || c.Role == "manager"
}

// CanAccessAllTenants 全テナントアクセス権限判定
func (c *Claims) CanAccessAllTenants() bool {
	return c.Role == "super_admin"
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	// Email メールアドレス
	Email string
	// Password パスワード
	Password string
}

// LoginResult ログイン結果
type LoginResult struct {
	// User ユーザー情報
	UserID   sharedDomain.ID
	Email    string
	FullName string
	Role     string
	// Tokens トークン
	TokenPair TokenPair
}
