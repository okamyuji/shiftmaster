// Package infrastructure 勤務希望インフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shiftmaster/internal/modules/request/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// RequestPeriodModel 受付期間DBモデル
type RequestPeriodModel struct {
	bun.BaseModel `bun:"table:request_periods"`

	ID                  uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID      uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	TargetYear          int       `bun:"target_year,notnull"`
	TargetMonth         int       `bun:"target_month,notnull"`
	StartDate           time.Time `bun:"start_date,type:date,notnull"`
	EndDate             time.Time `bun:"end_date,type:date,notnull"`
	MaxRequestsPerStaff int       `bun:"max_requests_per_staff,notnull"`
	MaxRequestsPerDay   int       `bun:"max_requests_per_day,notnull"`
	IsOpen              bool      `bun:"is_open,notnull"`
	CreatedAt           time.Time `bun:"created_at,notnull"`
	UpdatedAt           time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *RequestPeriodModel) ToDomain() *domain.RequestPeriod {
	return &domain.RequestPeriod{
		ID:                  m.ID,
		OrganizationID:      m.OrganizationID,
		TargetYear:          m.TargetYear,
		TargetMonth:         m.TargetMonth,
		StartDate:           m.StartDate,
		EndDate:             m.EndDate,
		MaxRequestsPerStaff: m.MaxRequestsPerStaff,
		MaxRequestsPerDay:   m.MaxRequestsPerDay,
		IsOpen:              m.IsOpen,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

// PostgresRequestPeriodRepository PostgreSQL受付期間リポジトリ
type PostgresRequestPeriodRepository struct {
	db *bun.DB
}

// NewPostgresRequestPeriodRepository リポジトリ生成
func NewPostgresRequestPeriodRepository(db *bun.DB) *PostgresRequestPeriodRepository {
	return &PostgresRequestPeriodRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresRequestPeriodRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.RequestPeriod, error) {
	model := &RequestPeriodModel{}
	err := r.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindAll 全件取得
func (r *PostgresRequestPeriodRepository) FindAll(ctx context.Context) ([]domain.RequestPeriod, error) {
	var models []RequestPeriodModel
	err := r.db.NewSelect().
		Model(&models).
		Order("target_year DESC", "target_month DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	periods := make([]domain.RequestPeriod, len(models))
	for i, m := range models {
		periods[i] = *m.ToDomain()
	}

	return periods, nil
}

// FindByOrganizationID 組織IDで検索
func (r *PostgresRequestPeriodRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.RequestPeriod, error) {
	var models []RequestPeriodModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Order("target_year DESC", "target_month DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	periods := make([]domain.RequestPeriod, len(models))
	for i, m := range models {
		periods[i] = *m.ToDomain()
	}

	return periods, nil
}

// FindByTargetMonth 対象年月で検索
func (r *PostgresRequestPeriodRepository) FindByTargetMonth(ctx context.Context, organizationID sharedDomain.ID, year, month int) (*domain.RequestPeriod, error) {
	model := &RequestPeriodModel{}
	err := r.db.NewSelect().
		Model(model).
		Where("organization_id = ? AND target_year = ? AND target_month = ?", organizationID, year, month).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindActive 受付中の期間を取得
func (r *PostgresRequestPeriodRepository) FindActive(ctx context.Context, organizationID sharedDomain.ID) ([]domain.RequestPeriod, error) {
	now := time.Now()
	var models []RequestPeriodModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ? AND is_open = ? AND start_date <= ? AND end_date >= ?", organizationID, true, now, now).
		Order("target_year DESC", "target_month DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	periods := make([]domain.RequestPeriod, len(models))
	for i, m := range models {
		periods[i] = *m.ToDomain()
	}

	return periods, nil
}

// Save 保存
func (r *PostgresRequestPeriodRepository) Save(ctx context.Context, period *domain.RequestPeriod) error {
	model := &RequestPeriodModel{
		ID:                  period.ID,
		OrganizationID:      period.OrganizationID,
		TargetYear:          period.TargetYear,
		TargetMonth:         period.TargetMonth,
		StartDate:           period.StartDate,
		EndDate:             period.EndDate,
		MaxRequestsPerStaff: period.MaxRequestsPerStaff,
		MaxRequestsPerDay:   period.MaxRequestsPerDay,
		IsOpen:              period.IsOpen,
		CreatedAt:           period.CreatedAt,
		UpdatedAt:           period.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("start_date = EXCLUDED.start_date").
		Set("end_date = EXCLUDED.end_date").
		Set("max_requests_per_staff = EXCLUDED.max_requests_per_staff").
		Set("max_requests_per_day = EXCLUDED.max_requests_per_day").
		Set("is_open = EXCLUDED.is_open").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresRequestPeriodRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*RequestPeriodModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// ShiftRequestModel 勤務希望DBモデル
type ShiftRequestModel struct {
	bun.BaseModel `bun:"table:shift_requests"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid"`
	PeriodID    uuid.UUID  `bun:"period_id,type:uuid,notnull"`
	StaffID     uuid.UUID  `bun:"staff_id,type:uuid,notnull"`
	TargetDate  time.Time  `bun:"target_date,type:date,notnull"`
	ShiftTypeID *uuid.UUID `bun:"shift_type_id,type:uuid"`
	RequestType string     `bun:"request_type,notnull"`
	Priority    string     `bun:"priority,notnull"`
	Comment     string     `bun:"comment"`
	CreatedAt   time.Time  `bun:"created_at,notnull"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *ShiftRequestModel) ToDomain() *domain.ShiftRequest {
	var shiftTypeID *sharedDomain.ID
	if m.ShiftTypeID != nil {
		id := *m.ShiftTypeID
		shiftTypeID = &id
	}

	return &domain.ShiftRequest{
		ID:          m.ID,
		PeriodID:    m.PeriodID,
		StaffID:     m.StaffID,
		TargetDate:  m.TargetDate,
		ShiftTypeID: shiftTypeID,
		RequestType: domain.RequestType(m.RequestType),
		Priority:    domain.RequestPriority(m.Priority),
		Comment:     m.Comment,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// PostgresShiftRequestRepository PostgreSQL勤務希望リポジトリ
type PostgresShiftRequestRepository struct {
	db *bun.DB
}

// NewPostgresShiftRequestRepository リポジトリ生成
func NewPostgresShiftRequestRepository(db *bun.DB) *PostgresShiftRequestRepository {
	return &PostgresShiftRequestRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresShiftRequestRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.ShiftRequest, error) {
	model := &ShiftRequestModel{}
	err := r.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByPeriodID 期間IDで検索
func (r *PostgresShiftRequestRepository) FindByPeriodID(ctx context.Context, periodID sharedDomain.ID) ([]domain.ShiftRequest, error) {
	var models []ShiftRequestModel
	err := r.db.NewSelect().
		Model(&models).
		Where("period_id = ?", periodID).
		Order("staff_id ASC", "target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]domain.ShiftRequest, len(models))
	for i, m := range models {
		requests[i] = *m.ToDomain()
	}

	return requests, nil
}

// FindByStaffID スタッフIDで検索
func (r *PostgresShiftRequestRepository) FindByStaffID(ctx context.Context, staffID sharedDomain.ID) ([]domain.ShiftRequest, error) {
	var models []ShiftRequestModel
	err := r.db.NewSelect().
		Model(&models).
		Where("staff_id = ?", staffID).
		Order("target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]domain.ShiftRequest, len(models))
	for i, m := range models {
		requests[i] = *m.ToDomain()
	}

	return requests, nil
}

// FindByPeriodAndStaff 期間とスタッフで検索
func (r *PostgresShiftRequestRepository) FindByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) ([]domain.ShiftRequest, error) {
	var models []ShiftRequestModel
	err := r.db.NewSelect().
		Model(&models).
		Where("period_id = ? AND staff_id = ?", periodID, staffID).
		Order("target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]domain.ShiftRequest, len(models))
	for i, m := range models {
		requests[i] = *m.ToDomain()
	}

	return requests, nil
}

// FindByPeriodAndDate 期間と日付で検索
func (r *PostgresShiftRequestRepository) FindByPeriodAndDate(ctx context.Context, periodID sharedDomain.ID, date time.Time) ([]domain.ShiftRequest, error) {
	var models []ShiftRequestModel
	err := r.db.NewSelect().
		Model(&models).
		Where("period_id = ? AND target_date = ?", periodID, date).
		Order("staff_id ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	requests := make([]domain.ShiftRequest, len(models))
	for i, m := range models {
		requests[i] = *m.ToDomain()
	}

	return requests, nil
}

// CountByPeriodAndStaff スタッフの希望数カウント
func (r *PostgresShiftRequestRepository) CountByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) (int, error) {
	count, err := r.db.NewSelect().
		Model((*ShiftRequestModel)(nil)).
		Where("period_id = ? AND staff_id = ?", periodID, staffID).
		Count(ctx)
	return count, err
}

// CountByPeriodAndDate 日付の希望数カウント
func (r *PostgresShiftRequestRepository) CountByPeriodAndDate(ctx context.Context, periodID sharedDomain.ID, date time.Time) (int, error) {
	count, err := r.db.NewSelect().
		Model((*ShiftRequestModel)(nil)).
		Where("period_id = ? AND target_date = ?", periodID, date).
		Count(ctx)
	return count, err
}

// Save 保存
func (r *PostgresShiftRequestRepository) Save(ctx context.Context, request *domain.ShiftRequest) error {
	var shiftTypeID *uuid.UUID
	if request.ShiftTypeID != nil {
		id := *request.ShiftTypeID
		shiftTypeID = &id
	}

	model := &ShiftRequestModel{
		ID:          request.ID,
		PeriodID:    request.PeriodID,
		StaffID:     request.StaffID,
		TargetDate:  request.TargetDate,
		ShiftTypeID: shiftTypeID,
		RequestType: request.RequestType.String(),
		Priority:    request.Priority.String(),
		Comment:     request.Comment,
		CreatedAt:   request.CreatedAt,
		UpdatedAt:   request.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("shift_type_id = EXCLUDED.shift_type_id").
		Set("request_type = EXCLUDED.request_type").
		Set("priority = EXCLUDED.priority").
		Set("comment = EXCLUDED.comment").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresShiftRequestRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ShiftRequestModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DeleteByPeriodAndStaff 期間とスタッフで削除
func (r *PostgresShiftRequestRepository) DeleteByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) error {
	_, err := r.db.NewDelete().
		Model((*ShiftRequestModel)(nil)).
		Where("period_id = ? AND staff_id = ?", periodID, staffID).
		Exec(ctx)
	return err
}
