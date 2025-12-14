// Package presentation スタッフプレゼンテーション層
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"shiftmaster/internal/modules/staff/application"
	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"
)

// StaffHandler スタッフHTTPハンドラー
type StaffHandler struct {
	useCase   *application.StaffUseCase
	teamRepo  domain.TeamRepository
	templates *web.TemplateEngine
	logger    *slog.Logger
}

// NewStaffHandler ハンドラー生成
func NewStaffHandler(
	useCase *application.StaffUseCase,
	teamRepo domain.TeamRepository,
	templates *web.TemplateEngine,
	logger *slog.Logger,
) *StaffHandler {
	return &StaffHandler{
		useCase:   useCase,
		teamRepo:  teamRepo,
		templates: templates,
		logger:    logger,
	}
}

// RegisterRoutes ルート登録
func (h *StaffHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /staffs", h.List)
	mux.HandleFunc("GET /staffs/new", h.New)
	mux.HandleFunc("GET /staffs/{id}", h.Show)
	mux.HandleFunc("GET /staffs/{id}/edit", h.Edit)
	mux.HandleFunc("POST /staffs", h.Create)
	mux.HandleFunc("PUT /staffs/{id}", h.Update)
	mux.HandleFunc("DELETE /staffs/{id}", h.Delete)

	// HTMX用エンドポイント
	mux.HandleFunc("GET /api/staffs", h.ListJSON)
	mux.HandleFunc("GET /api/staffs/{id}", h.ShowJSON)
	mux.HandleFunc("POST /api/staffs", h.CreateJSON)
	mux.HandleFunc("PUT /api/staffs/{id}", h.UpdateJSON)
	mux.HandleFunc("DELETE /api/staffs/{id}", h.DeleteJSON)
}

// getOrganizationID コンテキストから組織IDを取得
func (h *StaffHandler) getOrganizationID(r *http.Request) string {
	claims := web.GetClaimsFromContext(r.Context())
	if claims != nil && claims.OrganizationID != nil {
		return claims.OrganizationID.String()
	}
	return ""
}

// List スタッフ一覧ページ
func (h *StaffHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		// 組織未選択の場合は空リストを返す
		data := map[string]any{
			"Title":            "スタッフ一覧",
			"Staffs":           []any{},
			"Total":            0,
			"Page":             1,
			"PerPage":          20,
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください",
		}
		if err := h.templates.Render(w, "pages/staffs/list.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	result, err := h.useCase.ListByOrganization(r.Context(), orgID, page, perPage)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":   "スタッフ一覧",
		"Staffs":  result.Staffs,
		"Total":   result.Total,
		"Page":    result.Page,
		"PerPage": result.PerPage,
	}

	if isHTMXRequest(r) {
		if err := h.templates.RenderPartial(w, "pages/staffs/list.html", "staff-list", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.Render(w, "pages/staffs/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// New 新規作成フォーム
func (h *StaffHandler) New(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		data := map[string]any{
			"Title":            "スタッフ追加",
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください",
		}
		if err := h.templates.Render(w, "pages/staffs/form.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// 組織に属するチーム一覧取得
	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	teams, err := h.teamRepo.FindByOrganizationID(r.Context(), organizationID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title": "スタッフ追加",
		"Teams": teams,
	}

	if err := h.templates.Render(w, "pages/staffs/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Show スタッフ詳細ページ
func (h *StaffHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationID(r)

	if orgID == "" {
		http.Error(w, "組織が選択されていません", http.StatusBadRequest)
		return
	}

	staff, err := h.useCase.GetByIDWithOrg(r.Context(), id, orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title": staff.FullName,
		"Staff": staff,
	}

	if err := h.templates.Render(w, "pages/staffs/show.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Edit 編集フォーム
func (h *StaffHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationID(r)

	if orgID == "" {
		http.Error(w, "組織が選択されていません", http.StatusBadRequest)
		return
	}

	staff, err := h.useCase.GetByIDWithOrg(r.Context(), id, orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	// 組織に属するチーム一覧取得
	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	teams, err := h.teamRepo.FindByOrganizationID(r.Context(), organizationID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title": "スタッフ編集",
		"Staff": staff,
		"Teams": teams,
	}

	if err := h.templates.Render(w, "pages/staffs/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create スタッフ作成 フォーム送信
func (h *StaffHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	input := &application.CreateStaffInput{
		TeamID:         r.FormValue("team_id"),
		EmployeeCode:   r.FormValue("employee_code"),
		FirstName:      r.FormValue("first_name"),
		LastName:       r.FormValue("last_name"),
		Email:          r.FormValue("email"),
		Phone:          r.FormValue("phone"),
		HireDate:       r.FormValue("hire_date"),
		EmploymentType: r.FormValue("employment_type"),
	}

	_, err := h.useCase.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/staffs", http.StatusSeeOther)
}

// Update スタッフ更新 フォーム送信
func (h *StaffHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	input := &application.UpdateStaffInput{
		ID:             id,
		TeamID:         r.FormValue("team_id"),
		EmployeeCode:   r.FormValue("employee_code"),
		FirstName:      r.FormValue("first_name"),
		LastName:       r.FormValue("last_name"),
		Email:          r.FormValue("email"),
		Phone:          r.FormValue("phone"),
		HireDate:       r.FormValue("hire_date"),
		EmploymentType: r.FormValue("employment_type"),
		IsActive:       r.FormValue("is_active") == "true",
	}

	_, err := h.useCase.Update(r.Context(), input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/staffs/"+id, http.StatusSeeOther)
}

// Delete スタッフ削除
func (h *StaffHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationID(r)

	if orgID == "" {
		http.Error(w, "組織が選択されていません", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteWithOrg(r.Context(), id, orgID); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/staffs", http.StatusSeeOther)
}

// ListJSON スタッフ一覧JSON
func (h *StaffHandler) ListJSON(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "組織IDが必要です"})
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	result, err := h.useCase.ListByOrganization(r.Context(), orgID, page, perPage)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// ShowJSON スタッフ詳細JSON
func (h *StaffHandler) ShowJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "組織IDが必要です"})
		return
	}

	staff, err := h.useCase.GetByIDWithOrg(r.Context(), id, orgID)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, staff)
}

// CreateJSON スタッフ作成JSON
func (h *StaffHandler) CreateJSON(w http.ResponseWriter, r *http.Request) {
	var input application.CreateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
		return
	}

	staff, err := h.useCase.Create(r.Context(), &input)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, staff)
}

// UpdateJSON スタッフ更新JSON
func (h *StaffHandler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input application.UpdateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
		return
	}
	input.ID = id

	staff, err := h.useCase.Update(r.Context(), &input)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, staff)
}

// DeleteJSON スタッフ削除JSON
func (h *StaffHandler) DeleteJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "組織IDが必要です"})
		return
	}

	if err := h.useCase.DeleteWithOrg(r.Context(), id, orgID); err != nil {
		h.handleJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError エラーハンドリング
func (h *StaffHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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
			http.Error(w, domainErr.Message, http.StatusForbidden)
			return
		}
	}

	h.logger.Error("ハンドラーエラー", "error", err, "method", r.Method, "path", r.URL.Path)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// handleJSONError JSONエラーハンドリング
func (h *StaffHandler) handleJSONError(w http.ResponseWriter, err error) {
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

	h.logger.Error("ハンドラーエラー", "error", err)
	h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "内部エラーが発生しました"})
}

// writeJSON JSONレスポンス書き込み
func (h *StaffHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSONエンコード失敗", "error", err)
	}
}

// isHTMXRequest HTMXリクエスト判定
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
