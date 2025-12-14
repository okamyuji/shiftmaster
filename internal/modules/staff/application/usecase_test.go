// Package application スタッフユースケーステスト
package application

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/shared/infrastructure"
)

// モックリポジトリ

type mockStaffRepository struct {
	staffs map[sharedDomain.ID]*domain.Staff
}

func newMockStaffRepository() *mockStaffRepository {
	return &mockStaffRepository{
		staffs: make(map[sharedDomain.ID]*domain.Staff),
	}
}

func (m *mockStaffRepository) FindByID(_ context.Context, id sharedDomain.ID) (*domain.Staff, error) {
	return m.staffs[id], nil
}

func (m *mockStaffRepository) FindByTeamID(_ context.Context, _ sharedDomain.ID) ([]domain.Staff, error) {
	return nil, nil
}

func (m *mockStaffRepository) FindAll(_ context.Context, _ infrastructure.Pagination) ([]domain.Staff, int, error) {
	result := make([]domain.Staff, 0, len(m.staffs))
	for _, staff := range m.staffs {
		result = append(result, *staff)
	}
	return result, len(result), nil
}

func (m *mockStaffRepository) FindActive(_ context.Context) ([]domain.Staff, error) {
	var result []domain.Staff
	for _, staff := range m.staffs {
		if staff.IsActive {
			result = append(result, *staff)
		}
	}
	return result, nil
}

func (m *mockStaffRepository) Save(_ context.Context, staff *domain.Staff) error {
	m.staffs[staff.ID] = staff
	return nil
}

func (m *mockStaffRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	delete(m.staffs, id)
	return nil
}

func (m *mockStaffRepository) FindByOrganizationID(_ context.Context, _ sharedDomain.ID, _ infrastructure.Pagination) ([]domain.Staff, int, error) {
	result := make([]domain.Staff, 0, len(m.staffs))
	for _, staff := range m.staffs {
		result = append(result, *staff)
	}
	return result, len(result), nil
}

func (m *mockStaffRepository) FindActiveByOrganizationID(_ context.Context, _ sharedDomain.ID) ([]domain.Staff, error) {
	var result []domain.Staff
	for _, staff := range m.staffs {
		if staff.IsActive {
			result = append(result, *staff)
		}
	}
	return result, nil
}

// モックチームリポジトリ

type mockTeamRepository struct {
	teams map[sharedDomain.ID]*domain.Team
}

func newMockTeamRepository() *mockTeamRepository {
	return &mockTeamRepository{
		teams: make(map[sharedDomain.ID]*domain.Team),
	}
}

func (m *mockTeamRepository) FindByID(_ context.Context, id sharedDomain.ID) (*domain.Team, error) {
	return m.teams[id], nil
}

func (m *mockTeamRepository) FindByDepartmentID(_ context.Context, _ sharedDomain.ID) ([]domain.Team, error) {
	return nil, nil
}

func (m *mockTeamRepository) FindAll(_ context.Context) ([]domain.Team, error) {
	result := make([]domain.Team, 0, len(m.teams))
	for _, team := range m.teams {
		result = append(result, *team)
	}
	return result, nil
}

func (m *mockTeamRepository) Save(_ context.Context, team *domain.Team) error {
	m.teams[team.ID] = team
	return nil
}

func (m *mockTeamRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	delete(m.teams, id)
	return nil
}

func (m *mockTeamRepository) FindByOrganizationID(_ context.Context, _ sharedDomain.ID) ([]domain.Team, error) {
	result := make([]domain.Team, 0, len(m.teams))
	for _, team := range m.teams {
		result = append(result, *team)
	}
	return result, nil
}

// モック部署リポジトリ

type mockDepartmentRepository struct {
	departments map[sharedDomain.ID]*domain.Department
}

func newMockDepartmentRepository() *mockDepartmentRepository {
	return &mockDepartmentRepository{
		departments: make(map[sharedDomain.ID]*domain.Department),
	}
}

func (m *mockDepartmentRepository) FindByID(_ context.Context, id sharedDomain.ID) (*domain.Department, error) {
	return m.departments[id], nil
}

func (m *mockDepartmentRepository) FindByOrganizationID(_ context.Context, _ sharedDomain.ID) ([]domain.Department, error) {
	result := make([]domain.Department, 0, len(m.departments))
	for _, dept := range m.departments {
		result = append(result, *dept)
	}
	return result, nil
}

func (m *mockDepartmentRepository) FindAll(_ context.Context) ([]domain.Department, error) {
	result := make([]domain.Department, 0, len(m.departments))
	for _, dept := range m.departments {
		result = append(result, *dept)
	}
	return result, nil
}

