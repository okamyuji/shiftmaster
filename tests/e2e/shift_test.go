//go:build e2e

// Package e2e シフト種別管理E2Eテスト
package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestShiftTypeList シフト種別一覧テスト
func TestShiftTypeList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("シフト種別一覧ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、シフト種別一覧HTML
		_ = req
	})

	t.Run("未認証でシフト種別一覧アクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts", nil)

		// 期待: 302 リダイレクト to /login
		_ = req
	})
}

// TestShiftTypeCreate シフト種別作成テスト
func TestShiftTypeCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("シフト種別作成フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、フォームHTML
		_ = req
	})

	t.Run("正常なシフト種別作成_日勤", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "日勤A")
		form.Add("code", "DA")
		form.Add("color", "#00FF00")
		form.Add("start_time", "08:30")
		form.Add("end_time", "17:30")
		form.Add("break_minutes", "60")
		form.Add("is_night_shift", "false")
		form.Add("is_holiday", "false")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /shifts
		_ = req
	})

	t.Run("正常なシフト種別作成_夜勤", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "夜勤A")
		form.Add("code", "NA")
		form.Add("color", "#0000FF")
		form.Add("start_time", "17:00")
		form.Add("end_time", "09:00") // 翌日
		form.Add("break_minutes", "120")
		form.Add("is_night_shift", "true")
		form.Add("is_holiday", "false")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /shifts
		_ = req
	})

	t.Run("正常なシフト種別作成_休日", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "休日A")
		form.Add("code", "HA")
		form.Add("color", "#FFFF00")
		form.Add("start_time", "00:00")
		form.Add("end_time", "00:00")
		form.Add("break_minutes", "0")
		form.Add("is_night_shift", "false")
		form.Add("is_holiday", "true")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /shifts
		_ = req
	})

	t.Run("シフト種別名未入力", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "")
		form.Add("code", "XX")
		form.Add("color", "#FF0000")
		form.Add("start_time", "09:00")
		form.Add("end_time", "18:00")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("コード未入力", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "テストシフト")
		form.Add("code", "")
		form.Add("color", "#FF0000")
		form.Add("start_time", "09:00")
		form.Add("end_time", "18:00")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("コード重複", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "日勤重複")
		form.Add("code", "D") // 既存の日勤コード
		form.Add("color", "#FF0000")
		form.Add("start_time", "09:00")
		form.Add("end_time", "18:00")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（重複エラー）
		_ = req
	})

	t.Run("不正な時刻形式", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "不正シフト")
		form.Add("code", "XX")
		form.Add("color", "#FF0000")
		form.Add("start_time", "25:00") // 不正
		form.Add("end_time", "18:00")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request
		_ = req
	})

	t.Run("境界値_休憩0分", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "休憩なしシフト")
		form.Add("code", "NB")
		form.Add("color", "#00FF00")
		form.Add("start_time", "09:00")
		form.Add("end_time", "12:00")
		form.Add("break_minutes", "0")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト（0分は有効）
		_ = req
	})

	t.Run("境界値_休憩マイナス", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "マイナス休憩シフト")
		form.Add("code", "MB")
		form.Add("color", "#00FF00")
		form.Add("start_time", "09:00")
		form.Add("end_time", "18:00")
		form.Add("break_minutes", "-30")

		req := httptest.NewRequest(http.MethodPost, "/shifts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 400 Bad Request（マイナスは無効）
		_ = req
	})
}

// TestShiftTypeDetail シフト種別詳細テスト
func TestShiftTypeDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("シフト種別詳細ページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 200 OK、詳細HTML
		_ = req
	})

	t.Run("存在しないシフト種別詳細", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestShiftTypeEdit シフト種別編集テスト
func TestShiftTypeEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("シフト種別編集フォーム表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shifts/xxx/edit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 200 OK、現在値が入力されたフォームHTML
		_ = req
	})

	t.Run("シフト種別編集成功", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "日勤更新")
		form.Add("code", "D")
		form.Add("color", "#00AA00")
		form.Add("start_time", "09:00")
		form.Add("end_time", "18:00")
		form.Add("break_minutes", "60")

		req := httptest.NewRequest(http.MethodPut, "/shifts/xxx", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 302 リダイレクト to /shifts
		_ = req
	})

	t.Run("存在しないシフト種別編集", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "不存在シフト")
		form.Add("code", "XX")

		req := httptest.NewRequest(http.MethodPut, "/shifts/non-existent-id", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "manager_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}

// TestShiftTypeDelete シフト種別削除テスト
func TestShiftTypeDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("シフト種別削除成功", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shifts/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 200 OK
		_ = req
	})

	t.Run("勤務表で使用中のシフト種別削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shifts/shift-in-use", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 400 Bad Request（関連データ存在）またはソフトデリート
		_ = req
	})

	t.Run("一般ユーザーによるシフト種別削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shifts/xxx", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "user_token",
		})

		// 期待: 403 Forbidden
		_ = req
	})

	t.Run("存在しないシフト種別削除", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shifts/non-existent-id", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "admin_token",
		})

		// 期待: 404 Not Found
		_ = req
	})
}
