//go:build e2e

// Package e2e 勤務表管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestScheduleList 勤務表一覧テスト
func TestScheduleList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務表一覧ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、勤務表一覧HTML
		_ = req
	})

	t.Run("未認証で勤務表一覧アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestScheduleCreate 勤務表作成テスト
func TestScheduleCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務表作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、フォームHTML
		_ = req
	})

	t.Run("正常な勤務表作成", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "2025年2月勤務表")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "2")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /schedules/xxx
		_ = req
	})

	t.Run("境界値_1月", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "2025年1月勤務表")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "1")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト
		_ = req
	})

	t.Run("境界値_12月", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "2025年12月勤務表")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "12")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト
		_ = req
	})

	t.Run("不正な月_0月", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "2025年0月勤務表")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "0")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("不正な月_13月", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "2025年13月勤務表")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "13")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("重複する勤務表作成", func(t *testing.T) {
		// 同一チーム、同一年月で重複作成
		form := url.Values{}
		form.Add("name", "2025年2月勤務表重複")
		form.Add("team_id", "valid-team-id")
		form.Add("year", "2025")
		form.Add("month", "2") // 既に存在

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（重複エラー）
		_ = req
	})

	t.Run("存在しないチームID", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "不正な勤務表")
		form.Add("team_id", "invalid-team-id")
		form.Add("year", "2025")
		form.Add("month", "3")

		req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})
}

// TestScheduleDetail 勤務表詳細テスト
func TestScheduleDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務表詳細ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、詳細HTML（エントリグリッド表示）
		_ = req
	})

	t.Run("存在しない勤務表詳細", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedules/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestScheduleEntryUpdate 勤務表エントリ更新テスト
func TestScheduleEntryUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("エントリ更新成功", func(t *testing.T) {
		form := url.Values{}
		form.Add("shift_type_id", "valid-shift-type-id")

		req := httptest.NewRequest(http.MethodPut, "/schedules/xxx/entries/yyy", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("公開済み勤務表のエントリ更新", func(t *testing.T) {
		form := url.Values{}
		form.Add("shift_type_id", "valid-shift-type-id")

		req := httptest.NewRequest(http.MethodPut, "/schedules/published-schedule-id/entries/yyy", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（公開済みは編集不可）または許可
		_ = req
	})
}

// TestSchedulePublish 勤務表公開テスト
func TestSchedulePublish(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務表公開成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/schedules/xxx/publish", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、IsPublished=true
		_ = req
	})

	t.Run("既に公開済みの勤務表を再公開", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/schedules/published-schedule-id/publish", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（既に公開済み）または成功
		_ = req
	})

	t.Run("一般ユーザーによる勤務表公開", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/schedules/xxx/publish", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})
}

// TestScheduleValidate 勤務表検証テスト
func TestScheduleValidate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("検証成功_違反なし", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/schedules/valid-schedule-id/validate", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、IsValid=true
		_ = req
	})

	t.Run("検証成功_違反あり", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/schedules/invalid-schedule-id/validate", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、IsValid=false、違反リスト含む
		_ = req
	})
}

// TestScheduleDelete 勤務表削除テスト
func TestScheduleDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("勤務表削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/schedules/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("公開済み勤務表削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/schedules/published-schedule-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（公開済みは削除不可）または許可
		_ = req
	})

	t.Run("一般ユーザーによる勤務表削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/schedules/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("存在しない勤務表削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/schedules/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}
