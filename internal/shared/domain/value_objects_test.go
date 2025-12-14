package domain

import (
	"testing"
)

// Email テスト
func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常系_通常メール", "test@example.com", "test@example.com", false},
		{"正常系_大文字含む", "Test@Example.COM", "test@example.com", false},
		{"正常系_サブドメイン", "test@sub.example.com", "test@sub.example.com", false},
		{"正常系_プラス記号", "test+tag@example.com", "test+tag@example.com", false},
		{"異常系_空文字", "", "", true},
		{"異常系_@なし", "testexample.com", "", true},
		{"異常系_ドメインなし", "test@", "", true},
		{"異常系_ローカルパートなし", "@example.com", "", true},
		{"異常系_TLDなし", "test@example", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && email.String() != tt.want {
				t.Errorf("NewEmail() = %v, want %v", email.String(), tt.want)
			}
		})
	}
}

func TestEmail_Parts(t *testing.T) {
	email := MustNewEmail("test@example.com")

	if email.LocalPart() != "test" {
		t.Errorf("LocalPart() = %v, want test", email.LocalPart())
	}
	if email.Domain() != "example.com" {
		t.Errorf("Domain() = %v, want example.com", email.Domain())
	}
}

// PhoneNumber テスト
func TestNewPhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常系_携帯", "090-1234-5678", "09012345678", false},
		{"正常系_固定", "03-1234-5678", "0312345678", false},
		{"正常系_数字のみ", "09012345678", "09012345678", false},
		{"正常系_空文字許可", "", "", false},
		{"異常系_短すぎ", "123456789", "", true},
		{"異常系_長すぎ", "1234567890123456", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhoneNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && phone.String() != tt.want {
				t.Errorf("NewPhoneNumber() = %v, want %v", phone.String(), tt.want)
			}
		})
	}
}

func TestPhoneNumber_Formatted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"携帯", "09012345678", "090-1234-5678"},
		{"東京", "0312345678", "03-1234-5678"},
		{"大阪", "0612345678", "06-1234-5678"},
		{"地方", "0123456789", "0123-45-6789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := MustNewPhoneNumber(tt.input)
			if got := phone.Formatted(); got != tt.want {
				t.Errorf("Formatted() = %v, want %v", got, tt.want)
			}
		})
	}
}

// PersonName テスト
func TestNewPersonName(t *testing.T) {
	tests := []struct {
		name      string
		lastName  string
		firstName string
		wantFull  string
		wantErr   bool
	}{
		{"正常系", "山田", "太郎", "山田 太郎", false},
		{"正常系_スペースあり", " 山田 ", " 太郎 ", "山田 太郎", false},
		{"異常系_姓なし", "", "太郎", "", true},
		{"異常系_名なし", "山田", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := NewPersonName(tt.lastName, tt.firstName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPersonName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && name.FullName() != tt.wantFull {
				t.Errorf("FullName() = %v, want %v", name.FullName(), tt.wantFull)
			}
		})
	}
}

func TestPersonName_Initials(t *testing.T) {
	name := MustNewPersonName("山田", "太郎")
	if got := name.Initials(); got != "山太" {
		t.Errorf("Initials() = %v, want 山太", got)
	}
}

// Color テスト
func TestNewColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常系_大文字", "#FF5733", "#FF5733", false},
		{"正常系_小文字", "#ff5733", "#FF5733", false},
		{"正常系_#なし", "FF5733", "#FF5733", false},
		{"正常系_空文字デフォルト", "", "#6B7280", false},
		{"異常系_短すぎ", "#FFF", "", true},
		{"異常系_長すぎ", "#FF57331", "", true},
		{"異常系_不正文字", "#GGGGGG", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := NewColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && color.String() != tt.want {
				t.Errorf("NewColor() = %v, want %v", color.String(), tt.want)
			}
		})
	}
}

func TestColor_RGB(t *testing.T) {
	color := MustNewColor("#FF5733")
	r, g, b := color.RGB()
	if r != 255 || g != 87 || b != 51 {
		t.Errorf("RGB() = %d,%d,%d, want 255,87,51", r, g, b)
	}
}

func TestColor_IsDark(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  bool
	}{
		{"黒", "#000000", true},
		{"白", "#FFFFFF", false},
		{"濃い青", "#000080", true},
		{"明るい黄", "#FFFF00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := MustNewColor(tt.color)
			if got := color.IsDark(); got != tt.want {
				t.Errorf("IsDark() = %v, want %v", got, tt.want)
			}
		})
	}
}

// EmployeeCode テスト
func TestNewEmployeeCode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常系_英数字", "EMP001", "EMP001", false},
		{"正常系_小文字変換", "emp001", "EMP001", false},
		{"正常系_ハイフン", "EMP-001", "EMP-001", false},
		{"異常系_空文字", "", "", true},
		{"異常系_長すぎ", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "", true},
		{"異常系_不正文字", "EMP@001", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := NewEmployeeCode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmployeeCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && code.String() != tt.want {
				t.Errorf("NewEmployeeCode() = %v, want %v", code.String(), tt.want)
			}
		})
	}
}

func TestEmployeeCode_Equals(t *testing.T) {
	c1 := MustNewEmployeeCode("EMP001")
	c2 := MustNewEmployeeCode("emp001") // 小文字→大文字変換されるので同じ
	c3 := MustNewEmployeeCode("EMP002")

	if !c1.Equals(c2) {
		t.Error("同一コードがEqualsでfalse")
	}
	if c1.Equals(c3) {
		t.Error("異なるコードがEqualsでtrue")
	}
}
