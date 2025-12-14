// Package domain 勤務表ドメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ============================================
// Schedule関連テスト
// ============================================

func TestSchedule_TargetPeriodLabel(t *testing.T) {
	tests := []struct {
		name        string
		targetYear  int
		targetMonth int
		expected    string
	}{
		{
			name:        "2025年1月",
			targetYear:  2025,
			targetMonth: 1,
			expected:    "2025年1月",
		},
		{
			name:        "2025年12月",
			targetYear:  2025,
			targetMonth: 12,
			expected:    "2025年12月",
		},
		{
			name:        "2024年6月",
			targetYear:  2024,
			targetMonth: 6,
			expected:    "2024年6月",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &Schedule{
				TargetYear:  tt.targetYear,
				TargetMonth: tt.targetMonth,
			}

			result := schedule.TargetPeriodLabel()
			if result != tt.expected {
				t.Errorf("TargetPeriodLabel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchedule_DaysInMonth(t *testing.T) {
	tests := []struct {
		name        string
		targetYear  int
		targetMonth int
		expected    int
	}{
		{
			name:        "1月_31日",
			targetYear:  2025,
			targetMonth: 1,
			expected:    31,
		},
		{
			name:        "2月_通常年28日",
			targetYear:  2025,
			targetMonth: 2,
			expected:    28,
		},
		{
			name:        "2月_閏年29日",
			targetYear:  2024,
			targetMonth: 2,
			expected:    29,
		},
		{
			name:        "4月_30日",
			targetYear:  2025,
			targetMonth: 4,
			expected:    30,
		},
		{
			name:        "12月_31日",
			targetYear:  2025,
			targetMonth: 12,
			expected:    31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &Schedule{
				TargetYear:  tt.targetYear,
				TargetMonth: tt.targetMonth,
			}

			result := schedule.DaysInMonth()
			if result != tt.expected {
				t.Errorf("DaysInMonth() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchedule_StartDate(t *testing.T) {
	t.Run("月の開始日", func(t *testing.T) {
		schedule := &Schedule{
			TargetYear:  2025,
			TargetMonth: 3,
		}

		result := schedule.StartDate()

		if result.Year() != 2025 || result.Month() != time.March || result.Day() != 1 {
			t.Errorf("StartDate() = %v, want 2025-03-01", result)
		}
	})
}

func TestSchedule_EndDate(t *testing.T) {
	tests := []struct {
		name        string
		targetYear  int
		targetMonth int
		expectedDay int
	}{
		{
			name:        "1月末日",
			targetYear:  2025,
			targetMonth: 1,
			expectedDay: 31,
		},
		{
			name:        "2月末日_通常年",
			targetYear:  2025,
			targetMonth: 2,
			expectedDay: 28,
		},
		{
			name:        "2月末日_閏年",
			targetYear:  2024,
			targetMonth: 2,
			expectedDay: 29,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &Schedule{
				TargetYear:  tt.targetYear,
				TargetMonth: tt.targetMonth,
			}

			result := schedule.EndDate()

			if result.Day() != tt.expectedDay {
				t.Errorf("EndDate().Day() = %v, want %v", result.Day(), tt.expectedDay)
			}
		})
	}
}

func TestSchedule_Structure(t *testing.T) {
	t.Run("Schedule構造体の完全な初期化", func(t *testing.T) {
		scheduleID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()
		publishedAt := now.Add(-1 * time.Hour)

		schedule := &Schedule{
			ID:             scheduleID,
			OrganizationID: orgID,
			TargetYear:     2025,
			TargetMonth:    3,
			Status:         StatusPublished,
			PublishedAt:    &publishedAt,
			Entries:        []ScheduleEntry{},
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if schedule.ID != scheduleID {
			t.Errorf("ID = %v, want %v", schedule.ID, scheduleID)
		}
		if schedule.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", schedule.OrganizationID, orgID)
		}
		if schedule.TargetYear != 2025 {
			t.Errorf("TargetYear = %v, want %v", schedule.TargetYear, 2025)
		}
		if schedule.TargetMonth != 3 {
			t.Errorf("TargetMonth = %v, want %v", schedule.TargetMonth, 3)
		}
		if schedule.Status != StatusPublished {
			t.Errorf("Status = %v, want %v", schedule.Status, StatusPublished)
		}
		if schedule.PublishedAt == nil {
			t.Error("PublishedAt should not be nil")
		}
		if len(schedule.Entries) != 0 {
			t.Errorf("Entries length = %v, want %v", len(schedule.Entries), 0)
		}
		if !schedule.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt = %v, want %v", schedule.CreatedAt, now)
		}
		if !schedule.UpdatedAt.Equal(now) {
			t.Errorf("UpdatedAt = %v, want %v", schedule.UpdatedAt, now)
		}
	})
}

// ============================================
// ScheduleStatus関連テスト
// ============================================

func TestScheduleStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   ScheduleStatus
		expected string
	}{
		{
			name:     "draft",
			status:   StatusDraft,
			expected: "draft",
		},
		{
			name:     "in_progress",
			status:   StatusInProgress,
			expected: "in_progress",
		},
		{
			name:     "completed",
			status:   StatusCompleted,
			expected: "completed",
		},
		{
			name:     "published",
			status:   StatusPublished,
			expected: "published",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScheduleStatus_Label(t *testing.T) {
	tests := []struct {
		name     string
		status   ScheduleStatus
		expected string
	}{
		{
			name:     "draft",
			status:   StatusDraft,
			expected: "下書き",
		},
		{
			name:     "in_progress",
			status:   StatusInProgress,
			expected: "作成中",
		},
		{
			name:     "completed",
			status:   StatusCompleted,
			expected: "作成完了",
		},
		{
			name:     "published",
			status:   StatusPublished,
			expected: "公開済み",
		},
		{
			name:     "不明なステータス",
			status:   ScheduleStatus("unknown"),
			expected: "不明",
		},
		{
			name:     "空のステータス",
			status:   ScheduleStatus(""),
			expected: "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// ScheduleEntry関連テスト
// ============================================

func TestScheduleEntry_Structure(t *testing.T) {
	t.Run("ScheduleEntry構造体の完全な初期化", func(t *testing.T) {
		entryID := sharedDomain.NewID()
		scheduleID := sharedDomain.NewID()
		staffID := sharedDomain.NewID()
		shiftTypeID := sharedDomain.NewID()
		now := time.Now()
		targetDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

		entry := &ScheduleEntry{
			ID:          entryID,
			ScheduleID:  scheduleID,
			StaffID:     staffID,
			TargetDate:  targetDate,
			ShiftTypeID: &shiftTypeID,
			IsConfirmed: true,
			Note:        "特記事項あり",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if entry.ID != entryID {
			t.Errorf("ID = %v, want %v", entry.ID, entryID)
		}
		if entry.ScheduleID != scheduleID {
			t.Errorf("ScheduleID = %v, want %v", entry.ScheduleID, scheduleID)
		}
		if entry.StaffID != staffID {
			t.Errorf("StaffID = %v, want %v", entry.StaffID, staffID)
		}
		if !entry.TargetDate.Equal(targetDate) {
			t.Errorf("TargetDate = %v, want %v", entry.TargetDate, targetDate)
		}
		if entry.ShiftTypeID == nil || *entry.ShiftTypeID != shiftTypeID {
			t.Error("ShiftTypeID not set correctly")
		}
		if entry.Note != "特記事項あり" {
			t.Errorf("Note = %v, want %v", entry.Note, "特記事項あり")
		}
		if !entry.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt = %v, want %v", entry.CreatedAt, now)
		}
		if !entry.UpdatedAt.Equal(now) {
			t.Errorf("UpdatedAt = %v, want %v", entry.UpdatedAt, now)
		}
		if !entry.IsConfirmed {
			t.Error("IsConfirmed should be true")
		}
	})

	t.Run("ShiftTypeIDがnilの場合", func(t *testing.T) {
		entry := &ScheduleEntry{
			ShiftTypeID: nil,
		}

		if entry.ShiftTypeID != nil {
			t.Error("ShiftTypeID should be nil")
		}
	})
}

// ============================================
// ActualRecord関連テスト
// ============================================

func TestActualRecord_ActualWorkingMinutes(t *testing.T) {
	baseTime := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		startTime    *time.Time
		endTime      *time.Time
		breakMinutes int
		expected     int
	}{
		{
			name:         "通常勤務8時間_休憩60分",
			startTime:    ptrTime(baseTime.Add(9 * time.Hour)),
			endTime:      ptrTime(baseTime.Add(18 * time.Hour)),
			breakMinutes: 60,
			expected:     480, // 9時間 - 1時間 = 8時間 = 480分
		},
		{
			name:         "短時間勤務4時間_休憩なし",
			startTime:    ptrTime(baseTime.Add(9 * time.Hour)),
			endTime:      ptrTime(baseTime.Add(13 * time.Hour)),
			breakMinutes: 0,
			expected:     240, // 4時間 = 240分
		},
		{
			name:         "開始時刻がnilの場合",
			startTime:    nil,
			endTime:      ptrTime(baseTime.Add(18 * time.Hour)),
			breakMinutes: 60,
			expected:     0,
		},
		{
			name:         "終了時刻がnilの場合",
			startTime:    ptrTime(baseTime.Add(9 * time.Hour)),
			endTime:      nil,
			breakMinutes: 60,
			expected:     0,
		},
		{
			name:         "両方nilの場合",
			startTime:    nil,
			endTime:      nil,
			breakMinutes: 0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &ActualRecord{
				ActualStartTime:    tt.startTime,
				ActualEndTime:      tt.endTime,
				ActualBreakMinutes: tt.breakMinutes,
			}

			result := record.ActualWorkingMinutes()
			if result != tt.expected {
				t.Errorf("ActualWorkingMinutes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestActualRecord_Structure(t *testing.T) {
	t.Run("ActualRecord構造体の初期化", func(t *testing.T) {
		recordID := sharedDomain.NewID()
		entryID := sharedDomain.NewID()
		now := time.Now()
		startTime := time.Date(2025, 3, 15, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(2025, 3, 15, 18, 0, 0, 0, time.UTC)

		record := &ActualRecord{
			ID:                 recordID,
			ScheduleEntryID:    entryID,
			ActualStartTime:    &startTime,
			ActualEndTime:      &endTime,
			ActualBreakMinutes: 60,
			OvertimeMinutes:    30,
			Note:               "残業30分",
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		if record.ID != recordID {
			t.Errorf("ID = %v, want %v", record.ID, recordID)
		}
		if record.ScheduleEntryID != entryID {
			t.Errorf("ScheduleEntryID = %v, want %v", record.ScheduleEntryID, entryID)
		}
		if record.ActualStartTime == nil || !record.ActualStartTime.Equal(startTime) {
			t.Errorf("ActualStartTime = %v, want %v", record.ActualStartTime, startTime)
		}
		if record.ActualEndTime == nil || !record.ActualEndTime.Equal(endTime) {
			t.Errorf("ActualEndTime = %v, want %v", record.ActualEndTime, endTime)
		}
		if record.ActualBreakMinutes != 60 {
			t.Errorf("ActualBreakMinutes = %v, want %v", record.ActualBreakMinutes, 60)
		}
		if record.OvertimeMinutes != 30 {
			t.Errorf("OvertimeMinutes = %v, want %v", record.OvertimeMinutes, 30)
		}
		if record.Note != "残業30分" {
			t.Errorf("Note = %v, want %v", record.Note, "残業30分")
		}
		if !record.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt = %v, want %v", record.CreatedAt, now)
		}
		if !record.UpdatedAt.Equal(now) {
			t.Errorf("UpdatedAt = %v, want %v", record.UpdatedAt, now)
		}
	})
}

// ============================================
// OptimizeInput/OptimizeResult/Constraint関連テスト
// ============================================

func TestOptimizeInput_Structure(t *testing.T) {
	t.Run("OptimizeInput構造体の初期化", func(t *testing.T) {
		scheduleID := sharedDomain.NewID()
		startDate := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
		dateRange := sharedDomain.NewDateRange(startDate, endDate)

		input := OptimizeInput{
			ScheduleID: scheduleID,
			DateRange:  dateRange,
			Constraints: []Constraint{
				{Type: "min_staff", Priority: 100},
			},
			Options: OptimizeOptions{
				MaxIterations:      1000,
				TimeoutSeconds:     30,
				PrioritizeRequests: true,
			},
		}

		if input.ScheduleID != scheduleID {
			t.Errorf("ScheduleID = %v, want %v", input.ScheduleID, scheduleID)
		}
		if input.DateRange.Start != startDate {
			t.Errorf("DateRange.Start = %v, want %v", input.DateRange.Start, startDate)
		}
		if input.DateRange.End != endDate {
			t.Errorf("DateRange.End = %v, want %v", input.DateRange.End, endDate)
		}
		if len(input.Constraints) != 1 {
			t.Errorf("Constraints length = %d, want 1", len(input.Constraints))
		}
		if !input.Options.PrioritizeRequests {
			t.Error("PrioritizeRequests should be true")
		}
	})
}

func TestOptimizeResult_Structure(t *testing.T) {
	t.Run("成功したOptimizeResult", func(t *testing.T) {
		result := OptimizeResult{
			Success:    true,
			Entries:    []ScheduleEntry{},
			Score:      0.95,
			Violations: []ConstraintViolation{},
			Duration:   5 * time.Second,
		}

		if !result.Success {
			t.Error("Success should be true")
		}
		if len(result.Entries) != 0 {
			t.Errorf("Entries length = %d, want 0", len(result.Entries))
		}
		if result.Score != 0.95 {
			t.Errorf("Score = %v, want %v", result.Score, 0.95)
		}
		if len(result.Violations) != 0 {
			t.Errorf("Violations length = %d, want 0", len(result.Violations))
		}
		if result.Duration != 5*time.Second {
			t.Errorf("Duration = %v, want %v", result.Duration, 5*time.Second)
		}
	})

	t.Run("違反を含むOptimizeResult", func(t *testing.T) {
		staffID := sharedDomain.NewID()
		violationDate := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)

		result := OptimizeResult{
			Success: true,
			Score:   0.8,
			Violations: []ConstraintViolation{
				{
					ConstraintType: "consecutive",
					Message:        "連続勤務が6日を超えています",
					StaffID:        &staffID,
					Date:           &violationDate,
					Severity:       "warning",
				},
			},
		}

		if !result.Success {
			t.Error("Success should be true")
		}
		if result.Score != 0.8 {
			t.Errorf("Score = %v, want %v", result.Score, 0.8)
		}
		if len(result.Violations) != 1 {
			t.Errorf("Violations length = %d, want 1", len(result.Violations))
		}
		if result.Violations[0].Severity != "warning" {
			t.Errorf("Severity = %v, want %v", result.Violations[0].Severity, "warning")
		}
	})
}

func TestConstraint_Structure(t *testing.T) {
	t.Run("Constraint構造体の初期化", func(t *testing.T) {
		constraint := Constraint{
			Type:     "min_staff",
			Priority: 100,
			Config:   `{"shift_type_id": "xxx", "min_count": 3}`,
		}

		if constraint.Type != "min_staff" {
			t.Errorf("Type = %v, want %v", constraint.Type, "min_staff")
		}
		if constraint.Priority != 100 {
			t.Errorf("Priority = %v, want %v", constraint.Priority, 100)
		}
		if constraint.Config != `{"shift_type_id": "xxx", "min_count": 3}` {
			t.Errorf("Config = %v, want %v", constraint.Config, `{"shift_type_id": "xxx", "min_count": 3}`)
		}
	})
}

func TestConstraintViolation_Structure(t *testing.T) {
	t.Run("ConstraintViolation構造体の初期化", func(t *testing.T) {
		staffID := sharedDomain.NewID()
		date := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

		violation := ConstraintViolation{
			ConstraintType: "night_limit",
			Message:        "月間夜勤回数が上限を超えています",
			StaffID:        &staffID,
			Date:           &date,
			Severity:       "error",
		}

		if violation.ConstraintType != "night_limit" {
			t.Errorf("ConstraintType = %v, want %v", violation.ConstraintType, "night_limit")
		}
		if violation.Message != "月間夜勤回数が上限を超えています" {
			t.Errorf("Message = %v, want %v", violation.Message, "月間夜勤回数が上限を超えています")
		}
		if violation.StaffID == nil || *violation.StaffID != staffID {
			t.Errorf("StaffID = %v, want %v", violation.StaffID, staffID)
		}
		if violation.Date == nil || !violation.Date.Equal(date) {
			t.Errorf("Date = %v, want %v", violation.Date, date)
		}
		if violation.Severity != "error" {
			t.Errorf("Severity = %v, want %v", violation.Severity, "error")
		}
	})
}

// ============================================
// 定数テスト
// ============================================

func TestScheduleStatusConstants(t *testing.T) {
	t.Run("ステータス定数の値", func(t *testing.T) {
		if StatusDraft != "draft" {
			t.Errorf("StatusDraft = %v, want %v", StatusDraft, "draft")
		}
		if StatusInProgress != "in_progress" {
			t.Errorf("StatusInProgress = %v, want %v", StatusInProgress, "in_progress")
		}
		if StatusCompleted != "completed" {
			t.Errorf("StatusCompleted = %v, want %v", StatusCompleted, "completed")
		}
		if StatusPublished != "published" {
			t.Errorf("StatusPublished = %v, want %v", StatusPublished, "published")
		}
	})

	t.Run("ステータス定数の一意性", func(t *testing.T) {
		statuses := []ScheduleStatus{
			StatusDraft,
			StatusInProgress,
			StatusCompleted,
			StatusPublished,
		}

		seen := make(map[ScheduleStatus]bool)
		for _, s := range statuses {
			if seen[s] {
				t.Errorf("Duplicate status constant: %v", s)
			}
			seen[s] = true
		}
	})
}

// ============================================
// ヘルパー関数
// ============================================

func ptrTime(t time.Time) *time.Time {
	return &t
}
