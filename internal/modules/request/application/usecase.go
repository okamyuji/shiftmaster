// Package application 勤務希望アプリケーション層
package application

import (
	"context"
	"log/slog"
	"time"

	"shiftmaster/internal/modules/request/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// RequestPeriodUseCase 受付期間ユースケース
type RequestPeriodUseCase struct {
	periodRepo  domain.RequestPeriodRepository
	requestRepo domain.ShiftRequestRepository
	logger      *slog.Logger
}

// NewRequestPeriodUseCase 受付期間ユースケース生成
func NewRequestPeriodUseCase(
	periodRepo domain.RequestPeriodRepository,
	requestRepo domain.ShiftRequestRepository,
	logger *slog.Logger,
) *RequestPeriodUseCase {
	return &RequestPeriodUseCase{
		periodRepo:  periodRepo,
		requestRepo: requestRepo,
		logger:      logger,
	}
}

// CreatePeriod 受付期間作成
func (u *RequestPeriodUseCase) CreatePeriod(ctx context.Context, input *CreateRequestPeriodInput) (*RequestPeriodOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	orgID, err := sharedDomain.ParseID(input.OrganizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付開始日の形式が不正です")
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付終了日の形式が不正です")
	}

	if endDate.Before(startDate) {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付終了日は開始日より後である必要があります")
	}

	// 既存期間チェック
	existing, err := u.periodRepo.FindByTargetMonth(ctx, orgID, input.TargetYear, input.TargetMonth)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeConflict, "同じ対象月の受付期間が既に存在します")
	}

	maxPerStaff := input.MaxRequestsPerStaff
	if maxPerStaff <= 0 {
		maxPerStaff = 5
	}

	maxPerDay := input.MaxRequestsPerDay
	if maxPerDay <= 0 {
		maxPerDay = 3
	}

	now := time.Now()
	period := &domain.RequestPeriod{
		ID:                  sharedDomain.NewID(),
		OrganizationID:      orgID,
		TargetYear:          input.TargetYear,
		TargetMonth:         input.TargetMonth,
		StartDate:           startDate,
		EndDate:             endDate,
		MaxRequestsPerStaff: maxPerStaff,
		MaxRequestsPerDay:   maxPerDay,
		IsOpen:              false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := u.periodRepo.Save(ctx, period); err != nil {
		u.logger.Error("受付期間作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("受付期間作成完了", "period_id", period.ID)
	return ToRequestPeriodOutput(period), nil
}

// OpenPeriod 受付開始
func (u *RequestPeriodUseCase) OpenPeriod(ctx context.Context, id string) (*RequestPeriodOutput, error) {
	periodID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, sharedDomain.ErrNotFound
	}

	period.IsOpen = true
	period.UpdatedAt = time.Now()

	if err := u.periodRepo.Save(ctx, period); err != nil {
		return nil, err
	}

	u.logger.Info("受付開始", "period_id", period.ID)
	return ToRequestPeriodOutput(period), nil
}

// ClosePeriod 受付終了
func (u *RequestPeriodUseCase) ClosePeriod(ctx context.Context, id string) (*RequestPeriodOutput, error) {
	periodID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, sharedDomain.ErrNotFound
	}

	period.IsOpen = false
	period.UpdatedAt = time.Now()

	if err := u.periodRepo.Save(ctx, period); err != nil {
		return nil, err
	}

	u.logger.Info("受付終了", "period_id", period.ID)
	return ToRequestPeriodOutput(period), nil
}

// GetPeriodByID IDで受付期間取得
func (u *RequestPeriodUseCase) GetPeriodByID(ctx context.Context, id string) (*RequestPeriodOutput, error) {
	periodID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToRequestPeriodOutput(period), nil
}

// GetPeriodByIDWithOrg 組織IDを検証してIDで受付期間取得（マルチテナント対応）
func (u *RequestPeriodUseCase) GetPeriodByIDWithOrg(ctx context.Context, id, orgID string) (*RequestPeriodOutput, error) {
	periodID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, sharedDomain.ErrNotFound
	}

	// 組織チェック（super_adminは全組織アクセス可能）
	if orgID != "" {
		organizationID, parseErr := sharedDomain.ParseID(orgID)
		if parseErr != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
		}
		if period.OrganizationID != organizationID {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeForbidden, "この受付期間へのアクセス権限がありません")
		}
	}

	return ToRequestPeriodOutput(period), nil
}

// ListPeriods 受付期間一覧取得
// Deprecated: ListPeriodsByOrganization を使用してください
func (u *RequestPeriodUseCase) ListPeriods(ctx context.Context) ([]RequestPeriodOutput, error) {
	periods, err := u.periodRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]RequestPeriodOutput, len(periods))
	for i, p := range periods {
		outputs[i] = *ToRequestPeriodOutput(&p)
	}

	return outputs, nil
}

