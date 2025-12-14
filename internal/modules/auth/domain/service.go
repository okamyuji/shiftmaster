// Package domain 認証ドメイン層
package domain

// TokenService トークンサービスインターフェース
type TokenService interface {
	// GenerateTokenPair トークンペア生成
	GenerateTokenPair(claims *Claims) (*TokenPair, error)
	// ValidateAccessToken アクセストークン検証
	ValidateAccessToken(token string) (*Claims, error)
	// ValidateRefreshToken リフレッシュトークン検証
	ValidateRefreshToken(token string) (*Claims, error)
	// HashToken トークンハッシュ化
	HashToken(token string) string
}

// PasswordService パスワードサービスインターフェース
type PasswordService interface {
	// Hash パスワードハッシュ化
	Hash(password string) (string, error)
	// Verify パスワード検証
	Verify(hashedPassword, password string) error
}
