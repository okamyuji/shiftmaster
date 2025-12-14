//go:build e2e

// Package e2e ユーザー管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestUserList ユーザー一覧テスト
func TestUserList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("ユーザー一覧ページ表示_管理者", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK、ユーザー一覧HTML
		_ = req
	})

	t.Run("一般ユーザーがユーザー一覧にアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("未認証でユーザー一覧アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestUserCreate ユーザー作成テスト
func TestUserCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("ユーザー作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK、フォームHTML
		_ = req
	})

	t.Run("正常なユーザー作成_管理者", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "newadmin@example.com")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "新規")
		form.Add("last_name", "管理者")
		form.Add("role", "admin")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 302 リダイレクト to /admin/users
		_ = req
	})

	t.Run("正常なユーザー作成_マネージャー", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "newmanager@example.com")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "新規")
		form.Add("last_name", "マネージャー")
		form.Add("role", "manager")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 302 リダイレクト to /admin/users
		_ = req
	})

	t.Run("正常なユーザー作成_一般ユーザー", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "newuser@example.com")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "新規")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 302 リダイレクト to /admin/users
		_ = req
	})

	t.Run("重複するメールアドレス", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com") // 既存
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "重複")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（重複エラー）
		_ = req
	})

	t.Run("必須項目未入力_メールアドレス", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "テスト")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("必須項目未入力_パスワード", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "testuser@example.com")
		form.Add("password", "")
		form.Add("first_name", "テスト")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("不正なメールアドレス形式", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "invalid-email")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "テスト")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("弱いパスワード", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "weakpass@example.com")
		form.Add("password", "123") // 弱すぎる
		form.Add("first_name", "テスト")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（パスワードポリシー違反）
		_ = req
	})

	t.Run("不正なロール", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "invalidrole@example.com")
		form.Add("password", "SecurePass123!")
		form.Add("first_name", "テスト")
		form.Add("last_name", "ユーザー")
		form.Add("role", "invalid_role")

		req := httptest.NewRequest(http.MethodPost, "/admin/users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（不正なロール）
		_ = req
	})
}

// TestUserEdit ユーザー編集テスト
func TestUserEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("ユーザー編集フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users/xxx/edit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK、現在値が入力されたフォームHTML
		_ = req
	})

	t.Run("ユーザー編集成功", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "updated@example.com")
		form.Add("first_name", "更新")
		form.Add("last_name", "ユーザー")
		form.Add("role", "manager")

		req := httptest.NewRequest(http.MethodPut, "/admin/users/xxx", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 302 リダイレクト to /admin/users
		_ = req
	})

	t.Run("存在しないユーザー編集", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "notfound@example.com")
		form.Add("first_name", "不存在")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user")

		req := httptest.NewRequest(http.MethodPut, "/admin/users/non-existent-id", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})

	t.Run("自分自身の管理者権限を剥奪", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com")
		form.Add("first_name", "管理者")
		form.Add("last_name", "ユーザー")
		form.Add("role", "user") // 自分の権限を下げる

		req := httptest.NewRequest(http.MethodPut, "/admin/users/current-admin-id", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（自分自身の権限変更禁止）または許可
		_ = req
	})
}

// TestUserDelete ユーザー削除テスト
func TestUserDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("ユーザー削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/users/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("自分自身を削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/users/current-admin-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（自分自身は削除不可）
		_ = req
	})

	t.Run("マネージャーがユーザー削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/users/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 403 Forbidden（管理者のみ削除可能）
		_ = req
	})

	t.Run("存在しないユーザー削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/users/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}
