// Package domain 共有ドメイン型定義
package domain

import (
	"regexp"
	"strings"
)

// EmployeeCode 社員番号 Value Object
type EmployeeCode struct {
	value string
}

// employeeCodeRegex 社員番号形式チェック用正規表現（英数字とハイフンのみ）
var employeeCodeRegex = regexp.MustCompile(`^[A-Za-z0-9\-]{1,20}$`)

// NewEmployeeCode 社員番号生成
func NewEmployeeCode(code string) (EmployeeCode, error) {
	if code == "" {
		return EmployeeCode{}, NewDomainError(ErrCodeValidation, "社員番号は必須です")
	}

	code = strings.TrimSpace(strings.ToUpper(code))

	if !employeeCodeRegex.MatchString(code) {
		return EmployeeCode{}, NewDomainError(ErrCodeValidation, "社員番号は英数字とハイフンのみ20文字以内で入力してください")
	}

	return EmployeeCode{value: code}, nil
}

// MustNewEmployeeCode 生成失敗時パニック テスト用
func MustNewEmployeeCode(code string) EmployeeCode {
	c, err := NewEmployeeCode(code)
	if err != nil {
		panic(err)
	}
	return c
}

// String 文字列取得
func (c EmployeeCode) String() string {
	return c.value
}

// Equals 等価判定
func (c EmployeeCode) Equals(other EmployeeCode) bool {
	return c.value == other.value
}

// IsZero ゼロ値判定
func (c EmployeeCode) IsZero() bool {
	return c.value == ""
}
