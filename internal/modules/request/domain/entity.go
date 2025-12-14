// Package domain 勤務希望ドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// RequestPeriod 勤務希望受付期間エンティティ
type RequestPeriod struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// TargetYear 対象年
	TargetYear int
	// TargetMonth 対象月
	TargetMonth int
	// StartDate 受付開始日
	StartDate time.Time
	// EndDate 受付終了日
	EndDate time.Time
	// MaxRequestsPerStaff スタッフあたり最大希望数
	MaxRequestsPerStaff int
	// MaxRequestsPerDay 日あたり最大希望数
	MaxRequestsPerDay int
	// IsOpen 受付中フラグ
	IsOpen bool
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// IsActive 受付期間内判定
// 日付のみで比較（タイムゾーンの影響を受けないように）
func (p *RequestPeriod) IsActive() bool {
	if !p.IsOpen {
		return false
	}

	// 現在日付をUTCで取得し、日付部分のみで比較
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// StartDate/EndDateも日付部分のみで比較
	startDate := time.Date(p.StartDate.Year(), p.StartDate.Month(), p.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(p.EndDate.Year(), p.EndDate.Month(), p.EndDate.Day(), 23, 59, 59, 0, time.UTC)

	return !today.Before(startDate) && !today.After(endDate)
}

// TargetPeriodLabel 対象期間ラベル
func (p *RequestPeriod) TargetPeriodLabel() string {
	return time.Date(p.TargetYear, time.Month(p.TargetMonth), 1, 0, 0, 0, 0, time.Local).Format("2006年1月")
}

// ShiftRequest 勤務希望エンティティ
type ShiftRequest struct {
	// ID 一意識別子
	ID domain.ID
	// PeriodID 受付期間ID
	PeriodID domain.ID
	// StaffID スタッフID
	StaffID domain.ID
	// TargetDate 対象日
	TargetDate time.Time
	// ShiftTypeID シフト種別ID 指定なしの場合nil
	ShiftTypeID *domain.ID
	// RequestType 希望種別
	RequestType RequestType
	// Priority 優先度
	Priority RequestPriority
	// Comment コメント
	Comment string
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// RequestType 希望種別
type RequestType string

const (
	// RequestTypePreferred 希望 割り当てて欲しい
	RequestTypePreferred RequestType = "preferred"
	// RequestTypeAvoided 回避 割り当てたくない
	RequestTypeAvoided RequestType = "avoided"
	// RequestTypeFixed 固定 必ず割り当てる
	RequestTypeFixed RequestType = "fixed"
)

// String 文字列変換
func (t RequestType) String() string {
	return string(t)
}

// Label 表示ラベル
func (t RequestType) Label() string {
	switch t {
	case RequestTypePreferred:
		return "希望"
	case RequestTypeAvoided:
		return "回避"
	case RequestTypeFixed:
		return "固定"
	default:
		return "不明"
	}
}

// RequestPriority 優先度
type RequestPriority string

const (
	// PriorityRequired 必須
	PriorityRequired RequestPriority = "required"
	// PriorityOptional できれば
	PriorityOptional RequestPriority = "optional"
)

// String 文字列変換
func (p RequestPriority) String() string {
	return string(p)
}

// Label 表示ラベル
func (p RequestPriority) Label() string {
	switch p {
	case PriorityRequired:
		return "必須"
	case PriorityOptional:
		return "できれば"
	default:
		return "不明"
	}
}
