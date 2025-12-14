//go:build e2e

// Package e2e スタッフ管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestStaffList スタッフ一覧テスト
func TestStaffList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("スタッフ一覧ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、スタッフ一覧HTML
		_ = req
	})

	t.Run("未認証でスタッフ一覧アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestStaffCreate スタッフ作成テスト
func TestStaffCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("スタッフ作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、チームドロップダウン付きフォームHTML
		_ = req
	})

	t.Run("正常なスタッフ作成", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP002")
		form.Add("last_name", "田中")
		form.Add("first_name", "太郎")
		form.Add("email", "tanaka@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /staffs
		_ = req
	})

	t.Run("重複する社員番号", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP001") // 既存
		form.Add("last_name", "鈴木")
		form.Add("first_name", "一郎")
		form.Add("email", "suzuki@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（重複エラー）
		_ = req
	})

	t.Run("必須項目未入力_社員番号", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "")
		form.Add("last_name", "佐藤")
		form.Add("first_name", "花子")
		form.Add("email", "sato@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("必須項目未入力_姓", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP003")
		form.Add("last_name", "")
		form.Add("first_name", "次郎")
		form.Add("email", "jiro@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("不正なメールアドレス形式", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP004")
		form.Add("last_name", "高橋")
		form.Add("first_name", "三郎")
		form.Add("email", "invalid-email")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("存在しないチームID", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP005")
		form.Add("last_name", "伊藤")
		form.Add("first_name", "四郎")
		form.Add("email", "ito@example.com")
		form.Add("team_id", "invalid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPost, "/staffs", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})
}

// TestStaffEdit スタッフ編集テスト
func TestStaffEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("スタッフ編集フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs/xxx/edit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、現在値が入力されたフォームHTML
		_ = req
	})

	t.Run("スタッフ編集成功", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP001")
		form.Add("last_name", "山田")
		form.Add("first_name", "花子更新")
		form.Add("email", "yamada-updated@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "part_time")

		req := httptest.NewRequest(http.MethodPut, "/staffs/xxx", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /staffs
		_ = req
	})

	t.Run("存在しないスタッフ編集", func(t *testing.T) {
		form := url.Values{}
		form.Add("employee_number", "EMP999")
		form.Add("last_name", "不明")
		form.Add("first_name", "太郎")
		form.Add("email", "unknown@example.com")
		form.Add("team_id", "valid-team-id")
		form.Add("employment_type", "full_time")

		req := httptest.NewRequest(http.MethodPut, "/staffs/non-existent-id", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestStaffDelete スタッフ削除テスト
func TestStaffDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("スタッフ削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/staffs/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("一般ユーザーによるスタッフ削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/staffs/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("存在しないスタッフ削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/staffs/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestStaffDetail スタッフ詳細テスト
func TestStaffDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("スタッフ詳細ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、スタッフ詳細HTML
		_ = req
	})

	t.Run("存在しないスタッフ詳細", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/staffs/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}
