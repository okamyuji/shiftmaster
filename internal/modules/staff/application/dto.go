// Package application スタッフアプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// CreateStaffInput スタッフ作成入力
type CreateStaffInput struct {
	// TeamID 所属チームID
	TeamID string `json:"team_id"`
	// EmployeeCode 社員番号
	EmployeeCode string `json:"employee_code"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// Email メールアドレス
	Email string `json:"email"`
	// Phone 電話番号
	Phone string `json:"phone"`
	// HireDate 入社日
	HireDate string `json:"hire_date"`
	// EmploymentType 雇用形態
	EmploymentType string `json:"employment_type"`
}

// Validate 入力検証
func (i *CreateStaffInput) Validate() error {
	if i.TeamID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "チームIDは必須です")
	}
	if i.FirstName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "名は必須です")
	}
	if i.LastName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "姓は必須です")
	}
	return nil
}

// UpdateStaffInput スタッフ更新入力
type UpdateStaffInput struct {
	// ID スタッフID
	ID string `json:"id"`
	// TeamID 所属チームID
	TeamID string `json:"team_id"`
	// EmployeeCode 社員番号
	EmployeeCode string `json:"employee_code"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// Email メールアドレス
	Email string `json:"email"`
	// Phone 電話番号
	Phone string `json:"phone"`
	// HireDate 入社日
	HireDate string `json:"hire_date"`
	// EmploymentType 雇用形態
	EmploymentType string `json:"employment_type"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
}

// Validate 入力検証
func (i *UpdateStaffInput) Validate() error {
	if i.ID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDは必須です")
	}
	if i.TeamID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "チームIDは必須です")
	}
	if i.FirstName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "名は必須です")
	}
	if i.LastName == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "姓は必須です")
	}
	return nil
}

// StaffOutput スタッフ出力
type StaffOutput struct {
	// ID スタッフID
	ID string `json:"id"`
	// TeamID 所属チームID
	TeamID string `json:"team_id"`
	// EmployeeCode 社員番号
	EmployeeCode string `json:"employee_code"`
	// FirstName 名
	FirstName string `json:"first_name"`
	// LastName 姓
	LastName string `json:"last_name"`
	// FullName フルネーム
	FullName string `json:"full_name"`
	// Email メールアドレス
	Email string `json:"email"`
	// Phone 電話番号
	Phone string `json:"phone"`
	// HireDate 入社日
	HireDate string `json:"hire_date"`
	// EmploymentType 雇用形態
	EmploymentType string `json:"employment_type"`
	// EmploymentTypeLabel 雇用形態ラベル
	EmploymentTypeLabel string `json:"employment_type_label"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
	// Skills スキル
	Skills []StaffSkillOutput `json:"skills"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// StaffSkillOutput スタッフスキル出力
type StaffSkillOutput struct {
	// SkillID スキルID
	SkillID string `json:"skill_id"`
	// Level 習熟度
	Level int `json:"level"`
	// AcquiredAt 取得日
	AcquiredAt string `json:"acquired_at"`
}

// ToStaffOutput ドメインエンティティから出力DTOへ変換
func ToStaffOutput(staff *domain.Staff) *StaffOutput {
	hireDate := ""
	if staff.HireDate != nil {
		hireDate = staff.HireDate.Format("2006-01-02")
	}

	skills := make([]StaffSkillOutput, len(staff.Skills))
	for i, s := range staff.Skills {
		acquiredAt := ""
		if s.AcquiredAt != nil {
			acquiredAt = s.AcquiredAt.Format("2006-01-02")
		}
		skills[i] = StaffSkillOutput{
			SkillID:    s.SkillID.String(),
			Level:      s.Level,
			AcquiredAt: acquiredAt,
		}
	}

	return &StaffOutput{
		ID:                  staff.ID.String(),
		TeamID:              staff.TeamID.String(),
		EmployeeCode:        staff.EmployeeCode,
		FirstName:           staff.FirstName,
		LastName:            staff.LastName,
		FullName:            staff.FullName(),
		Email:               staff.Email,
		Phone:               staff.Phone,
		HireDate:            hireDate,
		EmploymentType:      staff.EmploymentType.String(),
		EmploymentTypeLabel: staff.EmploymentType.Label(),
		IsActive:            staff.IsActive,
		Skills:              skills,
		CreatedAt:           staff.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           staff.UpdatedAt.Format(time.RFC3339),
	}
}

// StaffListOutput スタッフ一覧出力
type StaffListOutput struct {
	// Staffs スタッフ一覧
	Staffs []StaffOutput `json:"staffs"`
	// Total 総件数
	Total int `json:"total"`
	// Page ページ番号
	Page int `json:"page"`
	// PerPage 1ページあたり件数
	PerPage int `json:"per_page"`
}
