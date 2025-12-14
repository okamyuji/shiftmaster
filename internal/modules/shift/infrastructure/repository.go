// Package infrastructure シフトインフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shiftmaster/internal/modules/shift/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ShiftTypeModel シフト種別DBモデル
type ShiftTypeModel struct {
	bun.BaseModel `bun:"table:shift_types"`

	ID              uuid.UUID     `bun:"id,pk,type:uuid"`
	OrganizationID  uuid.UUID     `bun:"organization_id,type:uuid,notnull"`
	ShiftPatternID  uuid.NullUUID `bun:"shift_pattern_id,type:uuid"`
	Name            string        `bun:"name,notnull"`
	Code            string        `bun:"code,notnull"`
	Color           string        `bun:"color,notnull"`
	StartTime       string        `bun:"start_time,type:time,notnull"`
	EndTime         string        `bun:"end_time,type:time,notnull"`
	BreakMinutes    int           `bun:"break_minutes,notnull"`
	HandoverMinutes int           `bun:"handover_minutes,notnull"`
	IsNightShift    bool          `bun:"is_night_shift,notnull"`
	IsHoliday       bool          `bun:"is_holiday,notnull"`
	SortOrder       int           `bun:"sort_order,notnull"`
	CreatedAt       time.Time     `bun:"created_at,notnull"`
	UpdatedAt       time.Time     `bun:"updated_at,notnull"`
}

