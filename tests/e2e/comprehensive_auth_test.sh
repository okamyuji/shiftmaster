#!/bin/bash

# ============================================
# ShiftMaster 包括的認証・認可テスト
# ============================================

set -e

BASE_URL="http://localhost:8080"
PASS_COUNT=0
FAIL_COUNT=0

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ログイン関数
login() {
    local email=$1
    local password=$2
    local response
    response=$(curl -s -X POST "$BASE_URL/api/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"password\":\"$password\"}")
    echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4
}

# テスト関数
test_request() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local token=$4
    local expected_status=$5
    local data=$6

    local status
    if [ -n "$data" ]; then
        status=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        status=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $token")
    fi

    if [ "$status" == "$expected_status" ]; then
        echo -e "${GREEN}✓${NC} $test_name (status: $status)"
        ((PASS_COUNT++))
    else
        echo -e "${RED}✗${NC} $test_name (expected: $expected_status, got: $status)"
        ((FAIL_COUNT++))
    fi
}

# ヘッダー表示
echo -e "\n${BLUE}============================================${NC}"
echo -e "${BLUE}  ShiftMaster 包括的認証・認可テスト${NC}"
echo -e "${BLUE}============================================${NC}\n"

# ============================================
# 1. 認証テスト（正常系）
# ============================================
echo -e "${YELLOW}=== 1. 認証テスト（正常系） ===${NC}"

# super_admin ログイン
SUPER_ADMIN_TOKEN=$(login "superadmin@example.com" "SuperAdmin123\$!")
if [ -n "$SUPER_ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} super_admin ログイン成功"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} super_admin ログイン失敗"
    ((FAIL_COUNT++))
fi

# admin ログイン (サンプル病院)
ADMIN_TOKEN=$(login "admin@example.com" "Password123\$!")
if [ -n "$ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} admin ログイン成功"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} admin ログイン失敗"
    ((FAIL_COUNT++))
fi

# manager ログイン
MANAGER_TOKEN=$(login "manager@example.com" "Manager123\$!")
if [ -n "$MANAGER_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} manager ログイン成功"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} manager ログイン失敗"
    ((FAIL_COUNT++))
fi

# user ログイン
USER_TOKEN=$(login "user@example.com" "User123\$!")
if [ -n "$USER_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} user ログイン成功"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} user ログイン失敗"
    ((FAIL_COUNT++))
fi

# other org admin ログイン (テスト医療センター)
OTHER_ADMIN_TOKEN=$(login "otheradmin@example.com" "OtherAdmin123\$!")
if [ -n "$OTHER_ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} other org admin ログイン成功"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} other org admin ログイン失敗"
    ((FAIL_COUNT++))
fi

# ============================================
# 2. 認証テスト（異常系）
# ============================================
echo -e "\n${YELLOW}=== 2. 認証テスト（異常系） ===${NC}"

# 不正なパスワード
INVALID_TOKEN=$(login "admin@example.com" "WrongPassword!")
if [ -z "$INVALID_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} 不正なパスワードでログイン拒否"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} 不正なパスワードでログインが成功した"
    ((FAIL_COUNT++))
fi

# 存在しないユーザー
NONEXIST_TOKEN=$(login "nonexist@example.com" "Password123\$!")
if [ -z "$NONEXIST_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} 存在しないユーザーでログイン拒否"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} 存在しないユーザーでログインが成功した"
    ((FAIL_COUNT++))
fi

# 空のメールアドレス
EMPTY_EMAIL_TOKEN=$(login "" "Password123\$!")
if [ -z "$EMPTY_EMAIL_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} 空のメールアドレスでログイン拒否"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} 空のメールアドレスでログインが成功した"
    ((FAIL_COUNT++))
fi

# 空のパスワード
EMPTY_PASS_TOKEN=$(login "admin@example.com" "")
if [ -z "$EMPTY_PASS_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} 空のパスワードでログイン拒否"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} 空のパスワードでログインが成功した"
    ((FAIL_COUNT++))
fi

# ============================================
# 3. super_admin 認可テスト
# ============================================
echo -e "\n${YELLOW}=== 3. super_admin 認可テスト ===${NC}"

test_request "super_admin: ダッシュボードアクセス" "GET" "/" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: ユーザー管理アクセス" "GET" "/admin/users" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: スタッフ一覧アクセス" "GET" "/staffs" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: シフト一覧アクセス" "GET" "/shifts" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: チーム一覧アクセス" "GET" "/teams" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: 勤務表一覧アクセス" "GET" "/schedules" "$SUPER_ADMIN_TOKEN" "200"
test_request "super_admin: /api/auth/me アクセス" "GET" "/api/auth/me" "$SUPER_ADMIN_TOKEN" "200"

# ============================================
# 4. admin 認可テスト（サンプル病院）
# ============================================
echo -e "\n${YELLOW}=== 4. admin 認可テスト（サンプル病院） ===${NC}"

