// Package application ユーザーアプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// CreateUserInput ユーザー作成入力
type CreateUserInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Email メールアドレス
	Email string `json:"email"`
	// Password パスワード
	Password string `json:"password"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// Role ロール
	Role string `json:"role"`
}

// Validate 入力検証
func (i *CreateUserInput) Validate() error {
	if i.Email == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "メールアドレスは必須です")
	}
	if i.Password == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パスワードは必須です")
	}
	if len(i.Password) < 8 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パスワードは8文字以上である必要があります")
	}
	if i.FirstName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "名は必須です")
	}
	if i.LastName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "姓は必須です")
	}
	role := domain.UserRole(i.Role)
	if !role.IsValid() {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "ロールが不正です")
	}
	return nil
}

// UpdateUserInput ユーザー更新入力
type UpdateUserInput struct {
	// ID ユーザーID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Email メールアドレス
	Email string `json:"email"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// Role ロール
	Role string `json:"role"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
}

// Validate 入力検証
func (i *UpdateUserInput) Validate() error {
	if i.ID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDは必須です")
	}
	if i.Email == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "メールアドレスは必須です")
	}
	if i.FirstName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "名は必須です")
	}
	if i.LastName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "姓は必須です")
	}
	role := domain.UserRole(i.Role)
	if !role.IsValid() {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "ロールが不正です")
	}
	return nil
}

// ChangePasswordInput パスワード変更入力
type ChangePasswordInput struct {
	// UserID ユーザーID
	UserID string `json:"user_id"`
	// CurrentPassword 現在のパスワード
	CurrentPassword string `json:"current_password"`
	// NewPassword 新しいパスワード
	NewPassword string `json:"new_password"`
}

// Validate 入力検証
func (i *ChangePasswordInput) Validate() error {
	if i.UserID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "ユーザーIDは必須です")
	}
	if i.CurrentPassword == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "現在のパスワードは必須です")
	}
	if i.NewPassword == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "新しいパスワードは必須です")
	}
	if len(i.NewPassword) < 8 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "新しいパスワードは8文字以上である必要があります")
	}
	return nil
}

// UserOutput ユーザー出力
type UserOutput struct {
	// ID ユーザーID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Email メールアドレス
	Email string `json:"email"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// FullName フルネーム
	FullName string `json:"full_name"`
	// Role ロール
	Role string `json:"role"`
	// RoleLabel ロールラベル
	RoleLabel string `json:"role_label"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
	// IsAdmin 管理者フラグ
	IsAdmin bool `json:"is_admin"`
	// LastLoginAt 最終ログイン日時
	LastLoginAt string `json:"last_login_at"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToUserOutput ドメインエンティティから出力DTOへ変換
func ToUserOutput(u *domain.User) *UserOutput {
	organizationID := ""
	if u.OrganizationID != nil {
		organizationID = u.OrganizationID.String()
	}

	lastLoginAt := ""
	if u.LastLoginAt != nil {
		lastLoginAt = u.LastLoginAt.Format(time.RFC3339)
	}

	return &UserOutput{
		ID:             u.ID.String(),
		OrganizationID: organizationID,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FullName:       u.FullName(),
		Role:           u.Role.String(),
		RoleLabel:      u.Role.Label(),
		IsActive:       u.IsActive,
		IsAdmin:        u.IsAdmin(),
		LastLoginAt:    lastLoginAt,
		CreatedAt:      u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      u.UpdatedAt.Format(time.RFC3339),
	}
}

// UserListOutput ユーザー一覧出力
type UserListOutput struct {
	// Users ユーザー一覧
	Users []UserOutput `json:"users"`
	// Total 総件数
	Total int `json:"total"`
}
