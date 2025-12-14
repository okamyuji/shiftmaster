// Package infrastructure スタッフインフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/shared/infrastructure"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// StaffModel スタッフDBモデル
type StaffModel struct {
	bun.BaseModel `bun:"table:staffs"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	TeamID         uuid.UUID  `bun:"team_id,type:uuid,notnull"`
	EmployeeCode   string     `bun:"employee_code"`
	FirstName      string     `bun:"first_name,notnull"`
	LastName       string     `bun:"last_name,notnull"`
	Email          string     `bun:"email"`
	Phone          string     `bun:"phone"`
	HireDate       *time.Time `bun:"hire_date,type:date"`
	EmploymentType string     `bun:"employment_type,notnull"`
	IsActive       bool       `bun:"is_active,notnull"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *StaffModel) ToDomain() *domain.Staff {
	return &domain.Staff{
		ID:             m.ID,
		TeamID:         m.TeamID,
		EmployeeCode:   m.EmployeeCode,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Email:          m.Email,
		Phone:          m.Phone,
		HireDate:       m.HireDate,
		EmploymentType: domain.EmploymentType(m.EmploymentType),
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromDomain ドメインエンティティからDBモデルへ変換
func (m *StaffModel) FromDomain(staff *domain.Staff) {
	m.ID = staff.ID
	m.TeamID = staff.TeamID
	m.EmployeeCode = staff.EmployeeCode
	m.FirstName = staff.FirstName
	m.LastName = staff.LastName
	m.Email = staff.Email
	m.Phone = staff.Phone
	m.HireDate = staff.HireDate
	m.EmploymentType = staff.EmploymentType.String()
	m.IsActive = staff.IsActive
	m.CreatedAt = staff.CreatedAt
	m.UpdatedAt = staff.UpdatedAt
}

// PostgresStaffRepository PostgreSQLスタッフリポジトリ
type PostgresStaffRepository struct {
	db *bun.DB
}

// NewPostgresStaffRepository リポジトリ生成
func NewPostgresStaffRepository(db *bun.DB) *PostgresStaffRepository {
	return &PostgresStaffRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresStaffRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.Staff, error) {
	model := &StaffModel{}
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
func (r *PostgresStaffRepository) FindAll(ctx context.Context, pagination infrastructure.Pagination) ([]domain.Staff, int, error) {
	var models []StaffModel

	count, err := r.db.NewSelect().
		Model(&models).
		Order("last_name ASC", "first_name ASC").
		Limit(pagination.Limit()).
		Offset(pagination.Offset()).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	staffs := make([]domain.Staff, len(models))
	for i, m := range models {
		staffs[i] = *m.ToDomain()
	}

	return staffs, count, nil
}

// FindByOrganizationID 組織IDで検索（マルチテナント対応）
// チーム→部署→組織の階層をJOINして組織IDでフィルタリング
func (r *PostgresStaffRepository) FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID, pagination infrastructure.Pagination) ([]domain.Staff, int, error) {
	var models []StaffModel

	// サブクエリで対象チームIDを取得
	teamSubquery := r.db.NewSelect().
		TableExpr("teams AS t").
		Column("t.id").
		Join("INNER JOIN departments AS d ON d.id = t.department_id").
		Where("d.organization_id = ?", orgID)

	count, err := r.db.NewSelect().
		Model(&models).
		Where("team_id IN (?)", teamSubquery).
		Order("last_name ASC", "first_name ASC").
		Limit(pagination.Limit()).
		Offset(pagination.Offset()).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	staffs := make([]domain.Staff, len(models))
	for i, m := range models {
		staffs[i] = *m.ToDomain()
	}

	return staffs, count, nil
}

// FindByTeamID チームIDで検索
func (r *PostgresStaffRepository) FindByTeamID(ctx context.Context, teamID sharedDomain.ID) ([]domain.Staff, error) {
	var models []StaffModel
	err := r.db.NewSelect().
		Model(&models).
		Where("team_id = ?", teamID).
		Order("last_name ASC", "first_name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	staffs := make([]domain.Staff, len(models))
	for i, m := range models {
		staffs[i] = *m.ToDomain()
	}

	return staffs, nil
}

// FindActive 有効スタッフのみ取得
func (r *PostgresStaffRepository) FindActive(ctx context.Context) ([]domain.Staff, error) {
	var models []StaffModel
	err := r.db.NewSelect().
		Model(&models).
		Where("is_active = ?", true).
		Order("last_name ASC", "first_name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	staffs := make([]domain.Staff, len(models))
	for i, m := range models {
		staffs[i] = *m.ToDomain()
	}

	return staffs, nil
}

// FindActiveByOrganizationID 組織IDで有効スタッフのみ取得
func (r *PostgresStaffRepository) FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]domain.Staff, error) {
	var models []StaffModel

	// サブクエリで対象チームIDを取得
	teamSubquery := r.db.NewSelect().
		TableExpr("teams AS t").
		Column("t.id").
		Join("INNER JOIN departments AS d ON d.id = t.department_id").
		Where("d.organization_id = ?", orgID)

	err := r.db.NewSelect().
		Model(&models).
		Where("team_id IN (?)", teamSubquery).
		Where("is_active = ?", true).
		Order("last_name ASC", "first_name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	staffs := make([]domain.Staff, len(models))
	for i, m := range models {
		staffs[i] = *m.ToDomain()
	}

	return staffs, nil
}

// Save 保存
func (r *PostgresStaffRepository) Save(ctx context.Context, staff *domain.Staff) error {
	model := &StaffModel{}
	model.FromDomain(staff)

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("team_id = EXCLUDED.team_id").
		Set("employee_code = EXCLUDED.employee_code").
		Set("first_name = EXCLUDED.first_name").
		Set("last_name = EXCLUDED.last_name").
		Set("email = EXCLUDED.email").
		Set("phone = EXCLUDED.phone").
		Set("hire_date = EXCLUDED.hire_date").
		Set("employment_type = EXCLUDED.employment_type").
		Set("is_active = EXCLUDED.is_active").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresStaffRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*StaffModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// TeamModel チームDBモデル
type TeamModel struct {
	bun.BaseModel `bun:"table:teams"`

	ID           uuid.UUID `bun:"id,pk,type:uuid"`
	DepartmentID uuid.UUID `bun:"department_id,type:uuid,notnull"`
	Name         string    `bun:"name,notnull"`
	Code         string    `bun:"code"`
	SortOrder    int       `bun:"sort_order,notnull"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	UpdatedAt    time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *TeamModel) ToDomain() *domain.Team {
	return &domain.Team{
		ID:           m.ID,
		DepartmentID: m.DepartmentID,
		Name:         m.Name,
		Code:         m.Code,
		SortOrder:    m.SortOrder,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// PostgresTeamRepository PostgreSQLチームリポジトリ
type PostgresTeamRepository struct {
	db *bun.DB
}

// NewPostgresTeamRepository リポジトリ生成
func NewPostgresTeamRepository(db *bun.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresTeamRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.Team, error) {
	model := &TeamModel{}
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
func (r *PostgresTeamRepository) FindAll(ctx context.Context) ([]domain.Team, error) {
	var models []TeamModel
	err := r.db.NewSelect().
		Model(&models).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]domain.Team, len(models))
	for i, m := range models {
		teams[i] = *m.ToDomain()
	}

	return teams, nil
}

// FindByOrganizationID 組織IDで検索（マルチテナント対応）
func (r *PostgresTeamRepository) FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]domain.Team, error) {
	var models []TeamModel

	// サブクエリで対象部署IDを取得
	deptSubquery := r.db.NewSelect().
		TableExpr("departments AS d").
		Column("d.id").
		Where("d.organization_id = ?", orgID)

	err := r.db.NewSelect().
		Model(&models).
		Where("department_id IN (?)", deptSubquery).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]domain.Team, len(models))
	for i, m := range models {
		teams[i] = *m.ToDomain()
	}

	return teams, nil
}

// FindByDepartmentID 部門IDで検索
func (r *PostgresTeamRepository) FindByDepartmentID(ctx context.Context, departmentID sharedDomain.ID) ([]domain.Team, error) {
	var models []TeamModel
	err := r.db.NewSelect().
		Model(&models).
		Where("department_id = ?", departmentID).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]domain.Team, len(models))
	for i, m := range models {
		teams[i] = *m.ToDomain()
	}

	return teams, nil
}

// Save 保存
func (r *PostgresTeamRepository) Save(ctx context.Context, team *domain.Team) error {
	model := &TeamModel{
		ID:           team.ID,
		DepartmentID: team.DepartmentID,
		Name:         team.Name,
		Code:         team.Code,
		SortOrder:    team.SortOrder,
		CreatedAt:    team.CreatedAt,
		UpdatedAt:    team.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("department_id = EXCLUDED.department_id").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("sort_order = EXCLUDED.sort_order").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresTeamRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*TeamModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DepartmentModel 部門DBモデル
type DepartmentModel struct {
	bun.BaseModel `bun:"table:departments"`

	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code"`
	SortOrder      int       `bun:"sort_order,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *DepartmentModel) ToDomain() *domain.Department {
	return &domain.Department{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
		SortOrder:      m.SortOrder,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// PostgresDepartmentRepository PostgreSQL部門リポジトリ
type PostgresDepartmentRepository struct {
	db *bun.DB
}

// NewPostgresDepartmentRepository リポジトリ生成
func NewPostgresDepartmentRepository(db *bun.DB) *PostgresDepartmentRepository {
	return &PostgresDepartmentRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresDepartmentRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.Department, error) {
	model := &DepartmentModel{}
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
func (r *PostgresDepartmentRepository) FindAll(ctx context.Context) ([]domain.Department, error) {
	var models []DepartmentModel
	err := r.db.NewSelect().
		Model(&models).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	departments := make([]domain.Department, len(models))
	for i, m := range models {
		departments[i] = *m.ToDomain()
	}

	return departments, nil
}

// FindByOrganizationID 組織IDで検索
func (r *PostgresDepartmentRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.Department, error) {
	var models []DepartmentModel
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Order("sort_order ASC", "name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	departments := make([]domain.Department, len(models))
	for i, m := range models {
		departments[i] = *m.ToDomain()
	}

	return departments, nil
}

// Save 保存
func (r *PostgresDepartmentRepository) Save(ctx context.Context, department *domain.Department) error {
	model := &DepartmentModel{
		ID:             department.ID,
		OrganizationID: department.OrganizationID,
		Name:           department.Name,
		Code:           department.Code,
		SortOrder:      department.SortOrder,
		CreatedAt:      department.CreatedAt,
		UpdatedAt:      department.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("organization_id = EXCLUDED.organization_id").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("sort_order = EXCLUDED.sort_order").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresDepartmentRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*DepartmentModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// OrganizationModel 組織DBモデル
type OrganizationModel struct {
	bun.BaseModel `bun:"table:organizations"`

	ID        uuid.UUID `bun:"id,pk,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Code      string    `bun:"code"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *OrganizationModel) ToDomain() *domain.Organization {
	return &domain.Organization{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// PostgresOrganizationRepository PostgreSQL組織リポジトリ
type PostgresOrganizationRepository struct {
	db *bun.DB
}

// NewPostgresOrganizationRepository リポジトリ生成
func NewPostgresOrganizationRepository(db *bun.DB) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

// FindByID IDで検索
func (r *PostgresOrganizationRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.Organization, error) {
	model := &OrganizationModel{}
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
func (r *PostgresOrganizationRepository) FindAll(ctx context.Context) ([]domain.Organization, error) {
	var models []OrganizationModel
	err := r.db.NewSelect().
		Model(&models).
		Order("name ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	orgs := make([]domain.Organization, len(models))
	for i, m := range models {
		orgs[i] = *m.ToDomain()
	}

	return orgs, nil
}

// Save 保存
func (r *PostgresOrganizationRepository) Save(ctx context.Context, org *domain.Organization) error {
	model := &OrganizationModel{
		ID:        org.ID,
		Name:      org.Name,
		Code:      org.Code,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(model).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("code = EXCLUDED.code").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

// Delete 削除
func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*OrganizationModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