func (m *mockDepartmentRepository) Save(_ context.Context, dept *domain.Department) error {
	m.departments[dept.ID] = dept
	return nil
}

func (m *mockDepartmentRepository) Delete(_ context.Context, id sharedDomain.ID) error {
	delete(m.departments, id)
	return nil
}

// テスト

func TestNewStaffUseCase(t *testing.T) {
	staffRepo := newMockStaffRepository()
	teamRepo := newMockTeamRepository()
	deptRepo := newMockDepartmentRepository()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := NewStaffUseCase(staffRepo, teamRepo, deptRepo, logger)

	if useCase == nil {
		t.Fatal("NewStaffUseCase returned nil")
	}
	if useCase.staffRepo == nil {
		t.Error("staffRepo should not be nil")
	}
	if useCase.teamRepo == nil {
		t.Error("teamRepo should not be nil")
	}
	if useCase.deptRepo == nil {
		t.Error("deptRepo should not be nil")
	}
	if useCase.logger == nil {
		t.Error("logger should not be nil")
	}
}

func TestStaffUseCase_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name    string
		input   *CreateStaffInput
		setup   func(*mockStaffRepository, *mockTeamRepository)
		wantErr bool
		errCode string
	}{
		{
			name: "正常なスタッフ作成",
			input: &CreateStaffInput{
				TeamID:         "", // setupで設定
				EmployeeCode:   "EMP001",
				FirstName:      "太郎",
				LastName:       "田中",
				Email:          "taro@example.com",
				Phone:          "090-1234-5678",
				EmploymentType: "full_time",
			},
			setup: func(_ *mockStaffRepository, teamRepo *mockTeamRepository) {
				teamID := sharedDomain.NewID()
				team := &domain.Team{
					ID:   teamID,
					Name: "テストチーム",
				}
				_ = teamRepo.Save(context.Background(), team)
			},
			wantErr: false,
		},
		{
			name: "チーム未検出",
			input: &CreateStaffInput{
				TeamID:         sharedDomain.NewID().String(),
				EmployeeCode:   "EMP002",
				FirstName:      "花子",
				LastName:       "山田",
				Email:          "hanako@example.com",
				EmploymentType: "full_time",
			},
			setup:   func(_ *mockStaffRepository, _ *mockTeamRepository) {},
			wantErr: true,
			errCode: sharedDomain.ErrCodeNotFound,
		},
		{
			name: "社員コード未入力",
			input: &CreateStaffInput{
				TeamID:         "",
				EmployeeCode:   "",
				FirstName:      "一郎",
				LastName:       "鈴木",
				Email:          "ichiro@example.com",
				EmploymentType: "full_time",
			},
			setup:   func(_ *mockStaffRepository, _ *mockTeamRepository) {},
			wantErr: true,
			errCode: sharedDomain.ErrCodeValidation,
		},
		{
			name: "名前未入力",
			input: &CreateStaffInput{
				TeamID:         "",
				EmployeeCode:   "EMP003",
				FirstName:      "",
				LastName:       "",
				Email:          "noname@example.com",
				EmploymentType: "full_time",
			},
			setup:   func(_ *mockStaffRepository, _ *mockTeamRepository) {},
			wantErr: true,
			errCode: sharedDomain.ErrCodeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staffRepo := newMockStaffRepository()
			teamRepo := newMockTeamRepository()
			tt.setup(staffRepo, teamRepo)

			// チームが存在する場合、そのIDを入力に設定
			if tt.input.TeamID == "" && len(teamRepo.teams) > 0 {
				for id := range teamRepo.teams {
					tt.input.TeamID = id.String()
					break
				}
			}

			useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
			output, err := useCase.Create(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if domainErr, ok := err.(*sharedDomain.DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("expected error code %s but got %s", tt.errCode, domainErr.Code)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if output == nil {
					t.Error("expected output but got nil")
				} else {
					if output.EmployeeCode != tt.input.EmployeeCode {
						t.Errorf("expected employee code %s but got %s", tt.input.EmployeeCode, output.EmployeeCode)
					}
				}
			}
		})
	}
}

