// Package main ShiftMaster サーバーエントリポイント
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shiftmaster/internal/config"
	"shiftmaster/internal/di"
	"shiftmaster/internal/web"
)

func main() {
	// ロガー初期化
	logger := setupLogger()

	// 設定読み込み
	cfg := config.Load()

	logger.Info("ShiftMaster 起動中",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)

	// DIコンテナ初期化
	container, err := di.NewContainer(cfg, logger)
	if err != nil {
		logger.Error("DIコンテナ初期化失敗", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			logger.Error("リソース解放失敗", "error", err)
		}
	}()

	// ミドルウェアチェーン
	handler := web.Chain(
		container.Router,
		web.SecurityHeaders(),
		web.CORS([]string{"*"}),
		web.Recover(logger),
		web.Logger(logger),
	)

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// グレースフルシャットダウン
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("シャットダウン開始")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("シャットダウン失敗", "error", err)
		}
		close(done)
	}()

	// サーバー起動
	logger.Info("サーバー起動完了", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("サーバー起動失敗", "error", err)
		os.Exit(1)
	}

	<-done
	logger.Info("シャットダウン完了")
}

// setupLogger slogロガー初期化
func setupLogger() *slog.Logger {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
