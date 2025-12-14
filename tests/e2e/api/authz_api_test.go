//go:build e2e

// Package api 認可APIのE2Eテスト
package api

import (
	"net/http"
	"testing"
)

// TestAuthorizationAPI 認可テスト
func TestAuthorizationAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E API test in short mode")
	}

	t.Run("admin_シフト種別一覧取得_許可", func(t *testing.T) {
		token := getAccessToken(t, "admin@example.com", "Password123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/shifts", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("super_admin_シフト種別一覧取得_許可", func(t *testing.T) {
		token := getAccessToken(t, "superadmin@example.com", "SuperAdmin123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/shifts", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("認証なし_シフト種別一覧取得_拒否", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/shifts", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("期待: %d, 実際: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}

// TestOrganizationSwitchAPI 組織切り替えAPIテスト
func TestOrganizationSwitchAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E API test in short mode")
	}

	t.Run("super_admin_組織切り替え_許可", func(t *testing.T) {
		token := getAccessToken(t, "superadmin@example.com", "SuperAdmin123$!")

		// 存在する組織IDを取得（サンプル病院）
		// 実際のテストではAPIから組織一覧を取得して使うべき
		orgID := "109e0070-a30d-46d7-8593-fdfa364a8e37"

		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/switch-organization/"+orgID, nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		// リダイレクトを追跡しない
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 302リダイレクトを期待
		if resp.StatusCode != http.StatusFound {
			t.Errorf("期待: %d, 実際: %d", http.StatusFound, resp.StatusCode)
		}

		// Cookieが設定されていることを確認
		cookies := resp.Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == "selected_organization_id" {
				found = true
				if c.Value != orgID {
					t.Errorf("期待する組織ID: %s, 実際: %s", orgID, c.Value)
				}
			}
		}
		if !found {
			t.Error("selected_organization_id Cookieが設定されていない")
		}
	})

	t.Run("admin_組織切り替え_拒否", func(t *testing.T) {
		token := getAccessToken(t, "admin@example.com", "Password123$!")

		orgID := "109e0070-a30d-46d7-8593-fdfa364a8e37"

		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/switch-organization/"+orgID, nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 403 Forbiddenを期待
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("期待: %d, 実際: %d", http.StatusForbidden, resp.StatusCode)
		}
	})

	t.Run("認証なし_組織切り替え_拒否", func(t *testing.T) {
		orgID := "109e0070-a30d-46d7-8593-fdfa364a8e37"

		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/switch-organization/"+orgID, nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 401 Unauthorizedを期待
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("期待: %d, 実際: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}

// TestAdminUserManagementAPI ユーザー管理API認可テスト
func TestAdminUserManagementAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E API test in short mode")
	}

	t.Run("admin_ユーザー一覧取得_許可", func(t *testing.T) {
		token := getAccessToken(t, "admin@example.com", "Password123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/users", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// HTMLページが返るため200を期待
		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("super_admin_ユーザー一覧取得_許可", func(t *testing.T) {
		token := getAccessToken(t, "superadmin@example.com", "SuperAdmin123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/users", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("認証なし_ユーザー一覧取得_拒否", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/admin/users", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("リクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("期待: %d, 実際: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}
