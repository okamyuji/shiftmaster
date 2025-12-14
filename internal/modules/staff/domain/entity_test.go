// Package domain スタッフドメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ============================================
// Staff関連テスト
// ============================================

func TestStaff_FullName(t *testing.T) {
	tests := []struct {
		name      string
		lastName  string
		firstName string
		expected  string
	}{
		{
			name:      "正常系_日本語名前",
			lastName:  "山田",
			firstName: "太郎",
			expected:  "山田 太郎",
		},
		{
			name:      "正常系_英語名前",
			lastName:  "Smith",
			firstName: "John",
			expected:  "Smith John",
		},
		{
			name:      "姓のみ",
			lastName:  "山田",
			firstName: "",
			expected:  "山田 ",
		},
		{
			name:      "名のみ",
			lastName:  "",
			firstName: "太郎",
			expected:  " 太郎",
		},
		{
			name:      "空の名前",
			lastName:  "",
			firstName: "",
			expected:  " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			staff := &Staff{
				LastName:  tt.lastName,
				FirstName: tt.firstName,
			}

			result := staff.FullName()
			if result != tt.expected {
				t.Errorf("FullName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStaff_HasSkill(t *testing.T) {
	skillID1 := sharedDomain.NewID()
	skillID2 := sharedDomain.NewID()
	skillID3 := sharedDomain.NewID()

	staff := &Staff{
		Skills: []StaffSkill{
			{SkillID: skillID1, Level: 3},
			{SkillID: skillID2, Level: 5},
		},
	}

	t.Run("保有しているスキル", func(t *testing.T) {
		if !staff.HasSkill(skillID1) {
			t.Error("HasSkill() should return true for owned skill")
		}
		if !staff.HasSkill(skillID2) {
			t.Error("HasSkill() should return true for owned skill")
		}
	})

	t.Run("保有していないスキル", func(t *testing.T) {
		if staff.HasSkill(skillID3) {
			t.Error("HasSkill() should return false for not owned skill")
		}
	})

	t.Run("スキルなしのスタッフ", func(t *testing.T) {
		emptyStaff := &Staff{Skills: nil}
		if emptyStaff.HasSkill(skillID1) {
			t.Error("HasSkill() should return false for staff with no skills")
		}
	})

	t.Run("空のスキルスライス", func(t *testing.T) {
		emptyStaff := &Staff{Skills: []StaffSkill{}}
		if emptyStaff.HasSkill(skillID1) {
			t.Error("HasSkill() should return false for staff with empty skills slice")
		}
	})
}

func TestStaff_Structure(t *testing.T) {
	t.Run("Staff構造体の完全な初期化", func(t *testing.T) {
		staffID := sharedDomain.NewID()
		teamID := sharedDomain.NewID()
		now := time.Now()
		hireDate := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

		staff := &Staff{
			ID:             staffID,
			TeamID:         teamID,
			EmployeeCode:   "EMP001",
			FirstName:      "太郎",
			LastName:       "山田",
			Email:          "yamada@example.com",
			Phone:          "090-1234-5678",
			HireDate:       &hireDate,
			EmploymentType: EmploymentFullTime,
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if staff.ID != staffID {
			t.Errorf("ID = %v, want %v", staff.ID, staffID)
		}
		if staff.TeamID != teamID {
			t.Errorf("TeamID = %v, want %v", staff.TeamID, teamID)
		}
		if staff.EmployeeCode != "EMP001" {
			t.Errorf("EmployeeCode = %v, want %v", staff.EmployeeCode, "EMP001")
		}
		if staff.FirstName != "太郎" {
			t.Errorf("FirstName = %v, want %v", staff.FirstName, "太郎")
		}
		if staff.LastName != "山田" {
			t.Errorf("LastName = %v, want %v", staff.LastName, "山田")
		}
		if staff.Email != "yamada@example.com" {
			t.Errorf("Email = %v, want %v", staff.Email, "yamada@example.com")
		}
		if staff.Phone != "090-1234-5678" {
			t.Errorf("Phone = %v, want %v", staff.Phone, "090-1234-5678")
		}
		if !staff.IsActive {
			t.Error("IsActive should be true")
		}
		if staff.EmploymentType != EmploymentFullTime {
			t.Errorf("EmploymentType = %v, want %v", staff.EmploymentType, EmploymentFullTime)
		}
		if staff.HireDate != &hireDate {
			t.Errorf("HireDate = %v, want %v", staff.HireDate, hireDate)
		}
		if staff.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", staff.CreatedAt, now)
		}
		if staff.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", staff.UpdatedAt, now)
		}
	})
}

// ============================================
// EmploymentType関連テスト
// ============================================

func TestEmploymentType_String(t *testing.T) {
	tests := []struct {
		name           string
		employmentType EmploymentType
		expected       string
	}{
		{
			name:           "正社員",
			employmentType: EmploymentFullTime,
			expected:       "full_time",
		},
		{
			name:           "パートタイム",
			employmentType: EmploymentPartTime,
			expected:       "part_time",
		},
		{
			name:           "契約社員",
			employmentType: EmploymentContract,
			expected:       "contract",
		},
		{
			name:           "派遣社員",
			employmentType: EmploymentTemporary,
			expected:       "temporary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.employmentType.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEmploymentType_Label(t *testing.T) {
	tests := []struct {
		name           string
		employmentType EmploymentType
		expected       string
	}{
		{
			name:           "正社員",
			employmentType: EmploymentFullTime,
			expected:       "正社員",
		},
		{
			name:           "パートタイム",
			employmentType: EmploymentPartTime,
			expected:       "パートタイム",
		},
		{
			name:           "契約社員",
			employmentType: EmploymentContract,
			expected:       "契約社員",
		},
		{
			name:           "派遣社員",
			employmentType: EmploymentTemporary,
			expected:       "派遣社員",
		},
		{
			name:           "不明な雇用形態",
			employmentType: EmploymentType("unknown"),
			expected:       "不明",
		},
		{
			name:           "空の雇用形態",
			employmentType: EmploymentType(""),
			expected:       "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.employmentType.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// StaffSkill関連テスト
// ============================================

func TestStaffSkill_Structure(t *testing.T) {
	t.Run("StaffSkill構造体の初期化", func(t *testing.T) {
		staffID := sharedDomain.NewID()
		skillID := sharedDomain.NewID()
		acquiredAt := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)

		skill := &StaffSkill{
			StaffID:    staffID,
			SkillID:    skillID,
			Level:      4,
			AcquiredAt: &acquiredAt,
		}

		if skill.StaffID != staffID {
			t.Errorf("StaffID = %v, want %v", skill.StaffID, staffID)
		}
		if skill.SkillID != skillID {
			t.Errorf("SkillID = %v, want %v", skill.SkillID, skillID)
		}
		if skill.Level != 4 {
			t.Errorf("Level = %v, want %v", skill.Level, 4)
		}
		if skill.AcquiredAt != &acquiredAt {
			t.Errorf("AcquiredAt = %v, want %v", skill.AcquiredAt, acquiredAt)
		}
	})

	t.Run("レベル境界値", func(t *testing.T) {
		tests := []struct {
			name  string
			level int
		}{
			{"最小レベル", 1},
			{"最大レベル", 5},
			{"中間レベル", 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				skill := &StaffSkill{Level: tt.level}
				if skill.Level != tt.level {
					t.Errorf("Level = %v, want %v", skill.Level, tt.level)
				}
			})
		}
	})
}

// ============================================
// Team関連テスト
// ============================================

func TestTeam_Structure(t *testing.T) {
	t.Run("Team構造体の初期化", func(t *testing.T) {
		teamID := sharedDomain.NewID()
		deptID := sharedDomain.NewID()
		now := time.Now()

		team := &Team{
			ID:           teamID,
			DepartmentID: deptID,
			Name:         "看護1チーム",
			Code:         "TEAM1",
			SortOrder:    1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if team.ID != teamID {
			t.Errorf("ID = %v, want %v", team.ID, teamID)
		}
		if team.DepartmentID != deptID {
			t.Errorf("DepartmentID = %v, want %v", team.DepartmentID, deptID)
		}
		if team.Name != "看護1チーム" {
			t.Errorf("Name = %v, want %v", team.Name, "看護1チーム")
		}
		if team.Code != "TEAM1" {
			t.Errorf("Code = %v, want %v", team.Code, "TEAM1")
		}
		if team.SortOrder != 1 {
			t.Errorf("SortOrder = %v, want %v", team.SortOrder, 1)
		}
		if team.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", team.CreatedAt, now)
		}
		if team.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", team.UpdatedAt, now)
		}
	})
}

// ============================================
// Department関連テスト
// ============================================

func TestDepartment_Structure(t *testing.T) {
	t.Run("Department構造体の初期化", func(t *testing.T) {
		deptID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()

		dept := &Department{
			ID:             deptID,
			OrganizationID: orgID,
			Name:           "看護部",
			Code:           "NURS",
			SortOrder:      1,
			Teams:          []Team{},
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if dept.ID != deptID {
			t.Errorf("ID = %v, want %v", dept.ID, deptID)
		}
		if dept.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", dept.OrganizationID, orgID)
		}
		if dept.Name != "看護部" {
			t.Errorf("Name = %v, want %v", dept.Name, "看護部")
		}
		if dept.Code != "NURS" {
			t.Errorf("Code = %v, want %v", dept.Code, "NURS")
		}
		if dept.SortOrder != 1 {
			t.Errorf("SortOrder = %v, want %v", dept.SortOrder, 1)
		}
		if len(dept.Teams) != 0 {
			t.Errorf("Teams length = %d, want 0", len(dept.Teams))
		}
		if dept.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", dept.CreatedAt, now)
		}
		if dept.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", dept.UpdatedAt, now)
		}
	})

	t.Run("チームを含む部門", func(t *testing.T) {
		dept := &Department{
			Name: "看護部",
			Teams: []Team{
				{Name: "チーム1"},
				{Name: "チーム2"},
			},
		}

		if dept.Name != "看護部" {
			t.Errorf("Name = %v, want %v", dept.Name, "看護部")
		}
		if len(dept.Teams) != 2 {
			t.Errorf("Teams length = %d, want 2", len(dept.Teams))
		}
	})
}

// ============================================
// Organization関連テスト
// ============================================

func TestOrganization_Structure(t *testing.T) {
	t.Run("Organization構造体の初期化", func(t *testing.T) {
		orgID := sharedDomain.NewID()
		now := time.Now()

		org := &Organization{
			ID:          orgID,
			Name:        "サンプル病院",
			Code:        "HOSP1",
			Departments: []Department{},
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if org.ID != orgID {
			t.Errorf("ID = %v, want %v", org.ID, orgID)
		}
		if org.Name != "サンプル病院" {
			t.Errorf("Name = %v, want %v", org.Name, "サンプル病院")
		}
		if org.Code != "HOSP1" {
			t.Errorf("Code = %v, want %v", org.Code, "HOSP1")
		}
		if len(org.Departments) != 0 {
			t.Errorf("Departments length = %d, want 0", len(org.Departments))
		}
		if org.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", org.CreatedAt, now)
		}
		if org.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", org.UpdatedAt, now)
		}
	})

	t.Run("部門を含む組織", func(t *testing.T) {
		org := &Organization{
			Name: "サンプル病院",
			Departments: []Department{
				{Name: "看護部"},
				{Name: "医事課"},
			},
		}

		if org.Name != "サンプル病院" {
			t.Errorf("Name = %v, want %v", org.Name, "サンプル病院")
		}
		if len(org.Departments) != 2 {
			t.Errorf("Departments length = %d, want 2", len(org.Departments))
		}
	})
}

// ============================================
// Skill関連テスト
// ============================================

func TestSkill_Structure(t *testing.T) {
	t.Run("Skill構造体の初期化", func(t *testing.T) {
		skillID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()

		skill := &Skill{
			ID:             skillID,
			OrganizationID: orgID,
			Name:           "注射",
			Description:    "静脈注射が可能",
			Color:          "#FF5722",
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if skill.ID != skillID {
			t.Errorf("ID = %v, want %v", skill.ID, skillID)
		}
		if skill.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", skill.OrganizationID, orgID)
		}
		if skill.Name != "注射" {
			t.Errorf("Name = %v, want %v", skill.Name, "注射")
		}
		if skill.Description != "静脈注射が可能" {
			t.Errorf("Description = %v, want %v", skill.Description, "静脈注射が可能")
		}
		if skill.Color != "#FF5722" {
			t.Errorf("Color = %v, want %v", skill.Color, "#FF5722")
		}
		if skill.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", skill.CreatedAt, now)
		}
		if skill.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", skill.UpdatedAt, now)
		}
	})
}

// ============================================
// 定数テスト
// ============================================

func TestEmploymentTypeConstants(t *testing.T) {
	t.Run("雇用形態定数の値", func(t *testing.T) {
		if EmploymentFullTime != "full_time" {
			t.Errorf("EmploymentFullTime = %v, want %v", EmploymentFullTime, "full_time")
		}
		if EmploymentPartTime != "part_time" {
			t.Errorf("EmploymentPartTime = %v, want %v", EmploymentPartTime, "part_time")
		}
		if EmploymentContract != "contract" {
			t.Errorf("EmploymentContract = %v, want %v", EmploymentContract, "contract")
		}
		if EmploymentTemporary != "temporary" {
			t.Errorf("EmploymentTemporary = %v, want %v", EmploymentTemporary, "temporary")
		}
	})

	t.Run("雇用形態定数の一意性", func(t *testing.T) {
		types := []EmploymentType{
			EmploymentFullTime,
			EmploymentPartTime,
			EmploymentContract,
			EmploymentTemporary,
		}

		seen := make(map[EmploymentType]bool)
		for _, et := range types {
			if seen[et] {
				t.Errorf("Duplicate employment type constant: %v", et)
			}
			seen[et] = true
		}
	})
}
