// Package application シフトアプリケーション層
package application

import (
	"context"
	"log/slog"
	"time"

	"shiftmaster/internal/modules/shift/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// ShiftTypeUseCase シフト種別ユースケース
type ShiftTypeUseCase struct {
	repo   domain.ShiftTypeRepository
	logger *slog.Logger
}

// NewShiftTypeUseCase シフト種別ユースケース生成
func NewShiftTypeUseCase(repo domain.ShiftTypeRepository, logger *slog.Logger) *ShiftTypeUseCase {
	return &ShiftTypeUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Create シフト種別作成
func (u *ShiftTypeUseCase) Create(ctx context.Context, input *CreateShiftTypeInput) (*ShiftTypeOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	orgID, err := sharedDomain.ParseID(input.OrganizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	var shiftPatternID *sharedDomain.ID
	if input.ShiftPatternID != "" {
		id, parseErr := sharedDomain.ParseID(input.ShiftPatternID)
		if parseErr != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフトパターンIDが不正です")
		}
		shiftPatternID = &id
	}

	startTime, err := parseTime(input.StartTime)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "開始時刻の形式が不正です")
	}

	endTime, err := parseTime(input.EndTime)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "終了時刻の形式が不正です")
	}

	color := input.Color
	if color == "" {
		color = "#3B82F6"
	}

	now := time.Now()
	st := &domain.ShiftType{
		ID:              sharedDomain.NewID(),
		OrganizationID:  orgID,
		ShiftPatternID:  shiftPatternID,
		Name:            input.Name,
		Code:            input.Code,
		Color:           color,
		StartTime:       startTime,
		EndTime:         endTime,
		BreakMinutes:    input.BreakMinutes,
		HandoverMinutes: input.HandoverMinutes,
		IsNightShift:    input.IsNightShift,
		IsHoliday:       input.IsHoliday,
		SortOrder:       input.SortOrder,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.repo.Save(ctx, st); err != nil {
		u.logger.Error("シフト種別作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("シフト種別作成完了", "shift_type_id", st.ID, "name", st.Name)
	return ToShiftTypeOutput(st), nil
}

// Update シフト種別更新
func (u *ShiftTypeUseCase) Update(ctx context.Context, input *UpdateShiftTypeInput) (*ShiftTypeOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	id, err := sharedDomain.ParseID(input.ID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	st, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, sharedDomain.ErrNotFound
	}

	// シフトパターンID
	if input.ShiftPatternID != "" {
		patternID, parseErr := sharedDomain.ParseID(input.ShiftPatternID)
		if parseErr != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフトパターンIDが不正です")
		}
		st.ShiftPatternID = &patternID
	} else {
		st.ShiftPatternID = nil
	}

	if input.StartTime != "" {
		startTime, parseErr := parseTime(input.StartTime)
		if parseErr != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "開始時刻の形式が不正です")
		}
		st.StartTime = startTime
	}

	if input.EndTime != "" {
		endTime, parseErr := parseTime(input.EndTime)
		if parseErr != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "終了時刻の形式が不正です")
		}
		st.EndTime = endTime
	}

	st.Name = input.Name
	st.Code = input.Code
	if input.Color != "" {
		st.Color = input.Color
	}
	st.BreakMinutes = input.BreakMinutes
	st.HandoverMinutes = input.HandoverMinutes
	st.IsNightShift = input.IsNightShift
	st.IsHoliday = input.IsHoliday
	st.SortOrder = input.SortOrder
	st.UpdatedAt = time.Now()

	if err := u.repo.Save(ctx, st); err != nil {
		u.logger.Error("シフト種別更新失敗", "error", err)
		return nil, err
	}

	u.logger.Info("シフト種別更新完了", "shift_type_id", st.ID)
	return ToShiftTypeOutput(st), nil
}

// GetByID IDでシフト種別取得
func (u *ShiftTypeUseCase) GetByID(ctx context.Context, id string) (*ShiftTypeOutput, error) {
	stID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	st, err := u.repo.FindByID(ctx, stID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToShiftTypeOutput(st), nil
}

// List シフト種別一覧取得
// Deprecated: ListByOrganization を使用してください
func (u *ShiftTypeUseCase) List(ctx context.Context) (*ShiftTypeListOutput, error) {
	shiftTypes, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]ShiftTypeOutput, len(shiftTypes))
	for i, st := range shiftTypes {
		outputs[i] = *ToShiftTypeOutput(&st)
	}

	return &ShiftTypeListOutput{
		ShiftTypes: outputs,
		Total:      len(outputs),
	}, nil
}

// ListByOrganization 組織IDでシフト種別一覧取得（マルチテナント対応）
func (u *ShiftTypeUseCase) ListByOrganization(ctx context.Context, orgID string) (*ShiftTypeListOutput, error) {
	if orgID == "" {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが必要です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	shiftTypes, err := u.repo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ShiftTypeOutput, len(shiftTypes))
	for i, st := range shiftTypes {
		outputs[i] = *ToShiftTypeOutput(&st)
	}

	return &ShiftTypeListOutput{
		ShiftTypes: outputs,
		Total:      len(outputs),
	}, nil
}

// Delete シフト種別削除
func (u *ShiftTypeUseCase) Delete(ctx context.Context, id string) error {
	stID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	st, err := u.repo.FindByID(ctx, stID)
	if err != nil {
		return err
	}
	if st == nil {
		return sharedDomain.ErrNotFound
	}

	if err := u.repo.Delete(ctx, stID); err != nil {
		u.logger.Error("シフト種別削除失敗", "error", err)
		return err
	}

	u.logger.Info("シフト種別削除完了", "shift_type_id", stID)
	return nil
}

// parseTime 時刻文字列パース HH:MM形式
func parseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}
