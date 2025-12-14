// Package domain 共有ドメイン型定義
package domain

import "strings"

// PersonName 人名 Value Object
type PersonName struct {
	lastName  string
	firstName string
}

// NewPersonName 人名生成
func NewPersonName(lastName, firstName string) (PersonName, error) {
	lastName = strings.TrimSpace(lastName)
	firstName = strings.TrimSpace(firstName)

	if lastName == "" {
		return PersonName{}, NewDomainError(ErrCodeValidation, "姓は必須です")
	}

	if firstName == "" {
		return PersonName{}, NewDomainError(ErrCodeValidation, "名は必須です")
	}

	if len(lastName) > 50 {
		return PersonName{}, NewDomainError(ErrCodeValidation, "姓は50文字以内で入力してください")
	}

	if len(firstName) > 50 {
		return PersonName{}, NewDomainError(ErrCodeValidation, "名は50文字以内で入力してください")
	}

	return PersonName{lastName: lastName, firstName: firstName}, nil
}

// MustNewPersonName 生成失敗時パニック テスト用
func MustNewPersonName(lastName, firstName string) PersonName {
	n, err := NewPersonName(lastName, firstName)
	if err != nil {
		panic(err)
	}
	return n
}

// LastName 姓取得
func (n PersonName) LastName() string {
	return n.lastName
}

// FirstName 名取得
func (n PersonName) FirstName() string {
	return n.firstName
}

// FullName フルネーム取得 "姓 名"
func (n PersonName) FullName() string {
	return n.lastName + " " + n.firstName
}

// FullNameReversed フルネーム取得 "名 姓"（英語式）
func (n PersonName) FullNameReversed() string {
	return n.firstName + " " + n.lastName
}

// Initials イニシャル取得
func (n PersonName) Initials() string {
	var initials string
	if len(n.lastName) > 0 {
		initials += string([]rune(n.lastName)[0])
	}
	if len(n.firstName) > 0 {
		initials += string([]rune(n.firstName)[0])
	}
	return initials
}

// Equals 等価判定
func (n PersonName) Equals(other PersonName) bool {
	return n.lastName == other.lastName && n.firstName == other.firstName
}

// IsZero ゼロ値判定
func (n PersonName) IsZero() bool {
	return n.lastName == "" && n.firstName == ""
}
