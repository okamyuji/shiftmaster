#!/usr/bin/env bash
# test.sh - Go テストスクリプト
# すべてのテストをシャッフルモードで実行する

set -e

# 色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

echo ""
echo "=========================================="
echo "  ShiftMaster Test Suite"
echo "=========================================="
echo ""

print_step "Running tests with shuffle mode..."
echo ""

if go test -shuffle=on -count=1 ./... 2>&1; then
    echo ""
    print_success "すべてのテストがPassしました！"
    echo "=========================================="
    exit 0
else
    echo ""
    print_error "テストが失敗しました"
    echo "=========================================="
    exit 1
fi

