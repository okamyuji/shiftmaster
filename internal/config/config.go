// Package config アプリケーション設定管理
package config

import (
	"os"
	"strconv"
)

// Config アプリケーション設定
type Config struct {
	// Server サーバー設定
	Server ServerConfig
	// Database データベース設定
	Database DatabaseConfig
	// Log ログ設定
	Log LogConfig
}

// ServerConfig サーバー設定
type ServerConfig struct {
	// Host バインドホスト
	Host string
	// Port リッスンポート
	Port int
}

// DatabaseConfig データベース設定
type DatabaseConfig struct {
	// URL 接続URL
	URL string
}

// LogConfig ログ設定
type LogConfig struct {
	// Level ログレベル debug info warn error
	Level string
}

// Load 環境変数から設定を読み込む
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://shiftmaster:shiftmaster_secret@localhost:5432/shiftmaster?sslmode=disable"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv 環境変数取得 デフォルト値付き
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 環境変数取得 整数型
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
