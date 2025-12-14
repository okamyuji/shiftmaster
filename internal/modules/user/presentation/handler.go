// Package presentation ユーザープレゼンテーション層
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"shiftmaster/internal/modules/user/application"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"
)

// UserHandler ユーザーハンドラー
type UserHandler struct {
	useCase   *application.UserUseCase
	templates web.TemplateRenderer
	logger    *slog.Logger
}

// NewUserHandler ユーザーハンドラー生成
func NewUserHandler(
	useCase *application.UserUseCase,
	templates web.TemplateRenderer,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		useCase:   useCase,
		templates: templates,
		logger:    logger,
	}
}

// ListUsers ユーザー一覧ページ
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	claims := web.GetClaimsFromContext(r.Context())

	var result *application.UserListOutput
	var err error

	// super_adminは全ユーザー表示、それ以外は自組織のユーザーのみ
	if claims != nil && claims.IsSuperAdmin() {
		result, err = h.useCase.List(r.Context())
	} else if claims != nil && claims.OrganizationID != nil {
		result, err = h.useCase.ListByOrganization(r.Context(), claims.OrganizationID.String())
	} else {
		result = &application.UserListOutput{Users: []application.UserOutput{}, Total: 0}
	}

	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":       "ユーザー管理",
		"Users":       result.Users,
		"Total":       result.Total,
		"CurrentUser": claims,
	}

	if err := h.templates.Render(w, "pages/admin/users.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewUserForm ユーザー作成フォーム
func (h *UserHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	claims := web.GetClaimsFromContext(r.Context())
	roles := getRoleOptions()
	if claims != nil && claims.IsSuperAdmin() {
		roles = getRoleOptionsForSuperAdmin()
	}

	data := map[string]any{
		"Title": "ユーザー追加",
		"IsNew": true,
		"User":  nil,
		"Roles": roles,
	}

	if err := h.templates.Render(w, "pages/admin/user_form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateUser ユーザー作成
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	input := &application.CreateUserInput{
		OrganizationID: r.FormValue("organization_id"),
		Email:          r.FormValue("email"),
		Password:       r.FormValue("password"),
		FirstName:      r.FormValue("first_name"),
		LastName:       r.FormValue("last_name"),
		Role:           r.FormValue("role"),
	}

	_, err := h.useCase.Create(r.Context(), input)
	if err != nil {
		h.handleFormError(w, r, err)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/users")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

// EditUserForm ユーザー編集フォーム
func (h *UserHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "IDが必要です", http.StatusBadRequest)
		return
	}

	user, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	claims := web.GetClaimsFromContext(r.Context())
	roles := getRoleOptions()
	if claims != nil && claims.IsSuperAdmin() {
		roles = getRoleOptionsForSuperAdmin()
	}

	data := map[string]any{
		"Title": "ユーザー編集",
		"IsNew": false,
		"User":  user,
		"Roles": roles,
	}

	if err := h.templates.Render(w, "pages/admin/user_form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateUser ユーザー更新
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "IDが必要です", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	isActive := r.FormValue("is_active") == "on" || r.FormValue("is_active") == "true"

	input := &application.UpdateUserInput{
		ID:             id,
		OrganizationID: r.FormValue("organization_id"),
		Email:          r.FormValue("email"),
		FirstName:      r.FormValue("first_name"),
		LastName:       r.FormValue("last_name"),
		Role:           r.FormValue("role"),
		IsActive:       isActive,
	}

	_, err := h.useCase.Update(r.Context(), input)
	if err != nil {
		h.handleFormError(w, r, err)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/users")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

// DeleteUser ユーザー削除
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "IDが必要です"})
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		h.handleJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError エラーハンドリング
func (h *UserHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *sharedDomain.DomainError
	if de, ok := err.(*sharedDomain.DomainError); ok {
		domainErr = de
	}

	if domainErr != nil {
		switch domainErr.Code {
		case sharedDomain.ErrCodeNotFound:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		case sharedDomain.ErrCodeValidation:
			http.Error(w, domainErr.Message, http.StatusBadRequest)
			return
		case sharedDomain.ErrCodeForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	h.logger.Error("ハンドラーエラー", "error", err, "method", r.Method, "path", r.URL.Path)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// handleFormError フォームエラーハンドリング
func (h *UserHandler) handleFormError(w http.ResponseWriter, r *http.Request, err error) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`<div class="text-red-400 text-sm">` + err.Error() + `</div>`))
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// handleJSONError JSONエラーハンドリング
func (h *UserHandler) handleJSONError(w http.ResponseWriter, err error) {
	var domainErr *sharedDomain.DomainError
	if de, ok := err.(*sharedDomain.DomainError); ok {
		domainErr = de
	}

	if domainErr != nil {
		switch domainErr.Code {
		case sharedDomain.ErrCodeNotFound:
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": domainErr.Message})
			return
		case sharedDomain.ErrCodeValidation:
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": domainErr.Message})
			return
		case sharedDomain.ErrCodeForbidden:
			h.writeJSON(w, http.StatusForbidden, map[string]string{"error": domainErr.Message})
			return
		}
	}

	h.logger.Error("JSONエラー", "error", err)
	h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
}

// writeJSON JSONレスポンス出力
func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSON書き込み失敗", "error", err)
	}
}

// RoleOption ロール選択肢
type RoleOption struct {
	Value string
	Label string
}

// getRoleOptions ロール選択肢取得（テナント管理者用）
func getRoleOptions() []RoleOption {
	return []RoleOption{
		{Value: "admin", Label: "テナント管理者"},
		{Value: "manager", Label: "マネージャー"},
		{Value: "user", Label: "一般ユーザー"},
	}
}

// getRoleOptionsForSuperAdmin ロール選択肢取得（スーパー管理者用）
func getRoleOptionsForSuperAdmin() []RoleOption {
	return []RoleOption{
		{Value: "super_admin", Label: "全体管理者"},
		{Value: "admin", Label: "テナント管理者"},
		{Value: "manager", Label: "マネージャー"},
		{Value: "user", Label: "一般ユーザー"},
	}
}
