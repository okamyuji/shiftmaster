// Package config アプリケーション設定テスト
package config

import (
	"os"
	"testing"
)

// ============================================
// Load関連テスト
// ============================================

func TestLoad(t *testing.T) {
	t.Run("デフォルト値での読み込み", func(t *testing.T) {
		// 環境変数をクリア
		_ = os.Unsetenv("SERVER_HOST")
		_ = os.Unsetenv("SERVER_PORT")
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Unsetenv("LOG_LEVEL")

		cfg := Load()

		if cfg.Server.Host != "0.0.0.0" {
			t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "0.0.0.0")
		}
		if cfg.Server.Port != 8080 {
			t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 8080)
		}
		if cfg.Database.URL == "" {
			t.Error("Database.URL should not be empty")
		}
		if cfg.Log.Level != "info" {
			t.Errorf("Log.Level = %v, want %v", cfg.Log.Level, "info")
		}
	})

	t.Run("環境変数からの読み込み", func(t *testing.T) {
		// テスト用の環境変数を設定
		_ = os.Setenv("SERVER_HOST", "127.0.0.1")
		_ = os.Setenv("SERVER_PORT", "3000")
		_ = os.Setenv("DATABASE_URL", "postgres://test:test@testdb:5432/testdb")
		_ = os.Setenv("LOG_LEVEL", "debug")
		defer func() {
			_ = os.Unsetenv("SERVER_HOST")
			_ = os.Unsetenv("SERVER_PORT")
			_ = os.Unsetenv("DATABASE_URL")
			_ = os.Unsetenv("LOG_LEVEL")
		}()

		cfg := Load()

		if cfg.Server.Host != "127.0.0.1" {
			t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "127.0.0.1")
		}
		if cfg.Server.Port != 3000 {
			t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 3000)
		}
		if cfg.Database.URL != "postgres://test:test@testdb:5432/testdb" {
			t.Errorf("Database.URL = %v, want %v", cfg.Database.URL, "postgres://test:test@testdb:5432/testdb")
		}
		if cfg.Log.Level != "debug" {
			t.Errorf("Log.Level = %v, want %v", cfg.Log.Level, "debug")
		}
	})
}

// ============================================
// getEnv関連テスト
// ============================================

func TestGetEnv(t *testing.T) {
	t.Run("環境変数が設定されている場合", func(t *testing.T) {
		_ = os.Setenv("TEST_VAR", "test_value")
		defer func() { _ = os.Unsetenv("TEST_VAR") }()

		result := getEnv("TEST_VAR", "default")
		if result != "test_value" {
			t.Errorf("getEnv() = %v, want %v", result, "test_value")
		}
	})

	t.Run("環境変数が設定されていない場合", func(t *testing.T) {
		_ = os.Unsetenv("NONEXISTENT_VAR")

		result := getEnv("NONEXISTENT_VAR", "default_value")
		if result != "default_value" {
			t.Errorf("getEnv() = %v, want %v", result, "default_value")
		}
	})

	t.Run("空の環境変数", func(t *testing.T) {
		_ = os.Setenv("EMPTY_VAR", "")
		defer func() { _ = os.Unsetenv("EMPTY_VAR") }()

		result := getEnv("EMPTY_VAR", "default")
		if result != "default" {
			t.Errorf("getEnv() = %v, want %v (empty string should use default)", result, "default")
		}
	})

	t.Run("空のデフォルト値", func(t *testing.T) {
		_ = os.Unsetenv("NONEXISTENT_VAR")

		result := getEnv("NONEXISTENT_VAR", "")
		if result != "" {
			t.Errorf("getEnv() = %v, want empty string", result)
		}
	})

	t.Run("特殊文字を含む値", func(t *testing.T) {
		specialValue := "postgres://user:p@ss=word@host:5432/db?sslmode=disable"
		_ = os.Setenv("SPECIAL_VAR", specialValue)
		defer func() { _ = os.Unsetenv("SPECIAL_VAR") }()

		result := getEnv("SPECIAL_VAR", "default")
		if result != specialValue {
			t.Errorf("getEnv() = %v, want %v", result, specialValue)
		}
	})

	t.Run("日本語を含む値", func(t *testing.T) {
		japaneseValue := "日本語テスト"
		_ = os.Setenv("JAPANESE_VAR", japaneseValue)
		defer func() { _ = os.Unsetenv("JAPANESE_VAR") }()

		result := getEnv("JAPANESE_VAR", "default")
		if result != japaneseValue {
			t.Errorf("getEnv() = %v, want %v", result, japaneseValue)
		}
	})
}

// ============================================
// getEnvInt関連テスト
// ============================================

