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

// PositionModel 職位DBモデル
type PositionModel struct {
	bun.BaseModel `bun:"table:positions"`

	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code,notnull"`
	Description    string    `bun:"description"`
	Level          int       `bun:"level,notnull"`
	SortOrder      int       `bun:"sort_order,notnull"`
	IsActive       bool      `bun:"is_active,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *PositionModel) ToDomain() *domain.Position {
	return &domain.Position{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
		Description:    m.Description,
		Level:          m.Level,
		SortOrder:      m.SortOrder,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromDomain ドメインエンティティからDBモデルへ変換
func (m *PositionModel) FromDomain(p *domain.Position) {
	m.ID = p.ID
	m.OrganizationID = p.OrganizationID
	m.Name = p.Name
	m.Code = p.Code
	m.Description = p.Description
	m.Level = p.Level
	m.SortOrder = p.SortOrder
	m.IsActive = p.IsActive
	m.CreatedAt = p.CreatedAt
	m.UpdatedAt = p.UpdatedAt
}

// PostgresPositionRepository PostgreSQL職位リポジトリ
type PostgresPositionRepository struct {
	db *bun.DB
}

// NewPostgresPositionRepository リポジトリ生成
func NewPostgresPositionRepository(db *bun.DB) *PostgresPositionRepository {
	return &PostgresPositionRepository{db: db}
}

// FindByID ID検索
func (r *PostgresPositionRepository) FindByID(ctx interface{}, id sharedDomain.ID) (*domain.Position, error) {
	c := ctx.(context.Context)
	var model PositionModel
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
func (r *PostgresPositionRepository) FindByOrganizationID(ctx interface{}, orgID sharedDomain.ID) ([]domain.Position, error) {
	c := ctx.(context.Context)
	var models []PositionModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", orgID).
		Where("is_active = ?", true).
		Order("level ASC", "sort_order ASC", "name ASC").
		Scan(c)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Position, len(models))
	for i, m := range models {
		result[i] = *m.ToDomain()
	}
	return result, nil
}

// FindByCode コード検索
func (r *PostgresPositionRepository) FindByCode(ctx interface{}, orgID sharedDomain.ID, code string) (*domain.Position, error) {
	c := ctx.(context.Context)
	var model PositionModel
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
func (r *PostgresPositionRepository) Save(ctx interface{}, position *domain.Position) error {
	c := ctx.(context.Context)
	model := &PositionModel{}
	model.FromDomain(position)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("description = EXCLUDED.description").
		Set("level = EXCLUDED.level").
		Set("sort_order = EXCLUDED.sort_order").
		Set("is_active = EXCLUDED.is_active").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(c)
	return err
}

// Delete 削除
func (r *PostgresPositionRepository) Delete(ctx interface{}, id sharedDomain.ID) error {
	c := ctx.(context.Context)
	_, err := r.db.NewDelete().Model((*PositionModel)(nil)).Where("id = ?", id).Exec(c)
	return err
}
