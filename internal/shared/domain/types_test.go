// Package domain 共有ドメイン型テスト
package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// ============================================
// ID関連テスト
// ============================================

func TestNewID(t *testing.T) {
	t.Run("新規ID生成_正常系", func(t *testing.T) {
		id := NewID()

		// UUIDが生成されていることを確認
		if id == uuid.Nil {
			t.Error("NewID() returned nil UUID")
		}
	})

	t.Run("新規ID生成_一意性確認", func(t *testing.T) {
		id1 := NewID()
		id2 := NewID()

		if id1 == id2 {
			t.Error("NewID() generated duplicate IDs")
		}
	})
}

func TestParseID(t *testing.T) {
	t.Run("有効なUUID文字列_正常系", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		id, err := ParseID(validUUID)

		if err != nil {
			t.Errorf("ParseID() error = %v, want nil", err)
		}
		if id.String() != validUUID {
			t.Errorf("ParseID() = %v, want %v", id.String(), validUUID)
		}
	})

	t.Run("無効なUUID文字列_異常系", func(t *testing.T) {
		invalidUUIDs := []string{
			"",
			"invalid",
			"550e8400-e29b-41d4-a716",
			"550e8400-e29b-41d4-a716-44665544000g",
			"not-a-uuid-at-all",
		}

		for _, invalid := range invalidUUIDs {
			_, err := ParseID(invalid)
			if err == nil {
				t.Errorf("ParseID(%q) expected error, got nil", invalid)
			}
		}
	})

	t.Run("境界値_空文字列", func(t *testing.T) {
		_, err := ParseID("")
		if err == nil {
			t.Error("ParseID(\"\") expected error, got nil")
		}
	})

	t.Run("境界値_nil UUID文字列", func(t *testing.T) {
		nilUUID := "00000000-0000-0000-0000-000000000000"
		id, err := ParseID(nilUUID)

		if err != nil {
			t.Errorf("ParseID() error = %v, want nil", err)
		}
		if id != uuid.Nil {
			t.Errorf("ParseID() = %v, want uuid.Nil", id)
		}
	})
}

// ============================================
// DateRange関連テスト
// ============================================

func TestNewDateRange(t *testing.T) {
	t.Run("正常な日付範囲生成", func(t *testing.T) {
		start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

		dr := NewDateRange(start, end)

		if !dr.Start.Equal(start) {
			t.Errorf("DateRange.Start = %v, want %v", dr.Start, start)
		}
		if !dr.End.Equal(end) {
			t.Errorf("DateRange.End = %v, want %v", dr.End, end)
		}
	})

	t.Run("同じ日付の範囲", func(t *testing.T) {
		date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		dr := NewDateRange(date, date)

		if !dr.Start.Equal(date) || !dr.End.Equal(date) {
			t.Error("DateRange should allow same start and end date")
		}
	})
}

func TestDateRange_Contains(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
	dr := NewDateRange(start, end)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "範囲内の日付_正常系",
			date:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "境界値_開始日",
			date:     start,
			expected: true,
		},
		{
			name:     "境界値_終了日",
			date:     end,
			expected: true,
		},
		{
			name:     "範囲外_開始日前",
			date:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "範囲外_終了日後",
			date:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "境界値_開始日の1日前",
			date:     time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			expected: false,
		},
		{
			name:     "境界値_終了日の1ナノ秒後",
			date:     time.Date(2025, 1, 31, 0, 0, 0, 1, time.UTC),
			expected: false, // time.Afterは厳密比較のため範囲外
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dr.Contains(tt.date)
			if result != tt.expected {
				t.Errorf("DateRange.Contains(%v) = %v, want %v", tt.date, result, tt.expected)
			}
		})
	}
}