// ShiftPatternModel シフトパターン（直）DBモデル
type ShiftPatternModel struct {
	bun.BaseModel `bun:"table:shift_patterns"`

	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code,notnull"`
	Description    string    `bun:"description"`
	RotationType   string    `bun:"rotation_type,notnull"`
	Color          string    `bun:"color,notnull"`
	SortOrder      int       `bun:"sort_order,notnull"`
	IsActive       bool      `bun:"is_active,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

// ToDomain ShiftPatternModel DBモデルからドメインエンティティへ変換
func (m *ShiftPatternModel) ToDomain() *domain.ShiftPattern {
	return &domain.ShiftPattern{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
		Description:    m.Description,
		RotationType:   domain.RotationType(m.RotationType),
		Color:          m.Color,
		SortOrder:      m.SortOrder,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromDomain ShiftPatternModel ドメインエンティティからDBモデルへ変換
func (m *ShiftPatternModel) FromDomain(sp *domain.ShiftPattern) {
	m.ID = sp.ID
	m.OrganizationID = sp.OrganizationID
	m.Name = sp.Name
	m.Code = sp.Code
	m.Description = sp.Description
	m.RotationType = sp.RotationType.String()
	m.Color = sp.Color
	m.SortOrder = sp.SortOrder
	m.IsActive = sp.IsActive
	m.CreatedAt = sp.CreatedAt
	m.UpdatedAt = sp.UpdatedAt
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *ShiftTypeModel) ToDomain() *domain.ShiftType {
	// 時刻文字列をtime.Timeに変換
	startTime, _ := time.Parse("15:04:05", m.StartTime)
	endTime, _ := time.Parse("15:04:05", m.EndTime)

	var shiftPatternID *sharedDomain.ID
	if m.ShiftPatternID.Valid {
		id := sharedDomain.ID(m.ShiftPatternID.UUID)
		shiftPatternID = &id
	}

	return &domain.ShiftType{
		ID:              m.ID,
		OrganizationID:  m.OrganizationID,
		ShiftPatternID:  shiftPatternID,
		Name:            m.Name,
		Code:            m.Code,
		Color:           m.Color,
		StartTime:       startTime,
		EndTime:         endTime,
		BreakMinutes:    m.BreakMinutes,
		HandoverMinutes: m.HandoverMinutes,
		IsNightShift:    m.IsNightShift,
		IsHoliday:       m.IsHoliday,
		SortOrder:       m.SortOrder,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// FromDomain ドメインエンティティからDBモデルへ変換
func (m *ShiftTypeModel) FromDomain(st *domain.ShiftType) {
	m.ID = st.ID
	m.OrganizationID = st.OrganizationID
	if st.ShiftPatternID != nil {
		m.ShiftPatternID = uuid.NullUUID{UUID: uuid.UUID(*st.ShiftPatternID), Valid: true}
	}
	m.Name = st.Name
	m.Code = st.Code
	m.Color = st.Color
	// time.Timeを時刻文字列に変換
	m.StartTime = st.StartTime.Format("15:04:05")
	m.EndTime = st.EndTime.Format("15:04:05")
	m.BreakMinutes = st.BreakMinutes
	m.HandoverMinutes = st.HandoverMinutes
	m.IsNightShift = st.IsNightShift
	m.IsHoliday = st.IsHoliday
	m.SortOrder = st.SortOrder
	m.CreatedAt = st.CreatedAt
	m.UpdatedAt = st.UpdatedAt
}

// PostgresShiftTypeRepository PostgreSQLシフト種別リポジトリ
type PostgresShiftTypeRepository struct {
	db *bun.DB
}

// NewPostgresShiftTypeRepository リポジトリ生成
func NewPostgresShiftTypeRepository(db *bun.DB) *PostgresShiftTypeRepository {
	return &PostgresShiftTypeRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresShiftTypeRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.ShiftType, error) {
	model := &ShiftTypeModel{}
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
func (r *PostgresShiftTypeRepository) FindAll(ctx context.Context) ([]domain.ShiftType, error) {
	var models []ShiftTypeModel
	err := r.db.NewSelect().
		Model(&models).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	shiftTypes := make([]domain.ShiftType, len(models))
	for i, m := range models {
		shiftTypes[i] = *m.ToDomain()
	}

	return shiftTypes, nil
}

// FindByOrganizationID 組織IDで検索
func (r *PostgresShiftTypeRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.ShiftType, error) {
	var models []ShiftTypeModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	shiftTypes := make([]domain.ShiftType, len(models))
	for i, m := range models {
		shiftTypes[i] = *m.ToDomain()
	}

	return shiftTypes, nil
}

// FindWorkShifts 勤務シフトのみ取得
func (r *PostgresShiftTypeRepository) FindWorkShifts(ctx context.Context, organizationID sharedDomain.ID) ([]domain.ShiftType, error) {
	var models []ShiftTypeModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Where("is_holiday = ?", false).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	shiftTypes := make([]domain.ShiftType, len(models))
	for i, m := range models {
		shiftTypes[i] = *m.ToDomain()
	}

	return shiftTypes, nil
}

// Save 保存
func (r *PostgresShiftTypeRepository) Save(ctx context.Context, st *domain.ShiftType) error {
	model := &ShiftTypeModel{}
	model.FromDomain(st)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("shift_pattern_id = EXCLUDED.shift_pattern_id").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("color = EXCLUDED.color").
		Set("start_time = EXCLUDED.start_time").
		Set("end_time = EXCLUDED.end_time").
		Set("break_minutes = EXCLUDED.break_minutes").
		Set("handover_minutes = EXCLUDED.handover_minutes").
		Set("is_night_shift = EXCLUDED.is_night_shift").
		Set("is_holiday = EXCLUDED.is_holiday").
		Set("sort_order = EXCLUDED.sort_order").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresShiftTypeRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ShiftTypeModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// PostgresShiftPatternRepository PostgreSQLシフトパターンリポジトリ
type PostgresShiftPatternRepository struct {
	db *bun.DB
}

// NewPostgresShiftPatternRepository リポジトリ生成
func NewPostgresShiftPatternRepository(db *bun.DB) *PostgresShiftPatternRepository {
	return &PostgresShiftPatternRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresShiftPatternRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.ShiftPattern, error) {
	model := &ShiftPatternModel{}
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
func (r *PostgresShiftPatternRepository) FindAll(ctx context.Context) ([]domain.ShiftPattern, error) {
	var models []ShiftPatternModel
	err := r.db.NewSelect().
		Model(&models).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	patterns := make([]domain.ShiftPattern, len(models))
	for i, m := range models {
		patterns[i] = *m.ToDomain()
	}

	return patterns, nil
}

// FindByOrganizationID 組織IDで検索
func (r *PostgresShiftPatternRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.ShiftPattern, error) {
	var models []ShiftPatternModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	patterns := make([]domain.ShiftPattern, len(models))
	for i, m := range models {
		patterns[i] = *m.ToDomain()
	}

	return patterns, nil
}

// FindActiveByOrganizationID 有効パターンのみ取得
func (r *PostgresShiftPatternRepository) FindActiveByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.ShiftPattern, error) {
	var models []ShiftPatternModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Where("is_active = ?", true).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	patterns := make([]domain.ShiftPattern, len(models))
	for i, m := range models {
		patterns[i] = *m.ToDomain()
	}

	return patterns, nil
}

// Save 保存
func (r *PostgresShiftPatternRepository) Save(ctx context.Context, sp *domain.ShiftPattern) error {
	model := &ShiftPatternModel{}
	model.FromDomain(sp)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("description = EXCLUDED.description").
		Set("rotation_type = EXCLUDED.rotation_type").
		Set("color = EXCLUDED.color").
		Set("sort_order = EXCLUDED.sort_order").
		Set("is_active = EXCLUDED.is_active").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresShiftPatternRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*ShiftPatternModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
