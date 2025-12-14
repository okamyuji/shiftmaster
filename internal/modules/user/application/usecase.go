// Package application ユーザーアプリケーション層
package application

import (
	"context"
	"log/slog"
	"time"

	"shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"golang.org/x/crypto/bcrypt"
)

// UserUseCase ユーザーユースケース
type UserUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.RefreshTokenRepository
	logger    *slog.Logger
}

// NewUserUseCase ユーザーユースケース生成
func NewUserUseCase(
	userRepo domain.UserRepository,
	tokenRepo domain.RefreshTokenRepository,
	logger *slog.Logger,
) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

// Create ユーザー作成
func (u *UserUseCase) Create(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// メールアドレス重複チェック
	existing, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeConflict, "このメールアドレスは既に使用されています")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("パスワードハッシュ化失敗", "error", err)
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "パスワードの処理に失敗しました")
	}

	var orgID *sharedDomain.ID
	if input.OrganizationID != "" {
		id, err := sharedDomain.ParseID(input.OrganizationID)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
		}
		orgID = &id
	}

	now := time.Now()
	user := &domain.User{
		ID:             sharedDomain.NewID(),
		OrganizationID: orgID,
		Email:          input.Email,
		PasswordHash:   string(hashedPassword),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Role:           domain.UserRole(input.Role),
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.userRepo.Save(ctx, user); err != nil {
		u.logger.Error("ユーザー作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("ユーザー作成完了", "user_id", user.ID, "email", user.Email)
	return ToUserOutput(user), nil
}

// GetByID IDでユーザー取得
func (u *UserUseCase) GetByID(ctx context.Context, id string) (*UserOutput, error) {
	userID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToUserOutput(user), nil
}

// GetByEmail メールアドレスでユーザー取得
func (u *UserUseCase) GetByEmail(ctx context.Context, email string) (*UserOutput, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToUserOutput(user), nil
}

// List ユーザー一覧取得
// Deprecated: ListByOrganization を使用してください
func (u *UserUseCase) List(ctx context.Context) (*UserListOutput, error) {
	users, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]UserOutput, len(users))
	for i, user := range users {
		outputs[i] = *ToUserOutput(&user)
	}

	return &UserListOutput{
		Users: outputs,
		Total: len(outputs),
	}, nil
}

// ListByOrganization 組織IDでユーザー一覧取得（マルチテナント対応）
func (u *UserUseCase) ListByOrganization(ctx context.Context, orgID string) (*UserListOutput, error) {
	if orgID == "" {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが必要です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	users, err := u.userRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	outputs := make([]UserOutput, len(users))
	for i, user := range users {
		outputs[i] = *ToUserOutput(&user)
	}

	return &UserListOutput{
		Users: outputs,
		Total: len(outputs),
	}, nil
}

// Update ユーザー更新
func (u *UserUseCase) Update(ctx context.Context, input *UpdateUserInput) (*UserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	userID, err := sharedDomain.ParseID(input.ID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, sharedDomain.ErrNotFound
	}

	// メールアドレス重複チェック 自分以外
	if user.Email != input.Email {
		existing, err := u.userRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != user.ID {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeConflict, "このメールアドレスは既に使用されています")
		}
	}

	var orgID *sharedDomain.ID
	if input.OrganizationID != "" {
		id, err := sharedDomain.ParseID(input.OrganizationID)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
		}
		orgID = &id
	}

	user.OrganizationID = orgID
	user.Email = input.Email
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Role = domain.UserRole(input.Role)
	user.IsActive = input.IsActive
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Save(ctx, user); err != nil {
		u.logger.Error("ユーザー更新失敗", "error", err)
		return nil, err
	}

	u.logger.Info("ユーザー更新完了", "user_id", user.ID)
	return ToUserOutput(user), nil
}

// ChangePassword パスワード変更
func (u *UserUseCase) ChangePassword(ctx context.Context, input *ChangePasswordInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	userID, err := sharedDomain.ParseID(input.UserID)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "ユーザーIDが不正です")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return sharedDomain.ErrNotFound
	}

	// 現在のパスワード検証
	if pwErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); pwErr != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeUnauthorized, "現在のパスワードが正しくありません")
	}

	// 新しいパスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("パスワードハッシュ化失敗", "error", err)
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeInternal, "パスワードの処理に失敗しました")
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Save(ctx, user); err != nil {
		u.logger.Error("パスワード変更失敗", "error", err)
		return err
	}

	// リフレッシュトークン削除 全セッション無効化
	if err := u.tokenRepo.DeleteByUserID(ctx, userID); err != nil {
		u.logger.Error("リフレッシュトークン削除失敗", "error", err)
	}

	u.logger.Info("パスワード変更完了", "user_id", user.ID)
	return nil
}

// Delete ユーザー削除
func (u *UserUseCase) Delete(ctx context.Context, id string) error {
	userID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return sharedDomain.ErrNotFound
	}

	// 管理者の最後の1人は削除不可
	if user.Role == domain.RoleAdmin {
		admins, err := u.userRepo.FindByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return err
		}
		if len(admins) <= 1 {
			return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "最後の管理者は削除できません")
		}
	}

	// リフレッシュトークン削除
	if err := u.tokenRepo.DeleteByUserID(ctx, userID); err != nil {
		u.logger.Error("リフレッシュトークン削除失敗", "error", err)
	}

	if err := u.userRepo.Delete(ctx, userID); err != nil {
		u.logger.Error("ユーザー削除失敗", "error", err)
		return err
	}

	u.logger.Info("ユーザー削除完了", "user_id", userID)
	return nil
}
