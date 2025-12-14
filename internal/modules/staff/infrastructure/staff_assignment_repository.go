// Package infrastructure スタッフインフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// StaffAssignmentModel スタッフ所属DBモデル
type StaffAssignmentModel struct {
	bun.BaseModel `bun:"table:staff_assignments"`

	ID         uuid.UUID  `bun:"id,pk,type:uuid"`
	StaffID    uuid.UUID  `bun:"staff_id,type:uuid,notnull"`
	TeamID     *uuid.UUID `bun:"team_id,type:uuid"`
	JobTypeID  *uuid.UUID `bun:"job_type_id,type:uuid"`
	PositionID *uuid.UUID `bun:"position_id,type:uuid"`
	IsPrimary  bool       `bun:"is_primary,notnull"`
	StartDate  *time.Time `bun:"start_date,type:date"`
	EndDate    *time.Time `bun:"end_date,type:date"`
	CreatedAt  time.Time  `bun:"created_at,notnull"`
	UpdatedAt  time.Time  `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *StaffAssignmentModel) ToDomain() *domain.StaffAssignment {
	var teamID, jobTypeID, positionID *sharedDomain.ID
	if m.TeamID != nil {
		t := sharedDomain.ID(*m.TeamID)
		teamID = &t
	}
	if m.JobTypeID != nil {
		j := sharedDomain.ID(*m.JobTypeID)
		jobTypeID = &j
	}
	if m.PositionID != nil {
		p := sharedDomain.ID(*m.PositionID)
		positionID = &p
	}

	return &domain.StaffAssignment{
		ID:         m.ID,
		StaffID:    m.StaffID,
		TeamID:     teamID,
		JobTypeID:  jobTypeID,
		PositionID: positionID,
		IsPrimary:  m.IsPrimary,
		StartDate:  m.StartDate,
		EndDate:    m.EndDate,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FromDomain ドメインエンティティからDBモデルへ変換
func (m *StaffAssignmentModel) FromDomain(a *domain.StaffAssignment) {
	m.ID = a.ID
	m.StaffID = a.StaffID
	if a.TeamID != nil {
		t := uuid.UUID(*a.TeamID)
		m.TeamID = &t
	}
	if a.JobTypeID != nil {
		j := uuid.UUID(*a.JobTypeID)
		m.JobTypeID = &j
	}
	if a.PositionID != nil {
		p := uuid.UUID(*a.PositionID)
		m.PositionID = &p
	}
	m.IsPrimary = a.IsPrimary
	m.StartDate = a.StartDate
	m.EndDate = a.EndDate
	m.CreatedAt = a.CreatedAt
	m.UpdatedAt = a.UpdatedAt
}

// PostgresStaffAssignmentRepository PostgreSQLスタッフ所属リポジトリ
type PostgresStaffAssignmentRepository struct {
	db *bun.DB
}

// NewPostgresStaffAssignmentRepository リポジトリ生成
func NewPostgresStaffAssignmentRepository(db *bun.DB) *PostgresStaffAssignmentRepository {
	return &PostgresStaffAssignmentRepository{db: db}
}

// FindByID ID検索
func (r *PostgresStaffAssignmentRepository) FindByID(ctx interface{}, id sharedDomain.ID) (*domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var model StaffAssignmentModel
	err := r.db.NewSelect().Model(&model).Where("id = ?", id).Scan(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByStaffID スタッフID検索
func (r *PostgresStaffAssignmentRepository) FindByStaffID(ctx interface{}, staffID sharedDomain.ID) ([]domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var models []StaffAssignmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("staff_id = ?", staffID).
		Order("is_primary DESC", "created_at ASC").
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StaffAssignment, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindActiveByStaffID スタッフIDで有効な所属を検索
func (r *PostgresStaffAssignmentRepository) FindActiveByStaffID(ctx interface{}, staffID sharedDomain.ID, date time.Time) ([]domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var models []StaffAssignmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("staff_id = ?", staffID).
		Where("(start_date IS NULL OR start_date <= ?)", date).
		Where("(end_date IS NULL OR end_date >= ?)", date).
		Order("is_primary DESC", "created_at ASC").
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StaffAssignment, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindByTeamID チームID検索
func (r *PostgresStaffAssignmentRepository) FindByTeamID(ctx interface{}, teamID sharedDomain.ID) ([]domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var models []StaffAssignmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("team_id = ?", teamID).
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StaffAssignment, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindByJobTypeID 職種ID検索
func (r *PostgresStaffAssignmentRepository) FindByJobTypeID(ctx interface{}, jobTypeID sharedDomain.ID) ([]domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var models []StaffAssignmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("job_type_id = ?", jobTypeID).
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StaffAssignment, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindByPositionID 職位ID検索
func (r *PostgresStaffAssignmentRepository) FindByPositionID(ctx interface{}, positionID sharedDomain.ID) ([]domain.StaffAssignment, error) {
	c := ctx.(context.Context)
	var models []StaffAssignmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("position_id = ?", positionID).
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StaffAssignment, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// Save 保存
func (r *PostgresStaffAssignmentRepository) Save(ctx interface{}, assignment *domain.StaffAssignment) error {
	c := ctx.(context.Context)
	model := &StaffAssignmentModel{}
	model.FromDomain(assignment)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("team_id = EXCLUDED.team_id").
		Set("job_type_id = EXCLUDED.job_type_id").
		Set("position_id = EXCLUDED.position_id").
		Set("is_primary = EXCLUDED.is_primary").
		Set("start_date = EXCLUDED.start_date").
		Set("end_date = EXCLUDED.end_date").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(c)
	return err
}

// Delete 削除
func (r *PostgresStaffAssignmentRepository) Delete(ctx interface{}, id sharedDomain.ID) error {
	c := ctx.(context.Context)
	_, err := r.db.NewDelete().Model((*StaffAssignmentModel)(nil)).Where("id = ?", id).Exec(c)
	return err
}
