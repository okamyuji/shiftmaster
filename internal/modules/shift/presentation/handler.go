// Package presentation シフトプレゼンテーション層
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"shiftmaster/internal/modules/shift/application"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"
)

// ShiftTypeHandler シフト種別HTTPハンドラー
type ShiftTypeHandler struct {
	useCase   *application.ShiftTypeUseCase
	templates *web.TemplateEngine
	logger    *slog.Logger
}

// NewShiftTypeHandler ハンドラー生成
func NewShiftTypeHandler(
	useCase *application.ShiftTypeUseCase,
	templates *web.TemplateEngine,
	logger *slog.Logger,
) *ShiftTypeHandler {
	return &ShiftTypeHandler{
		useCase:   useCase,
		templates: templates,
		logger:    logger,
	}
}

// RegisterRoutes ルート登録
func (h *ShiftTypeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /shifts", h.List)
	mux.HandleFunc("GET /shifts/new", h.New)
	mux.HandleFunc("GET /shifts/{id}", h.Show)
	mux.HandleFunc("GET /shifts/{id}/edit", h.Edit)
	mux.HandleFunc("POST /shifts", h.Create)
	mux.HandleFunc("PUT /shifts/{id}", h.Update)
	mux.HandleFunc("DELETE /shifts/{id}", h.Delete)

	// API
	mux.HandleFunc("GET /api/shifts", h.ListJSON)
	mux.HandleFunc("GET /api/shifts/{id}", h.ShowJSON)
	mux.HandleFunc("POST /api/shifts", h.CreateJSON)
	mux.HandleFunc("PUT /api/shifts/{id}", h.UpdateJSON)
	mux.HandleFunc("DELETE /api/shifts/{id}", h.DeleteJSON)
}

// getOrganizationID コンテキストから組織IDを取得
func (h *ShiftTypeHandler) getOrganizationID(r *http.Request) string {
	claims := web.GetClaimsFromContext(r.Context())
	if claims != nil && claims.OrganizationID != nil {
		return claims.OrganizationID.String()
	}
	return ""
}

// List シフト種別一覧ページ
func (h *ShiftTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		data := map[string]any{
			"Title":            "シフト種別一覧",
			"ShiftTypes":       []any{},
			"Total":            0,
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください",
		}
		if err := h.templates.Render(w, "pages/shifts/list.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	result, err := h.useCase.ListByOrganization(r.Context(), orgID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := map[string]any{
		"Title":      "シフト種別一覧",
		"ShiftTypes": result.ShiftTypes,
		"Total":      result.Total,
	}

	if err := h.templates.Render(w, "pages/shifts/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// New 新規作成フォーム
func (h *ShiftTypeHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "シフト種別追加",
	}

	if err := h.templates.Render(w, "pages/shifts/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Show シフト種別詳細ページ
func (h *ShiftTypeHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	shiftType, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := map[string]any{
		"Title":     shiftType.Name,
		"ShiftType": shiftType,
	}

	if err := h.templates.Render(w, "pages/shifts/show.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Edit 編集フォーム
func (h *ShiftTypeHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	shiftType, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	data := map[string]any{
		"Title":     "シフト種別編集",
		"ShiftType": shiftType,
	}

	if err := h.templates.Render(w, "pages/shifts/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create シフト種別作成
func (h *ShiftTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 認証ユーザーからOrganizationID取得
	claims := web.GetClaimsFromContext(r.Context())
	var orgID string
	if claims != nil && claims.OrganizationID != nil {
		orgID = claims.OrganizationID.String()
	}

	breakMinutes, _ := strconv.Atoi(r.FormValue("break_minutes"))
	handoverMinutes, _ := strconv.Atoi(r.FormValue("handover_minutes"))
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))

	input := &application.CreateShiftTypeInput{
		OrganizationID:  orgID,
		ShiftPatternID:  r.FormValue("shift_pattern_id"),
		Name:            r.FormValue("name"),
		Code:            r.FormValue("code"),
		Color:           r.FormValue("color"),
		StartTime:       r.FormValue("start_time"),
		EndTime:         r.FormValue("end_time"),
		BreakMinutes:    breakMinutes,
		HandoverMinutes: handoverMinutes,
		IsNightShift:    r.FormValue("is_night_shift") == "true" || r.FormValue("is_night_shift") == "on",
		IsHoliday:       r.FormValue("is_holiday") == "true" || r.FormValue("is_holiday") == "on",
		SortOrder:       sortOrder,
	}

	_, err := h.useCase.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	http.Redirect(w, r, "/shifts", http.StatusSeeOther)
}

// Update シフト種別更新
func (h *ShiftTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	breakMinutes, _ := strconv.Atoi(r.FormValue("break_minutes"))
	handoverMinutes, _ := strconv.Atoi(r.FormValue("handover_minutes"))
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))

	input := &application.UpdateShiftTypeInput{
		ID:              id,
		ShiftPatternID:  r.FormValue("shift_pattern_id"),
		Name:            r.FormValue("name"),
		Code:            r.FormValue("code"),
		Color:           r.FormValue("color"),
		StartTime:       r.FormValue("start_time"),
		EndTime:         r.FormValue("end_time"),
		BreakMinutes:    breakMinutes,
		HandoverMinutes: handoverMinutes,
		IsNightShift:    r.FormValue("is_night_shift") == "true" || r.FormValue("is_night_shift") == "on",
		IsHoliday:       r.FormValue("is_holiday") == "true" || r.FormValue("is_holiday") == "on",
		SortOrder:       sortOrder,
	}

	_, err := h.useCase.Update(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	http.Redirect(w, r, "/shifts/"+id, http.StatusSeeOther)
}

// Delete シフト種別削除
func (h *ShiftTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/shifts", http.StatusSeeOther)
}

// ListJSON シフト種別一覧JSON
func (h *ShiftTypeHandler) ListJSON(w http.ResponseWriter, r *http.Request) {
	result, err := h.useCase.List(r.Context())
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// ShowJSON シフト種別詳細JSON
func (h *ShiftTypeHandler) ShowJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	shiftType, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, shiftType)
}

// CreateJSON シフト種別作成JSON
func (h *ShiftTypeHandler) CreateJSON(w http.ResponseWriter, r *http.Request) {
	var input application.CreateShiftTypeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
		return
	}

	shiftType, err := h.useCase.Create(r.Context(), &input)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, shiftType)
}

// UpdateJSON シフト種別更新JSON
func (h *ShiftTypeHandler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input application.UpdateShiftTypeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
		return
	}
	input.ID = id

	shiftType, err := h.useCase.Update(r.Context(), &input)
	if err != nil {
		h.handleJSONError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, shiftType)
}

// DeleteJSON シフト種別削除JSON
func (h *ShiftTypeHandler) DeleteJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		h.handleJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError エラーハンドリング
func (h *ShiftTypeHandler) handleError(w http.ResponseWriter, err error) {
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
		}
	}

	h.logger.Error("ハンドラーエラー", "error", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// handleJSONError JSONエラーハンドリング
func (h *ShiftTypeHandler) handleJSONError(w http.ResponseWriter, err error) {
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
		}
	}

	h.logger.Error("ハンドラーエラー", "error", err)
	h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "内部エラーが発生しました"})
}

// writeJSON JSONレスポンス書き込み
func (h *ShiftTypeHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSONエンコード失敗", "error", err)
	}
}