func TestGetEnvInt(t *testing.T) {
	t.Run("正常な整数値", func(t *testing.T) {
		_ = os.Setenv("INT_VAR", "12345")
		defer func() { _ = os.Unsetenv("INT_VAR") }()

		result := getEnvInt("INT_VAR", 0)
		if result != 12345 {
			t.Errorf("getEnvInt() = %v, want %v", result, 12345)
		}
	})

	t.Run("環境変数が設定されていない場合", func(t *testing.T) {
		_ = os.Unsetenv("NONEXISTENT_INT_VAR")

		result := getEnvInt("NONEXISTENT_INT_VAR", 9999)
		if result != 9999 {
			t.Errorf("getEnvInt() = %v, want %v", result, 9999)
		}
	})

	t.Run("不正な整数値", func(t *testing.T) {
		_ = os.Setenv("INVALID_INT_VAR", "not_a_number")
		defer func() { _ = os.Unsetenv("INVALID_INT_VAR") }()

		result := getEnvInt("INVALID_INT_VAR", 100)
		if result != 100 {
			t.Errorf("getEnvInt() = %v, want %v (invalid string should use default)", result, 100)
		}
	})

	t.Run("空文字列", func(t *testing.T) {
		_ = os.Setenv("EMPTY_INT_VAR", "")
		defer func() { _ = os.Unsetenv("EMPTY_INT_VAR") }()

		result := getEnvInt("EMPTY_INT_VAR", 50)
		if result != 50 {
			t.Errorf("getEnvInt() = %v, want %v (empty string should use default)", result, 50)
		}
	})

	t.Run("境界値_ゼロ", func(t *testing.T) {
		_ = os.Setenv("ZERO_VAR", "0")
		defer func() { _ = os.Unsetenv("ZERO_VAR") }()

		result := getEnvInt("ZERO_VAR", 100)
		if result != 0 {
			t.Errorf("getEnvInt() = %v, want %v", result, 0)
		}
	})

	t.Run("境界値_負数", func(t *testing.T) {
		_ = os.Setenv("NEGATIVE_VAR", "-100")
		defer func() { _ = os.Unsetenv("NEGATIVE_VAR") }()

		result := getEnvInt("NEGATIVE_VAR", 50)
		if result != -100 {
			t.Errorf("getEnvInt() = %v, want %v", result, -100)
		}
	})

	t.Run("境界値_大きな正数", func(t *testing.T) {
		_ = os.Setenv("LARGE_VAR", "2147483647") // int32 max
		defer func() { _ = os.Unsetenv("LARGE_VAR") }()

		result := getEnvInt("LARGE_VAR", 0)
		if result != 2147483647 {
			t.Errorf("getEnvInt() = %v, want %v", result, 2147483647)
		}
	})

	t.Run("小数点を含む値", func(t *testing.T) {
		_ = os.Setenv("FLOAT_VAR", "3.14")
		defer func() { _ = os.Unsetenv("FLOAT_VAR") }()

		result := getEnvInt("FLOAT_VAR", 10)
		if result != 10 {
			t.Errorf("getEnvInt() = %v, want %v (float should use default)", result, 10)
		}
	})

	t.Run("先頭に空白を含む値", func(t *testing.T) {
		_ = os.Setenv("SPACE_VAR", " 123")
		defer func() { _ = os.Unsetenv("SPACE_VAR") }()

		result := getEnvInt("SPACE_VAR", 10)
		if result != 10 {
			t.Errorf("getEnvInt() = %v, want %v (string with leading space should use default)", result, 10)
		}
	})
}

// ============================================
// 構造体テスト
// ============================================

func TestConfigStruct(t *testing.T) {
	t.Run("Config構造体の初期化", func(t *testing.T) {
		cfg := Config{
			Server: ServerConfig{
				Host: "localhost",
				Port: 8080,
			},
			Database: DatabaseConfig{
				URL: "postgres://localhost/test",
			},
			Log: LogConfig{
				Level: "debug",
			},
		}

		if cfg.Server.Host != "localhost" {
			t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "localhost")
		}
		if cfg.Server.Port != 8080 {
			t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 8080)
		}
		if cfg.Database.URL != "postgres://localhost/test" {
			t.Errorf("Database.URL = %v, want %v", cfg.Database.URL, "postgres://localhost/test")
		}
		if cfg.Log.Level != "debug" {
			t.Errorf("Log.Level = %v, want %v", cfg.Log.Level, "debug")
		}
	})

	t.Run("ゼロ値の構造体", func(t *testing.T) {
		var cfg Config

		if cfg.Server.Host != "" {
			t.Errorf("Server.Host zero value should be empty string")
		}
		if cfg.Server.Port != 0 {
			t.Errorf("Server.Port zero value should be 0")
		}
		if cfg.Database.URL != "" {
			t.Errorf("Database.URL zero value should be empty string")
		}
		if cfg.Log.Level != "" {
			t.Errorf("Log.Level zero value should be empty string")
		}
	})
}
