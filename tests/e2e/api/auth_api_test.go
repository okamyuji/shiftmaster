//go:build e2e

// Package api 認証APIのE2Eテスト
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

var baseURL string

func init() {
	baseURL = os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse ログインレスポンス
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		Role      string `json:"role"`
		RoleLabel string `json:"role_label"`
		IsAdmin   bool   `json:"is_admin"`
	} `json:"user"`
}

// MeResponse /api/auth/me レスポンス
type MeResponse struct {
	UserID         string  `json:"user_id"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	OrganizationID *string `json:"organization_id"`
}

// TestLoginAPI ログインAPIテスト
func TestLoginAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E API test in short mode")
	}

	t.Run("正常なadminログイン", func(t *testing.T) {
		resp, err := login("admin@example.com", "Password123$!")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}

		var result LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("JSONデコード失敗: %v", err)
		}

		if result.AccessToken == "" {
			t.Error("アクセストークンが空")
		}
		if result.RefreshToken == "" {
			t.Error("リフレッシュトークンが空")
		}
		if result.User.Role != "admin" {
			t.Errorf("期待するロール: admin, 実際: %s", result.User.Role)
		}
		if result.User.RoleLabel != "テナント管理者" {
			t.Errorf("期待するロールラベル: テナント管理者, 実際: %s", result.User.RoleLabel)
		}
	})

	t.Run("正常なsuper_adminログイン", func(t *testing.T) {
		resp, err := login("superadmin@example.com", "SuperAdmin123$!")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期待: %d, 実際: %d", http.StatusOK, resp.StatusCode)
		}

		var result LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("JSONデコード失敗: %v", err)
		}

		if result.User.Role != "super_admin" {
			t.Errorf("期待するロール: super_admin, 実際: %s", result.User.Role)
		}
		if result.User.RoleLabel != "全体管理者" {
			t.Errorf("期待するロールラベル: 全体管理者, 実際: %s", result.User.RoleLabel)
		}
	})

	t.Run("不正なメールアドレス", func(t *testing.T) {
		resp, err := login("notexist@example.com", "Password123$!")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("期待: %d, 実際: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("不正なパスワード", func(t *testing.T) {
		resp, err := login("admin@example.com", "WrongPassword!")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("期待: %d, 実際: %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("空のメールアドレス", func(t *testing.T) {
		resp, err := login("", "Password123$!")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 実装では空のメールアドレスも401を返す
		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("期待: 400または401, 実際: %d", resp.StatusCode)
		}
	})

	t.Run("空のパスワード", func(t *testing.T) {
		resp, err := login("admin@example.com", "")
		if err != nil {
			t.Fatalf("ログインリクエスト失敗: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 実装では空のパスワードも401を返す
		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("期待: 400または401, 実際: %d", resp.StatusCode)
		}
	})
}

// TestMeAPI /api/auth/me テスト
func TestMeAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E API test in short mode")
	}

	t.Run("認証済みユーザー情報取得_admin", func(t *testing.T) {
		token := getAccessToken(t, "admin@example.com", "Password123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/auth/me", nil)
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

		var result MeResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("JSONデコード失敗: %v", err)
		}

		if result.Email != "admin@example.com" {
			t.Errorf("期待するEmail: admin@example.com, 実際: %s", result.Email)
		}
		if result.Role != "admin" {
			t.Errorf("期待するロール: admin, 実際: %s", result.Role)
		}
		if result.OrganizationID == nil {
			t.Error("OrganizationIDがnull（adminはテナントに属するべき）")
		}
	})

	t.Run("認証済みユーザー情報取得_super_admin", func(t *testing.T) {
		token := getAccessToken(t, "superadmin@example.com", "SuperAdmin123$!")

		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/auth/me", nil)
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

		var result MeResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("JSONデコード失敗: %v", err)
		}

		if result.Email != "superadmin@example.com" {
			t.Errorf("期待するEmail: superadmin@example.com, 実際: %s", result.Email)
		}
		if result.Role != "super_admin" {
			t.Errorf("期待するロール: super_admin, 実際: %s", result.Role)
		}
		if result.OrganizationID != nil {
			t.Error("OrganizationIDが設定されている（super_adminは組織に属さない）")
		}
	})

	t.Run("未認証でユーザー情報取得", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/auth/me", nil)
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

	t.Run("無効なトークンでユーザー情報取得", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/auth/me", nil)
		if err != nil {
			t.Fatalf("リクエスト生成失敗: %v", err)
		}
		req.Header.Set("Authorization", "Bearer invalid_token_here")

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

// login ログインヘルパー
func login(email, password string) (*http.Response, error) {
	reqBody := LoginRequest{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(reqBody)

	return http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
}

// getAccessToken ログインしてアクセストークンを取得
func getAccessToken(t *testing.T, email, password string) string {
	t.Helper()

	resp, err := login(email, password)
	if err != nil {
		t.Fatalf("ログイン失敗: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ログイン失敗: ステータス %d", resp.StatusCode)
	}

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("JSONデコード失敗: %v", err)
	}

	return result.AccessToken
}
