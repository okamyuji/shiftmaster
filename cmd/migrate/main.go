// Package main マイグレーションツール
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"shiftmaster/internal/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Load()

	// データベース接続
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.URL)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer func() {
		_ = db.Close()
	}()

	ctx := context.Background()

	// マイグレーションテーブル作成
	if err := createMigrationTable(ctx, db); err != nil {
		logger.Error("マイグレーションテーブル作成失敗", "error", err)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "up":
		if err := migrateUp(ctx, db, logger); err != nil {
			logger.Error("マイグレーション失敗", "error", err)
			os.Exit(1)
		}
	case "down":
		if err := migrateDown(ctx, db, logger); err != nil {
			logger.Error("ロールバック失敗", "error", err)
			os.Exit(1)
		}
	case "status":
		if err := showStatus(ctx, db, logger); err != nil {
			logger.Error("ステータス取得失敗", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// Migration マイグレーション履歴
type Migration struct {
	bun.BaseModel `bun:"table:schema_migrations"`
	Version       string `bun:"version,pk"`
	AppliedAt     string `bun:"applied_at"`
}

// createMigrationTable マイグレーション管理テーブル作成
func createMigrationTable(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

// getAppliedMigrations 適用済みマイグレーション取得
func getAppliedMigrations(ctx context.Context, db *bun.DB) (map[string]bool, error) {
	var migrations []Migration
	err := db.NewSelect().Model(&migrations).Scan(ctx)
	if err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, m := range migrations {
		applied[m.Version] = true
	}
	return applied, nil
}

// getMigrationFiles マイグレーションファイル一覧取得
func getMigrationFiles(suffix string) ([]string, error) {
	pattern := filepath.Join("migrations", fmt.Sprintf("*.%s.sql", suffix))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// extractVersion ファイル名からバージョン抽出
func extractVersion(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, ".")
	if len(parts) >= 2 {
		return parts[0]
	}
	return base
}

// migrateUp マイグレーション実行
func migrateUp(ctx context.Context, db *bun.DB, logger *slog.Logger) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	files, err := getMigrationFiles("up")
	if err != nil {
		return err
	}

	for _, file := range files {
		version := extractVersion(file)
		if applied[version] {
			continue
		}

		logger.Info("マイグレーション実行", "version", version)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("マイグレーション失敗 %s: %w", version, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		logger.Info("マイグレーション完了", "version", version)
	}

	return nil
}

// migrateDown ロールバック実行
func migrateDown(ctx context.Context, db *bun.DB, logger *slog.Logger) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	files, err := getMigrationFiles("down")
	if err != nil {
		return err
	}

	// 逆順でロールバック
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	for _, file := range files {
		version := extractVersion(file)
		if !applied[version] {
			continue
		}

		logger.Info("ロールバック実行", "version", version)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("ロールバック失敗 %s: %w", version, err)
		}

		if _, err := tx.ExecContext(ctx, "DELETE FROM schema_migrations WHERE version = ?", version); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		logger.Info("ロールバック完了", "version", version)
		break // 1つずつロールバック
	}

	return nil
}

// showStatus マイグレーションステータス表示
func showStatus(ctx context.Context, db *bun.DB, logger *slog.Logger) error {
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	files, err := getMigrationFiles("up")
	if err != nil {
		return err
	}

	logger.Info("マイグレーションステータス表示")
	fmt.Println("Migration Status:")
	fmt.Println("=================")
	for _, file := range files {
		version := extractVersion(file)
		status := "pending"
		if applied[version] {
			status = "applied"
		}
		logger.Debug("マイグレーションステータス", "version", version, "status", status)
		fmt.Printf("  %s: %s\n", version, status)
	}

	return nil
}
