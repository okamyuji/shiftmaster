// Package domain シフトドメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ============================================
// ShiftType関連テスト
// ============================================

func TestShiftType_WorkingMinutes(t *testing.T) {
	tests := []struct {
		name         string
		startTime    time.Time
		endTime      time.Time
		breakMinutes int
		expected     int
	}{
		{
			name:         "通常日勤_8時間勤務1時間休憩",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
			breakMinutes: 60,
			expected:     480, // 9時間 - 1時間 = 8時間 = 480分
		},
		{
			name:         "休憩なし",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			breakMinutes: 0,
			expected:     180, // 3時間 = 180分
		},
		{
			name:         "夜勤_日跨ぎ",
			startTime:    time.Date(2025, 1, 1, 17, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC),
			breakMinutes: 120,
			expected:     840, // 16時間 - 2時間 = 14時間 = 840分
		},
		{
			name:         "同じ時刻",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			breakMinutes: 0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime:    tt.startTime,
				EndTime:      tt.endTime,
				BreakMinutes: tt.breakMinutes,
			}

			result := shift.WorkingMinutes()
			if result != tt.expected {
				t.Errorf("WorkingMinutes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_TotalMinutes(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		expected  int
	}{
		{
			name:      "通常勤務_9時間",
			startTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
			expected:  540, // 9時間 = 540分
		},
		{
			name:      "短時間勤務_4時間",
			startTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC),
			expected:  240, // 4時間 = 240分
		},
		{
			name:      "日跨ぎ_夜勤16時間",
			startTime: time.Date(2025, 1, 1, 17, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC),
			expected:  960, // 16時間 = 960分
		},
		{
			name:      "日跨ぎ_23時から7時",
			startTime: time.Date(2025, 1, 1, 23, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 2, 7, 0, 0, 0, time.UTC),
			expected:  480, // 8時間 = 480分
		},
		{
			name:      "同じ時刻",
			startTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  0,
		},
		{
			name:      "分単位を含む",
			startTime: time.Date(2025, 1, 1, 8, 30, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 1, 17, 15, 0, 0, time.UTC),
			expected:  525, // 8時間45分 = 525分
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}

			result := shift.TotalMinutes()
			if result != tt.expected {
				t.Errorf("TotalMinutes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_WorkingHours(t *testing.T) {
	tests := []struct {
		name         string
		startTime    time.Time
		endTime      time.Time
		breakMinutes int
		expected     float64
	}{
		{
			name:         "8時間勤務",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
			breakMinutes: 60,
			expected:     8.0,
		},
		{
			name:         "7.5時間勤務",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 17, 30, 0, 0, time.UTC),
			breakMinutes: 60,
			expected:     7.5,
		},
		{
			name:         "ゼロ時間",
			startTime:    time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:      time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			breakMinutes: 0,
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime:    tt.startTime,
				EndTime:      tt.endTime,
				BreakMinutes: tt.breakMinutes,
			}

			result := shift.WorkingHours()
			if result != tt.expected {
				t.Errorf("WorkingHours() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_StartTimeString(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		expected  string
	}{
		{
			name:      "9時0分",
			startTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			expected:  "09:00",
		},
		{
			name:      "17時30分",
			startTime: time.Date(2025, 1, 1, 17, 30, 0, 0, time.UTC),
			expected:  "17:30",
		},
		{
			name:      "0時0分",
			startTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  "00:00",
		},
		{
			name:      "23時59分",
			startTime: time.Date(2025, 1, 1, 23, 59, 0, 0, time.UTC),
			expected:  "23:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{StartTime: tt.startTime}

			result := shift.StartTimeString()
			if result != tt.expected {
				t.Errorf("StartTimeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_EndTimeString(t *testing.T) {
	tests := []struct {
		name     string
		endTime  time.Time
		expected string
	}{
		{
			name:     "18時0分",
			endTime:  time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
			expected: "18:00",
		},
		{
			name:     "8時45分",
			endTime:  time.Date(2025, 1, 1, 8, 45, 0, 0, time.UTC),
			expected: "08:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{EndTime: tt.endTime}

			result := shift.EndTimeString()
			if result != tt.expected {
				t.Errorf("EndTimeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// 申し送り時間（HandoverMinutes）関連テスト
// ============================================

func TestShiftType_EffectiveWorkingMinutes(t *testing.T) {
	tests := []struct {
		name            string
		startTime       time.Time
		endTime         time.Time
		breakMinutes    int
		handoverMinutes int
		expected        int
	}{
		{
			name:            "申し送りなし_通常勤務",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:         time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC),
			breakMinutes:    60,
			handoverMinutes: 0,
			expected:        480, // 9時間 - 1時間 = 8時間
		},
		{
			name:            "申し送り15分_病院標準",
			startTime:       time.Date(2025, 1, 1, 8, 30, 0, 0, time.UTC),
			endTime:         time.Date(2025, 1, 1, 17, 0, 0, 0, time.UTC),
			breakMinutes:    60,
			handoverMinutes: 15,
			expected:        435, // 8.5時間 - 1時間 - 15分 = 7時間15分 = 435分
		},
		{
			name:            "申し送り30分_看護師夜勤",
			startTime:       time.Date(2025, 1, 1, 16, 30, 0, 0, time.UTC),
			endTime:         time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC),
			breakMinutes:    120,
			handoverMinutes: 30,
			expected:        840, // 16.5時間 - 2時間 - 30分 = 14時間 = 840分
		},
		{
			name:            "申し送りが実働時間より長い場合",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:         time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			breakMinutes:    30,
			handoverMinutes: 60,
			expected:        0, // 1時間 - 30分 = 30分 < 60分申し送り = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime:       tt.startTime,
				EndTime:         tt.endTime,
				BreakMinutes:    tt.breakMinutes,
				HandoverMinutes: tt.handoverMinutes,
			}

			result := shift.EffectiveWorkingMinutes()
			if result != tt.expected {
				t.Errorf("EffectiveWorkingMinutes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_HasHandover(t *testing.T) {
	tests := []struct {
		name            string
		handoverMinutes int
		expected        bool
	}{
		{
			name:            "申し送りあり_15分",
			handoverMinutes: 15,
			expected:        true,
		},
		{
			name:            "申し送りあり_30分",
			handoverMinutes: 30,
			expected:        true,
		},
		{
			name:            "申し送りなし",
			handoverMinutes: 0,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{HandoverMinutes: tt.handoverMinutes}

			result := shift.HasHandover()
			if result != tt.expected {
				t.Errorf("HasHandover() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_ActualStartTime(t *testing.T) {
	tests := []struct {
		name            string
		startTime       time.Time
		handoverMinutes int
		expected        time.Time
	}{
		{
			name:            "申し送りなし",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			handoverMinutes: 0,
			expected:        time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			name:            "申し送り15分",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			handoverMinutes: 15,
			expected:        time.Date(2025, 1, 1, 8, 45, 0, 0, time.UTC),
		},
		{
			name:            "申し送り30分",
			startTime:       time.Date(2025, 1, 1, 17, 0, 0, 0, time.UTC),
			handoverMinutes: 30,
			expected:        time.Date(2025, 1, 1, 16, 30, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime:       tt.startTime,
				HandoverMinutes: tt.handoverMinutes,
			}

			result := shift.ActualStartTime()
			if !result.Equal(tt.expected) {
				t.Errorf("ActualStartTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftType_ActualStartTimeString(t *testing.T) {
	tests := []struct {
		name            string
		startTime       time.Time
		handoverMinutes int
		expected        string
	}{
		{
			name:            "申し送りなし_9時開始",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			handoverMinutes: 0,
			expected:        "09:00",
		},
		{
			name:            "申し送り15分_9時開始",
			startTime:       time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			handoverMinutes: 15,
			expected:        "08:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift := &ShiftType{
				StartTime:       tt.startTime,
				HandoverMinutes: tt.handoverMinutes,
			}

			result := shift.ActualStartTimeString()
			if result != tt.expected {
				t.Errorf("ActualStartTimeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// RotationType（ローテーション種別）関連テスト
// ============================================

func TestRotationType_String(t *testing.T) {
	tests := []struct {
		name         string
		rotationType RotationType
		expected     string
	}{
		{
			name:         "二交代制",
			rotationType: RotationTypeTwoShift,
			expected:     "two_shift",
		},
		{
			name:         "三交代制",
			rotationType: RotationTypeThreeShift,
			expected:     "three_shift",
		},
		{
			name:         "四直三交代制",
			rotationType: RotationTypeFourShiftThreeTeam,
			expected:     "four_shift_three_team",
		},
		{
			name:         "カスタム",
			rotationType: RotationTypeCustom,
			expected:     "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rotationType.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRotationType_Label(t *testing.T) {
	tests := []struct {
		name         string
		rotationType RotationType
		expected     string
	}{
		{
			name:         "二交代制",
			rotationType: RotationTypeTwoShift,
			expected:     "二交代制",
		},
		{
			name:         "三交代制",
			rotationType: RotationTypeThreeShift,
			expected:     "三交代制",
		},
		{
			name:         "四直三交代制",
			rotationType: RotationTypeFourShiftThreeTeam,
			expected:     "四直三交代制",
		},
		{
			name:         "カスタム",
			rotationType: RotationTypeCustom,
			expected:     "カスタム",
		},
		{
			name:         "不明なタイプ",
			rotationType: RotationType("unknown"),
			expected:     "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rotationType.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRotationType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		rotationType RotationType
		expected     bool
	}{
		{
			name:         "二交代制_有効",
			rotationType: RotationTypeTwoShift,
			expected:     true,
		},
		{
			name:         "三交代制_有効",
			rotationType: RotationTypeThreeShift,
			expected:     true,
		},
		{
			name:         "四直三交代制_有効",
			rotationType: RotationTypeFourShiftThreeTeam,
			expected:     true,
		},
		{
			name:         "カスタム_有効",
			rotationType: RotationTypeCustom,
			expected:     true,
		},
		{
			name:         "不明_無効",
			rotationType: RotationType("unknown"),
			expected:     false,
		},
		{
			name:         "空文字_無効",
			rotationType: RotationType(""),
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rotationType.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// ShiftPattern（直）関連テスト
// ============================================

func TestShiftPattern_Structure(t *testing.T) {
	t.Run("ShiftPattern構造体の完全な初期化", func(t *testing.T) {
		patternID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()

		pattern := &ShiftPattern{
			ID:             patternID,
			OrganizationID: orgID,
			Name:           "1直",
			Code:           "D1",
			Description:    "第1シフトグループ（日勤担当）",
			RotationType:   RotationTypeThreeShift,
			Color:          "#4CAF50",
			SortOrder:      1,
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if pattern.ID != patternID {
			t.Errorf("ID = %v, want %v", pattern.ID, patternID)
		}
		if pattern.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", pattern.OrganizationID, orgID)
		}
		if pattern.Name != "1直" {
			t.Errorf("Name = %v, want %v", pattern.Name, "1直")
		}
		if pattern.Code != "D1" {
			t.Errorf("Code = %v, want %v", pattern.Code, "D1")
		}
		if pattern.Description != "第1シフトグループ（日勤担当）" {
			t.Errorf("Description = %v, want %v", pattern.Description, "第1シフトグループ（日勤担当）")
		}
		if pattern.RotationType != RotationTypeThreeShift {
			t.Errorf("RotationType = %v, want %v", pattern.RotationType, RotationTypeThreeShift)
		}
		if pattern.Color != "#4CAF50" {
			t.Errorf("Color = %v, want %v", pattern.Color, "#4CAF50")
		}
		if pattern.SortOrder != 1 {
			t.Errorf("SortOrder = %v, want %v", pattern.SortOrder, 1)
		}
		if !pattern.IsActive {
			t.Error("IsActive should be true")
		}
		if pattern.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", pattern.CreatedAt, now)
		}
		if pattern.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", pattern.UpdatedAt, now)
		}
	})

	t.Run("二交代制パターン", func(t *testing.T) {
		pattern := &ShiftPattern{
			Name:         "日勤チーム",
			RotationType: RotationTypeTwoShift,
		}
		if pattern.Name != "日勤チーム" {
			t.Errorf("Name = %v, want %v", pattern.Name, "日勤チーム")
		}
		if pattern.RotationType != RotationTypeTwoShift {
			t.Errorf("RotationType = %v, want %v", pattern.RotationType, RotationTypeTwoShift)
		}
	})

	t.Run("四直三交代制パターン（工場向け）", func(t *testing.T) {
		pattern := &ShiftPattern{
			Name:         "A班",
			Code:         "A",
			RotationType: RotationTypeFourShiftThreeTeam,
			Description:  "工場24時間稼働用4班体制",
		}
		if pattern.Name != "A班" {
			t.Errorf("Name = %v, want %v", pattern.Name, "A班")
		}
		if pattern.Code != "A" {
			t.Errorf("Code = %v, want %v", pattern.Code, "A")
		}
		if pattern.Description != "工場24時間稼働用4班体制" {
			t.Errorf("Description = %v, want %v", pattern.Description, "工場24時間稼働用4班体制")
		}
		if pattern.RotationType != RotationTypeFourShiftThreeTeam {
			t.Errorf("RotationType = %v, want %v", pattern.RotationType, RotationTypeFourShiftThreeTeam)
		}
	})
}

func TestShiftType_Structure(t *testing.T) {
	t.Run("ShiftType構造体の完全な初期化", func(t *testing.T) {
		shiftID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()
		startTime := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)

		shift := &ShiftType{
			ID:             shiftID,
			OrganizationID: orgID,
			Name:           "日勤",
			Code:           "D",
			Color:          "#4CAF50",
			StartTime:      startTime,
			EndTime:        endTime,
			BreakMinutes:   60,
			IsNightShift:   false,
			IsHoliday:      false,
			SortOrder:      1,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if shift.ID != shiftID {
			t.Errorf("ID = %v, want %v", shift.ID, shiftID)
		}
		if shift.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", shift.OrganizationID, orgID)
		}
		if shift.Name != "日勤" {
			t.Errorf("Name = %v, want %v", shift.Name, "日勤")
		}
		if shift.Code != "D" {
			t.Errorf("Code = %v, want %v", shift.Code, "D")
		}
		if shift.Color != "#4CAF50" {
			t.Errorf("Color = %v, want %v", shift.Color, "#4CAF50")
		}
		if shift.StartTime != startTime {
			t.Errorf("StartTime = %v, want %v", shift.StartTime, startTime)
		}
		if shift.EndTime != endTime {
			t.Errorf("EndTime = %v, want %v", shift.EndTime, endTime)
		}
		if shift.BreakMinutes != 60 {
			t.Errorf("BreakMinutes = %v, want %v", shift.BreakMinutes, 60)
		}
		if shift.SortOrder != 1 {
			t.Errorf("SortOrder = %v, want %v", shift.SortOrder, 1)
		}
		if shift.IsNightShift {
			t.Errorf("IsNightShift = %v, want %v", shift.IsNightShift, false)
		}
		if shift.IsHoliday {
			t.Errorf("IsHoliday = %v, want %v", shift.IsHoliday, false)
		}
		if shift.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", shift.CreatedAt, now)
		}
		if shift.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", shift.UpdatedAt, now)
		}
		if shift.IsNightShift {
			t.Error("IsNightShift should be false")
		}
	})

	t.Run("夜勤シフト", func(t *testing.T) {
		shift := &ShiftType{
			Name:         "夜勤",
			IsNightShift: true,
		}
		if shift.Name != "夜勤" {
			t.Errorf("Name = %v, want %v", shift.Name, "夜勤")
		}
		if !shift.IsNightShift {
			t.Error("IsNightShift should be true for night shift")
		}
	})

	t.Run("休日シフト", func(t *testing.T) {
		shift := &ShiftType{
			Name:      "公休",
			IsHoliday: true,
		}
		if shift.Name != "公休" {
			t.Errorf("Name = %v, want %v", shift.Name, "公休")
		}
		if !shift.IsHoliday {
			t.Error("IsHoliday should be true for holiday shift")
		}
	})
}

// ============================================
// ShiftRule関連テスト
// ============================================

func TestShiftRule_Structure(t *testing.T) {
	t.Run("ShiftRule構造体の初期化", func(t *testing.T) {
		ruleID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()

		rule := &ShiftRule{
			ID:             ruleID,
			OrganizationID: orgID,
			Name:           "最小配置ルール",
			Description:    "日勤は最低3名必要",
			RuleType:       RuleTypeMinStaff,
			Priority:       100,
			IsActive:       true,
			Config:         `{"min_staff": 3}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if rule.ID != ruleID {
			t.Errorf("ID = %v, want %v", rule.ID, ruleID)
		}
		if rule.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", rule.OrganizationID, orgID)
		}
		if rule.Name != "最小配置ルール" {
			t.Errorf("Name = %v, want %v", rule.Name, "最小配置ルール")
		}
		if rule.Description != "日勤は最低3名必要" {
			t.Errorf("Description = %v, want %v", rule.Description, "日勤は最低3名必要")
		}
		if rule.Priority != 100 {
			t.Errorf("Priority = %v, want %v", rule.Priority, 100)
		}
		if rule.RuleType != RuleTypeMinStaff {
			t.Errorf("RuleType = %v, want %v", rule.RuleType, RuleTypeMinStaff)
		}
		if !rule.IsActive {
			t.Error("IsActive should be true")
		}
		if rule.Config != `{"min_staff": 3}` {
			t.Errorf("Config = %v, want %v", rule.Config, `{"min_staff": 3}`)
		}
		if rule.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", rule.CreatedAt, now)
		}
		if rule.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", rule.UpdatedAt, now)
		}
	})
}

// ============================================
// ShiftRuleType関連テスト
// ============================================

func TestShiftRuleType_String(t *testing.T) {
	tests := []struct {
		name     string
		ruleType ShiftRuleType
		expected string
	}{
		{
			name:     "min_staff",
			ruleType: RuleTypeMinStaff,
			expected: "min_staff",
		},
		{
			name:     "max_staff",
			ruleType: RuleTypeMaxStaff,
			expected: "max_staff",
		},
		{
			name:     "consecutive",
			ruleType: RuleTypeConsecutive,
			expected: "consecutive",
		},
		{
			name:     "interval",
			ruleType: RuleTypeInterval,
			expected: "interval",
		},
		{
			name:     "skill_required",
			ruleType: RuleTypeSkillRequired,
			expected: "skill_required",
		},
		{
			name:     "night_limit",
			ruleType: RuleTypeNightLimit,
			expected: "night_limit",
		},
		{
			name:     "weekly_hours",
			ruleType: RuleTypeWeeklyHours,
			expected: "weekly_hours",
		},
		{
			name:     "monthly_hours",
			ruleType: RuleTypeMonthlyHours,
			expected: "monthly_hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ruleType.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShiftRuleType_Label(t *testing.T) {
	tests := []struct {
		name     string
		ruleType ShiftRuleType
		expected string
	}{
		{
			name:     "min_staff",
			ruleType: RuleTypeMinStaff,
			expected: "最小配置人数",
		},
		{
			name:     "max_staff",
			ruleType: RuleTypeMaxStaff,
			expected: "最大配置人数",
		},
		{
			name:     "consecutive",
			ruleType: RuleTypeConsecutive,
			expected: "連続勤務制限",
		},
		{
			name:     "interval",
			ruleType: RuleTypeInterval,
			expected: "シフト間隔",
		},
		{
			name:     "skill_required",
			ruleType: RuleTypeSkillRequired,
			expected: "必須スキル",
		},
		{
			name:     "night_limit",
			ruleType: RuleTypeNightLimit,
			expected: "夜勤回数制限",
		},
		{
			name:     "weekly_hours",
			ruleType: RuleTypeWeeklyHours,
			expected: "週間労働時間",
		},
		{
			name:     "monthly_hours",
			ruleType: RuleTypeMonthlyHours,
			expected: "月間労働時間",
		},
		{
			name:     "不明なタイプ",
			ruleType: ShiftRuleType("unknown"),
			expected: "不明",
		},
		{
			name:     "空のタイプ",
			ruleType: ShiftRuleType(""),
			expected: "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ruleType.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// 定数テスト
// ============================================

func TestRuleTypeConstants(t *testing.T) {
	t.Run("ルール種別定数の一意性", func(t *testing.T) {
		ruleTypes := []ShiftRuleType{
			RuleTypeMinStaff,
			RuleTypeMaxStaff,
			RuleTypeConsecutive,
			RuleTypeInterval,
			RuleTypeSkillRequired,
			RuleTypeNightLimit,
			RuleTypeWeeklyHours,
			RuleTypeMonthlyHours,
		}

		seen := make(map[ShiftRuleType]bool)
		for _, rt := range ruleTypes {
			if seen[rt] {
				t.Errorf("Duplicate rule type constant: %v", rt)
			}
			seen[rt] = true
		}
	})
}
