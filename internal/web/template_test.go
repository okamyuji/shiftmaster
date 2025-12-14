// Package web テンプレートエンジンテスト
package web

import (
	"html/template"
	"testing"
	"time"
)

// ============================================
// NewTemplateEngine関連テスト
// ============================================

func TestNewTemplateEngine(t *testing.T) {
	t.Run("正常系_テンプレートエンジン生成", func(t *testing.T) {
		engine := NewTemplateEngine("/path/to/templates")

		if engine == nil {
			t.Fatal("NewTemplateEngine() returned nil")
		}
		if engine.templates == nil {
			t.Error("templates map should not be nil")
		}
		if engine.funcMap == nil {
			t.Error("funcMap should not be nil")
		}
		if engine.baseDir != "/path/to/templates" {
			t.Errorf("baseDir = %v, want %v", engine.baseDir, "/path/to/templates")
		}
	})

	t.Run("空のベースディレクトリ", func(t *testing.T) {
		engine := NewTemplateEngine("")

		if engine == nil {
			t.Fatal("NewTemplateEngine() returned nil")
		}
		if engine.baseDir != "" {
			t.Errorf("baseDir = %v, want empty string", engine.baseDir)
		}
	})
}

// ============================================
// TemplateError関連テスト
// ============================================

func TestTemplateError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *TemplateError
		expected string
	}{
		{
			name: "正常なエラーメッセージ",
			err: &TemplateError{
				Name:    "pages/test.html",
				Message: "テンプレートが見つかりません",
			},
			expected: "テンプレートエラー [pages/test.html]: テンプレートが見つかりません",
		},
		{
			name: "空の名前",
			err: &TemplateError{
				Name:    "",
				Message: "エラー",
			},
			expected: "テンプレートエラー []: エラー",
		},
		{
			name: "空のメッセージ",
			err: &TemplateError{
				Name:    "test.html",
				Message: "",
			},
			expected: "テンプレートエラー [test.html]: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============================================
// defaultFuncMap関連テスト
// ============================================

func TestDefaultFuncMap(t *testing.T) {
	funcMap := defaultFuncMap()

	t.Run("funcMapの存在確認", func(t *testing.T) {
		expectedFuncs := []string{
			"safeHTML",
			"formatDate",
			"formatDateInput",
			"formatDateTime",
			"add",
			"sub",
			"mul",
			"seq",
			"iterate",
		}

		for _, name := range expectedFuncs {
			if _, ok := funcMap[name]; !ok {
				t.Errorf("funcMap should contain %s", name)
			}
		}
	})
}

func TestFuncMap_SafeHTML(t *testing.T) {
	funcMap := defaultFuncMap()
	safeHTMLFunc := funcMap["safeHTML"].(func(string) template.HTML)

	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "通常のHTML",
			input:    "<div>Hello</div>",
			expected: template.HTML("<div>Hello</div>"),
		},
		{
			name:     "スクリプトタグ",
			input:    "<script>alert('test')</script>",
			expected: template.HTML("<script>alert('test')</script>"),
		},
		{
			name:     "空文字列",
			input:    "",
			expected: template.HTML(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeHTMLFunc(tt.input)
			if result != tt.expected {
				t.Errorf("safeHTML() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFuncMap_FormatDate(t *testing.T) {
	funcMap := defaultFuncMap()
	formatDateFunc := funcMap["formatDate"].(func(any) string)

	t.Run("time.Time型", func(t *testing.T) {
		date := time.Date(2025, 3, 15, 10, 30, 0, 0, time.UTC)
		result := formatDateFunc(date)
		if result != "2025/03/15" {
			t.Errorf("formatDate() = %v, want %v", result, "2025/03/15")
		}
	})

	t.Run("RFC3339形式の文字列", func(t *testing.T) {
		result := formatDateFunc("2025-03-15T10:30:00Z")
		if result != "2025/03/15" {
			t.Errorf("formatDate() = %v, want %v", result, "2025/03/15")
		}
	})

	t.Run("YYYY-MM-DD形式の文字列", func(t *testing.T) {
		result := formatDateFunc("2025-03-15")
		if result != "2025/03/15" {
			t.Errorf("formatDate() = %v, want %v", result, "2025/03/15")
		}
	})

	t.Run("不正な形式の文字列", func(t *testing.T) {
		result := formatDateFunc("invalid-date")
		if result != "invalid-date" {
			t.Errorf("formatDate() = %v, want %v", result, "invalid-date")
		}
	})

	t.Run("空の文字列", func(t *testing.T) {
		result := formatDateFunc("")
		if result != "" {
			t.Errorf("formatDate() = %v, want empty string", result)
		}
	})

	t.Run("その他の型", func(t *testing.T) {
		result := formatDateFunc(12345)
		if result != "" {
			t.Errorf("formatDate() = %v, want empty string for int type", result)
		}
	})
}

func TestFuncMap_FormatDateInput(t *testing.T) {
	funcMap := defaultFuncMap()
	formatDateInputFunc := funcMap["formatDateInput"].(func(any) string)

	t.Run("time.Time型", func(t *testing.T) {
		date := time.Date(2025, 3, 15, 10, 30, 0, 0, time.UTC)
		result := formatDateInputFunc(date)
		if result != "2025-03-15" {
			t.Errorf("formatDateInput() = %v, want %v", result, "2025-03-15")
		}
	})

	t.Run("RFC3339形式の文字列", func(t *testing.T) {
		result := formatDateInputFunc("2025-03-15T10:30:00Z")
		if result != "2025-03-15" {
			t.Errorf("formatDateInput() = %v, want %v", result, "2025-03-15")
		}
	})

	t.Run("その他の型", func(t *testing.T) {
		result := formatDateInputFunc(nil)
		if result != "" {
			t.Errorf("formatDateInput() = %v, want empty string for nil", result)
		}
	})
}

func TestFuncMap_FormatDateTime(t *testing.T) {
	funcMap := defaultFuncMap()
	formatDateTimeFunc := funcMap["formatDateTime"].(func(any) string)

	t.Run("time.Time型", func(t *testing.T) {
		// UTC 10:30 はJST 19:30
		datetime := time.Date(2025, 3, 15, 10, 30, 0, 0, time.UTC)
		result := formatDateTimeFunc(datetime)
		if result != "2025/03/15 19:30" {
			t.Errorf("formatDateTime() = %v, want %v", result, "2025/03/15 19:30")
		}
	})

	t.Run("RFC3339形式の文字列", func(t *testing.T) {
		// UTC 14:45 はJST 23:45
		result := formatDateTimeFunc("2025-03-15T14:45:00Z")
		if result != "2025/03/15 23:45" {
			t.Errorf("formatDateTime() = %v, want %v", result, "2025/03/15 23:45")
		}
	})
}

func TestFuncMap_Add(t *testing.T) {
	funcMap := defaultFuncMap()
	addFunc := funcMap["add"].(func(int, int) int)

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正の数", 5, 3, 8},
		{"負の数", -5, 3, -2},
		{"ゼロ", 0, 0, 0},
		{"大きな数", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addFunc(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFuncMap_Sub(t *testing.T) {
	funcMap := defaultFuncMap()
	subFunc := funcMap["sub"].(func(int, int) int)

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正の数", 5, 3, 2},
		{"負の結果", 3, 5, -2},
		{"ゼロ", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subFunc(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("sub(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFuncMap_Mul(t *testing.T) {
	funcMap := defaultFuncMap()
	mulFunc := funcMap["mul"].(func(int, int) int)

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正の数", 5, 3, 15},
		{"負の数", -5, 3, -15},
		{"ゼロ", 5, 0, 0},
		{"両方負", -5, -3, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mulFunc(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("mul(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFuncMap_Seq(t *testing.T) {
	funcMap := defaultFuncMap()
	seqFunc := funcMap["seq"].(func(int, int) []int)

	tests := []struct {
		name       string
		start, end int
		expected   []int
	}{
		{"1から5", 1, 5, []int{1, 2, 3, 4, 5}},
		{"0から3", 0, 3, []int{0, 1, 2, 3}},
		{"同じ値", 5, 5, []int{5}},
		{"負の範囲", -2, 2, []int{-2, -1, 0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seqFunc(tt.start, tt.end)

			if len(result) != len(tt.expected) {
				t.Errorf("seq(%d, %d) length = %d, want %d", tt.start, tt.end, len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("seq(%d, %d)[%d] = %d, want %d", tt.start, tt.end, i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestFuncMap_Iterate(t *testing.T) {
	funcMap := defaultFuncMap()
	iterateFunc := funcMap["iterate"].(func(int, int) []int)

	tests := []struct {
		name       string
		start, end int
		expected   []int
	}{
		{"1から5まで（5は含まない）", 1, 5, []int{1, 2, 3, 4}},
		{"0から3まで（3は含まない）", 0, 3, []int{0, 1, 2}},
		{"同じ値", 5, 5, []int{}},
		{"月のイテレート", 1, 13, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := iterateFunc(tt.start, tt.end)

			if len(result) != len(tt.expected) {
				t.Errorf("iterate(%d, %d) length = %d, want %d", tt.start, tt.end, len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("iterate(%d, %d)[%d] = %d, want %d", tt.start, tt.end, i, v, tt.expected[i])
				}
			}
		})
	}
}
