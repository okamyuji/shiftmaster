// Package domain 共有ドメイン型定義
package domain

import (
	"regexp"
	"strings"
)

// Email メールアドレス Value Object
type Email struct {
	value string
}

// emailRegex メールアドレス形式チェック用正規表現
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail メールアドレス生成
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, NewDomainError(ErrCodeValidation, "メールアドレスは必須です")
	}

	email = strings.TrimSpace(strings.ToLower(email))

	if !emailRegex.MatchString(email) {
		return Email{}, NewDomainError(ErrCodeValidation, "メールアドレスの形式が不正です")
	}

	return Email{value: email}, nil
}

// MustNewEmail 生成失敗時パニック テスト用
func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

// String 文字列取得
func (e Email) String() string {
	return e.value
}

// LocalPart ローカルパート取得 @より前
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// Domain ドメイン取得 @より後
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// Equals 等価判定
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsZero ゼロ値判定
func (e Email) IsZero() bool {
	return e.value == ""
}
