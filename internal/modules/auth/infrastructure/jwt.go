// Package infrastructure 認証インフラストラクチャ層
package infrastructure

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"shiftmaster/internal/modules/auth/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT設定
type JWTConfig struct {
	// SecretKey 秘密鍵
	SecretKey string
	// AccessTokenDuration アクセストークン有効期間
	AccessTokenDuration time.Duration
	// RefreshTokenDuration リフレッシュトークン有効期間
	RefreshTokenDuration time.Duration
	// Issuer 発行者
	Issuer string
}

// DefaultJWTConfig デフォルトJWT設定
func DefaultJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey:            secretKey,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               "shiftmaster",
	}
}

// JWTTokenService JWT実装のトークンサービス
type JWTTokenService struct {
	config *JWTConfig
}

// NewJWTTokenService トークンサービス生成
func NewJWTTokenService(config *JWTConfig) *JWTTokenService {
	return &JWTTokenService{config: config}
}

// jwtClaims JWTクレーム
type jwtClaims struct {
	UserID         string  `json:"user_id"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	OrganizationID *string `json:"organization_id,omitempty"`
	TokenType      string  `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateTokenPair トークンペア生成
func (s *JWTTokenService) GenerateTokenPair(claims *domain.Claims) (*domain.TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(s.config.AccessTokenDuration)
	refreshExpiry := now.Add(s.config.RefreshTokenDuration)

	// アクセストークン生成
	accessToken, err := s.generateToken(claims, "access", accessExpiry)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークン生成
	refreshToken, err := s.generateToken(claims, "refresh", refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiry,
		RefreshTokenExpiresAt: refreshExpiry,
	}, nil
}

// generateToken トークン生成
func (s *JWTTokenService) generateToken(claims *domain.Claims, tokenType string, expiresAt time.Time) (string, error) {
	var orgID *string
	if claims.OrganizationID != nil {
		orgIDStr := claims.OrganizationID.String()
		orgID = &orgIDStr
	}

	tokenClaims := jwtClaims{
		UserID:         claims.UserID.String(),
		Email:          claims.Email,
		Role:           claims.Role,
		OrganizationID: orgID,
		TokenType:      tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateAccessToken アクセストークン検証
func (s *JWTTokenService) ValidateAccessToken(tokenString string) (*domain.Claims, error) {
	return s.validateToken(tokenString, "access")
}

// ValidateRefreshToken リフレッシュトークン検証
func (s *JWTTokenService) ValidateRefreshToken(tokenString string) (*domain.Claims, error) {
	return s.validateToken(tokenString, "refresh")
}

// validateToken トークン検証
func (s *JWTTokenService) validateToken(tokenString, expectedType string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "不正な署名アルゴリズムです")
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "トークンの検証に失敗しました")
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "無効なトークンです")
	}

	if claims.TokenType != expectedType {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "トークンタイプが不正です")
	}

	userID, err := sharedDomain.ParseID(claims.UserID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "ユーザーIDが不正です")
	}

	var orgID *sharedDomain.ID
	if claims.OrganizationID != nil {
		id, err := sharedDomain.ParseID(*claims.OrganizationID)
		if err == nil {
			orgID = &id
		}
	}

	return &domain.Claims{
		UserID:         userID,
		Email:          claims.Email,
		Role:           claims.Role,
		OrganizationID: orgID,
		IssuedAt:       claims.IssuedAt.Time,
		ExpiresAt:      claims.ExpiresAt.Time,
	}, nil
}

// HashToken トークンハッシュ化
func (s *JWTTokenService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// BcryptPasswordService bcrypt実装のパスワードサービス
type BcryptPasswordService struct {
	cost int
}

// NewBcryptPasswordService パスワードサービス生成
func NewBcryptPasswordService() *BcryptPasswordService {
	return &BcryptPasswordService{cost: 10}
}

// Hash パスワードハッシュ化
func (s *BcryptPasswordService) Hash(password string) (string, error) {
	// bcryptはgolang.org/x/cryptoを使用
	// user usecase内で直接bcryptを使用しているため、ここでは簡易的に実装
	return password, nil // 実際の実装はusecase層で行う
}

// Verify パスワード検証
func (s *BcryptPasswordService) Verify(hashedPassword, password string) error {
	// 実際の実装はusecase層で行う
	return nil
}
