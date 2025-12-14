// Package application レポートアプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/report/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// GenerateReportInput レポート生成入力
type GenerateReportInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Type レポート種別
	Type string `json:"type"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
}

// Validate 入力検証
func (i *GenerateReportInput) Validate() error {
	if i.OrganizationID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDは必須です")
	}
	if i.Type == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "レポート種別は必須です")
	}
	if i.TargetYear < 2020 || i.TargetYear > 2100 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象年が不正です")
	}
	if i.TargetMonth < 1 || i.TargetMonth > 12 {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象月が不正です")
	}
	return nil
}

// ReportOutput レポート出力
type ReportOutput struct {
	// ID レポートID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Type レポート種別
	Type string `json:"type"`
	// TypeLabel レポート種別ラベル
	TypeLabel string `json:"type_label"`
	// Title タイトル
	Title string `json:"title"`
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
	// Status 状態
	Status string `json:"status"`
	// StatusLabel 状態ラベル
	StatusLabel string `json:"status_label"`
	// GeneratedAt 生成日時
	GeneratedAt string `json:"generated_at"`
	// FilePath ファイルパス
	FilePath string `json:"file_path"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToReportOutput ドメインエンティティから出力DTOへ変換
func ToReportOutput(r *domain.Report) *ReportOutput {
	generatedAt := ""
	if r.GeneratedAt != nil {
		generatedAt = r.GeneratedAt.Format(time.RFC3339)
	}

	return &ReportOutput{
		ID:             r.ID.String(),
		OrganizationID: r.OrganizationID.String(),
		Type:           r.Type.String(),
		TypeLabel:      r.Type.Label(),
		Title:          r.Title,
		TargetYear:     r.TargetYear,
		TargetMonth:    r.TargetMonth,
		Status:         r.Status.String(),
		StatusLabel:    r.Status.Label(),
		GeneratedAt:    generatedAt,
		FilePath:       r.FilePath,
		CreatedAt:      r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      r.UpdatedAt.Format(time.RFC3339),
	}
}

// MonthlySummaryOutput 月次集計出力
type MonthlySummaryOutput struct {
	// TargetYear 対象年
	TargetYear int `json:"target_year"`
	// TargetMonth 対象月
	TargetMonth int `json:"target_month"`
	// TargetPeriodLabel 対象期間ラベル
	TargetPeriodLabel string `json:"target_period_label"`
	// StaffSummaries スタッフ別集計
	StaffSummaries []StaffSummaryOutput `json:"staff_summaries"`
	// DailySummaries 日別集計
	DailySummaries []DailySummaryOutput `json:"daily_summaries"`
	// ShiftTypeSummaries シフト種別集計
	ShiftTypeSummaries []ShiftTypeSummaryOutput `json:"shift_type_summaries"`
	// TotalWorkDays 総勤務日数
	TotalWorkDays int `json:"total_work_days"`
	// TotalWorkHours 総勤務時間
	TotalWorkHours float64 `json:"total_work_hours"`
	// AverageWorkHoursPerDay 日平均勤務時間
	AverageWorkHoursPerDay float64 `json:"average_work_hours_per_day"`
}

// StaffSummaryOutput スタッフ集計出力
type StaffSummaryOutput struct {
	// StaffID スタッフID
	StaffID string `json:"staff_id"`
	// StaffName スタッフ名
	StaffName string `json:"staff_name"`
	// TotalWorkDays 総勤務日数
	TotalWorkDays int `json:"total_work_days"`
	// TotalWorkHours 総勤務時間
	TotalWorkHours float64 `json:"total_work_hours"`
	// NightShiftCount 夜勤回数
	NightShiftCount int `json:"night_shift_count"`
	// HolidayCount 休日数
	HolidayCount int `json:"holiday_count"`
	// OvertimeHours 残業時間
	OvertimeHours float64 `json:"overtime_hours"`
	// PaidLeaveUsed 有給休暇使用日数
	PaidLeaveUsed float64 `json:"paid_leave_used"`
}

