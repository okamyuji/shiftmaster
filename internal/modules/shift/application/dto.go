// Package application シフトアプリケーション層
package application

import (
	"time"

	"shiftmaster/internal/modules/shift/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// CreateShiftTypeInput シフト種別作成入力
type CreateShiftTypeInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// ShiftPatternID シフトパターンID（直）オプション
	ShiftPatternID string `json:"shift_pattern_id"`
	// Name シフト名
	Name string `json:"name"`
	// Code シフトコード
	Code string `json:"code"`
	// Color 表示色
	Color string `json:"color"`
	// StartTime 開始時刻 HH:MM形式
	StartTime string `json:"start_time"`
	// EndTime 終了時刻 HH:MM形式
	EndTime string `json:"end_time"`
	// BreakMinutes 休憩時間
	BreakMinutes int `json:"break_minutes"`
	// HandoverMinutes 申し送り時間
	HandoverMinutes int `json:"handover_minutes"`
	// IsNightShift 夜勤フラグ
	IsNightShift bool `json:"is_night_shift"`
	// IsHoliday 休日フラグ
	IsHoliday bool `json:"is_holiday"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
}

// Validate 入力検証
func (i *CreateShiftTypeInput) Validate() error {
	if i.OrganizationID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDは必須です")
	}
	if i.Name == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフト名は必須です")
	}
	if i.Code == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフトコードは必須です")
	}
	if i.StartTime == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "開始時刻は必須です")
	}
	if i.EndTime == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "終了時刻は必須です")
	}
	return nil
}

// UpdateShiftTypeInput シフト種別更新入力
type UpdateShiftTypeInput struct {
	// ID シフト種別ID
	ID string `json:"id"`
	// ShiftPatternID シフトパターンID（直）オプション
	ShiftPatternID string `json:"shift_pattern_id"`
	// Name シフト名
	Name string `json:"name"`
	// Code シフトコード
	Code string `json:"code"`
	// Color 表示色
	Color string `json:"color"`
	// StartTime 開始時刻
	StartTime string `json:"start_time"`
	// EndTime 終了時刻
	EndTime string `json:"end_time"`
	// BreakMinutes 休憩時間
	BreakMinutes int `json:"break_minutes"`
	// HandoverMinutes 申し送り時間
	HandoverMinutes int `json:"handover_minutes"`
	// IsNightShift 夜勤フラグ
	IsNightShift bool `json:"is_night_shift"`
	// IsHoliday 休日フラグ
	IsHoliday bool `json:"is_holiday"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
}

// Validate 入力検証
func (i *UpdateShiftTypeInput) Validate() error {
	if i.ID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDは必須です")
	}
	if i.Name == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフト名は必須です")
	}
	if i.Code == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフトコードは必須です")
	}
	return nil
}

// ShiftTypeOutput シフト種別出力
type ShiftTypeOutput struct {
	// ID シフト種別ID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// ShiftPatternID シフトパターンID（直）
	ShiftPatternID string `json:"shift_pattern_id,omitempty"`
	// Name シフト名
	Name string `json:"name"`
	// Code シフトコード
	Code string `json:"code"`
	// Color 表示色
	Color string `json:"color"`
	// StartTime 開始時刻
	StartTime string `json:"start_time"`
	// EndTime 終了時刻
	EndTime string `json:"end_time"`
	// ActualStartTime 実際の出勤開始時刻（申し送り含む）
	ActualStartTime string `json:"actual_start_time"`
	// BreakMinutes 休憩時間
	BreakMinutes int `json:"break_minutes"`
	// HandoverMinutes 申し送り時間
	HandoverMinutes int `json:"handover_minutes"`
	// WorkingMinutes 実働時間
	WorkingMinutes int `json:"working_minutes"`
	// EffectiveWorkingMinutes 正味実働時間（申し送り除く）
	EffectiveWorkingMinutes int `json:"effective_working_minutes"`
	// WorkingHours 実働時間 時間単位
	WorkingHours float64 `json:"working_hours"`
	// HasHandover 申し送りありフラグ
	HasHandover bool `json:"has_handover"`
	// IsNightShift 夜勤フラグ
	IsNightShift bool `json:"is_night_shift"`
	// IsHoliday 休日フラグ
	IsHoliday bool `json:"is_holiday"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToShiftTypeOutput ドメインエンティティから出力DTOへ変換
func ToShiftTypeOutput(st *domain.ShiftType) *ShiftTypeOutput {
	var shiftPatternID string
	if st.ShiftPatternID != nil {
		shiftPatternID = st.ShiftPatternID.String()
	}

	return &ShiftTypeOutput{
		ID:                      st.ID.String(),
		OrganizationID:          st.OrganizationID.String(),
		ShiftPatternID:          shiftPatternID,
		Name:                    st.Name,
		Code:                    st.Code,
		Color:                   st.Color,
		StartTime:               st.StartTimeString(),
		EndTime:                 st.EndTimeString(),
		ActualStartTime:         st.ActualStartTimeString(),
		BreakMinutes:            st.BreakMinutes,
		HandoverMinutes:         st.HandoverMinutes,
		WorkingMinutes:          st.WorkingMinutes(),
		EffectiveWorkingMinutes: st.EffectiveWorkingMinutes(),
		WorkingHours:            st.WorkingHours(),
		HasHandover:             st.HasHandover(),
		IsNightShift:            st.IsNightShift,
		IsHoliday:               st.IsHoliday,
		SortOrder:               st.SortOrder,
		CreatedAt:               st.CreatedAt.Format(time.RFC3339),
		UpdatedAt:               st.UpdatedAt.Format(time.RFC3339),
	}
}

// ShiftTypeListOutput シフト種別一覧出力
type ShiftTypeListOutput struct {
	// ShiftTypes シフト種別一覧
	ShiftTypes []ShiftTypeOutput `json:"shift_types"`
	// Total 総件数
	Total int `json:"total"`
}

// CreateShiftPatternInput シフトパターン（直）作成入力
type CreateShiftPatternInput struct {
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Name パターン名
	Name string `json:"name"`
	// Code パターンコード
	Code string `json:"code"`
	// Description 説明
	Description string `json:"description"`
	// RotationType ローテーション種別
	RotationType string `json:"rotation_type"`
	// Color 表示色
	Color string `json:"color"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
}

