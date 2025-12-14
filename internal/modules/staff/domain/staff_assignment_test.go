package domain

import (
	"testing"
	"time"

	"shiftmaster/internal/shared/domain"
)

func TestNewStaffAssignment(t *testing.T) {
	staffID := domain.NewID()
	teamID := domain.NewID()
	jobTypeID := domain.NewID()
	positionID := domain.NewID()

	tests := []struct {
		name       string
		staffID    domain.ID
		teamID     *domain.ID
		jobTypeID  *domain.ID
		positionID *domain.ID
		isPrimary  bool
		wantErr    bool
	}{
		{"正常系_チームのみ", staffID, &teamID, nil, nil, true, false},
		{"正常系_職種のみ", staffID, nil, &jobTypeID, nil, false, false},
		{"正常系_職位のみ", staffID, nil, nil, &positionID, false, false},
		{"正常系_全指定", staffID, &teamID, &jobTypeID, &positionID, true, false},
		{"異常系_すべてnil", staffID, nil, nil, nil, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa, err := NewStaffAssignment(tt.staffID, tt.teamID, tt.jobTypeID, tt.positionID, tt.isPrimary)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStaffAssignment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if sa.StaffID != tt.staffID {
					t.Errorf("StaffID = %v, want %v", sa.StaffID, tt.staffID)
				}
				if sa.IsPrimary != tt.isPrimary {
					t.Errorf("IsPrimary = %v, want %v", sa.IsPrimary, tt.isPrimary)
				}
			}
		})
	}
}

func TestStaffAssignment_SetDateRange(t *testing.T) {
	staffID := domain.NewID()
	teamID := domain.NewID()
	sa, _ := NewStaffAssignment(staffID, &teamID, nil, nil, true)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name      string
		startDate *time.Time
		endDate   *time.Time
		wantErr   bool
	}{
		{"正常系_両方指定", &yesterday, &tomorrow, false},
		{"正常系_開始のみ", &yesterday, nil, false},
		{"正常系_終了のみ", nil, &tomorrow, false},
		{"正常系_両方nil", nil, nil, false},
		{"異常系_終了が開始より前", &tomorrow, &yesterday, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sa.SetDateRange(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffAssignment_IsActiveOn(t *testing.T) {
	staffID := domain.NewID()
	teamID := domain.NewID()
	sa, _ := NewStaffAssignment(staffID, &teamID, nil, nil, true)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	tests := []struct {
		name      string
		startDate *time.Time
		endDate   *time.Time
		checkDate time.Time
		want      bool
	}{
		{"範囲なし_常に有効", nil, nil, now, true},
		{"開始後_有効", &yesterday, nil, now, true},
		{"開始前_無効", &tomorrow, nil, now, false},
		{"終了前_有効", nil, &tomorrow, now, true},
		{"終了後_無効", nil, &yesterday, now, false},
		{"範囲内_有効", &yesterday, &nextWeek, now, true},
		{"範囲外_開始前", &tomorrow, &nextWeek, now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sa.SetDateRange(tt.startDate, tt.endDate)
			if got := sa.IsActiveOn(tt.checkDate); got != tt.want {
				t.Errorf("IsActiveOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaffAssignment_SetAsPrimary(t *testing.T) {
	staffID := domain.NewID()
	teamID := domain.NewID()
	sa, _ := NewStaffAssignment(staffID, &teamID, nil, nil, false)

	if sa.IsPrimary {
		t.Error("初期状態でIsPrimaryがtrue")
	}

	sa.SetAsPrimary()
	if !sa.IsPrimary {
		t.Error("SetAsPrimary後もIsPrimaryがfalse")
	}

	sa.UnsetAsPrimary()
	if sa.IsPrimary {
		t.Error("UnsetAsPrimary後もIsPrimaryがtrue")
	}
}

func TestStaffAssignment_End(t *testing.T) {
	staffID := domain.NewID()
	teamID := domain.NewID()
	sa, _ := NewStaffAssignment(staffID, &teamID, nil, nil, true)

	if sa.EndDate != nil {
		t.Error("初期状態でEndDateがnilでない")
	}

	endDate := time.Now()
	sa.End(endDate)

	if sa.EndDate == nil {
		t.Error("End後もEndDateがnil")
	}
	if !sa.EndDate.Equal(endDate) {
		t.Errorf("EndDate = %v, want %v", sa.EndDate, endDate)
	}
}
