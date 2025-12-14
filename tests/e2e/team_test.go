//go:build e2e

// Package e2e チーム管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestTeamList チーム一覧テスト
func TestTeamList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("チーム一覧ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/teams", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、チーム一覧HTML
		_ = req
	})

	t.Run("未認証でチーム一覧アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/teams", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestTeamCreate チーム作成テスト
func TestTeamCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("チーム作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/teams/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、部署選択付きフォームHTML
		_ = req
	})

	t.Run("正常なチーム作成", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "看護2チーム")
		form.Add("department_id", "valid-department-id")
		form.Add("description", "看護部の2番目のチーム")

		req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /teams
		_ = req
	})

	t.Run("チーム名未入力", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "")
		form.Add("department_id", "valid-department-id")

		req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("部署未選択", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "チーム名")
		form.Add("department_id", "")

		req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("存在しない部署ID", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "チーム名")
		form.Add("department_id", "invalid-department-id")

		req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（外部キー違反）
		_ = req
	})

	t.Run("同一部署内でチーム名重複", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "看護1チーム") // 既存
		form.Add("department_id", "valid-department-id")

		req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（重複エラー）
		_ = req
	})
}

// TestTeamEdit チーム編集テスト
func TestTeamEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("チーム編集フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/teams/xxx/edit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、現在値が入力されたフォームHTML
		_ = req
	})

	t.Run("チーム編集成功", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "看護1チーム更新")
		form.Add("department_id", "valid-department-id")
		form.Add("description", "更新された説明")

		req := httptest.NewRequest(http.MethodPut, "/teams/xxx", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /teams
		_ = req
	})

	t.Run("存在しないチーム編集", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "不存在チーム")
		form.Add("department_id", "valid-department-id")

		req := httptest.NewRequest(http.MethodPut, "/teams/non-existent-id", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestTeamDelete チーム削除テスト
func TestTeamDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("チーム削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/teams/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("スタッフが所属するチーム削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/teams/team-with-staff", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（関連データ存在）または CASCADE削除
		_ = req
	})

	t.Run("一般ユーザーによるチーム削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/teams/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("存在しないチーム削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/teams/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}
