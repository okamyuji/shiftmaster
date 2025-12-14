#!/bin/bash
# lint.sh - Go Formatter/Linter スクリプト
# すべてのFormatter/Linterを実行し、結果を表示する

set -e

# 色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 結果追跡
ERRORS=0

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}!${NC} $1"
}

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

echo ""
echo "=========================================="
echo "  ShiftMaster Lint & Format Check"
echo "=========================================="
echo ""

# 1. go fmt
print_step "Running go fmt..."
FMT_OUTPUT=$(gofmt -l . 2>&1 || true)
if [ -n "$FMT_OUTPUT" ]; then
    print_error "go fmt: フォーマットが必要なファイルがあります"
    echo "$FMT_OUTPUT"
    ERRORS=$((ERRORS + 1))
else
    print_success "go fmt: OK"
fi

# 2. go vet
print_step "Running go vet..."
if go vet ./... 2>&1; then
    print_success "go vet: OK"
else
    print_error "go vet: 問題が検出されました"
    ERRORS=$((ERRORS + 1))
fi

# 3. staticcheck (インストールされている場合)
print_step "Running staticcheck..."
if command -v staticcheck &> /dev/null; then
    if staticcheck ./... 2>&1; then
        print_success "staticcheck: OK"
    else
        print_error "staticcheck: 問題が検出されました"
        ERRORS=$((ERRORS + 1))
    fi
else
    print_warning "staticcheck: インストールされていません (go install honnef.co/go/tools/cmd/staticcheck@latest)"
fi

# 4. golangci-lint
print_step "Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run ./... 2>&1; then
        print_success "golangci-lint: OK"
    else
        print_error "golangci-lint: 問題が検出されました"
        ERRORS=$((ERRORS + 1))
    fi
else
    print_warning "golangci-lint: インストールされていません (https://golangci-lint.run/usage/install/)"
fi

# 5. go mod tidy チェック
print_step "Checking go mod tidy..."
cp go.mod go.mod.backup
cp go.sum go.sum.backup 2>/dev/null || true
go mod tidy
if ! diff -q go.mod go.mod.backup > /dev/null 2>&1; then
    print_error "go mod tidy: go.mod に変更が必要です"
    mv go.mod.backup go.mod
    mv go.sum.backup go.sum 2>/dev/null || true
    ERRORS=$((ERRORS + 1))
else
    print_success "go mod tidy: OK"
    rm go.mod.backup
    rm go.sum.backup 2>/dev/null || true
fi

# 6. go build
print_step "Running go build..."
if go build ./... 2>&1; then
    print_success "go build: OK"
else
    print_error "go build: ビルドに失敗しました"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "=========================================="

if [ $ERRORS -eq 0 ]; then
    print_success "すべてのチェックがPassしました！"
    echo "=========================================="
    exit 0
else
    print_error "$ERRORS 個のチェックが失敗しました"
    echo "=========================================="
    exit 1
fi

