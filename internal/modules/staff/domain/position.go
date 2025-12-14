// Package domain スタッフドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// Position 職位エンティティ
// 師長、主任、リーダー、一般など
type Position struct {
	// ID 一意識別子
	ID domain.ID
	// OrganizationID 組織ID
	OrganizationID domain.ID
	// Name 職位名
	Name string
	// Code 職位コード
	Code string
	// Description 説明
	Description string
	// Level 階層レベル 数値が小さいほど上位
	Level int
	// SortOrder 表示順
	SortOrder int
	// IsActive 有効フラグ
	IsActive bool
	// CreatedAt 作成日時
	CreatedAt time.Time
	// UpdatedAt 更新日時
	UpdatedAt time.Time
}

// NewPosition 職位生成
func NewPosition(orgID domain.ID, name, code string, level int) (*Position, error) {
	if name == "" {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "職位名は必須です")
	}
	if code == "" {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "職位コードは必須です")
	}
	if level < 0 {
		return nil, domain.NewDomainError(domain.ErrCodeValidation, "レベルは0以上で指定してください")
	}

	now := time.Now()
	return &Position{
		ID:             domain.NewID(),
		OrganizationID: orgID,
		Name:           name,
		Code:           code,
		Level:          level,
		SortOrder:      0,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Update 職位更新
func (p *Position) Update(name, code, description string, level, sortOrder int) error {
	if name == "" {
		return domain.NewDomainError(domain.ErrCodeValidation, "職位名は必須です")
	}
	if code == "" {
		return domain.NewDomainError(domain.ErrCodeValidation, "職位コードは必須です")
	}
	if level < 0 {
		return domain.NewDomainError(domain.ErrCodeValidation, "レベルは0以上で指定してください")
	}

	p.Name = name
	p.Code = code
	p.Description = description
	p.Level = level
	p.SortOrder = sortOrder
	p.UpdatedAt = time.Now()
	return nil
}

// IsHigherThan 指定職位より上位か判定
func (p *Position) IsHigherThan(other *Position) bool {
	if other == nil {
		return true
	}
	return p.Level < other.Level
}

// IsEqualOrHigherThan 指定職位と同等以上か判定
func (p *Position) IsEqualOrHigherThan(other *Position) bool {
	if other == nil {
		return true
	}
	return p.Level <= other.Level
}

// Deactivate 無効化
func (p *Position) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// Activate 有効化
func (p *Position) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// PositionRepository 職位リポジトリインターフェース
type PositionRepository interface {
	// FindByID ID検索
	FindByID(ctx interface{}, id domain.ID) (*Position, error)
	// FindByOrganizationID 組織ID検索
	FindByOrganizationID(ctx interface{}, orgID domain.ID) ([]Position, error)
	// FindByCode コード検索
	FindByCode(ctx interface{}, orgID domain.ID, code string) (*Position, error)
	// Save 保存
	Save(ctx interface{}, position *Position) error
	// Delete 削除
	Delete(ctx interface{}, id domain.ID) error
}