func TestDateRange_Days(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "1日間",
			start:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "31日間_1月",
			start:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		{
			name:     "28日間_2月",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		{
			name:     "29日間_閏年2月",
			start:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		{
			name:     "365日間_1年",
			start:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		{
			name:     "366日間_閏年",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 366,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := NewDateRange(tt.start, tt.end)
			result := dr.Days()
			if result != tt.expected {
				t.Errorf("DateRange.Days() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// TimeRange関連テスト
// ============================================

func TestNewTimeRange(t *testing.T) {
	t.Run("通常の時間範囲", func(t *testing.T) {
		tr := NewTimeRange(9, 0, 17, 30)

		expectedStart := 9*60 + 0
		expectedEnd := 17*60 + 30

		if tr.Start != expectedStart {
			t.Errorf("TimeRange.Start = %v, want %v", tr.Start, expectedStart)
		}
		if tr.End != expectedEnd {
			t.Errorf("TimeRange.End = %v, want %v", tr.End, expectedEnd)
		}
	})

	t.Run("境界値_深夜0時", func(t *testing.T) {
		tr := NewTimeRange(0, 0, 0, 0)

		if tr.Start != 0 || tr.End != 0 {
			t.Errorf("TimeRange = {%v, %v}, want {0, 0}", tr.Start, tr.End)
		}
	})

	t.Run("境界値_23時59分", func(t *testing.T) {
		tr := NewTimeRange(23, 59, 23, 59)

		expected := 23*60 + 59

		if tr.Start != expected || tr.End != expected {
			t.Errorf("TimeRange = {%v, %v}, want {%v, %v}", tr.Start, tr.End, expected, expected)
		}
	})
}

func TestTimeRange_Duration(t *testing.T) {
	tests := []struct {
		name                                 string
		startHour, startMin, endHour, endMin int
		expected                             int
	}{
		{
			name:      "通常勤務_8時間",
			startHour: 9, startMin: 0, endHour: 17, endMin: 0,
			expected: 480,
		},
		{
			name:      "通常勤務_8時間30分",
			startHour: 9, startMin: 0, endHour: 17, endMin: 30,
			expected: 510,
		},
		{
			name:      "同じ時刻_0分",
			startHour: 12, startMin: 0, endHour: 12, endMin: 0,
			expected: 0,
		},
		{
			name:      "1分間",
			startHour: 12, startMin: 0, endHour: 12, endMin: 1,
			expected: 1,
		},
		{
			name:      "日跨ぎ_夜勤_16時間",
			startHour: 17, startMin: 0, endHour: 9, endMin: 0,
			expected: 960, // 17:00 -> 翌9:00 = 16時間 = 960分
		},
		{
			name:      "日跨ぎ_23時から1時",
			startHour: 23, startMin: 0, endHour: 1, endMin: 0,
			expected: 120, // 2時間 = 120分
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTimeRange(tt.startHour, tt.startMin, tt.endHour, tt.endMin)
			result := tr.Duration()
			if result != tt.expected {
				t.Errorf("TimeRange.Duration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// DomainError関連テスト
// ============================================

func TestDomainError_Error(t *testing.T) {
	t.Run("エラーメッセージ取得", func(t *testing.T) {
		err := &DomainError{
			Code:    "TEST_CODE",
			Message: "テストエラーメッセージ",
		}

		if err.Error() != "テストエラーメッセージ" {
			t.Errorf("DomainError.Error() = %v, want %v", err.Error(), "テストエラーメッセージ")
		}
	})

	t.Run("空のメッセージ", func(t *testing.T) {
		err := &DomainError{
			Code:    "EMPTY",
			Message: "",
		}

		if err.Error() != "" {
			t.Errorf("DomainError.Error() = %v, want empty string", err.Error())
		}
	})
}

func TestDomainError_Unwrap(t *testing.T) {
	t.Run("元エラーあり", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := &DomainError{
			Code:    "WRAP",
			Message: "wrapped error",
			Err:     originalErr,
		}

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("DomainError.Unwrap() = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("元エラーなし", func(t *testing.T) {
		err := &DomainError{
			Code:    "NO_WRAP",
			Message: "no wrapped error",
		}

		unwrapped := err.Unwrap()
		if unwrapped != nil {
			t.Errorf("DomainError.Unwrap() = %v, want nil", unwrapped)
		}
	})

	t.Run("errors.Is_互換性", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := &DomainError{
			Code:    "WRAP",
			Message: "wrapped error",
			Err:     originalErr,
		}

		if !errors.Is(err, originalErr) {
			t.Error("errors.Is() should return true for wrapped error")
		}
	})
}

func TestNewDomainError(t *testing.T) {
	t.Run("ドメインエラー生成", func(t *testing.T) {
		err := NewDomainError("TEST_CODE", "テストメッセージ")

		if err.Code != "TEST_CODE" {
			t.Errorf("DomainError.Code = %v, want %v", err.Code, "TEST_CODE")
		}
		if err.Message != "テストメッセージ" {
			t.Errorf("DomainError.Message = %v, want %v", err.Message, "テストメッセージ")
		}
		if err.Err != nil {
			t.Errorf("DomainError.Err = %v, want nil", err.Err)
		}
	})

	t.Run("空のコードとメッセージ", func(t *testing.T) {
		err := NewDomainError("", "")

		if err.Code != "" || err.Message != "" {
			t.Error("NewDomainError should allow empty code and message")
		}
	})
}

func TestWrapDomainError(t *testing.T) {
	t.Run("エラーラップ", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := WrapDomainError("WRAP_CODE", "ラップメッセージ", originalErr)

		if err.Code != "WRAP_CODE" {
			t.Errorf("DomainError.Code = %v, want %v", err.Code, "WRAP_CODE")
		}
		if err.Message != "ラップメッセージ" {
			t.Errorf("DomainError.Message = %v, want %v", err.Message, "ラップメッセージ")
		}
		if err.Err != originalErr {
			t.Errorf("DomainError.Err = %v, want %v", err.Err, originalErr)
		}
	})

	t.Run("nilエラーのラップ", func(t *testing.T) {
		err := WrapDomainError("NIL_WRAP", "メッセージ", nil)

		if err.Err != nil {
			t.Errorf("DomainError.Err = %v, want nil", err.Err)
		}
	})
}

// ============================================
// 定義済みエラーテスト
// ============================================

func TestPredefinedErrors(t *testing.T) {
	t.Run("ErrNotFound", func(t *testing.T) {
		if ErrNotFound.Code != ErrCodeNotFound {
			t.Errorf("ErrNotFound.Code = %v, want %v", ErrNotFound.Code, ErrCodeNotFound)
		}
		if ErrNotFound.Message == "" {
			t.Error("ErrNotFound.Message should not be empty")
		}
	})

	t.Run("ErrInvalidInput", func(t *testing.T) {
		if ErrInvalidInput.Code != ErrCodeInvalidInput {
			t.Errorf("ErrInvalidInput.Code = %v, want %v", ErrInvalidInput.Code, ErrCodeInvalidInput)
		}
		if ErrInvalidInput.Message == "" {
			t.Error("ErrInvalidInput.Message should not be empty")
		}
	})

	t.Run("ErrConflict", func(t *testing.T) {
		if ErrConflict.Code != ErrCodeConflict {
			t.Errorf("ErrConflict.Code = %v, want %v", ErrConflict.Code, ErrCodeConflict)
		}
		if ErrConflict.Message == "" {
			t.Error("ErrConflict.Message should not be empty")
		}
	})
}

func TestErrorCodes(t *testing.T) {
	expectedCodes := map[string]string{
		"ErrCodeNotFound":     ErrCodeNotFound,
		"ErrCodeInvalidInput": ErrCodeInvalidInput,
		"ErrCodeConflict":     ErrCodeConflict,
		"ErrCodeInternal":     ErrCodeInternal,
		"ErrCodeUnauthorized": ErrCodeUnauthorized,
		"ErrCodeForbidden":    ErrCodeForbidden,
		"ErrCodeValidation":   ErrCodeValidation,
	}

	for name, code := range expectedCodes {
		t.Run(name+"_非空", func(t *testing.T) {
			if code == "" {
				t.Errorf("%s should not be empty", name)
			}
		})
	}

	// 一意性チェック
	t.Run("エラーコードの一意性", func(t *testing.T) {
		codes := []string{
			ErrCodeNotFound,
			ErrCodeInvalidInput,
			ErrCodeConflict,
			ErrCodeInternal,
			ErrCodeUnauthorized,
			ErrCodeForbidden,
			ErrCodeValidation,
		}

		seen := make(map[string]bool)
		for _, code := range codes {
			if seen[code] {
				t.Errorf("Duplicate error code: %s", code)
			}
			seen[code] = true
		}
	})
}
