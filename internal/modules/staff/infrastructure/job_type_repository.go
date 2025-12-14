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

// JobTypeModel 職種DBモデル
type JobTypeModel struct {
	bun.BaseModel `bun:"table:job_types"`

	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code,notnull"`
	Description    string    `bun:"description"`
	Color          string    `bun:"color"`
	SortOrder      int       `bun:"sort_order,notnull"`
	IsActive       bool      `bun:"is_active,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *JobTypeModel) ToDomain() *domain.JobType {
	return &domain.JobType{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
		Description:    m.Description,
		Color:          m.Color,
		SortOrder:      m.SortOrder,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromDomain ドメインエンティティからDBモデルへ変換
func (m *JobTypeModel) FromDomain(j *domain.JobType) {
	m.ID = j.ID
	m.OrganizationID = j.OrganizationID
	m.Name = j.Name
	m.Code = j.Code
	m.Description = j.Description
	m.Color = j.Color
	m.SortOrder = j.SortOrder
	m.IsActive = j.IsActive
	m.CreatedAt = j.CreatedAt
	m.UpdatedAt = j.UpdatedAt
}

// PostgresJobTypeRepository PostgreSQL職種リポジトリ
type PostgresJobTypeRepository struct {
	db *bun.DB
}

// NewPostgresJobTypeRepository リポジトリ生成
func NewPostgresJobTypeRepository(db *bun.DB) *PostgresJobTypeRepository {
	return &PostgresJobTypeRepository{db: db}
}

// FindByID ID検索
func (r *PostgresJobTypeRepository) FindByID(ctx interface{}, id sharedDomain.ID) (*domain.JobType, error) {
	c := ctx.(context.Context)
	var model JobTypeModel
	err := r.db.NewSelect().Model(&model).Where("id = ?", id).Scan(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByOrganizationID 組織ID検索
func (r *PostgresJobTypeRepository) FindByOrganizationID(ctx interface{}, orgID sharedDomain.ID) ([]domain.JobType, error) {
	c := ctx.(context.Context)
	var models []JobTypeModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", orgID).
		Where("is_active = ?", true).
		Order("sort_order ASC", "name ASC").
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.JobType, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindByCode コード検索
func (r *PostgresJobTypeRepository) FindByCode(ctx interface{}, orgID sharedDomain.ID, code string) (*domain.JobType, error) {
	c := ctx.(context.Context)
	var model JobTypeModel
	err := r.db.NewSelect().
		Model(&model).
		Where("organization_id = ?", orgID).
		Where("code = ?", code).
		Scan(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// Save 保存
func (r *PostgresJobTypeRepository) Save(ctx interface{}, jobType *domain.JobType) error {
	c := ctx.(context.Context)
	model := &JobTypeModel{}
	model.FromDomain(jobType)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("description = EXCLUDED.description").
		Set("color = EXCLUDED.color").
		Set("sort_order = EXCLUDED.sort_order").
		Set("is_active = EXCLUDED.is_active").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(c)
	return err
}

// Delete 削除
func (r *PostgresJobTypeRepository) Delete(ctx interface{}, id sharedDomain.ID) error {
	c := ctx.(context.Context)
	_, err := r.db.NewDelete().Model((*JobTypeModel)(nil)).Where("id = ?", id).Exec(c)
	return err
}
