// Package presentation 認証ハンドラーテスト
package presentation

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// テスト用のモックテンプレートレンダラー
type mockTemplateRenderer struct {
	renderCalled           bool
	renderWithLayoutCalled bool
	lastTemplate           string
	lastLayout             string
	lastData               any
}

func (m *mockTemplateRenderer) Render(w http.ResponseWriter, template string, data any) error {
	m.renderCalled = true
	m.lastTemplate = template
	m.lastData = data
	return nil
}

func (m *mockTemplateRenderer) RenderWithLayout(w http.ResponseWriter, template, layout string, data any) error {
	m.renderWithLayoutCalled = true
	m.lastTemplate = template
	m.lastLayout = layout
	m.lastData = data
	return nil
}

func TestNewAuthHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templates := &mockTemplateRenderer{}

	handler := NewAuthHandler(nil, templates, logger)

	if handler == nil {
		t.Fatal("NewAuthHandler returned nil")
	}
	if handler.templates == nil {
		t.Error("templates should not be nil")
	}
	if handler.logger == nil {
		t.Error("logger should not be nil")
	}
}

func TestAuthHandler_LoginPage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	templates := &mockTemplateRenderer{}

	handler := NewAuthHandler(nil, templates, logger)

	t.Run("ログインページ表示", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()

		handler.LoginPage(w, req)

		if !templates.renderWithLayoutCalled {
			t.Error("RenderWithLayout should be called")
		}
		if templates.lastTemplate != "pages/auth/login.html" {
			t.Errorf("expected template pages/auth/login.html but got %s", templates.lastTemplate)
		}
		if templates.lastLayout != "auth" {
			t.Errorf("expected layout auth but got %s", templates.lastLayout)
		}
		dataMap, ok := templates.lastData.(map[string]any)
		if !ok {
			t.Fatal("lastData should be map[string]any")
		}
		if dataMap["Title"] != "ログイン" {
			t.Errorf("expected title ログイン but got %v", dataMap["Title"])
		}
	})

	t.Run("エラーパラメータ付きログインページ表示", func(t *testing.T) {
		templates := &mockTemplateRenderer{}
		handler := NewAuthHandler(nil, templates, logger)

		req := httptest.NewRequest(http.MethodGet, "/login?error=test_error", nil)
		w := httptest.NewRecorder()

		handler.LoginPage(w, req)

		dataMap, ok := templates.lastData.(map[string]any)
		if !ok {
			t.Fatal("lastData should be map[string]any")
		}
		if dataMap["Error"] != "test_error" {
			t.Errorf("expected error test_error but got %v", dataMap["Error"])
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	templates := &mockTemplateRenderer{}

	// useCaseがnilでもLogoutは動作する（トークン削除はエラーを無視）
	handler := NewAuthHandler(nil, templates, logger)

	t.Run("通常のログアウト", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		resp := w.Result()
		defer func() { _ = resp.Body.Close() }()

		// リダイレクト確認
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected status 302 but got %d", resp.StatusCode)
		}
		if resp.Header.Get("Location") != "/login" {
			t.Errorf("expected redirect to /login but got %s", resp.Header.Get("Location"))
		}
	})

	t.Run("HTMXリクエストでのログアウト", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("HX-Request", "true")
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		resp := w.Result()
		defer func() { _ = resp.Body.Close() }()

		// HX-Redirectヘッダー確認
		if resp.Header.Get("HX-Redirect") != "/login" {
			t.Errorf("expected HX-Redirect /login but got %s", resp.Header.Get("HX-Redirect"))
		}
	})
}
