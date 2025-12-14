//go:build e2e

// Package e2e 勤務希望管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestRequestPeriodList 受付期間一覧テスト
func TestRequestPeriodList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("受付期間一覧ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/requests", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、受付期間一覧HTML
		_ = req
	})
}

// TestRequestPeriodCreate 受付期間作成テスト
func TestRequestPeriodCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("受付期間作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/requests/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、フォームHTML
		_ = req
	})

	t.Run("正常な受付期間作成", func(t *testing.T) {
		form := url.Values{}
		form.Add("target_year", "2025")
		form.Add("target_month", "2")
		form.Add("start_date", "2025-01-15")
		form.Add("end_date", "2025-01-25")
		form.Add("max_requests_per_staff", "10")

		req := httptest.NewRequest(http.MethodPost, "/requests", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト
		_ = req
	})

	t.Run("境界値_最大希望数0", func(t *testing.T) {
		form := url.Values{}
		form.Add("target_year", "2025")
		form.Add("target_month", "3")
		form.Add("start_date", "2025-02-15")
		form.Add("end_date", "2025-02-25")
		form.Add("max_requests_per_staff", "0")

		req := httptest.NewRequest(http.MethodPost, "/requests", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（0は無効）または許可
		_ = req
	})

	t.Run("終了日が開始日より前", func(t *testing.T) {
		form := url.Values{}
		form.Add("target_year", "2025")
		form.Add("target_month", "4")
		form.Add("start_date", "2025-03-25")
		form.Add("end_date", "2025-03-15") // 開始日より前

		req := httptest.NewRequest(http.MethodPost, "/requests", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})
}

// TestRequestPeriodDetail 受付期間詳細テスト
func TestRequestPeriodDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("受付期間詳細ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/requests/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、詳細HTML + 勤務希望一覧
		_ = req
	})
}

// TestRequestPeriodStatus 受付期間ステータス変更テスト
func TestRequestPeriodStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("受付開始", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/open", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、IsOpen=true
		_ = req
	})

	t.Run("受付終了", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/close", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、IsOpen=false
		_ = req
	})
}

// TestShiftRequestCreate 勤務希望作成テスト
func TestShiftRequestCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務希望作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/requests/xxx/entries/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、フォームHTML
		_ = req
	})

	t.Run("正常な勤務希望作成_希望シフト", func(t *testing.T) {
		form := url.Values{}
		form.Add("staff_id", "valid-staff-id")
		form.Add("target_date", "2025-02-10")
		form.Add("shift_type_id", "valid-shift-type-id")
		form.Add("request_type", "preferred")
		form.Add("priority", "required")

		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/entries", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 302 リダイレクト
		_ = req
	})

	t.Run("正常な勤務希望作成_回避", func(t *testing.T) {
		form := url.Values{}
		form.Add("staff_id", "valid-staff-id")
		form.Add("target_date", "2025-02-15")
		form.Add("request_type", "avoided")
		form.Add("priority", "optional")
		form.Add("comment", "この日は家族の予定があります")

		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/entries", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 302 リダイレクト
		_ = req
	})

	t.Run("対象期間外の日付", func(t *testing.T) {
		form := url.Values{}
		form.Add("staff_id", "valid-staff-id")
		form.Add("target_date", "2025-01-10") // 2月の期間なのに1月の日付
		form.Add("request_type", "preferred")

		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/entries", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("受付期間終了後の登録", func(t *testing.T) {
		// 受付期間がクローズされている場合
		form := url.Values{}
		form.Add("staff_id", "valid-staff-id")
		form.Add("target_date", "2025-02-10")
		form.Add("request_type", "preferred")

		req := httptest.NewRequest(http.MethodPost, "/requests/closed-period-id/entries", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 400 Bad Request（受付終了）
		_ = req
	})

	t.Run("希望上限超過", func(t *testing.T) {
		// 既に上限まで登録済みの場合
		form := url.Values{}
		form.Add("staff_id", "staff-with-max-requests")
		form.Add("target_date", "2025-02-20")
		form.Add("request_type", "preferred")

		req := httptest.NewRequest(http.MethodPost, "/requests/xxx/entries", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 400 Bad Request（上限超過）
		_ = req
	})
}

// TestShiftRequestDelete 勤務希望削除テスト
func TestShiftRequestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("正常な勤務希望削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/requests/entries/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("他人の勤務希望削除試行", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/requests/entries/other-user-entry", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token", // 別のユーザー
		})

		// 期待: 403 Forbidden（管理者以外は自分のみ削除可能）
		_ = req
	})

	t.Run("管理者による他人の勤務希望削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/requests/entries/other-user-entry", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK（管理者は削除可能）
		_ = req
	})
}