func TestStaffUseCase_GetByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("存在するスタッフ取得", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		staffID := sharedDomain.NewID()
		staff := &domain.Staff{
			ID:           staffID,
			EmployeeCode: "EMP001",
			FirstName:    "太郎",
			LastName:     "田中",
		}
		_ = staffRepo.Save(context.Background(), staff)

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		output, err := useCase.GetByID(context.Background(), staffID.String())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.ID != staffID.String() {
			t.Errorf("expected ID %s but got %s", staffID.String(), output.ID)
		}
	})

	t.Run("存在しないスタッフ取得", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		_, err := useCase.GetByID(context.Background(), sharedDomain.NewID().String())

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})

	t.Run("不正なID形式", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		_, err := useCase.GetByID(context.Background(), "invalid-id")

		if err == nil {
			t.Error("expected error but got nil")
		}
		if domainErr, ok := err.(*sharedDomain.DomainError); ok {
			if domainErr.Code != sharedDomain.ErrCodeValidation {
				t.Errorf("expected ErrCodeValidation but got %s", domainErr.Code)
			}
		}
	})
}

func TestStaffUseCase_List(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("スタッフ一覧取得_空", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		output, err := useCase.List(context.Background(), 1, 10)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.Total != 0 {
			t.Errorf("expected total 0 but got %d", output.Total)
		}
	})

	t.Run("スタッフ一覧取得_複数", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		// スタッフを追加
		for i := 0; i < 3; i++ {
			staff := &domain.Staff{
				ID:           sharedDomain.NewID(),
				EmployeeCode: "EMP00" + string(rune('1'+i)),
				FirstName:    "名前" + string(rune('1'+i)),
				LastName:     "苗字",
			}
			_ = staffRepo.Save(context.Background(), staff)
		}

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		output, err := useCase.List(context.Background(), 1, 10)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.Total != 3 {
			t.Errorf("expected total 3 but got %d", output.Total)
		}
	})
}

func TestStaffUseCase_Update(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なスタッフ更新", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		teamID := sharedDomain.NewID()
		team := &domain.Team{
			ID:   teamID,
			Name: "テストチーム",
		}
		_ = teamRepo.Save(context.Background(), team)

		staffID := sharedDomain.NewID()
		staff := &domain.Staff{
			ID:           staffID,
			TeamID:       teamID,
			EmployeeCode: "EMP001",
			FirstName:    "太郎",
			LastName:     "田中",
			IsActive:     true,
		}
		_ = staffRepo.Save(context.Background(), staff)

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		output, err := useCase.Update(context.Background(), &UpdateStaffInput{
			ID:             staffID.String(),
			TeamID:         teamID.String(),
			EmployeeCode:   "EMP001-UPDATED",
			FirstName:      "次郎",
			LastName:       "田中",
			EmploymentType: "part_time",
			IsActive:       true,
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected output but got nil")
		}
		if output.EmployeeCode != "EMP001-UPDATED" {
			t.Errorf("expected employee code EMP001-UPDATED but got %s", output.EmployeeCode)
		}
	})

	t.Run("存在しないスタッフ更新", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		teamID := sharedDomain.NewID()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		_, err := useCase.Update(context.Background(), &UpdateStaffInput{
			ID:             sharedDomain.NewID().String(),
			TeamID:         teamID.String(),
			EmployeeCode:   "EMP999",
			FirstName:      "不明",
			LastName:       "太郎",
			EmploymentType: "full_time",
			IsActive:       true,
		})

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})
}

func TestStaffUseCase_Delete(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("正常なスタッフ削除", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		staffID := sharedDomain.NewID()
		staff := &domain.Staff{
			ID:           staffID,
			EmployeeCode: "EMP001",
			FirstName:    "太郎",
			LastName:     "田中",
		}
		_ = staffRepo.Save(context.Background(), staff)

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		err := useCase.Delete(context.Background(), staffID.String())

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// 削除確認
		deleted, _ := staffRepo.FindByID(context.Background(), staffID)
		if deleted != nil {
			t.Error("staff should be deleted")
		}
	})

	t.Run("存在しないスタッフ削除", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		err := useCase.Delete(context.Background(), sharedDomain.NewID().String())

		if err != sharedDomain.ErrNotFound {
			t.Errorf("expected ErrNotFound but got %v", err)
		}
	})

	t.Run("不正なID形式", func(t *testing.T) {
		staffRepo := newMockStaffRepository()
		teamRepo := newMockTeamRepository()

		useCase := NewStaffUseCase(staffRepo, teamRepo, newMockDepartmentRepository(), logger)
		err := useCase.Delete(context.Background(), "invalid-id")

		if err == nil {
			t.Error("expected error but got nil")
		}
		if domainErr, ok := err.(*sharedDomain.DomainError); ok {
			if domainErr.Code != sharedDomain.ErrCodeValidation {
				t.Errorf("expected ErrCodeValidation but got %s", domainErr.Code)
			}
		}
	})
}
