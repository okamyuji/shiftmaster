// Package domain 共有ドメイン型定義
package domain

import (
	"regexp"
	"strings"
)

// PhoneNumber 電話番号 Value Object
type PhoneNumber struct {
	value string
}

// phoneCleanRegex 電話番号から不要文字を除去する正規表現
var phoneCleanRegex = regexp.MustCompile(`[^0-9]`)

// NewPhoneNumber 電話番号生成
func NewPhoneNumber(phone string) (PhoneNumber, error) {
	if phone == "" {
		// 電話番号は任意項目として空を許可
		return PhoneNumber{}, nil
	}

	// 数字のみ抽出
	cleaned := phoneCleanRegex.ReplaceAllString(phone, "")

	if len(cleaned) < 10 {
		return PhoneNumber{}, NewDomainError(ErrCodeValidation, "電話番号は10桁以上で入力してください")
	}

	if len(cleaned) > 15 {
		return PhoneNumber{}, NewDomainError(ErrCodeValidation, "電話番号は15桁以内で入力してください")
	}

	return PhoneNumber{value: cleaned}, nil
}

// MustNewPhoneNumber 生成失敗時パニック テスト用
func MustNewPhoneNumber(phone string) PhoneNumber {
	p, err := NewPhoneNumber(phone)
	if err != nil {
		panic(err)
	}
	return p
}

// String 文字列取得（数字のみ）
func (p PhoneNumber) String() string {
	return p.value
}

// Formatted ハイフン区切りフォーマット
func (p PhoneNumber) Formatted() string {
	if p.value == "" {
		return ""
	}

	// 日本の電話番号形式でフォーマット
	if len(p.value) == 11 && strings.HasPrefix(p.value, "0") {
		// 携帯電話 090-1234-5678
		return p.value[:3] + "-" + p.value[3:7] + "-" + p.value[7:]
	}
	if len(p.value) == 10 && strings.HasPrefix(p.value, "0") {
		// 固定電話 03-1234-5678 or 0123-45-6789
		if strings.HasPrefix(p.value, "03") || strings.HasPrefix(p.value, "06") {
			return p.value[:2] + "-" + p.value[2:6] + "-" + p.value[6:]
		}
		return p.value[:4] + "-" + p.value[4:6] + "-" + p.value[6:]
	}

	return p.value
}

// Equals 等価判定
func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

// IsZero ゼロ値判定
func (p PhoneNumber) IsZero() bool {
	return p.value == ""
}
