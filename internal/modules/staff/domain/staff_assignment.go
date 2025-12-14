// Package domain スタッフドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// StaffAssignment スタッフ所属エンティティ
// スタッフと チーム/職種/職位 の多対多関係を表現
type StaffAssignment struct {
	// ID 一意識別子
	ID domain.ID
	// StaffID スタッフID
	StaffID domain.ID
	// TeamID チームID 任意
	TeamID *domain.ID
	// JobTypeID 職種ID 任意
	JobTypeID *domain.ID
	// PositionID 職位ID 任意
	PositionID *domain.ID
	// IsPrimary 主たる所属かどうか
	IsPrimary bool
	// StartDate 所属開始日
	StartDate *time.Time
	// EndDate 所属終了日 異動対応
	EndDate *time.Time
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// NewStaffAssignment スタッフ所属生成
func NewStaffAssignment(staffID domain.ID, teamID, jobTypeID, positionID *domain.ID, isPrimary bool) (*StaffAssignment, error) {
	// 少なくとも1つは指定されている必要がある
	if teamID == nil && jobTypeID == nil && positionID == nil {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "チーム、職種、職位のいずれかを指定してください")
	}

	now := time.Now()
	return &StaffAssignment{
		ID:         domain.NewID(),
		StaffID:    staffID,
		TeamID:     teamID,
		JobTypeID:  jobTypeID,
		PositionID: positionID,
		IsPrimary:  isPrimary,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// SetDateRange 所属期間設定
func (a *StaffAssignment) SetDateRange(startDate, endDate *time.Time) error {
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return domain.NewDomainError(domain.ErrCodeValidation, "終了日は開始日より後に設定してください")
	}
	a.StartDate = startDate
	a.EndDate = endDate
	a.UpdatedAt = time.Now()
	return nil
}

// IsActiveOn 指定日に有効か判定
func (a *StaffAssignment) IsActiveOn(date time.Time) bool {
	// 開始日が設定されていて、指定日より後なら無効
	if a.StartDate != nil && date.Before(*a.StartDate) {
		return false
	}
	// 終了日が設定されていて、指定日より前なら無効
	if a.EndDate != nil && date.After(*a.EndDate) {
		return false
	}
	return true
}

// IsCurrentlyActive 現在有効か判定
func (a *StaffAssignment) IsCurrentlyActive() bool {
	return a.IsActiveOn(time.Now())
}

// SetAsPrimary 主たる所属に設定
func (a *StaffAssignment) SetAsPrimary() {
	a.IsPrimary = true
	a.UpdatedAt = time.Now()
}

// UnsetAsPrimary 主たる所属を解除
func (a *StaffAssignment) UnsetAsPrimary() {
	a.IsPrimary = false
	a.UpdatedAt = time.Now()
}

// End 所属終了
func (a *StaffAssignment) End(endDate time.Time) {
	a.EndDate = &endDate
	a.UpdatedAt = time.Now()
}

// StaffAssignmentRepository スタッフ所属リポジトリインターフェース
type StaffAssignmentRepository interface {
	// FindByID ID検索
	FindByID(ctx interface{}, id domain.ID) (*StaffAssignment, error)
	// FindByStaffID スタッフID検索
	FindByStaffID(ctx interface{}, staffID domain.ID) ([]StaffAssignment, error)
	// FindActiveByStaffID スタッフIDで有効な所属を検索
	FindActiveByStaffID(ctx interface{}, staffID domain.ID, date time.Time) ([]StaffAssignment, error)
	// FindByTeamID チームID検索
	FindByTeamID(ctx interface{}, teamID domain.ID) ([]StaffAssignment, error)
	// FindByJobTypeID 職種ID検索
	FindByJobTypeID(ctx interface{}, jobTypeID domain.ID) ([]StaffAssignment, error)
	// FindByPositionID 職位ID検索
	FindByPositionID(ctx interface{}, positionID domain.ID) ([]StaffAssignment, error)
	// Save 保存
	Save(ctx interface{}, assignment *StaffAssignment) error
	// Delete 削除
	Delete(ctx interface{}, id domain.ID) error
}
