// Package presentation 認証プレゼンテーション層
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"shiftmaster/internal/modules/auth/application"
	"shiftmaster/internal/web"
)

// AuthHandler 認証ハンドラー
type AuthHandler struct {
	useCase   *application.AuthUseCase
	templates web.TemplateRenderer
	logger    *slog.Logger
}

// NewAuthHandler 認証ハンドラー生成
func NewAuthHandler(
	useCase *application.AuthUseCase,
	templates web.TemplateRenderer,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		useCase:   useCase,
		templates: templates,
		logger:    logger,
	}
}

// LoginPage ログインページ表示
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "ログイン",
		"Error": r.URL.Query().Get("error"),
	}

	// 認証ページ用のシンプルなレイアウトを使用
	if err := h.templates.RenderWithLayout(w, "pages/auth/login.html", "auth", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Login ログイン処理
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// フォームデータ取得
	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	input := &application.LoginInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	result, err := h.useCase.Login(r.Context(), input)
	if err != nil {
		h.logger.Warn("ログイン失敗", "email", input.Email, "error", err)
		// HTMX リクエストの場合
		if r.Header.Get("HX-Request") == "true" {
			h.handleError(w, r, "メールアドレスまたはパスワードが正しくありません", http.StatusUnauthorized)
			return
		}
		// 通常リクエスト
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusFound)
		return
	}

	// Cookie設定
	h.setAuthCookies(w, result)

	// HTMX リクエストの場合
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 通常リクエスト
	http.Redirect(w, r, "/", http.StatusFound)
}

// LoginAPI ログインAPI JSON
func (h *AuthHandler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	var input application.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
		return
	}

	result, err := h.useCase.Login(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// Refresh トークンリフレッシュAPI
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input application.RefreshInput

	// Cookie からリフレッシュトークン取得
	cookie, cookieErr := r.Cookie("refresh_token")
	if cookieErr == nil && cookie.Value != "" {
		input.RefreshToken = cookie.Value
	} else {
		// JSON ボディから取得
		if decodeErr := json.NewDecoder(r.Body).Decode(&input); decodeErr != nil {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
			return
		}
	}

	result, err := h.useCase.Refresh(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	// Cookie更新
	h.setAuthCookies(w, result)

	h.writeJSON(w, http.StatusOK, result)
}

// Logout ログアウト
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// リフレッシュトークン取得
	refreshToken := ""
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		refreshToken = cookie.Value
	}

	// ログアウト処理
	if err := h.useCase.Logout(r.Context(), refreshToken); err != nil {
		h.logger.Error("ログアウト失敗", "error", err)
	}

	// Cookie削除
	h.clearAuthCookies(w)

	// HTMX リクエストの場合
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// Me 現在のユーザー情報取得
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := web.GetClaimsFromContext(r.Context())
	if claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "認証が必要です"})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{
		"user_id":         claims.UserID.String(),
		"email":           claims.Email,
		"role":            claims.Role,
		"organization_id": claims.OrganizationID,
	})
}

// setAuthCookies 認証Cookie設定
func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, result *application.AuthOutput) {
	// アクセストークンCookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		HttpOnly: true,
		Secure:   false, // 本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
	})

	// リフレッシュトークンCookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // 本番環境ではtrue
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookies 認証Cookie削除
func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

// handleError エラーハンドリング
func (h *AuthHandler) handleError(w http.ResponseWriter, r *http.Request, message string, status int) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`<div class="text-red-400 text-sm">` + message + `</div>`))
		return
	}
	http.Error(w, message, status)
}

// writeJSON JSONレスポンス出力
func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSON書き込み失敗", "error", err)
	}
}
