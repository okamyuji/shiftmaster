package domain

import (
	"testing"
	"time"

	"shiftmaster/internal/shared/domain"
)

func TestScheduleValidator_ValidateNoTimeOverlap_SameDay(t *testing.T) {
	validator := NewScheduleValidator()
	staffID := domain.NewID()
	scheduleID := domain.NewID()
	shiftID1 := domain.NewID()
	shiftID2 := domain.NewID()

	// 同一日に重複するシフト
	targetDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name           string
		shift1         ShiftTimeInfo
		shift2         ShiftTimeInfo
		wantViolations int
	}{
		{
			name:           "重複あり_日勤同士",
			shift1:         ShiftTimeInfo{ID: shiftID1, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true},
			shift2:         ShiftTimeInfo{ID: shiftID2, StartTime: "12:00", EndTime: "20:00", IsWorkingShift: true},
			wantViolations: 1,
		},
		{
			name:           "重複なし_連続シフト",
			shift1:         ShiftTimeInfo{ID: shiftID1, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true},
			shift2:         ShiftTimeInfo{ID: shiftID2, StartTime: "17:00", EndTime: "22:00", IsWorkingShift: true},
			wantViolations: 0,
		},
		{
			name:           "重複あり_完全包含",
			shift1:         ShiftTimeInfo{ID: shiftID1, StartTime: "08:00", EndTime: "20:00", IsWorkingShift: true},
			shift2:         ShiftTimeInfo{ID: shiftID2, StartTime: "10:00", EndTime: "18:00", IsWorkingShift: true},
			wantViolations: 1,
		},
		{
			name:           "公休はスキップ",
			shift1:         ShiftTimeInfo{ID: shiftID1, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true},
			shift2:         ShiftTimeInfo{ID: shiftID2, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: false},
			wantViolations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := []ScheduleEntry{
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: targetDate, ShiftTypeID: &shiftID1},
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: targetDate, ShiftTypeID: &shiftID2},
			}
			shiftInfoMap := map[domain.ID]ShiftTimeInfo{
				shiftID1: tt.shift1,
				shiftID2: tt.shift2,
			}

			violations := validator.ValidateEntries(entries, shiftInfoMap)
			if len(violations) != tt.wantViolations {
				t.Errorf("ValidateEntries() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestScheduleValidator_ValidateNoTimeOverlap_ConsecutiveDays(t *testing.T) {
	validator := NewScheduleValidator()
	staffID := domain.NewID()
	scheduleID := domain.NewID()
	nightShiftID := domain.NewID()
	dayShiftID := domain.NewID()
	earlyShiftID := domain.NewID()

	day1 := time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)
	day2 := time.Date(2025, 1, 16, 0, 0, 0, 0, time.Local)

	nightShift := ShiftTimeInfo{ID: nightShiftID, StartTime: "17:00", EndTime: "09:00", IsWorkingShift: true}
	dayShift := ShiftTimeInfo{ID: dayShiftID, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true}
	earlyShift := ShiftTimeInfo{ID: earlyShiftID, StartTime: "06:00", EndTime: "14:00", IsWorkingShift: true}

	tests := []struct {
		name           string
		day1Shift      ShiftTimeInfo
		day2Shift      ShiftTimeInfo
		wantViolations int
	}{
		{
			name:           "夜勤後の早番_重複あり",
			day1Shift:      nightShift,
			day2Shift:      earlyShift,
			wantViolations: 1,
		},
		{
			name:           "夜勤後の日勤_重複なし",
			day1Shift:      nightShift,
			day2Shift:      dayShift,
			wantViolations: 0,
		},
		{
			name:           "日勤後の早番_重複なし",
			day1Shift:      dayShift,
			day2Shift:      earlyShift,
			wantViolations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := []ScheduleEntry{
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: day1, ShiftTypeID: &tt.day1Shift.ID},
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: day2, ShiftTypeID: &tt.day2Shift.ID},
			}
			shiftInfoMap := map[domain.ID]ShiftTimeInfo{
				tt.day1Shift.ID: tt.day1Shift,
				tt.day2Shift.ID: tt.day2Shift,
			}

			violations := validator.ValidateEntries(entries, shiftInfoMap)
			if len(violations) != tt.wantViolations {
				t.Errorf("ValidateEntries() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestScheduleValidator_ValidateNoTimeOverlap_DifferentStaff(t *testing.T) {
	validator := NewScheduleValidator()
	staffID1 := domain.NewID()
	staffID2 := domain.NewID()
	scheduleID := domain.NewID()
	shiftID := domain.NewID()

	targetDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)

	// 異なるスタッフなら同じ時間でも重複なし
	entries := []ScheduleEntry{
		{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID1, TargetDate: targetDate, ShiftTypeID: &shiftID},
		{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID2, TargetDate: targetDate, ShiftTypeID: &shiftID},
	}
	shiftInfoMap := map[domain.ID]ShiftTimeInfo{
		shiftID: {ID: shiftID, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true},
	}

	violations := validator.ValidateEntries(entries, shiftInfoMap)
	if len(violations) != 0 {
		t.Errorf("異なるスタッフなら重複なし、violations = %d, want 0", len(violations))
	}
}

func TestScheduleValidator_ValidateConsecutiveWorkDays(t *testing.T) {
	validator := NewScheduleValidator()
	staffID := domain.NewID()
	scheduleID := domain.NewID()
	shiftID := domain.NewID()

	shiftInfo := ShiftTimeInfo{ID: shiftID, StartTime: "09:00", EndTime: "17:00", IsWorkingShift: true}
	shiftInfoMap := map[domain.ID]ShiftTimeInfo{shiftID: shiftInfo}

	tests := []struct {
		name               string
		consecutiveDays    int
		maxConsecutiveDays int
		wantViolations     int
	}{
		{"5連勤_上限5_違反なし", 5, 5, 0},
		{"6連勤_上限5_違反あり", 6, 5, 1},
		{"7連勤_上限5_違反2回", 7, 5, 2},
		{"3連勤_上限5_違反なし", 3, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entries []ScheduleEntry
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
			for i := 0; i < tt.consecutiveDays; i++ {
				date := baseDate.AddDate(0, 0, i)
				entries = append(entries, ScheduleEntry{
					ID:          domain.NewID(),
					ScheduleID:  scheduleID,
					StaffID:     staffID,
					TargetDate:  date,
					ShiftTypeID: &shiftID,
				})
			}

			violations := validator.ValidateConsecutiveWorkDays(entries, shiftInfoMap, tt.maxConsecutiveDays)
			if len(violations) != tt.wantViolations {
				t.Errorf("ValidateConsecutiveWorkDays() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestScheduleValidator_ValidateNightShiftInterval(t *testing.T) {
	validator := NewScheduleValidator()
	staffID := domain.NewID()
	scheduleID := domain.NewID()
	nightShiftID := domain.NewID()

	nightShift := ShiftTimeInfo{ID: nightShiftID, StartTime: "17:00", EndTime: "09:00", IsWorkingShift: true}
	shiftInfoMap := map[domain.ID]ShiftTimeInfo{nightShiftID: nightShift}

	tests := []struct {
		name            string
		intervalDays    int
		minIntervalDays int
		wantViolations  int
	}{
		{"2日間隔_最小2日_違反なし", 2, 2, 0},
		{"1日間隔_最小2日_違反あり", 1, 2, 1},
		{"3日間隔_最小2日_違反なし", 3, 2, 0},
		{"連続夜勤_最小1日_違反なし", 1, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
			entries := []ScheduleEntry{
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: baseDate, ShiftTypeID: &nightShiftID},
				{ID: domain.NewID(), ScheduleID: scheduleID, StaffID: staffID, TargetDate: baseDate.AddDate(0, 0, tt.intervalDays), ShiftTypeID: &nightShiftID},
			}

			violations := validator.ValidateNightShiftInterval(entries, shiftInfoMap, tt.minIntervalDays)
			if len(violations) != tt.wantViolations {
				t.Errorf("ValidateNightShiftInterval() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestShiftTimeInfo_CrossesDate(t *testing.T) {
	tests := []struct {
		name      string
		startTime string
		endTime   string
		want      bool
	}{
		{"日勤", "09:00", "17:00", false},
		{"夜勤", "17:00", "09:00", true},
		{"深夜勤", "22:00", "06:00", true},
		{"早番", "06:00", "14:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ShiftTimeInfo{
				ID:             domain.NewID(),
				StartTime:      tt.startTime,
				EndTime:        tt.endTime,
				IsWorkingShift: true,
			}
			if got := info.CrossesDate(); got != tt.want {
				t.Errorf("CrossesDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
