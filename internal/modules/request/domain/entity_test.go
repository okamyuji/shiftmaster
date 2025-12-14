// Package domain 勤務希望ドメイン層テスト
package domain

import (
	"testing"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
)

// ============================================
// RequestPeriod関連テスト
// ============================================

func TestRequestPeriod_IsActive(t *testing.T) {
	// 日付ベースで比較（タイムゾーン非依存）
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		period   *RequestPeriod
		expected bool
	}{
		{
			name: "アクティブ_受付中かつ期間内",
			period: &RequestPeriod{
				IsOpen:    true,
				StartDate: today.AddDate(0, 0, -1), // 昨日
				EndDate:   today.AddDate(0, 0, 1),  // 明日
			},
			expected: true,
		},
		{
			name: "非アクティブ_IsOpenがfalse",
			period: &RequestPeriod{
				IsOpen:    false,
				StartDate: today.AddDate(0, 0, -1),
				EndDate:   today.AddDate(0, 0, 1),
			},
			expected: false,
		},
		{
			name: "非アクティブ_開始日前",
			period: &RequestPeriod{
				IsOpen:    true,
				StartDate: today.AddDate(0, 0, 1), // 明日
				EndDate:   today.AddDate(0, 0, 2), // 明後日
			},
			expected: false,
		},
		{
			name: "非アクティブ_終了日後",
			period: &RequestPeriod{
				IsOpen:    true,
				StartDate: today.AddDate(0, 0, -3), // 3日前
				EndDate:   today.AddDate(0, 0, -1), // 昨日
			},
			expected: false,
		},
		{
			name: "境界値_開始日当日",
			period: &RequestPeriod{
				IsOpen:    true,
				StartDate: today,                  // 今日
				EndDate:   today.AddDate(0, 0, 1), // 明日
			},
			expected: true,
		},
		{
			name: "境界値_終了日当日",
			period: &RequestPeriod{
				IsOpen:    true,
				StartDate: today.AddDate(0, 0, -1), // 昨日
				EndDate:   today,                   // 今日
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.period.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRequestPeriod_TargetPeriodLabel(t *testing.T) {
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
			period := &RequestPeriod{
				TargetYear:  tt.targetYear,
				TargetMonth: tt.targetMonth,
			}

			result := period.TargetPeriodLabel()
			if result != tt.expected {
				t.Errorf("TargetPeriodLabel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRequestPeriod_Structure(t *testing.T) {
	t.Run("RequestPeriod構造体の完全な初期化", func(t *testing.T) {
		periodID := sharedDomain.NewID()
		orgID := sharedDomain.NewID()
		now := time.Now()
		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		period := &RequestPeriod{
			ID:                  periodID,
			OrganizationID:      orgID,
			TargetYear:          2025,
			TargetMonth:         2,
			StartDate:           startDate,
			EndDate:             endDate,
			MaxRequestsPerStaff: 10,
			MaxRequestsPerDay:   3,
			IsOpen:              true,
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		if period.ID != periodID {
			t.Errorf("ID = %v, want %v", period.ID, periodID)
		}
		if period.OrganizationID != orgID {
			t.Errorf("OrganizationID = %v, want %v", period.OrganizationID, orgID)
		}
		if period.TargetYear != 2025 {
			t.Errorf("TargetYear = %v, want %v", period.TargetYear, 2025)
		}
		if period.TargetMonth != 2 {
			t.Errorf("TargetMonth = %v, want %v", period.TargetMonth, 2)
		}
		if period.StartDate != startDate {
			t.Errorf("StartDate = %v, want %v", period.StartDate, startDate)
		}
		if period.EndDate != endDate {
			t.Errorf("EndDate = %v, want %v", period.EndDate, endDate)
		}
		if period.MaxRequestsPerStaff != 10 {
			t.Errorf("MaxRequestsPerStaff = %v, want %v", period.MaxRequestsPerStaff, 10)
		}
		if period.MaxRequestsPerDay != 3 {
			t.Errorf("MaxRequestsPerDay = %v, want %v", period.MaxRequestsPerDay, 3)
		}
		if !period.IsOpen {
			t.Error("IsOpen should be true")
		}
		if period.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", period.CreatedAt, now)
		}
		if period.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", period.UpdatedAt, now)
		}
	})
}

// ============================================
// ShiftRequest関連テスト
// ============================================

func TestShiftRequest_Structure(t *testing.T) {
	t.Run("ShiftRequest構造体の完全な初期化", func(t *testing.T) {
		requestID := sharedDomain.NewID()
		periodID := sharedDomain.NewID()
		staffID := sharedDomain.NewID()
		shiftTypeID := sharedDomain.NewID()
		now := time.Now()
		targetDate := time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC)

		request := &ShiftRequest{
			ID:          requestID,
			PeriodID:    periodID,
			StaffID:     staffID,
			TargetDate:  targetDate,
			ShiftTypeID: &shiftTypeID,
			RequestType: RequestTypePreferred,
			Priority:    PriorityRequired,
			Comment:     "この日は午前勤務希望",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if request.ID != requestID {
			t.Errorf("ID = %v, want %v", request.ID, requestID)
		}
		if request.PeriodID != periodID {
			t.Errorf("PeriodID = %v, want %v", request.PeriodID, periodID)
		}
		if request.StaffID != staffID {
			t.Errorf("StaffID = %v, want %v", request.StaffID, staffID)
		}
		if request.TargetDate != targetDate {
			t.Errorf("TargetDate = %v, want %v", request.TargetDate, targetDate)
		}
		if request.ShiftTypeID == nil || *request.ShiftTypeID != shiftTypeID {
			t.Error("ShiftTypeID not set correctly")
		}
		if request.RequestType != RequestTypePreferred {
			t.Errorf("RequestType = %v, want %v", request.RequestType, RequestTypePreferred)
		}
		if request.Priority != PriorityRequired {
			t.Errorf("Priority = %v, want %v", request.Priority, PriorityRequired)
		}
		if request.Comment != "この日は午前勤務希望" {
			t.Errorf("Comment = %v, want %v", request.Comment, "この日は午前勤務希望")
		}
		if request.CreatedAt != now {
			t.Errorf("CreatedAt = %v, want %v", request.CreatedAt, now)
		}
		if request.UpdatedAt != now {
			t.Errorf("UpdatedAt = %v, want %v", request.UpdatedAt, now)
		}
		if request.ID == sharedDomain.ID(uuid.Nil) {
			t.Error("ID should not be nil")
		}
		if request.PeriodID == sharedDomain.ID(uuid.Nil) {
			t.Error("PeriodID should not be nil")
		}
		if request.StaffID == sharedDomain.ID(uuid.Nil) {
			t.Error("StaffID should not be nil")
		}
		if request.TargetDate.IsZero() {
			t.Error("TargetDate should not be zero time")
		}
		if request.ShiftTypeID == nil {
			t.Error("ShiftTypeID should not be nil")
		}
	})

	t.Run("ShiftTypeIDがnilの場合", func(t *testing.T) {
		request := &ShiftRequest{
			RequestType: RequestTypeAvoided,
			ShiftTypeID: nil,
		}
		if request.RequestType != RequestTypeAvoided {
			t.Errorf("RequestType = %v, want %v", request.RequestType, RequestTypeAvoided)
		}
		if request.ShiftTypeID != nil {
			t.Error("ShiftTypeID should be nil")
		}
	})
}

// ============================================
// RequestType関連テスト
// ============================================

func TestRequestType_String(t *testing.T) {
	tests := []struct {
		name        string
		requestType RequestType
		expected    string
	}{
		{
			name:        "preferred",
			requestType: RequestTypePreferred,
			expected:    "preferred",
		},
		{
			name:        "avoided",
			requestType: RequestTypeAvoided,
			expected:    "avoided",
		},
		{
			name:        "fixed",
			requestType: RequestTypeFixed,
			expected:    "fixed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.requestType.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRequestType_Label(t *testing.T) {
	tests := []struct {
		name        string
		requestType RequestType
		expected    string
	}{
		{
			name:        "preferred",
			requestType: RequestTypePreferred,
			expected:    "希望",
		},
		{
			name:        "avoided",
			requestType: RequestTypeAvoided,
			expected:    "回避",
		},
		{
			name:        "fixed",
			requestType: RequestTypeFixed,
			expected:    "固定",
		},
		{
			name:        "不明なタイプ",
			requestType: RequestType("unknown"),
			expected:    "不明",
		},
		{
			name:        "空のタイプ",
			requestType: RequestType(""),
			expected:    "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.requestType.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// RequestPriority関連テスト
// ============================================

func TestRequestPriority_String(t *testing.T) {
	tests := []struct {
		name     string
		priority RequestPriority
		expected string
	}{
		{
			name:     "required",
			priority: PriorityRequired,
			expected: "required",
		},
		{
			name:     "optional",
			priority: PriorityOptional,
			expected: "optional",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priority.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRequestPriority_Label(t *testing.T) {
	tests := []struct {
		name     string
		priority RequestPriority
		expected string
	}{
		{
			name:     "required",
			priority: PriorityRequired,
			expected: "必須",
		},
		{
			name:     "optional",
			priority: PriorityOptional,
			expected: "できれば",
		},
		{
			name:     "不明な優先度",
			priority: RequestPriority("unknown"),
			expected: "不明",
		},
		{
			name:     "空の優先度",
			priority: RequestPriority(""),
			expected: "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priority.Label()
			if result != tt.expected {
				t.Errorf("Label() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// 定数テスト
// ============================================

func TestRequestTypeConstants(t *testing.T) {
	t.Run("希望種別定数の値", func(t *testing.T) {
		if RequestTypePreferred != "preferred" {
			t.Errorf("RequestTypePreferred = %v, want %v", RequestTypePreferred, "preferred")
		}
		if RequestTypeAvoided != "avoided" {
			t.Errorf("RequestTypeAvoided = %v, want %v", RequestTypeAvoided, "avoided")
		}
		if RequestTypeFixed != "fixed" {
			t.Errorf("RequestTypeFixed = %v, want %v", RequestTypeFixed, "fixed")
		}
	})

	t.Run("希望種別定数の一意性", func(t *testing.T) {
		types := []RequestType{
			RequestTypePreferred,
			RequestTypeAvoided,
			RequestTypeFixed,
		}

		seen := make(map[RequestType]bool)
		for _, rt := range types {
			if seen[rt] {
				t.Errorf("Duplicate request type constant: %v", rt)
			}
			seen[rt] = true
		}
	})
}

func TestRequestPriorityConstants(t *testing.T) {
	t.Run("優先度定数の値", func(t *testing.T) {
		if PriorityRequired != "required" {
			t.Errorf("PriorityRequired = %v, want %v", PriorityRequired, "required")
		}
		if PriorityOptional != "optional" {
			t.Errorf("PriorityOptional = %v, want %v", PriorityOptional, "optional")
		}
	})

	t.Run("優先度定数の一意性", func(t *testing.T) {
		priorities := []RequestPriority{
			PriorityRequired,
			PriorityOptional,
		}

		seen := make(map[RequestPriority]bool)
		for _, p := range priorities {
			if seen[p] {
				t.Errorf("Duplicate priority constant: %v", p)
			}
			seen[p] = true
		}
	})
}