// Validate 入力検証
func (i *CreateShiftPatternInput) Validate() error {
	if i.OrganizationID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDは必須です")
	}
	if i.Name == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パターン名は必須です")
	}
	if i.Code == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パターンコードは必須です")
	}
	return nil
}

// UpdateShiftPatternInput シフトパターン（直）更新入力
type UpdateShiftPatternInput struct {
	// ID パターンID
	ID string `json:"id"`
	// Name パターン名
	Name string `json:"name"`
	// Code パターンコード
	Code string `json:"code"`
	// Description 説明
	Description string `json:"description"`
	// RotationType ローテーション種別
	RotationType string `json:"rotation_type"`
	// Color 表示色
	Color string `json:"color"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
}

// Validate 入力検証
func (i *UpdateShiftPatternInput) Validate() error {
	if i.ID == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDは必須です")
	}
	if i.Name == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パターン名は必須です")
	}
	if i.Code == "" {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "パターンコードは必須です")
	}
	return nil
}

// ShiftPatternOutput シフトパターン（直）出力
type ShiftPatternOutput struct {
	// ID パターンID
	ID string `json:"id"`
	// OrganizationID 組織ID
	OrganizationID string `json:"organization_id"`
	// Name パターン名
	Name string `json:"name"`
	// Code パターンコード
	Code string `json:"code"`
	// Description 説明
	Description string `json:"description"`
	// RotationType ローテーション種別
	RotationType string `json:"rotation_type"`
	// RotationTypeLabel ローテーション種別ラベル
	RotationTypeLabel string `json:"rotation_type_label"`
	// Color 表示色
	Color string `json:"color"`
	// SortOrder 表示順
	SortOrder int `json:"sort_order"`
	// IsActive 有効フラグ
	IsActive bool `json:"is_active"`
	// CreatedAt 作成日時
	CreatedAt string `json:"created_at"`
	// UpdatedAt 更新日時
	UpdatedAt string `json:"updated_at"`
}

// ToShiftPatternOutput ドメインエンティティから出力DTOへ変換
func ToShiftPatternOutput(sp *domain.ShiftPattern) *ShiftPatternOutput {
	return &ShiftPatternOutput{
		ID:                sp.ID.String(),
		OrganizationID:    sp.OrganizationID.String(),
		Name:              sp.Name,
		Code:              sp.Code,
		Description:       sp.Description,
		RotationType:      sp.RotationType.String(),
		RotationTypeLabel: sp.RotationType.Label(),
		Color:             sp.Color,
		SortOrder:         sp.SortOrder,
		IsActive:          sp.IsActive,
		CreatedAt:         sp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         sp.UpdatedAt.Format(time.RFC3339),
	}
}

// ShiftPatternListOutput シフトパターン一覧出力
type ShiftPatternListOutput struct {
	// ShiftPatterns シフトパターン一覧
	ShiftPatterns []ShiftPatternOutput `json:"shift_patterns"`
	// Total 総件数
	Total int `json:"total"`
}

// RotationTypeOption ローテーション種別選択肢
type RotationTypeOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetRotationTypeOptions ローテーション種別選択肢取得
func GetRotationTypeOptions() []RotationTypeOption {
	return []RotationTypeOption{
		{Value: string(domain.RotationTypeTwoShift), Label: domain.RotationTypeTwoShift.Label()},
		{Value: string(domain.RotationTypeThreeShift), Label: domain.RotationTypeThreeShift.Label()},
		{Value: string(domain.RotationTypeFourShiftThreeTeam), Label: domain.RotationTypeFourShiftThreeTeam.Label()},
		{Value: string(domain.RotationTypeCustom), Label: domain.RotationTypeCustom.Label()},
	}
}
