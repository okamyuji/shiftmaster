// Package domain スタッフドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// JobType 職種エンティティ
// 看護師、介護士、オペレーターなど
type JobType struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name 職種名
	Name string
	// Code 職種コード
	Code string
	// Description 説明
	Description string
	// Color 表示色
	Color string
	// SortOrder 表示順
	SortOrder int
	// IsActive 有効フラグ
	IsActive bool
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// NewJobType 職種生成
func NewJobType(orgID domain.ID, name, code string) (*JobType, error) {
	if name == "" {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "職種名は必須です")
	}
	if code == "" {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "職種コードは必須です")
	}

	now := time.Now()
	return &JobType{
		ID:             domain.NewID(),
		OrganizationID: orgID,
		Name:           name,
		Code:           code,
		Color:          "#6B7280",
		SortOrder:      0,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Update 職種更新
func (j *JobType) Update(name, code, description, color string, sortOrder int) error {
	if name == "" {
		return domain.NewDomainError(domain.ErrCodeValidation, "職種名は必須です")
	}
	if code == "" {
		return domain.NewDomainError(domain.ErrCodeValidation, "職種コードは必須です")
	}

	j.Name = name
	j.Code = code
	j.Description = description
	if color != "" {
		j.Color = color
	}
	j.SortOrder = sortOrder
	j.UpdatedAt = time.Now()
	return nil
}

// Deactivate 無効化
func (j *JobType) Deactivate() {
	j.IsActive = false
	j.UpdatedAt = time.Now()
}

// Activate 有効化
func (j *JobType) Activate() {
	j.IsActive = true
	j.UpdatedAt = time.Now()
}

// JobTypeRepository 職種リポジトリインターフェース
type JobTypeRepository interface {
	// FindByID ID検索
	FindByID(ctx interface{}, id domain.ID) (*JobType, error)
	// FindByOrganizationID 組織ID検索
	FindByOrganizationID(ctx interface{}, orgID domain.ID) ([]JobType, error)
	// FindByCode コード検索
	FindByCode(ctx interface{}, orgID domain.ID, code string) (*JobType, error)
	// Save 保存
	Save(ctx interface{}, jobType *JobType) error
	// Delete 削除
	Delete(ctx interface{}, id domain.ID) error
}
