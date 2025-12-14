// Package application 勤務表アプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/schedule/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// CreateScheduleInput 勤務表作成入力
type CreateScheduleInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
}

// Validate 入力検証
func (i *CreateScheduleInput) Validate() error {
	if i.OrganizationID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDは必須です")
	}
	if i.TargetYear < 2020 || i.TargetYear > 2100 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象年が不正です")
	}
	if i.TargetMonth < 1 || i.TargetMonth > 12 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象月が不正です")
	}
	return nil
}

// CreateEntryInput エントリ作成入力
type CreateEntryInput struct {
	// ScheduleID 勤務表ID
	ScheduleID string `json:"schedule_id"`
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// TargetDate 対象日（YYYY-MM-DD形式）
	TargetDate string `json:"target_date"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// Note 備考
	Note string `json:"note"`
}

// Validate 入力検証
func (i *CreateEntryInput) Validate() error {
	if i.ScheduleID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "勤務表IDは必須です")
	}
	if i.StaffID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "スタッフIDは必須です")
	}
	if i.TargetDate == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象日は必須です")
	}
	return nil
}

// UpdateEntryInput エントリ更新入力
type UpdateEntryInput struct {
	// ID エントリID
	ID string `json:"id"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// IsConfirmed 確定フラグ
	IsConfirmed bool `json:"is_confirmed"`
	// Note 備考
	Note string `json:"note"`
}

// Validate 入力検証
func (i *UpdateEntryInput) Validate() error {
	if i.ID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDは必須です")
	}
	return nil
}

// BulkUpdateEntriesInput 一括エントリ更新入力
type BulkUpdateEntriesInput struct {
	// ScheduleID 勤務表ID
	ScheduleID string `json:"schedule_id"`
	// Entries エントリ一覧
	Entries []EntryInput `json:"entries"`
}

// EntryInput エントリ入力
type EntryInput struct {
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// TargetDate 対象日
	TargetDate string `json:"target_date"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
}

// ScheduleOutput 勤務表出力
type ScheduleOutput struct {
	// ID 勤務表ID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
	// TargetPeriodLabel 対象期間ラベル
	TargetPeriodLabel string `json:"target_period_label"`
	// Status 状態
	Status string `json:"status"`
	// StatusLabel 状態ラベル
	StatusLabel string `json:"status_label"`
	// DaysInMonth 月の日数
	DaysInMonth int `json:"days_in_month"`
	// PublishedAt 公開日時
	PublishedAt string `json:"published_at"`
	// Entries エントリ一覧
	Entries []ScheduleEntryOutput `json:"entries"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToScheduleOutput ドメインエンティティから出力DTOへ変換
func ToScheduleOutput(s *domain.Schedule) *ScheduleOutput {
	publishedAt := ""
	if s.PublishedAt != nil {
		publishedAt = s.PublishedAt.Format(time.RFC3339)
	}

	entries := make([]ScheduleEntryOutput, len(s.Entries))
	for i, e := range s.Entries {
		entries[i] = *ToScheduleEntryOutput(&e)
	}

	return &ScheduleOutput{
		ID:                s.ID.String(),
		OrganizationID:    s.OrganizationID.String(),
		TargetYear:        s.TargetYear,
		TargetMonth:       s.TargetMonth,
		TargetPeriodLabel: s.TargetPeriodLabel(),
		Status:            s.Status.String(),
		StatusLabel:       s.Status.Label(),
		DaysInMonth:       s.DaysInMonth(),
		PublishedAt:       publishedAt,
		Entries:           entries,
		CreatedAt:         s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         s.UpdatedAt.Format(time.RFC3339),
	}
}

// ScheduleEntryOutput 勤務表エントリ出力
type ScheduleEntryOutput struct {
	// ID エントリID
	ID string `json:"id"`
	// ScheduleID 勤務表ID
	ScheduleID string `json:"schedule_id"`
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// StaffName スタッフ名
	StaffName string `json:"staff_name"`
	// TargetDate 対象日
	TargetDate string `json:"target_date"`
	// DayOfWeek 曜日
	DayOfWeek string `json:"day_of_week"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// ShiftTypeName シフト種別名
	ShiftTypeName string `json:"shift_type_name"`
	// ShiftTypeCode シフト種別コード
	ShiftTypeCode string `json:"shift_type_code"`
	// IsConfirmed 確定フラグ
	IsConfirmed bool `json:"is_confirmed"`
	// Note 備考
	Note string `json:"note"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToScheduleEntryOutput ドメインエンティティから出力DTOへ変換
func ToScheduleEntryOutput(e *domain.ScheduleEntry) *ScheduleEntryOutput {
	shiftTypeID := ""
	if e.ShiftTypeID != nil {
		shiftTypeID = e.ShiftTypeID.String()
	}

	weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}

	return &ScheduleEntryOutput{
		ID:          e.ID.String(),
		ScheduleID:  e.ScheduleID.String(),
		StaffID:     e.StaffID.String(),
		TargetDate:  e.TargetDate.Format("2006-01-02"),
		DayOfWeek:   weekdays[e.TargetDate.Weekday()],
		ShiftTypeID: shiftTypeID,
		IsConfirmed: e.IsConfirmed,
		Note:        e.Note,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   e.UpdatedAt.Format(time.RFC3339),
	}
}

// ScheduleListOutput 勤務表一覧出力
type ScheduleListOutput struct {
	// Schedules 勤務表一覧
	Schedules []ScheduleOutput `json:"schedules"`
	// Total 総件数
	Total int `json:"total"`
}

// ValidateResult 検証結果
type ValidateResult struct {
	// IsValid 有効フラグ
	IsValid bool `json:"is_valid"`
	// Violations 違反リスト
	Violations []ViolationOutput `json:"violations"`
}

// ViolationOutput 違反出力
type ViolationOutput struct {
	// Type 種別
	Type string `json:"type"`
	// Message メッセージ
	Message string `json:"message"`
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// Date 日付
	Date string `json:"date"`
	// Severity 重大度
	Severity string `json:"severity"`
}