// ToStaffSummaryOutput ドメインエンティティから出力DTOへ変換
func ToStaffSummaryOutput(s *domain.StaffSummary) *StaffSummaryOutput {
	return &StaffSummaryOutput{
		StaffID:         s.StaffID.String(),
		StaffName:       s.StaffName,
		TotalWorkDays:   s.TotalWorkDays,
		TotalWorkHours:  s.TotalWorkHours(),
		NightShiftCount: s.NightShiftCount,
		HolidayCount:    s.HolidayCount,
		OvertimeHours:   s.OvertimeHours(),
		PaidLeaveUsed:   s.PaidLeaveUsed,
	}
}

// DailySummaryOutput 日別集計出力
type DailySummaryOutput struct {
	// Date 日付
	Date string `json:"date"`
	// DayOfWeek 曜日
	DayOfWeek string `json:"day_of_week"`
	// TotalStaff 総スタッフ数
	TotalStaff int `json:"total_staff"`
	// DayShiftCount 日勤人数
	DayShiftCount int `json:"day_shift_count"`
	// NightShiftCount 夜勤人数
	NightShiftCount int `json:"night_shift_count"`
	// HolidayCount 休日人数
	HolidayCount int `json:"holiday_count"`
}

// ToDailySummaryOutput ドメインエンティティから出力DTOへ変換
func ToDailySummaryOutput(s *domain.DailySummary) *DailySummaryOutput {
	return &DailySummaryOutput{
		Date:            s.Date.Format("2006-01-02"),
		DayOfWeek:       s.DayOfWeek(),
		TotalStaff:      s.TotalStaff,
		DayShiftCount:   s.DayShiftCount,
		NightShiftCount: s.NightShiftCount,
		HolidayCount:    s.HolidayCount,
	}
}

// ShiftTypeSummaryOutput シフト種別集計出力
type ShiftTypeSummaryOutput struct {
	// ShiftTypeID シフト種別ID
	ShiftTypeID string `json:"shift_type_id"`
	// ShiftTypeName シフト種別名
	ShiftTypeName string `json:"shift_type_name"`
	// ShiftTypeCode シフト種別コード
	ShiftTypeCode string `json:"shift_type_code"`
	// Count 件数
	Count int `json:"count"`
	// TotalHours 合計時間
	TotalHours float64 `json:"total_hours"`
}

// ToShiftTypeSummaryOutput ドメインエンティティから出力DTOへ変換
func ToShiftTypeSummaryOutput(s *domain.ShiftTypeSummary) *ShiftTypeSummaryOutput {
	return &ShiftTypeSummaryOutput{
		ShiftTypeID:   s.ShiftTypeID.String(),
		ShiftTypeName: s.ShiftTypeName,
		ShiftTypeCode: s.ShiftTypeCode,
		Count:         s.Count,
		TotalHours:    s.TotalHours(),
	}
}

// ToMonthlySummaryOutput ドメインエンティティから出力DTOへ変換
func ToMonthlySummaryOutput(s *domain.MonthlySummary) *MonthlySummaryOutput {
	staffSummaries := make([]StaffSummaryOutput, len(s.StaffSummaries))
	for i, ss := range s.StaffSummaries {
		staffSummaries[i] = *ToStaffSummaryOutput(&ss)
	}

	dailySummaries := make([]DailySummaryOutput, len(s.DailySummaries))
	for i, ds := range s.DailySummaries {
		dailySummaries[i] = *ToDailySummaryOutput(&ds)
	}

	shiftTypeSummaries := make([]ShiftTypeSummaryOutput, len(s.ShiftTypeSummaries))
	for i, sts := range s.ShiftTypeSummaries {
		shiftTypeSummaries[i] = *ToShiftTypeSummaryOutput(&sts)
	}

	return &MonthlySummaryOutput{
		TargetYear:             s.TargetYear,
		TargetMonth:            s.TargetMonth,
		TargetPeriodLabel:      time.Date(s.TargetYear, time.Month(s.TargetMonth), 1, 0, 0, 0, 0, time.Local).Format("2006年1月"),
		StaffSummaries:         staffSummaries,
		DailySummaries:         dailySummaries,
		ShiftTypeSummaries:     shiftTypeSummaries,
		TotalWorkDays:          s.TotalWorkDays,
		TotalWorkHours:         float64(s.TotalWorkMinutes) / 60.0,
		AverageWorkHoursPerDay: s.AverageWorkMinutesPerDay / 60.0,
	}
}
