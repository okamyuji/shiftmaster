// Package application 認証アプリケーション層
package application

import (
	sharedDomain "shiftmaster/internal/shared/domain"
)

// LoginInput ログイン入力
type LoginInput struct {
	// Email メールアドレス
	Email string `json:"email"`
	// Password パスワード
	Password string `json:"password"`
}

// Validate 入力検証
func (i *LoginInput) Validate() error {
	if i.Email == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "メールアドレスは必須です")
	}
	if i.Password == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パスワードは必須です")
	}
	return nil
}

// RefreshInput トークンリフレッシュ入力
type RefreshInput struct {
	// RefreshToken リフレッシュトークン
	RefreshToken string `json:"refresh_token"`
}

// Validate 入力検証
func (i *RefreshInput) Validate() error {
	if i.RefreshToken == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "リフレッシュトークンは必須です")
	}
	return nil
}

// AuthOutput 認証結果出力
type AuthOutput struct {
	// AccessToken アクセストークン
	AccessToken string `json:"access_token"`
	// RefreshToken リフレッシュトークン
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn 有効期限（秒）
	ExpiresIn int64 `json:"expires_in"`
	// TokenType トークンタイプ
	TokenType string `json:"token_type"`
	// User ユーザー情報
	User AuthUserOutput `json:"user"`
}

// AuthUserOutput 認証ユーザー情報出力
type AuthUserOutput struct {
	// ID ユーザーID
	ID string `json:"id"`
	// Email メールアドレス
	Email string `json:"email"`
	// FullName フルネーム
	FullName string `json:"full_name"`
	// Role ロール
	Role string `json:"role"`
	// RoleLabel ロールラベル
	RoleLabel string `json:"role_label"`
	// IsAdmin 管理者フラグ
	IsAdmin bool `json:"is_admin"`
}