test_request "admin: ダッシュボードアクセス" "GET" "/" "$ADMIN_TOKEN" "200"
test_request "admin: ユーザー管理アクセス" "GET" "/admin/users" "$ADMIN_TOKEN" "200"
test_request "admin: スタッフ一覧アクセス" "GET" "/staffs" "$ADMIN_TOKEN" "200"
test_request "admin: シフト一覧アクセス" "GET" "/shifts" "$ADMIN_TOKEN" "200"
test_request "admin: チーム一覧アクセス" "GET" "/teams" "$ADMIN_TOKEN" "200"
test_request "admin: 勤務表一覧アクセス" "GET" "/schedules" "$ADMIN_TOKEN" "200"
test_request "admin: /api/auth/me アクセス" "GET" "/api/auth/me" "$ADMIN_TOKEN" "200"

# ============================================
# 5. manager 認可テスト
# ============================================
echo -e "\n${YELLOW}=== 5. manager 認可テスト ===${NC}"

test_request "manager: ダッシュボードアクセス" "GET" "/" "$MANAGER_TOKEN" "200"
test_request "manager: ユーザー管理アクセス拒否" "GET" "/admin/users" "$MANAGER_TOKEN" "403"
test_request "manager: スタッフ一覧アクセス" "GET" "/staffs" "$MANAGER_TOKEN" "200"
test_request "manager: シフト一覧アクセス" "GET" "/shifts" "$MANAGER_TOKEN" "200"
test_request "manager: チーム一覧アクセス" "GET" "/teams" "$MANAGER_TOKEN" "200"
test_request "manager: 勤務表一覧アクセス" "GET" "/schedules" "$MANAGER_TOKEN" "200"
test_request "manager: /api/auth/me アクセス" "GET" "/api/auth/me" "$MANAGER_TOKEN" "200"

# ============================================
# 6. user 認可テスト
# ============================================
echo -e "\n${YELLOW}=== 6. user 認可テスト ===${NC}"

test_request "user: ダッシュボードアクセス" "GET" "/" "$USER_TOKEN" "200"
test_request "user: ユーザー管理アクセス拒否" "GET" "/admin/users" "$USER_TOKEN" "403"
test_request "user: スタッフ一覧アクセス" "GET" "/staffs" "$USER_TOKEN" "200"
test_request "user: シフト一覧アクセス" "GET" "/shifts" "$USER_TOKEN" "200"
test_request "user: チーム一覧アクセス" "GET" "/teams" "$USER_TOKEN" "200"
test_request "user: 勤務表一覧アクセス" "GET" "/schedules" "$USER_TOKEN" "200"
test_request "user: /api/auth/me アクセス" "GET" "/api/auth/me" "$USER_TOKEN" "200"

# ============================================
# 7. 他組織admin 認可テスト（テスト医療センター）
# ============================================
echo -e "\n${YELLOW}=== 7. 他組織admin 認可テスト（テスト医療センター） ===${NC}"

test_request "other admin: ダッシュボードアクセス" "GET" "/" "$OTHER_ADMIN_TOKEN" "200"
test_request "other admin: ユーザー管理アクセス" "GET" "/admin/users" "$OTHER_ADMIN_TOKEN" "200"
test_request "other admin: /api/auth/me アクセス" "GET" "/api/auth/me" "$OTHER_ADMIN_TOKEN" "200"

# ============================================
# 8. 未認証アクセステスト
# ============================================
echo -e "\n${YELLOW}=== 8. 未認証アクセステスト ===${NC}"

test_request "未認証: ダッシュボードアクセス拒否" "GET" "/" "" "401"
test_request "未認証: ユーザー管理アクセス拒否" "GET" "/admin/users" "" "401"
test_request "未認証: スタッフ一覧アクセス拒否" "GET" "/staffs" "" "401"
test_request "未認証: /api/auth/me アクセス拒否" "GET" "/api/auth/me" "" "401"

# ============================================
# 9. 無効なトークンテスト
# ============================================
echo -e "\n${YELLOW}=== 9. 無効なトークンテスト ===${NC}"

test_request "無効トークン: ダッシュボードアクセス拒否" "GET" "/" "invalid_token_here" "401"
test_request "無効トークン: /api/auth/me アクセス拒否" "GET" "/api/auth/me" "invalid_token_here" "401"

# ============================================
# 10. 境界値テスト
# ============================================
echo -e "\n${YELLOW}=== 10. 境界値テスト ===${NC}"

# 非常に長いトークン
LONG_TOKEN=$(printf 'a%.0s' {1..10000})
test_request "長いトークン: アクセス拒否" "GET" "/" "$LONG_TOKEN" "401"

# 特殊文字を含むトークン
test_request "特殊文字トークン: アクセス拒否" "GET" "/" "token<>\"'&test" "401"

# ============================================
# 結果サマリー
# ============================================
echo -e "\n${BLUE}============================================${NC}"
echo -e "${BLUE}  テスト結果サマリー${NC}"
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}Pass: $PASS_COUNT${NC}"
echo -e "${RED}Fail: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ すべてのテストがPassしました！${NC}"
    exit 0
else
    echo -e "${RED}✗ $FAIL_COUNT 個のテストが失敗しました${NC}"
    exit 1
fi
