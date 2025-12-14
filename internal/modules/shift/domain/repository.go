// Package domain シフトドメイン層
package domain

import (
	"context"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ShiftTypeRepository シフト種別リポジトリインターフェース
type ShiftTypeRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ShiftType, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]ShiftType, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftType, error)
	// FindWorkShifts 勤務シフトのみ取得 休日除外
	FindWorkShifts(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftType, error)
	// Save 保存
	Save(ctx context.Context, shiftType *ShiftType) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// ShiftPatternRepository シフトパターン（直）リポジトリインターフェース
type ShiftPatternRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ShiftPattern, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]ShiftPattern, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftPattern, error)
	// FindActiveByOrganizationID 有効パターンのみ取得
	FindActiveByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftPattern, error)
	// Save 保存
	Save(ctx context.Context, pattern *ShiftPattern) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// ShiftRuleRepository シフトルールリポジトリインターフェース
type ShiftRuleRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ShiftRule, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]ShiftRule, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftRule, error)
	// FindActiveByOrganizationID 有効ルールのみ取得
	FindActiveByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]ShiftRule, error)
	// Save 保存
	Save(ctx context.Context, rule *ShiftRule) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}
