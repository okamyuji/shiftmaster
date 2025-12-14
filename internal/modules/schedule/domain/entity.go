// Package domain 勤務表ドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// Schedule 勤務表エンティティ
type Schedule struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// TargetYear 対象年
	TargetYear int
	// TargetMonth 対象月
	TargetMonth int
	// Status 状態
	Status ScheduleStatus
	// PublishedAt 公開日時
	PublishedAt *time.Time
	// Entries エントリ一覧
	Entries []ScheduleEntry
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// TargetPeriodLabel 対象期間ラベル
func (s *Schedule) TargetPeriodLabel() string {
	return time.Date(s.TargetYear, time.Month(s.TargetMonth), 1, 0, 0, 0, 0, time.Local).Format("2006年1月")
}

// DaysInMonth 月の日数
func (s *Schedule) DaysInMonth() int {
	return time.Date(s.TargetYear, time.Month(s.TargetMonth)+1, 0, 0, 0, 0, 0, time.Local).Day()
}

// StartDate 開始日
func (s *Schedule) StartDate() time.Time {
	return time.Date(s.TargetYear, time.Month(s.TargetMonth), 1, 0, 0, 0, 0, time.Local)
}

// EndDate 終了日
func (s *Schedule) EndDate() time.Time {
	return time.Date(s.TargetYear, time.Month(s.TargetMonth), s.DaysInMonth(), 0, 0, 0, 0, time.Local)
}

// ScheduleStatus 勤務表状態
type ScheduleStatus string

const (
	// StatusDraft 下書き
	StatusDraft ScheduleStatus = "draft"
	// StatusInProgress 作成中
	StatusInProgress ScheduleStatus = "in_progress"
	// StatusCompleted 作成完了
	StatusCompleted ScheduleStatus = "completed"
	// StatusPublished 公開済み
	StatusPublished ScheduleStatus = "published"
)

// String 文字列変換
func (s ScheduleStatus) String() string {
	return string(s)
}

// Label 表示ラベル
func (s ScheduleStatus) Label() string {
	switch s {
	case StatusDraft:
		return "下書き"
	case StatusInProgress:
		return "作成中"
	case StatusCompleted:
		return "作成完了"
	case StatusPublished:
		return "公開済み"
	default:
		return "不明"
	}
}

// ScheduleEntry 勤務表エントリエンティティ
type ScheduleEntry struct {
	// ID 一意識別子
	ID domain.ID
	// ScheduleID 勤務表ID
	ScheduleID domain.ID
	// StaffID スタッフID
	StaffID domain.ID
	// TargetDate 対象日
	TargetDate time.Time
	// ShiftTypeID シフト種別ID
	ShiftTypeID *domain.ID
	// IsConfirmed 確定フラグ
	IsConfirmed bool
	// Note 備考
	Note string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// ActualRecord 勤務実績エンティティ
type ActualRecord struct {
	// ID 一意識別子
	ID domain.ID
	// ScheduleEntryID 勤務表エントリID
	ScheduleEntryID domain.ID
	// ActualStartTime 実際の開始時刻
	ActualStartTime *time.Time
	// ActualEndTime 実際の終了時刻
	ActualEndTime *time.Time
	// ActualBreakMinutes 実際の休憩時間
	ActualBreakMinutes int
	// OvertimeMinutes 残業時間
	OvertimeMinutes int
	// Note 備考
	Note string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// ActualWorkingMinutes 実際の実働時間
func (r *ActualRecord) ActualWorkingMinutes() int {
	if r.ActualStartTime == nil || r.ActualEndTime == nil {
		return 0
	}
	total := int(r.ActualEndTime.Sub(*r.ActualStartTime).Minutes())
	return total - r.ActualBreakMinutes
}

// ScheduleOptimizer 勤務表最適化インターフェース AI拡張用
type ScheduleOptimizer interface {
	// Optimize 勤務表最適化
	Optimize(ctx interface{}, input OptimizeInput) (*OptimizeResult, error)
}

// OptimizeInput 最適化入力
type OptimizeInput struct {
	// ScheduleID 勤務表ID
	ScheduleID domain.ID
	// DateRange 対象日範囲
	DateRange domain.DateRange
	// Constraints 制約条件
	Constraints []Constraint
	// Options 最適化オプション
	Options OptimizeOptions
}

// OptimizeOptions 最適化オプション
type OptimizeOptions struct {
	// MaxIterations 最大反復回数
	MaxIterations int
	// TimeoutSeconds タイムアウト秒数
	TimeoutSeconds int
	// PrioritizeRequests 希望優先
	PrioritizeRequests bool
}

// OptimizeResult 最適化結果
type OptimizeResult struct {
	// Success 成功フラグ
	Success bool
	// Entries 生成エントリ
	Entries []ScheduleEntry
	// Score スコア
	Score float64
	// Violations 違反リスト
	Violations []ConstraintViolation
	// Duration 処理時間
	Duration time.Duration
}

// Constraint 制約条件
type Constraint struct {
	// Type 制約種別
	Type string
	// Priority 優先度
	Priority int
	// Config 設定JSON
	Config string
}

// ConstraintViolation 制約違反
type ConstraintViolation struct {
	// ConstraintType 制約種別
	ConstraintType string
	// Message メッセージ
	Message string
	// StaffID 関連スタッフID
	StaffID *domain.ID
	// Date 関連日付
	Date *time.Time
	// Severity 重大度
	Severity string
}
