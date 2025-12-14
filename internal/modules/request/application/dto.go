// Package application 勤務希望アプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/request/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// CreateRequestPeriodInput 受付期間作成入力
type CreateRequestPeriodInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
	// StartDate 受付開始日
	StartDate string `json:"start_date"`
	// EndDate 受付終了日
	EndDate string `json:"end_date"`
	// MaxRequestsPerStaff スタッフあたり最大希望数
	MaxRequestsPerStaff int `json:"max_requests_per_staff"`
	// MaxRequestsPerDay 日あたり最大希望数
	MaxRequestsPerDay int `json:"max_requests_per_day"`
}

// Validate 入力検証
func (i *CreateRequestPeriodInput) Validate() error {
	if i.OrganizationID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDは必須です")
	}
	if i.TargetYear < 2020 || i.TargetYear > 2100 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象年が不正です")
	}
	if i.TargetMonth < 1 || i.TargetMonth > 12 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象月が不正です")
	}
	if i.StartDate == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付開始日は必須です")
	}
	if i.EndDate == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付終了日は必須です")
	}
	return nil
}

// RequestPeriodOutput 受付期間出力
type RequestPeriodOutput struct {
	// ID 受付期間ID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
	// TargetPeriodLabel 対象期間ラベル
	TargetPeriodLabel string `json:"target_period_label"`
	// StartDate 受付開始日
	StartDate string `json:"start_date"`
	// EndDate 受付終了日
	EndDate string `json:"end_date"`
	// MaxRequestsPerStaff スタッフあたり最大希望数
	MaxRequestsPerStaff int `json:"max_requests_per_staff"`
	// MaxRequestsPerDay 日あたり最大希望数
	MaxRequestsPerDay int `json:"max_requests_per_day"`
	// IsOpen 受付中フラグ
	IsOpen bool `json:"is_open"`
	// IsActive 受付期間内フラグ
	IsActive bool `json:"is_active"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToRequestPeriodOutput ドメインエンティティから出力DTOへ変換
func ToRequestPeriodOutput(p *domain.RequestPeriod) *RequestPeriodOutput {
	return &RequestPeriodOutput{
		ID:                  p.ID.String(),
		OrganizationID:      p.OrganizationID.String(),
		TargetYear:          p.TargetYear,
		TargetMonth:         p.TargetMonth,
		TargetPeriodLabel:   p.TargetPeriodLabel(),
		StartDate:           p.StartDate.Format("2006-01-02"),
		EndDate:             p.EndDate.Format("2006-01-02"),
		MaxRequestsPerStaff: p.MaxRequestsPerStaff,
		MaxRequestsPerDay:   p.MaxRequestsPerDay,
		IsOpen:              p.IsOpen,
		IsActive:            p.IsActive(),
		CreatedAt:           p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           p.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateShiftRequestInput 勤務希望作成入力
type CreateShiftRequestInput struct {
	// PeriodID 受付期間ID
	PeriodID string `json:"period_id"`
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// TargetDate 対象日
	TargetDate string `json:"target_date"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// RequestType 希望種別
	RequestType string `json:"request_type"`
	// Priority 優先度
	Priority string `json:"priority"`
	// Comment コメント
	Comment string `json:"comment"`
}

// Validate 入力検証
func (i *CreateShiftRequestInput) Validate() error {
	if i.PeriodID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付期間IDは必須です")
	}
	if i.StaffID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "スタッフIDは必須です")
	}
	if i.TargetDate == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象日は必須です")
	}
	return nil
}

// ShiftRequestOutput 勤務希望出力
type ShiftRequestOutput struct {
	// ID 勤務希望ID
	ID string `json:"id"`
	// PeriodID 受付期間ID
	PeriodID string `json:"period_id"`
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// StaffName スタッフ名
	StaffName string `json:"staff_name"`
	// TargetDate 対象日
	TargetDate string `json:"target_date"`
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// ShiftTypeName シフト種別名
	ShiftTypeName string `json:"shift_type_name"`
	// RequestType 希望種別
	RequestType string `json:"request_type"`
	// RequestTypeLabel 希望種別ラベル
	RequestTypeLabel string `json:"request_type_label"`
	// Priority 優先度
	Priority string `json:"priority"`
	// PriorityLabel 優先度ラベル
	PriorityLabel string `json:"priority_label"`
	// Comment コメント
	Comment string `json:"comment"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToShiftRequestOutput ドメインエンティティから出力DTOへ変換
func ToShiftRequestOutput(r *domain.ShiftRequest) *ShiftRequestOutput {
	shiftTypeID := ""
	if r.ShiftTypeID != nil {
		shiftTypeID = r.ShiftTypeID.String()
	}

	return &ShiftRequestOutput{
		ID:               r.ID.String(),
		PeriodID:         r.PeriodID.String(),
		StaffID:          r.StaffID.String(),
		TargetDate:       r.TargetDate.Format("2006-01-02"),
		ShiftTypeID:      shiftTypeID,
		RequestType:      r.RequestType.String(),
		RequestTypeLabel: r.RequestType.Label(),
		Priority:         r.Priority.String(),
		PriorityLabel:    r.Priority.Label(),
		Comment:          r.Comment,
		CreatedAt:        r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        r.UpdatedAt.Format(time.RFC3339),
	}
}

// ShiftRequestListOutput 勤務希望一覧出力
type ShiftRequestListOutput struct {
	// Requests 勤務希望一覧
	Requests []ShiftRequestOutput `json:"requests"`
	// Total 総件数
	Total int `json:"total"`
}
