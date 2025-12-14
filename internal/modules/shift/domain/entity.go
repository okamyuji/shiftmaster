// Package domain シフトドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// ShiftType シフト種別エンティティ
type ShiftType struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// ShiftPatternID シフトパターンID（直）への参照 オプション
	ShiftPatternID *domain.ID
	// Name シフト名
	Name string
	// Code シフトコード 短縮表示用
	Code string
	// Color 表示色 HEX形式
	Color string
	// StartTime 開始時刻
	StartTime time.Time
	// EndTime 終了時刻
	EndTime time.Time
	// BreakMinutes 休憩時間 分
	BreakMinutes int
	// HandoverMinutes 申し送り時間 分（引き継ぎ用のオーバーラップ時間）
	HandoverMinutes int
	// IsNightShift 夜勤フラグ
	IsNightShift bool
	// IsHoliday 休日フラグ
	IsHoliday bool
	// SortOrder 表示順
	SortOrder int
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// WorkingMinutes 実働時間 分単位（申し送り時間を含む）
func (s *ShiftType) WorkingMinutes() int {
	total := s.TotalMinutes()
	if total <= 0 {
		return 0
	}
	return total - s.BreakMinutes
}

// EffectiveWorkingMinutes 正味実働時間 分単位（申し送り時間を除く）
func (s *ShiftType) EffectiveWorkingMinutes() int {
	working := s.WorkingMinutes()
	if working <= s.HandoverMinutes {
		return 0
	}
	return working - s.HandoverMinutes
}

// HasHandover 申し送り時間があるかどうか
func (s *ShiftType) HasHandover() bool {
	return s.HandoverMinutes > 0
}

// ActualStartTime 実際の出勤開始時刻（申し送り開始時刻）
// 申し送りがある場合、通常の開始時刻より早く出勤する必要がある
func (s *ShiftType) ActualStartTime() time.Time {
	if s.HandoverMinutes > 0 {
		return s.StartTime.Add(-time.Duration(s.HandoverMinutes) * time.Minute)
	}
	return s.StartTime
}

// ActualStartTimeString 実際の出勤開始時刻文字列
func (s *ShiftType) ActualStartTimeString() string {
	return s.ActualStartTime().Format("15:04")
}

// TotalMinutes 拘束時間 分単位
func (s *ShiftType) TotalMinutes() int {
	start := s.StartTime.Hour()*60 + s.StartTime.Minute()
	end := s.EndTime.Hour()*60 + s.EndTime.Minute()

	if end < start {
		// 日跨ぎ
		return (24*60 - start) + end
	}
	return end - start
}

// WorkingHours 実働時間 時間単位
func (s *ShiftType) WorkingHours() float64 {
	return float64(s.WorkingMinutes()) / 60.0
}

// StartTimeString 開始時刻文字列
func (s *ShiftType) StartTimeString() string {
	return s.StartTime.Format("15:04")
}

// EndTimeString 終了時刻文字列
func (s *ShiftType) EndTimeString() string {
	return s.EndTime.Format("15:04")
}

// ShiftPattern シフトパターン（直）エンティティ
// 1直、2直、3直などの交代勤務グループを表す
type ShiftPattern struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name パターン名 例: 1直、2直、3直
	Name string
	// Code パターンコード 例: D1, D2, D3
	Code string
	// Description 説明
	Description string
	// RotationType ローテーション種別 二交代/三交代/四直三交代など
	RotationType RotationType
	// Color 表示色 HEX形式
	Color string
	// SortOrder 表示順
	SortOrder int
	// IsActive 有効フラグ
	IsActive bool
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// RotationType ローテーション種別
type RotationType string

const (
	// RotationTypeTwoShift 二交代制
	RotationTypeTwoShift RotationType = "two_shift"
	// RotationTypeThreeShift 三交代制
	RotationTypeThreeShift RotationType = "three_shift"
	// RotationTypeFourShiftThreeTeam 四直三交代制（工場等で使用）
	RotationTypeFourShiftThreeTeam RotationType = "four_shift_three_team"
	// RotationTypeCustom カスタム
	RotationTypeCustom RotationType = "custom"
)

// String 文字列変換
func (t RotationType) String() string {
	return string(t)
}

// Label 表示ラベル
func (t RotationType) Label() string {
	switch t {
	case RotationTypeTwoShift:
		return "二交代制"
	case RotationTypeThreeShift:
		return "三交代制"
	case RotationTypeFourShiftThreeTeam:
		return "四直三交代制"
	case RotationTypeCustom:
		return "カスタム"
	default:
		return "不明"
	}
}

// IsValid 有効なローテーション種別かチェック
func (t RotationType) IsValid() bool {
	switch t {
	case RotationTypeTwoShift, RotationTypeThreeShift, RotationTypeFourShiftThreeTeam, RotationTypeCustom:
		return true
	default:
		return false
	}
}

// ShiftRule シフトルールエンティティ
type ShiftRule struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name ルール名
	Name string
	// Description 説明
	Description string
	// RuleType ルール種別
	RuleType ShiftRuleType
	// Priority 優先度 高いほど優先
	Priority int
	// IsActive 有効フラグ
	IsActive bool
	// Config ルール設定 JSON形式
	Config string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// ShiftRuleType シフトルール種別
type ShiftRuleType string

const (
	// RuleTypeMinStaff 最小配置人数
	RuleTypeMinStaff ShiftRuleType = "min_staff"
	// RuleTypeMaxStaff 最大配置人数
	RuleTypeMaxStaff ShiftRuleType = "max_staff"
	// RuleTypeConsecutive 連続勤務制限
	RuleTypeConsecutive ShiftRuleType = "consecutive"
	// RuleTypeInterval シフト間隔
	RuleTypeInterval ShiftRuleType = "interval"
	// RuleTypeSkillRequired 必須スキル
	RuleTypeSkillRequired ShiftRuleType = "skill_required"
	// RuleTypeNightLimit 夜勤回数制限
	RuleTypeNightLimit ShiftRuleType = "night_limit"
	// RuleTypeWeeklyHours 週間労働時間制限
	RuleTypeWeeklyHours ShiftRuleType = "weekly_hours"
	// RuleTypeMonthlyHours 月間労働時間制限
	RuleTypeMonthlyHours ShiftRuleType = "monthly_hours"
)

// String 文字列変換
func (t ShiftRuleType) String() string {
	return string(t)
}

// Label 表示ラベル
func (t ShiftRuleType) Label() string {
	switch t {
	case RuleTypeMinStaff:
		return "最小配置人数"
	case RuleTypeMaxStaff:
		return "最大配置人数"
	case RuleTypeConsecutive:
		return "連続勤務制限"
	case RuleTypeInterval:
		return "シフト間隔"
	case RuleTypeSkillRequired:
		return "必須スキル"
	case RuleTypeNightLimit:
		return "夜勤回数制限"
	case RuleTypeWeeklyHours:
		return "週間労働時間"
	case RuleTypeMonthlyHours:
		return "月間労働時間"
	default:
		return "不明"
	}
}