// ListPeriodsByOrganization 組織IDで受付期間一覧取得（マルチテナント対応）
func (u *RequestPeriodUseCase) ListPeriodsByOrganization(ctx context.Context, orgID string) ([]RequestPeriodOutput, error) {
	if orgID == "" {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが必要です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	periods, err := u.periodRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	outputs := make([]RequestPeriodOutput, len(periods))
	for i, p := range periods {
		outputs[i] = *ToRequestPeriodOutput(&p)
	}

	return outputs, nil
}

// Delete 受付期間削除
func (u *RequestPeriodUseCase) Delete(ctx context.Context, id string) error {
	periodID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return err
	}
	if period == nil {
		return sharedDomain.ErrNotFound
	}

	// 受付期間に紐づく勤務希望を削除
	requests, err := u.requestRepo.FindByPeriodID(ctx, periodID)
	if err != nil {
		return err
	}
	for _, r := range requests {
		if err := u.requestRepo.Delete(ctx, r.ID); err != nil {
			u.logger.Error("勤務希望削除失敗", "error", err)
			return err
		}
	}

	if err := u.periodRepo.Delete(ctx, periodID); err != nil {
		u.logger.Error("受付期間削除失敗", "error", err)
		return err
	}

	u.logger.Info("受付期間削除完了", "period_id", periodID)
	return nil
}

// ShiftRequestUseCase 勤務希望ユースケース
type ShiftRequestUseCase struct {
	requestRepo domain.ShiftRequestRepository
	periodRepo  domain.RequestPeriodRepository
	logger      *slog.Logger
}

// NewShiftRequestUseCase 勤務希望ユースケース生成
func NewShiftRequestUseCase(
	requestRepo domain.ShiftRequestRepository,
	periodRepo domain.RequestPeriodRepository,
	logger *slog.Logger,
) *ShiftRequestUseCase {
	return &ShiftRequestUseCase{
		requestRepo: requestRepo,
		periodRepo:  periodRepo,
		logger:      logger,
	}
}

// Create 勤務希望作成
func (u *ShiftRequestUseCase) Create(ctx context.Context, input *CreateShiftRequestInput) (*ShiftRequestOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	periodID, err := sharedDomain.ParseID(input.PeriodID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付期間IDが不正です")
	}

	staffID, err := sharedDomain.ParseID(input.StaffID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "スタッフIDが不正です")
	}

	targetDate, err := time.Parse("2006-01-02", input.TargetDate)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象日の形式が不正です")
	}

	// 受付期間チェック
	period, err := u.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeNotFound, "受付期間が見つかりません")
	}
	if !period.IsActive() {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付期間外です")
	}

	// スタッフあたりの希望数チェック
	staffCount, err := u.requestRepo.CountByPeriodAndStaff(ctx, periodID, staffID)
	if err != nil {
		return nil, err
	}
	if staffCount >= period.MaxRequestsPerStaff {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "希望数の上限に達しています")
	}

	// 日あたりの希望数チェック
	dateCount, err := u.requestRepo.CountByPeriodAndDate(ctx, periodID, targetDate)
	if err != nil {
		return nil, err
	}
	if dateCount >= period.MaxRequestsPerDay {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "この日の希望数上限に達しています")
	}

	var shiftTypeID *sharedDomain.ID
	if input.ShiftTypeID != "" {
		id, err := sharedDomain.ParseID(input.ShiftTypeID)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフト種別IDが不正です")
		}
		shiftTypeID = &id
	}

	requestType := domain.RequestTypePreferred
	if input.RequestType != "" {
		requestType = domain.RequestType(input.RequestType)
	}

	priority := domain.PriorityOptional
	if input.Priority != "" {
		priority = domain.RequestPriority(input.Priority)
	}

	now := time.Now()
	request := &domain.ShiftRequest{
		ID:          sharedDomain.NewID(),
		PeriodID:    periodID,
		StaffID:     staffID,
		TargetDate:  targetDate,
		ShiftTypeID: shiftTypeID,
		RequestType: requestType,
		Priority:    priority,
		Comment:     input.Comment,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.requestRepo.Save(ctx, request); err != nil {
		u.logger.Error("勤務希望作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("勤務希望作成完了", "request_id", request.ID)
	return ToShiftRequestOutput(request), nil
}

// GetByID IDで勤務希望取得
func (u *ShiftRequestUseCase) GetByID(ctx context.Context, id string) (*ShiftRequestOutput, error) {
	requestID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	request, err := u.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToShiftRequestOutput(request), nil
}

// ListByPeriod 受付期間の勤務希望一覧取得
func (u *ShiftRequestUseCase) ListByPeriod(ctx context.Context, periodID string) (*ShiftRequestListOutput, error) {
	pID, err := sharedDomain.ParseID(periodID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "受付期間IDが不正です")
	}

	requests, err := u.requestRepo.FindByPeriodID(ctx, pID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ShiftRequestOutput, len(requests))
	for i, r := range requests {
		outputs[i] = *ToShiftRequestOutput(&r)
	}

	return &ShiftRequestListOutput{
		Requests: outputs,
		Total:    len(outputs),
	}, nil
}

// Delete 勤務希望削除
func (u *ShiftRequestUseCase) Delete(ctx context.Context, id string) error {
	requestID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	request, err := u.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return sharedDomain.ErrNotFound
	}

	if err := u.requestRepo.Delete(ctx, requestID); err != nil {
		u.logger.Error("勤務希望削除失敗", "error", err)
		return err
	}

	u.logger.Info("勤務希望削除完了", "request_id", requestID)
	return nil
}
