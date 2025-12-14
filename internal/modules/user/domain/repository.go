// Package domain ユーザードメイン層
package domain

import (
	"context"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*User, error)
	// FindByEmail メールアドレスで検索
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]User, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]User, error)
	// FindByRole ロールで検索
	FindByRole(ctx context.Context, role UserRole) ([]User, error)
	// Save 保存
	Save(ctx context.Context, user *User) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
	// UpdateLastLogin 最終ログイン日時更新
	UpdateLastLogin(ctx context.Context, id sharedDomain.ID) error
}

// RefreshTokenRepository リフレッシュトークンリポジトリインターフェース
type RefreshTokenRepository interface {
	// FindByTokenHash トークンハッシュで検索
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	// FindByUserID ユーザーIDで検索
	FindByUserID(ctx context.Context, userID sharedDomain.ID) ([]RefreshToken, error)
	// Save 保存
	Save(ctx context.Context, token *RefreshToken) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
	// DeleteByUserID ユーザーIDで削除
	DeleteByUserID(ctx context.Context, userID sharedDomain.ID) error
	// DeleteExpired 期限切れトークン削除
	DeleteExpired(ctx context.Context) error
}
