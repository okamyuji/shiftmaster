// Package application スタッフアプリケーション層
package application

import (
	"context"
	"log/slog"
	"time"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/shared/infrastructure"
)

// StaffUseCase スタッフユースケース
type StaffUseCase struct {
	staffRepo domain.StaffRepository
	teamRepo  domain.TeamRepository
	deptRepo  domain.DepartmentRepository
	logger    *slog.Logger
}

// NewStaffUseCase スタッフユースケース生成
func NewStaffUseCase(
	staffRepo domain.StaffRepository,
	teamRepo domain.TeamRepository,
	deptRepo domain.DepartmentRepository,
	logger *slog.Logger,
) *StaffUseCase {
	return &StaffUseCase{
		staffRepo: staffRepo,
		teamRepo:  teamRepo,
		deptRepo:  deptRepo,
		logger:    logger,
	}
}

// verifyStaffBelongsToOrganization スタッフが指定組織に属しているか検証
func (u *StaffUseCase) verifyStaffBelongsToOrganization(ctx context.Context, staff *domain.Staff, orgID sharedDomain.ID) error {
	team, err := u.teamRepo.FindByID(ctx, staff.TeamID)
	if err != nil {
		return err
	}
	if team == nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeNotFound, "チームが見つかりません")
	}

	dept, err := u.deptRepo.FindByID(ctx, team.DepartmentID)
	if err != nil {
		return err
	}
	if dept == nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeNotFound, "部署が見つかりません")
	}

	if dept.OrganizationID != orgID {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeForbidden, "このスタッフへのアクセス権限がありません")
	}

	return nil
}

// Create スタッフ作成
func (u *StaffUseCase) Create(ctx context.Context, input *CreateStaffInput) (*StaffOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	teamID, err := sharedDomain.ParseID(input.TeamID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "チームIDが不正です")
	}

	// チーム存在確認
	team, err := u.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeNotFound, "チームが見つかりません")
	}

	var hireDate *time.Time
	if input.HireDate != "" {
		t, err := time.Parse("2006-01-02", input.HireDate)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "入社日の形式が不正です")
		}
		hireDate = &t
	}

	empType := domain.EmploymentFullTime
	if input.EmploymentType != "" {
		empType = domain.EmploymentType(input.EmploymentType)
	}

	now := time.Now()
	staff := &domain.Staff{
		ID:             sharedDomain.NewID(),
		TeamID:         teamID,
		EmployeeCode:   input.EmployeeCode,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Phone:          input.Phone,
		HireDate:       hireDate,
		EmploymentType: empType,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.staffRepo.Save(ctx, staff); err != nil {
		u.logger.Error("スタッフ作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("スタッフ作成完了", "staff_id", staff.ID, "name", staff.FullName())
	return ToStaffOutput(staff), nil
}

// Update スタッフ更新
func (u *StaffUseCase) Update(ctx context.Context, input *UpdateStaffInput) (*StaffOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	id, err := sharedDomain.ParseID(input.ID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	staff, err := u.staffRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, sharedDomain.ErrNotFound
	}

	teamID, err := sharedDomain.ParseID(input.TeamID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "チームIDが不正です")
	}

	var hireDate *time.Time
	if input.HireDate != "" {
		t, err := time.Parse("2006-01-02", input.HireDate)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "入社日の形式が不正です")
		}
		hireDate = &t
	}

	staff.TeamID = teamID
	staff.EmployeeCode = input.EmployeeCode
	staff.FirstName = input.FirstName
	staff.LastName = input.LastName
	staff.Email = input.Email
	staff.Phone = input.Phone
	staff.HireDate = hireDate
	staff.EmploymentType = domain.EmploymentType(input.EmploymentType)
	staff.IsActive = input.IsActive
	staff.UpdatedAt = time.Now()

	if err := u.staffRepo.Save(ctx, staff); err != nil {
		u.logger.Error("スタッフ更新失敗", "error", err)
		return nil, err
	}

	u.logger.Info("スタッフ更新完了", "staff_id", staff.ID)
	return ToStaffOutput(staff), nil
}

// GetByID IDでスタッフ取得
// Deprecated: GetByIDWithOrg を使用してください
func (u *StaffUseCase) GetByID(ctx context.Context, id string) (*StaffOutput, error) {
	staffID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	staff, err := u.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToStaffOutput(staff), nil
}

// GetByIDWithOrg IDと組織IDでスタッフ取得（マルチテナント対応）
func (u *StaffUseCase) GetByIDWithOrg(ctx context.Context, id, orgID string) (*StaffOutput, error) {
	staffID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	staff, err := u.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, sharedDomain.ErrNotFound
	}

	// 組織チェック
	if err := u.verifyStaffBelongsToOrganization(ctx, staff, organizationID); err != nil {
		return nil, err
	}

	return ToStaffOutput(staff), nil
}

// List スタッフ一覧取得
// Deprecated: ListByOrganization を使用してください
func (u *StaffUseCase) List(ctx context.Context, page, perPage int) (*StaffListOutput, error) {
	pagination := infrastructure.NewPagination(page, perPage)

	staffs, total, err := u.staffRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	outputs := make([]StaffOutput, len(staffs))
	for i, s := range staffs {
		outputs[i] = *ToStaffOutput(&s)
	}

	return &StaffListOutput{
		Staffs:  outputs,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.Limit(),
	}, nil
}

// ListByOrganization 組織IDでスタッフ一覧取得（マルチテナント対応）
func (u *StaffUseCase) ListByOrganization(ctx context.Context, orgID string, page, perPage int) (*StaffListOutput, error) {
	if orgID == "" {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが必要です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	pagination := infrastructure.NewPagination(page, perPage)

	staffs, total, err := u.staffRepo.FindByOrganizationID(ctx, organizationID, pagination)
	if err != nil {
		return nil, err
	}

	outputs := make([]StaffOutput, len(staffs))
	for i, s := range staffs {
		outputs[i] = *ToStaffOutput(&s)
	}

	return &StaffListOutput{
		Staffs:  outputs,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.Limit(),
	}, nil
}

// Delete スタッフ削除
// Deprecated: DeleteWithOrg を使用してください
func (u *StaffUseCase) Delete(ctx context.Context, id string) error {
	staffID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	staff, err := u.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		return err
	}
	if staff == nil {
		return sharedDomain.ErrNotFound
	}

	if err := u.staffRepo.Delete(ctx, staffID); err != nil {
		u.logger.Error("スタッフ削除失敗", "error", err)
		return err
	}

	u.logger.Info("スタッフ削除完了", "staff_id", staffID)
	return nil
}

// DeleteWithOrg 組織IDを検証してスタッフ削除（マルチテナント対応）
func (u *StaffUseCase) DeleteWithOrg(ctx context.Context, id, orgID string) error {
	staffID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	staff, err := u.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		return err
	}
	if staff == nil {
		return sharedDomain.ErrNotFound
	}

	// 組織チェック
	if err := u.verifyStaffBelongsToOrganization(ctx, staff, organizationID); err != nil {
		return err
	}

	if err := u.staffRepo.Delete(ctx, staffID); err != nil {
		u.logger.Error("スタッフ削除失敗", "error", err)
		return err
	}

	u.logger.Info("スタッフ削除完了", "staff_id", staffID, "organization_id", orgID)
	return nil
}
