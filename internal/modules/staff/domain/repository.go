// Package domain スタッフドメイン層
package domain

import (
	"context"

	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/shared/infrastructure"
)

// StaffRepository スタッフリポジトリインターフェース
type StaffRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Staff, error)
	// FindAll 全件取得
	FindAll(ctx context.Context, pagination infrastructure.Pagination) ([]Staff, int, error)
	// FindByOrganizationID 組織IDで検索（マルチテナント対応）
	FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID, pagination infrastructure.Pagination) ([]Staff, int, error)
	// FindByTeamID チームIDで検索
	FindByTeamID(ctx context.Context, teamID sharedDomain.ID) ([]Staff, error)
	// FindActive 有効スタッフのみ取得
	FindActive(ctx context.Context) ([]Staff, error)
	// FindActiveByOrganizationID 組織IDで有効スタッフのみ取得
	FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]Staff, error)
	// Save 保存 新規作成または更新
	Save(ctx context.Context, staff *Staff) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// TeamRepository チームリポジトリインターフェース
type TeamRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Team, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]Team, error)
	// FindByOrganizationID 組織IDで検索（マルチテナント対応）
	FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]Team, error)
	// FindByDepartmentID 部門IDで検索
	FindByDepartmentID(ctx context.Context, departmentID sharedDomain.ID) ([]Team, error)
	// Save 保存
	Save(ctx context.Context, team *Team) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// DepartmentRepository 部門リポジトリインターフェース
type DepartmentRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Department, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]Department, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]Department, error)
	// Save 保存
	Save(ctx context.Context, department *Department) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// OrganizationRepository 組織リポジトリインターフェース
type OrganizationRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Organization, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]Organization, error)
	// Save 保存
	Save(ctx context.Context, org *Organization) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// SkillRepository スキルリポジトリインターフェース
type SkillRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Skill, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]Skill, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]Skill, error)
	// Save 保存
	Save(ctx context.Context, skill *Skill) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}
