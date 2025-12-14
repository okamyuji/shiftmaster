// Package domain 共有ドメイン型定義
package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// Color 色 Value Object
type Color struct {
	value string
}

// colorHexRegex HEXカラー形式チェック用正規表現
var colorHexRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// DefaultColor デフォルト色
var DefaultColor = Color{value: "#6B7280"}

// NewColor 色生成
func NewColor(hex string) (Color, error) {
	if hex == "" {
		return DefaultColor, nil
	}

	hex = strings.TrimSpace(hex)

	// #がなければ付与
	if !strings.HasPrefix(hex, "#") {
		hex = "#" + hex
	}

	if !colorHexRegex.MatchString(hex) {
		return Color{}, NewDomainError(ErrCodeValidation, "色はHEX形式（#RRGGBB）で指定してください")
	}

	return Color{value: strings.ToUpper(hex)}, nil
}

// MustNewColor 生成失敗時パニック テスト用
func MustNewColor(hex string) Color {
	c, err := NewColor(hex)
	if err != nil {
		panic(err)
	}
	return c
}

// String HEX文字列取得
func (c Color) String() string {
	return c.value
}

// RGB RGB値取得
func (c Color) RGB() (r, g, b int) {
	if c.value == "" {
		return 0, 0, 0
	}
	_, _ = fmt.Sscanf(c.value, "#%02X%02X%02X", &r, &g, &b)
	return
}

// IsDark 暗い色か判定（テキスト色決定用）
func (c Color) IsDark() bool {
	r, g, b := c.RGB()
	// 輝度計算（YIQ式）
	luminance := (r*299 + g*587 + b*114) / 1000
	return luminance < 128
}

// ContrastColor コントラスト色取得（テキスト用）
func (c Color) ContrastColor() Color {
	if c.IsDark() {
		return Color{value: "#FFFFFF"}
	}
	return Color{value: "#000000"}
}

// Equals 等価判定
func (c Color) Equals(other Color) bool {
	return c.value == other.value
}

// IsZero ゼロ値判定
func (c Color) IsZero() bool {
	return c.value == ""
}
