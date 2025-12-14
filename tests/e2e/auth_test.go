//go:build e2e

// Package e2e 認証E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestLoginPage ログインページ表示テスト
func TestLoginPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("ログインページ正常表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)

		// 期待: 200 OK、ログインフォームHTML
		_ = req
	})

	t.Run("認証済みユーザーはダッシュボードへリダイレクト", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 302 リダイレクト to /
		_ = req
	})
}

// TestLogin ログイン処理テスト
func TestLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("正常なログイン", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com")
		form.Add("password", "Password123$!")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 302 リダイレクト to /, Cookieにトークン設定
		_ = req
	})

	t.Run("不正なメールアドレス", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "invalid@example.com")
		form.Add("password", "Password123$!")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 401 Unauthorized
		_ = req
	})

	t.Run("不正なパスワード", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com")
		form.Add("password", "wrongpassword")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 401 Unauthorized
		_ = req
	})

	t.Run("空のメールアドレス", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "")
		form.Add("password", "Password123$!")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("空のパスワード", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com")
		form.Add("password", "")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("SQLインジェクション試行", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", "admin@example.com' OR '1'='1")
		form.Add("password", "Password123$!")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 期待: 401 Unauthorized（インジェクション無効化）
		_ = req
	})
}

// TestLogout ログアウト処理テスト
func TestLogout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("正常なログアウト", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 302 リダイレクト to /login, Cookie削除
		_ = req
	})

	t.Run("未認証でのログアウト", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestAuthMiddleware 認証ミドルウェアテスト
func TestAuthMiddleware(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("未認証でダッシュボードアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})

	t.Run("無効なトークンでダッシュボードアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "invalid_token",
		})

		// 期待: 302 リダイレクト to /login
		_ = req
	})

	t.Run("期限切れトークンでダッシュボードアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "expired_token",
		})

		// 期待: 302 リダイレクト to /login
		_ = req
	})

	t.Run("有効なトークンでダッシュボードアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、ダッシュボードHTML
		_ = req
	})
}

// TestRoleMiddleware ロールミドルウェアテスト
func TestRoleMiddleware(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("一般ユーザーが管理者ページアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token", // roleがuser
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("管理者が管理者ページアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token", // roleがadmin
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("マネージャーが勤務表管理アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token", // roleがmanager
		})

		// 期待: 200 OK
		_ = req
	})
}
