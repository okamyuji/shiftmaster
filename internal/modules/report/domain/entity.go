// Package domain レポートドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// Report レポートエンティティ
type Report struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Type レポート種別
	Type ReportType
	// Title タイトル
	Title string
	// TargetYear 対象年
	TargetYear int
	// TargetMonth 対象月
	TargetMonth int
	// Status 状態
	Status ReportStatus
	// GeneratedAt 生成日時
	GeneratedAt *time.Time
	// FilePath ファイルパス
	FilePath string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// ReportType レポート種別
type ReportType string

const (
	// ReportTypeSchedule 勤務表レポート
	ReportTypeSchedule ReportType = "schedule"
	// ReportTypeSummary 集計レポート
	ReportTypeSummary ReportType = "summary"
	// ReportTypeActual 実績レポート
	ReportTypeActual ReportType = "actual"
	// ReportTypeStaffList スタッフ一覧
	ReportTypeStaffList ReportType = "staff_list"
)

// String 文字列変換
func (t ReportType) String() string {
	return string(t)
}

// Label 表示ラベル
func (t ReportType) Label() string {
	switch t {
	case ReportTypeSchedule:
		return "勤務表"
	case ReportTypeSummary:
		return "集計レポート"
	case ReportTypeActual:
		return "勤務実績"
	case ReportTypeStaffList:
		return "スタッフ一覧"
	default:
		return "不明"
	}
}

// ReportStatus レポート状態
type ReportStatus string

const (
	// ReportStatusPending 生成待ち
	ReportStatusPending ReportStatus = "pending"
	// ReportStatusGenerating 生成中
	ReportStatusGenerating ReportStatus = "generating"
	// ReportStatusCompleted 完了
	ReportStatusCompleted ReportStatus = "completed"
	// ReportStatusFailed 失敗
	ReportStatusFailed ReportStatus = "failed"
)

// String 文字列変換
func (s ReportStatus) String() string {
	return string(s)
}

// Label 表示ラベル
func (s ReportStatus) Label() string {
	switch s {
	case ReportStatusPending:
		return "生成待ち"
	case ReportStatusGenerating:
		return "生成中"
	case ReportStatusCompleted:
		return "完了"
	case ReportStatusFailed:
		return "失敗"
	default:
		return "不明"
	}
}

// StaffSummary スタッフ集計
type StaffSummary struct {
	// StaffID スタッフID
	StaffID domain.ID
	// StaffName スタッフ名
	StaffName string
	// TotalWorkDays 総勤務日数
	TotalWorkDays int
	// TotalWorkMinutes 総勤務時間 分
	TotalWorkMinutes int
	// NightShiftCount 夜勤回数
	NightShiftCount int
	// HolidayCount 休日数
	HolidayCount int
	// OvertimeMinutes 残業時間 分
	OvertimeMinutes int
	// PaidLeaveUsed 有給休暇使用日数
	PaidLeaveUsed float64
}

// TotalWorkHours 総勤務時間 時間単位
func (s *StaffSummary) TotalWorkHours() float64 {
	return float64(s.TotalWorkMinutes) / 60.0
}

// OvertimeHours 残業時間 時間単位
func (s *StaffSummary) OvertimeHours() float64 {
	return float64(s.OvertimeMinutes) / 60.0
}

// DailySummary 日別集計
type DailySummary struct {
	// Date 日付
	Date time.Time
	// TotalStaff 総スタッフ数
	TotalStaff int
	// DayShiftCount 日勤人数
	DayShiftCount int
	// NightShiftCount 夜勤人数
	NightShiftCount int
	// HolidayCount 休日人数
	HolidayCount int
}

// DayOfWeek 曜日
func (s *DailySummary) DayOfWeek() string {
	weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}
	return weekdays[s.Date.Weekday()]
}

// ShiftTypeSummary シフト種別集計
type ShiftTypeSummary struct {
	// ShiftTypeID シフト種別ID
	ShiftTypeID domain.ID
	// ShiftTypeName シフト種別名
	ShiftTypeName string
	// ShiftTypeCode シフト種別コード
	ShiftTypeCode string
	// Count 件数
	Count int
	// TotalMinutes 合計時間 分
	TotalMinutes int
}

// TotalHours 合計時間 時間単位
func (s *ShiftTypeSummary) TotalHours() float64 {
	return float64(s.TotalMinutes) / 60.0
}

// MonthlySummary 月次集計
type MonthlySummary struct {
	// TargetYear 対象年
	TargetYear int
	// TargetMonth 対象月
	TargetMonth int
	// StaffSummaries スタッフ別集計
	StaffSummaries []StaffSummary
	// DailySummaries 日別集計
	DailySummaries []DailySummary
	// ShiftTypeSummaries シフト種別集計
	ShiftTypeSummaries []ShiftTypeSummary
	// TotalWorkDays 総勤務日数
	TotalWorkDays int
	// TotalWorkMinutes 総勤務時間 分
	TotalWorkMinutes int
	// AverageWorkMinutesPerDay 日平均勤務時間 分
	AverageWorkMinutesPerDay float64
}
