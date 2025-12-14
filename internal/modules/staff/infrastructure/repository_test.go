// Package infrastructure スタッフインフラストラクチャ層テスト
package infrastructure

import (
	"testing"
	"time"

	"shiftmaster/internal/modules/staff/domain"

	"github.com/google/uuid"
)

func TestStaffModel_ToDomain(t *testing.T) {
	now := time.Now()
	hireDate := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

	t.Run("全フィールドが設定されている場合", func(t *testing.T) {
		model := &StaffModel{
			ID:             uuid.New(),
			TeamID:         uuid.New(),
			EmployeeCode:   "EMP001",
			FirstName:      "太郎",
			LastName:       "田中",
			Email:          "taro@example.com",
			Phone:          "090-1234-5678",
			HireDate:       &hireDate,
			EmploymentType: "full_time",
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		staff := model.ToDomain()

		if staff.ID != model.ID {
			t.Error("ID should match")
		}
		if staff.TeamID != model.TeamID {
			t.Error("TeamID should match")
		}
		if staff.EmployeeCode != model.EmployeeCode {
			t.Errorf("expected employee code %s but got %s", model.EmployeeCode, staff.EmployeeCode)
		}
		if staff.FirstName != model.FirstName {
			t.Errorf("expected first name %s but got %s", model.FirstName, staff.FirstName)
		}
		if staff.LastName != model.LastName {
			t.Errorf("expected last name %s but got %s", model.LastName, staff.LastName)
		}
		if staff.Email != model.Email {
			t.Errorf("expected email %s but got %s", model.Email, staff.Email)
		}
		if staff.Phone != model.Phone {
			t.Errorf("expected phone %s but got %s", model.Phone, staff.Phone)
		}
		if staff.HireDate == nil {
			t.Fatal("HireDate should not be nil")
		}
		if !staff.HireDate.Equal(hireDate) {
			t.Error("HireDate should match")
		}
		if staff.EmploymentType != domain.EmploymentFullTime {
			t.Errorf("expected employment type full_time but got %s", staff.EmploymentType)
		}
		if !staff.IsActive {
			t.Error("IsActive should be true")
		}
	})

	t.Run("HireDateがnilの場合", func(t *testing.T) {
		model := &StaffModel{
			ID:             uuid.New(),
			TeamID:         uuid.New(),
			EmployeeCode:   "EMP002",
			FirstName:      "花子",
			LastName:       "山田",
			Email:          "hanako@example.com",
			EmploymentType: "part_time",
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		staff := model.ToDomain()

		if staff.HireDate != nil {
			t.Error("HireDate should be nil")
		}
		if staff.EmploymentType != domain.EmploymentPartTime {
			t.Errorf("expected employment type part_time but got %s", staff.EmploymentType)
		}
	})
}

func TestStaffModel_FromDomain(t *testing.T) {
	now := time.Now()
	hireDate := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

	t.Run("全フィールドが設定されている場合", func(t *testing.T) {
		staff := &domain.Staff{
			ID:             uuid.New(),
			TeamID:         uuid.New(),
			EmployeeCode:   "EMP001",
			FirstName:      "太郎",
			LastName:       "田中",
			Email:          "taro@example.com",
			Phone:          "090-1234-5678",
			HireDate:       &hireDate,
			EmploymentType: domain.EmploymentFullTime,
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		model := &StaffModel{}
		model.FromDomain(staff)

		if model.ID != staff.ID {
			t.Error("ID should match")
		}
		if model.TeamID != staff.TeamID {
			t.Error("TeamID should match")
		}
		if model.EmployeeCode != staff.EmployeeCode {
			t.Errorf("expected employee code %s but got %s", staff.EmployeeCode, model.EmployeeCode)
		}
		if model.EmploymentType != "full_time" {
			t.Errorf("expected employment type full_time but got %s", model.EmploymentType)
		}
		if model.HireDate == nil {
			t.Fatal("HireDate should not be nil")
		}
	})

	t.Run("HireDateがnilの場合", func(t *testing.T) {
		staff := &domain.Staff{
			ID:             uuid.New(),
			TeamID:         uuid.New(),
			EmployeeCode:   "EMP002",
			FirstName:      "花子",
			LastName:       "山田",
			EmploymentType: domain.EmploymentPartTime,
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		model := &StaffModel{}
		model.FromDomain(staff)

		if model.HireDate != nil {
			t.Error("HireDate should be nil")
		}
		if model.EmploymentType != "part_time" {
			t.Errorf("expected employment type part_time but got %s", model.EmploymentType)
		}
	})
}

func TestTeamModel_ToDomain(t *testing.T) {
	now := time.Now()

	t.Run("正常な変換", func(t *testing.T) {
		model := &TeamModel{
			ID:           uuid.New(),
			DepartmentID: uuid.New(),
			Name:         "テストチーム",
			Code:         "TEAM01",
			SortOrder:    1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		team := model.ToDomain()

		if team.ID != model.ID {
			t.Error("ID should match")
		}
		if team.DepartmentID != model.DepartmentID {
			t.Error("DepartmentID should match")
		}
		if team.Name != model.Name {
			t.Errorf("expected name %s but got %s", model.Name, team.Name)
		}
		if team.Code != model.Code {
			t.Errorf("expected code %s but got %s", model.Code, team.Code)
		}
		if team.SortOrder != model.SortOrder {
			t.Errorf("expected sort order %d but got %d", model.SortOrder, team.SortOrder)
		}
	})
}

// 境界値テスト

func TestStaffModel_ToDomain_BoundaryValues(t *testing.T) {
	t.Run("空文字列のフィールド", func(t *testing.T) {
		model := &StaffModel{
			ID:             uuid.New(),
			TeamID:         uuid.New(),
			EmployeeCode:   "",
			FirstName:      "",
			LastName:       "",
			Email:          "",
			Phone:          "",
			EmploymentType: "",
			IsActive:       false,
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}

		staff := model.ToDomain()

		if staff.EmployeeCode != "" {
			t.Error("EmployeeCode should be empty")
		}
		if staff.FirstName != "" {
			t.Error("FirstName should be empty")
		}
		if staff.IsActive {
			t.Error("IsActive should be false")
		}
	})

	t.Run("ゼロ値のUUID", func(t *testing.T) {
		model := &StaffModel{
			ID:             uuid.Nil,
			TeamID:         uuid.Nil,
			EmployeeCode:   "EMP001",
			FirstName:      "Test",
			LastName:       "User",
			EmploymentType: "full_time",
			IsActive:       true,
		}

		staff := model.ToDomain()

		if staff.ID != uuid.Nil {
			t.Error("ID should be nil UUID")
		}
		if staff.TeamID != uuid.Nil {
			t.Error("TeamID should be nil UUID")
		}
	})
}
