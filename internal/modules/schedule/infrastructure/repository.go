// Package infrastructure 勤務表インフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shiftmaster/internal/modules/schedule/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ScheduleModel 勤務表DBモデル
type ScheduleModel struct {
	bun.BaseModel `bun:"table:schedules"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID  `bun:"organization_id,type:uuid,notnull"`
	TargetYear     int        `bun:"target_year,notnull"`
	TargetMonth    int        `bun:"target_month,notnull"`
	Status         string     `bun:"status,notnull"`
	PublishedAt    *time.Time `bun:"published_at"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *ScheduleModel) ToDomain() *domain.Schedule {
	return &domain.Schedule{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		TargetYear:     m.TargetYear,
		TargetMonth:    m.TargetMonth,
		Status:         domain.ScheduleStatus(m.Status),
		PublishedAt:    m.PublishedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// PostgresScheduleRepository PostgreSQL勤務表リポジトリ
type PostgresScheduleRepository struct {
	db *bun.DB
}

// NewPostgresScheduleRepository リポジトリ生成
func NewPostgresScheduleRepository(db *bun.DB) *PostgresScheduleRepository {
	return &PostgresScheduleRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresScheduleRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.Schedule, error) {
	model := &ScheduleModel{}
	err := r.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByIDWithEntries エントリ付きで検索
func (r *PostgresScheduleRepository) FindByIDWithEntries(ctx context.Context, id sharedDomain.ID) (*domain.Schedule, error) {
	schedule, err := r.FindByID(ctx, id)
	if err != nil || schedule == nil {
		return schedule, err
	}

	// エントリを別途取得
	var entryModels []ScheduleEntryModel
	err = r.db.NewSelect().
		Model(&entryModels).
		Where("schedule_id = ?", id).
		Order("staff_id ASC", "target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.ScheduleEntry, len(entryModels))
	for i, m := range entryModels {
		entries[i] = *m.ToDomain()
	}
	schedule.Entries = entries

	return schedule, nil
}

// FindByOrganizationID 組織IDで検索
func (r *PostgresScheduleRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.Schedule, error) {
	var models []ScheduleModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Order("target_year DESC", "target_month DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	schedules := make([]domain.Schedule, len(models))
	for i, m := range models {
		schedules[i] = *m.ToDomain()
	}

	return schedules, nil
}

// FindByTargetMonth 対象年月で検索
func (r *PostgresScheduleRepository) FindByTargetMonth(ctx context.Context, organizationID sharedDomain.ID, year, month int) (*domain.Schedule, error) {
	model := &ScheduleModel{}
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

// Save 保存
func (r *PostgresScheduleRepository) Save(ctx context.Context, schedule *domain.Schedule) error {
	model := &ScheduleModel{
		ID:             schedule.ID,
		OrganizationID: schedule.OrganizationID,
		TargetYear:     schedule.TargetYear,
		TargetMonth:    schedule.TargetMonth,
		Status:         schedule.Status.String(),
		PublishedAt:    schedule.PublishedAt,
		CreatedAt:      schedule.CreatedAt,
		UpdatedAt:      schedule.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("status = EXCLUDED.status").
		Set("published_at = EXCLUDED.published_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresScheduleRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ScheduleModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// ScheduleEntryModel 勤務表エントリDBモデル
type ScheduleEntryModel struct {
	bun.BaseModel `bun:"table:schedule_entries"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid"`
	ScheduleID  uuid.UUID  `bun:"schedule_id,type:uuid,notnull"`
	StaffID     uuid.UUID  `bun:"staff_id,type:uuid,notnull"`
	TargetDate  time.Time  `bun:"target_date,type:date,notnull"`
	ShiftTypeID *uuid.UUID `bun:"shift_type_id,type:uuid"`
	IsConfirmed bool       `bun:"is_confirmed,notnull"`
	Note        string     `bun:"note"`
	CreatedAt   time.Time  `bun:"created_at,notnull"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *ScheduleEntryModel) ToDomain() *domain.ScheduleEntry {
	var shiftTypeID *sharedDomain.ID
	if m.ShiftTypeID != nil {
		id := *m.ShiftTypeID
		shiftTypeID = &id
	}

	return &domain.ScheduleEntry{
		ID:          m.ID,
		ScheduleID:  m.ScheduleID,
		StaffID:     m.StaffID,
		TargetDate:  m.TargetDate,
		ShiftTypeID: shiftTypeID,
		IsConfirmed: m.IsConfirmed,
		Note:        m.Note,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// PostgresScheduleEntryRepository PostgreSQL勤務表エントリリポジトリ
type PostgresScheduleEntryRepository struct {
	db *bun.DB
}

// NewPostgresScheduleEntryRepository リポジトリ生成
func NewPostgresScheduleEntryRepository(db *bun.DB) *PostgresScheduleEntryRepository {
	return &PostgresScheduleEntryRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresScheduleEntryRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.ScheduleEntry, error) {
	model := &ScheduleEntryModel{}
	err := r.db.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByScheduleID 勤務表IDで検索
func (r *PostgresScheduleEntryRepository) FindByScheduleID(ctx context.Context, scheduleID sharedDomain.ID) ([]domain.ScheduleEntry, error) {
	var models []ScheduleEntryModel
	err := r.db.NewSelect().
		Model(&models).
		Where("schedule_id = ?", scheduleID).
		Order("staff_id ASC", "target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.ScheduleEntry, len(models))
	for i, m := range models {
		entries[i] = *m.ToDomain()
	}

	return entries, nil
}

// FindByScheduleAndStaff 勤務表とスタッフで検索
func (r *PostgresScheduleEntryRepository) FindByScheduleAndStaff(ctx context.Context, scheduleID, staffID sharedDomain.ID) ([]domain.ScheduleEntry, error) {
	var models []ScheduleEntryModel
	err := r.db.NewSelect().
		Model(&models).
		Where("schedule_id = ? AND staff_id = ?", scheduleID, staffID).
		Order("target_date ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.ScheduleEntry, len(models))
	for i, m := range models {
		entries[i] = *m.ToDomain()
	}

	return entries, nil
}

// FindByScheduleAndDate 勤務表と日付で検索
func (r *PostgresScheduleEntryRepository) FindByScheduleAndDate(ctx context.Context, scheduleID sharedDomain.ID, date time.Time) ([]domain.ScheduleEntry, error) {
	var models []ScheduleEntryModel
	err := r.db.NewSelect().
		Model(&models).
		Where("schedule_id = ? AND target_date = ?", scheduleID, date).
		Order("staff_id ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.ScheduleEntry, len(models))
	for i, m := range models {
		entries[i] = *m.ToDomain()
	}

	return entries, nil
}

// Save 保存
func (r *PostgresScheduleEntryRepository) Save(ctx context.Context, entry *domain.ScheduleEntry) error {
	var shiftTypeID *uuid.UUID
	if entry.ShiftTypeID != nil {
		id := *entry.ShiftTypeID
		shiftTypeID = &id
	}

	model := &ScheduleEntryModel{
		ID:          entry.ID,
		ScheduleID:  entry.ScheduleID,
		StaffID:     entry.StaffID,
		TargetDate:  entry.TargetDate,
		ShiftTypeID: shiftTypeID,
		IsConfirmed: entry.IsConfirmed,
		Note:        entry.Note,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("shift_type_id = EXCLUDED.shift_type_id").
		Set("is_confirmed = EXCLUDED.is_confirmed").
		Set("note = EXCLUDED.note").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// SaveBatch 一括保存
func (r *PostgresScheduleEntryRepository) SaveBatch(ctx context.Context, entries []domain.ScheduleEntry) error {
	if len(entries) == 0 {
		return nil
	}

	models := make([]ScheduleEntryModel, len(entries))
	for i, entry := range entries {
		var shiftTypeID *uuid.UUID
		if entry.ShiftTypeID != nil {
			id := *entry.ShiftTypeID
			shiftTypeID = &id
		}

		models[i] = ScheduleEntryModel{
			ID:          entry.ID,
			ScheduleID:  entry.ScheduleID,
			StaffID:     entry.StaffID,
			TargetDate:  entry.TargetDate,
			ShiftTypeID: shiftTypeID,
			IsConfirmed: entry.IsConfirmed,
			Note:        entry.Note,
			CreatedAt:   entry.CreatedAt,
			UpdatedAt:   entry.UpdatedAt,
		}
	}

	_, err := r.db.NewInsert().
		Model(&models).
		On("CONFLICT (id) DO UPDATE").
		Set("shift_type_id = EXCLUDED.shift_type_id").
		Set("is_confirmed = EXCLUDED.is_confirmed").
		Set("note = EXCLUDED.note").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresScheduleEntryRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ScheduleEntryModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DeleteBySchedule 勤務表のエントリ全削除
func (r *PostgresScheduleEntryRepository) DeleteBySchedule(ctx context.Context, scheduleID sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ScheduleEntryModel)(nil)).Where("schedule_id = ?", scheduleID).Exec(ctx)
	return err
}
