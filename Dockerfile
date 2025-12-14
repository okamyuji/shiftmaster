# ShiftMaster Dockerfile
# マルチステージビルド

# ビルドステージ
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 依存関係ファイルコピー
COPY go.mod go.sum ./

# 依存関係ダウンロード
RUN go mod download

# ソースコードコピー
COPY . .

# バイナリビルド
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# 開発環境ステージ
FROM golang:1.25-alpine AS development

WORKDIR /app

# 開発ツールインストール
RUN go install github.com/air-verse/air@latest

# ソースコードマウント用 実際のソースはボリュームでマウント
COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

# 本番環境ステージ
FROM alpine:3.20 AS production

WORKDIR /app

# セキュリティ用 非rootユーザー作成
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# 必要なファイルのみコピー
COPY --from=builder /app/server /app/server
COPY --from=builder /app/static /app/static
COPY --from=builder /app/internal/web/templates /app/internal/web/templates

# 所有者変更
RUN chown -R appuser:appgroup /app

# 非rootユーザーに切り替え
USER appuser

# ポート公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# サーバー起動
CMD ["/app/server"]

