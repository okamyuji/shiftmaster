#!/usr/bin/env bash
# check-all.sh - 全チェックスクリプト（lint + test）
# pre-commitフックから呼び出される

set -e

# 色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}"
    echo "╔════════════════════════════════════════╗"
    echo "║     ShiftMaster Pre-Commit Check       ║"
    echo "╚════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

print_header

ERRORS=0

# 1. Lint チェック
print_step "Step 1/2: Lint チェック実行中..."
echo ""
if ./scripts/lint.sh; then
    print_success "Lint チェック完了"
else
    print_error "Lint チェック失敗"
    ERRORS=$((ERRORS + 1))
fi

echo ""

# 2. テスト実行
print_step "Step 2/2: テスト実行中..."
echo ""
if ./scripts/test.sh; then
    print_success "テスト完了"
else
    print_error "テスト失敗"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "=========================================="

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════╗"
    echo "║   ✓ すべてのチェックがPassしました！   ║"
    echo "║     コミットを続行します               ║"
    echo "╚════════════════════════════════════════╝"
    echo -e "${NC}"
    exit 0
else
    echo -e "${RED}"
    echo "╔════════════════════════════════════════╗"
    echo "║   ✗ チェックが失敗しました             ║"
    echo "║     コミットを中止します               ║"
    echo "╚════════════════════════════════════════╝"
    echo -e "${NC}"
    echo ""
    echo "修正後、再度コミットしてください。"
    echo "チェックをスキップするには: git commit --no-verify"
    exit 1
fi

