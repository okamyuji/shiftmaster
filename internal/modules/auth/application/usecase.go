// Package application 認証アプリケーション層
package application

import (
	"context"
	"log/slog"
	"time"

	authDomain "shiftmaster/internal/modules/auth/domain"
	userDomain "shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase 認証ユースケース
type AuthUseCase struct {
	userRepo     userDomain.UserRepository
	tokenRepo    userDomain.RefreshTokenRepository
	tokenService authDomain.TokenService
	logger       *slog.Logger
}

// NewAuthUseCase 認証ユースケース生成
func NewAuthUseCase(
	userRepo userDomain.UserRepository,
	tokenRepo userDomain.RefreshTokenRepository,
	tokenService authDomain.TokenService,
	logger *slog.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		logger:       logger,
	}
}

// Login ログイン
func (u *AuthUseCase) Login(ctx context.Context, input *LoginInput) (*AuthOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// ユーザー検索
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		u.logger.Error("ユーザー検索失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "認証処理に失敗しました")
	}
	if user == nil {
		u.logger.Warn("ログイン失敗 ユーザー未検出", "email", input.Email)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "メールアドレスまたはパスワードが正しくありません")
	}

	// 有効ユーザーチェック
	if !user.IsActive {
		u.logger.Warn("ログイン失敗 無効ユーザー", "user_id", user.ID)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "このアカウントは無効化されています")
	}

	// パスワード検証
	if pwErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); pwErr != nil {
		u.logger.Warn("ログイン失敗 パスワード不一致", "email", input.Email)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "メールアドレスまたはパスワードが正しくありません")
	}

	// 既存のリフレッシュトークンを削除（重複防止）
	if delErr := u.tokenRepo.DeleteByUserID(ctx, user.ID); delErr != nil {
		u.logger.Warn("既存リフレッシュトークン削除失敗", "error", delErr, "user_id", user.ID)
		// エラーでも処理を続行
	}

	// トークン生成
	claims := &authDomain.Claims{
		UserID:         user.ID,
		Email:          user.Email,
		Role:           user.Role.String(),
		OrganizationID: user.OrganizationID,
		IssuedAt:       time.Now(),
	}

	tokenPair, err := u.tokenService.GenerateTokenPair(claims)
	if err != nil {
		u.logger.Error("トークン生成失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "トークンの生成に失敗しました")
	}

	// リフレッシュトークン保存
	refreshToken := &userDomain.RefreshToken{
		ID:        sharedDomain.NewID(),
		UserID:    user.ID,
		TokenHash: u.tokenService.HashToken(tokenPair.RefreshToken),
		ExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := u.tokenRepo.Save(ctx, refreshToken); err != nil {
		u.logger.Error("リフレッシュトークン保存失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "トークンの保存に失敗しました")
	}

	// 最終ログイン日時更新
	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		u.logger.Error("最終ログイン日時更新失敗", "error", err)
	}

	u.logger.Info("ログイン成功", "user_id", user.ID, "email", user.Email)

	expiresIn := int64(time.Until(tokenPair.AccessTokenExpiresAt).Seconds())

	return &AuthOutput{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User: AuthUserOutput{
			ID:        user.ID.String(),
			Email:     user.Email,
			FullName:  user.FullName(),
			Role:      user.Role.String(),
			RoleLabel: user.Role.Label(),
			IsAdmin:   user.IsAdmin(),
		},
	}, nil
}

// Refresh トークンリフレッシュ
func (u *AuthUseCase) Refresh(ctx context.Context, input *RefreshInput) (*AuthOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// リフレッシュトークン検証
	claims, err := u.tokenService.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		u.logger.Warn("リフレッシュトークン検証失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "無効なリフレッシュトークンです")
	}

	// トークンハッシュ確認
	tokenHash := u.tokenService.HashToken(input.RefreshToken)
	storedToken, err := u.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		u.logger.Error("リフレッシュトークン検索失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "トークンの検証に失敗しました")
	}
	if storedToken == nil {
		u.logger.Warn("リフレッシュトークン未検出", "user_id", claims.UserID)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "無効なリフレッシュトークンです")
	}

	// 期限切れチェック
	if storedToken.IsExpired() {
		u.logger.Warn("リフレッシュトークン期限切れ", "user_id", claims.UserID)
		// 期限切れトークン削除
		if delErr := u.tokenRepo.Delete(ctx, storedToken.ID); delErr != nil {
			u.logger.Error("期限切れトークン削除失敗", "error", delErr)
		}
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "リフレッシュトークンが期限切れです")
	}

	// ユーザー取得
	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		u.logger.Error("ユーザー検索失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "ユーザー情報の取得に失敗しました")
	}
	if user == nil || !user.IsActive {
		u.logger.Warn("リフレッシュ失敗 無効ユーザー", "user_id", claims.UserID)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "このアカウントは無効化されています")
	}

	// 古いトークン削除
	if delErr := u.tokenRepo.Delete(ctx, storedToken.ID); delErr != nil {
		u.logger.Error("古いトークン削除失敗", "error", delErr)
	}

	// 新しいトークン生成
	newClaims := &authDomain.Claims{
		UserID:         user.ID,
		Email:          user.Email,
		Role:           user.Role.String(),
		OrganizationID: user.OrganizationID,
		IssuedAt:       time.Now(),
	}

	tokenPair, err := u.tokenService.GenerateTokenPair(newClaims)
	if err != nil {
		u.logger.Error("トークン生成失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "トークンの生成に失敗しました")
	}

	// 新しいリフレッシュトークン保存
	newRefreshToken := &userDomain.RefreshToken{
		ID:        sharedDomain.NewID(),
		UserID:    user.ID,
		TokenHash: u.tokenService.HashToken(tokenPair.RefreshToken),
		ExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := u.tokenRepo.Save(ctx, newRefreshToken); err != nil {
		u.logger.Error("リフレッシュトークン保存失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "トークンの保存に失敗しました")
	}

	u.logger.Info("トークンリフレッシュ成功", "user_id", user.ID)

	expiresIn := int64(time.Until(tokenPair.AccessTokenExpiresAt).Seconds())

	return &AuthOutput{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User: AuthUserOutput{
			ID:        user.ID.String(),
			Email:     user.Email,
			FullName:  user.FullName(),
			Role:      user.Role.String(),
			RoleLabel: user.Role.Label(),
			IsAdmin:   user.IsAdmin(),
		},
	}, nil
}

// Logout ログアウト
func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	tokenHash := u.tokenService.HashToken(refreshToken)
	storedToken, err := u.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		u.logger.Error("リフレッシュトークン検索失敗", "error", err)
		return nil
	}
	if storedToken == nil {
		return nil
	}

	if err := u.tokenRepo.Delete(ctx, storedToken.ID); err != nil {
		u.logger.Error("リフレッシュトークン削除失敗", "error", err)
		return err
	}

	u.logger.Info("ログアウト成功", "user_id", storedToken.UserID)
	return nil
}

// ValidateToken トークン検証
func (u *AuthUseCase) ValidateToken(ctx context.Context, accessToken string) (*authDomain.Claims, error) {
	claims, err := u.tokenService.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	// ユーザーが有効か確認
	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "ユーザー情報の取得に失敗しました")
	}
	if user == nil || !user.IsActive {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "このアカウントは無効化されています")
	}

	return claims, nil
}
