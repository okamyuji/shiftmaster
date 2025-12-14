// Package domain スタッフドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// Staff スタッフエンティティ
type Staff struct {
	// ID 一意識別子
	ID domain.ID
	// TeamID 所属チームID
	TeamID domain.ID
	// EmployeeCode 社員番号
	EmployeeCode string
	// FirstName 名
	FirstName string
	// LastName 姓
	LastName string
	// Email メールアドレス
	Email string
	// Phone 電話番号
	Phone string
	// HireDate 入社日
	HireDate *time.Time
	// EmploymentType 雇用形態
	EmploymentType EmploymentType
	// IsActive 有効フラグ
	IsActive bool
	// Skills 保有スキル
	Skills []StaffSkill
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// FullName フルネーム取得
func (s *Staff) FullName() string {
	return s.LastName + " " + s.FirstName
}

// HasSkill スキル保有判定
func (s *Staff) HasSkill(skillID domain.ID) bool {
	for _, skill := range s.Skills {
		if skill.SkillID == skillID {
			return true
		}
	}
	return false
}

// EmploymentType 雇用形態
type EmploymentType string

const (
	// EmploymentFullTime 正社員
	EmploymentFullTime EmploymentType = "full_time"
	// EmploymentPartTime パートタイム
	EmploymentPartTime EmploymentType = "part_time"
	// EmploymentContract 契約社員
	EmploymentContract EmploymentType = "contract"
	// EmploymentTemporary 派遣社員
	EmploymentTemporary EmploymentType = "temporary"
)

// String 文字列変換
func (e EmploymentType) String() string {
	return string(e)
}

// Label 表示ラベル
func (e EmploymentType) Label() string {
	switch e {
	case EmploymentFullTime:
		return "正社員"
	case EmploymentPartTime:
		return "パートタイム"
	case EmploymentContract:
		return "契約社員"
	case EmploymentTemporary:
		return "派遣社員"
	default:
		return "不明"
	}
}

// StaffSkill スタッフスキル
type StaffSkill struct {
	// StaffID スタッフID
	StaffID domain.ID
	// SkillID スキルID
	SkillID domain.ID
	// Level 習熟度 1-5
	Level int
	// AcquiredAt 取得日
	AcquiredAt *time.Time
}

// Team チームエンティティ
type Team struct {
	// ID 一意識別子
	ID domain.ID
	// DepartmentID 部門ID
	DepartmentID domain.ID
	// Name チーム名
	Name string
	// Code チームコード
	Code string
	// SortOrder 表示順
	SortOrder int
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// Department 部門エンティティ
type Department struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name 部門名
	Name string
	// Code 部門コード
	Code string
	// SortOrder 表示順
	SortOrder int
	// Teams 所属チーム
	Teams []Team
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// Organization 組織エンティティ
type Organization struct {
	// ID 一意識別子
	ID domain.ID
	// Name 組織名
	Name string
	// Code 組織コード
	Code string
	// Departments 部門
	Departments []Department
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// Skill スキルエンティティ
type Skill struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name スキル名
	Name string
	// Description 説明
	Description string
	// Color 表示色
	Color string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}
