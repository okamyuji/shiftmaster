#!/usr/bin/env bash
# setup-hooks.sh - Git hooksのセットアップスクリプト
# プロジェクトのセットアップ時に一度実行してください

set -e

# 色定義
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

echo ""
echo "=========================================="
echo "  Git Hooks セットアップ"
echo "=========================================="
echo ""

# Git hooksディレクトリを設定
print_step "Git hooksディレクトリを設定中..."
git config core.hooksPath .githooks
print_success "core.hooksPath を .githooks に設定しました"

# 実行権限を付与
print_step "実行権限を付与中..."
chmod +x .githooks/pre-commit
chmod +x scripts/lint.sh
chmod +x scripts/test.sh
chmod +x scripts/check-all.sh
chmod +x scripts/setup-hooks.sh
print_success "実行権限を付与しました"

echo ""
echo "=========================================="
print_success "セットアップ完了！"
echo ""
echo "以下のコマンドが利用可能です:"
echo "  ./scripts/lint.sh      - Lintチェックのみ実行"
echo "  ./scripts/test.sh      - テストのみ実行"
echo "  ./scripts/check-all.sh - Lint + テストを実行"
echo ""
echo "git commit 時に自動的にチェックが実行されます。"
echo "チェックをスキップするには: git commit --no-verify"
echo "=========================================="

