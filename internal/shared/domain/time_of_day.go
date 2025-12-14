// Package domain 共有ドメイン型定義
package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// TimeOfDay 時刻 Value Object
type TimeOfDay struct {
	hour   int
	minute int
}

// NewTimeOfDay 時刻生成
func NewTimeOfDay(hour, minute int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "時間は0-23の範囲で指定してください")
	}
	if minute < 0 || minute > 59 {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "分は0-59の範囲で指定してください")
	}
	return TimeOfDay{hour: hour, minute: minute}, nil
}

// ParseTimeOfDay 文字列から時刻生成 "HH:MM" 形式
func ParseTimeOfDay(s string) (TimeOfDay, error) {
	if s == "" {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "時刻が指定されていません")
	}

	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "時刻は HH:MM 形式で指定してください")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "時間が数値ではありません")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return TimeOfDay{}, NewDomainError(ErrCodeValidation, "分が数値ではありません")
	}

	return NewTimeOfDay(hour, minute)
}

// MustParseTimeOfDay パース失敗時パニック テスト用
func MustParseTimeOfDay(s string) TimeOfDay {
	t, err := ParseTimeOfDay(s)
	if err != nil {
		panic(err)
	}
	return t
}

// Hour 時間取得
func (t TimeOfDay) Hour() int {
	return t.hour
}

// Minute 分取得
func (t TimeOfDay) Minute() int {
	return t.minute
}

// String 文字列変換 "HH:MM" 形式
func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.minute)
}

// ToMinutes 00:00からの経過分数
func (t TimeOfDay) ToMinutes() int {
	return t.hour*60 + t.minute
}

// IsBefore 指定時刻より前か判定
func (t TimeOfDay) IsBefore(other TimeOfDay) bool {
	return t.ToMinutes() < other.ToMinutes()
}

// IsAfter 指定時刻より後か判定
func (t TimeOfDay) IsAfter(other TimeOfDay) bool {
	return t.ToMinutes() > other.ToMinutes()
}

// Equals 等価判定
func (t TimeOfDay) Equals(other TimeOfDay) bool {
	return t.hour == other.hour && t.minute == other.minute
}

// IsZero ゼロ値判定
func (t TimeOfDay) IsZero() bool {
	return t.hour == 0 && t.minute == 0
}

// Add 分を加算 日跨ぎ対応
func (t TimeOfDay) Add(minutes int) TimeOfDay {
	total := t.ToMinutes() + minutes
	// 負の値対応
	for total < 0 {
		total += 1440
	}
	total = total % 1440
	return TimeOfDay{hour: total / 60, minute: total % 60}
}
