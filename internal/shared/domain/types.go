// Package domain 共有ドメイン型定義
package domain

import (
	"time"

	"github.com/google/uuid"
)

// ID エンティティ識別子
type ID = uuid.UUID

// NewID 新規ID生成
func NewID() ID {
	return uuid.New()
}

// ParseID 文字列からID変換
func ParseID(s string) (ID, error) {
	return uuid.Parse(s)
}

// DateRange 日付範囲
type DateRange struct {
	// Start 開始日
	Start time.Time
	// End 終了日
	End time.Time
}

// NewDateRange 日付範囲生成
func NewDateRange(start, end time.Time) DateRange {
	return DateRange{Start: start, End: end}
}

// Contains 指定日が範囲内か判定
func (r DateRange) Contains(t time.Time) bool {
	return !t.Before(r.Start) && !t.After(r.End)
}

// Days 範囲内の日数
func (r DateRange) Days() int {
	return int(r.End.Sub(r.Start).Hours()/24) + 1
}

// TimeRange 時間範囲
type TimeRange struct {
	// Start 開始時刻 分単位 0-1439
	Start int
	// End 終了時刻 分単位 0-1439 日跨ぎの場合は1440以上
	End int
}

// NewTimeRange 時間範囲生成
func NewTimeRange(startHour, startMin, endHour, endMin int) TimeRange {
	start := startHour*60 + startMin
	end := endHour*60 + endMin
	return TimeRange{Start: start, End: end}
}

// Duration 所要時間 分単位
func (r TimeRange) Duration() int {
	if r.End >= r.Start {
		return r.End - r.Start
	}
	// 日跨ぎの場合
	return (1440 - r.Start) + r.End
}

// DomainError ドメインエラー
type DomainError struct {
	// Code エラーコード
	Code string
	// Message エラーメッセージ
	Message string
	// Err 元エラー
	Err error
}

// Error エラーメッセージ取得
func (e *DomainError) Error() string {
	return e.Message
}

// Unwrap 元エラー取得
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError ドメインエラー生成
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// WrapDomainError エラーラップ
func WrapDomainError(code, message string, err error) *DomainError {
	return &DomainError{Code: code, Message: message, Err: err}
}

// 共通エラーコード
const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInvalidInput = "INVALID_INPUT"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeValidation   = "VALIDATION_ERROR"
)

// ErrNotFound リソース未検出エラー
var ErrNotFound = NewDomainError(ErrCodeNotFound, "リソースが見つかりません")

// ErrInvalidInput 入力値不正エラー
var ErrInvalidInput = NewDomainError(ErrCodeInvalidInput, "入力値が不正です")

// ErrConflict 競合エラー
var ErrConflict = NewDomainError(ErrCodeConflict, "リソースが競合しています")
