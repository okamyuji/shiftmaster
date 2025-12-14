// Package presentation チームプレゼンテーション層
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"

	"github.com/google/uuid"
)

// TeamHandler チームHTTPハンドラー
type TeamHandler struct {
	teamRepo       domain.TeamRepository
	departmentRepo domain.DepartmentRepository
	templates      *web.TemplateEngine
	logger         *slog.Logger
}

// NewTeamHandler ハンドラー生成
func NewTeamHandler(
	teamRepo domain.TeamRepository,
	departmentRepo domain.DepartmentRepository,
	templates *web.TemplateEngine,
	logger *slog.Logger,
) *TeamHandler {
	return &TeamHandler{
		teamRepo:       teamRepo,
		departmentRepo: departmentRepo,
		templates:      templates,
		logger:         logger,
	}
}

// RegisterRoutes ルート登録
func (h *TeamHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teams", h.List)
	mux.HandleFunc("GET /teams/new", h.New)
	mux.HandleFunc("GET /teams/{id}/edit", h.Edit)
	mux.HandleFunc("POST /teams", h.Create)
	mux.HandleFunc("PUT /teams/{id}", h.Update)
	mux.HandleFunc("DELETE /teams/{id}", h.Delete)

	// API用エンドポイント
	mux.HandleFunc("GET /api/teams", h.ListJSON)
}

// getOrganizationID コンテキストから組織IDを取得
func (h *TeamHandler) getOrganizationID(r *http.Request) string {
	claims := web.GetClaimsFromContext(r.Context())
	if claims != nil && claims.OrganizationID != nil {
		return claims.OrganizationID.String()
	}
	return ""
}

// List チーム一覧ページ
func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		data := map[string]any{
			"Title":            "チーム一覧",
			"Teams":            []domain.Team{},
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください",
		}
		if err := h.templates.Render(w, "pages/teams/list.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

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
		"Title": "チーム一覧",
		"Teams": teams,
	}

	if isHTMXRequest(r) {
		if err := h.templates.RenderPartial(w, "pages/teams/list.html", "team-list", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.Render(w, "pages/teams/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// New 新規作成フォーム
func (h *TeamHandler) New(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationID(r)
	if orgID == "" {
		data := map[string]any{
			"Title":            "チーム追加",
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください",
		}
		if err := h.templates.Render(w, "pages/teams/form.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	organizationID, err := sharedDomain.ParseID(orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	departments, err := h.departmentRepo.FindByOrganizationID(r.Context(), organizationID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":       "チーム追加",
		"Departments": departments,
	}

	if err := h.templates.Render(w, "pages/teams/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Edit 編集フォーム
func (h *TeamHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	teamID, err := sharedDomain.ParseID(id)
	if err != nil {
		http.Error(w, "不正なID", http.StatusBadRequest)
		return
	}

	team, err := h.teamRepo.FindByID(r.Context(), teamID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	if team == nil {
		http.Error(w, "チームが見つかりません", http.StatusNotFound)
		return
	}

	departments, err := h.departmentRepo.FindAll(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":       "チーム編集",
		"Team":        team,
		"Departments": departments,
	}

	if err := h.templates.Render(w, "pages/teams/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create チーム作成
func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 部門IDを取得
	departmentIDStr := r.FormValue("department_id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		h.logger.Error("部門ID解析失敗", "department_id", departmentIDStr, "error", err)
		http.Error(w, "部門を選択してください", http.StatusBadRequest)
		return
	}

	team := &domain.Team{
		ID:           uuid.New(),
		DepartmentID: departmentID,
		Name:         r.FormValue("name"),
		Code:         r.FormValue("code"),
		SortOrder:    0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.teamRepo.Save(r.Context(), team); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

// Update チーム更新
func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	teamID, err := sharedDomain.ParseID(id)
	if err != nil {
		http.Error(w, "不正なID", http.StatusBadRequest)
		return
	}

	if parseErr := r.ParseForm(); parseErr != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	team, err := h.teamRepo.FindByID(r.Context(), teamID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	if team == nil {
		http.Error(w, "チームが見つかりません", http.StatusNotFound)
		return
	}

	// 部門IDを取得
	departmentIDStr := r.FormValue("department_id")
	if departmentIDStr != "" {
		departmentID, err := uuid.Parse(departmentIDStr)
		if err == nil {
			team.DepartmentID = departmentID
		}
	}

	team.Name = r.FormValue("name")
	team.Code = r.FormValue("code")
	team.UpdatedAt = time.Now()

	if err := h.teamRepo.Save(r.Context(), team); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

// Delete チーム削除
func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	teamID, err := sharedDomain.ParseID(id)
	if err != nil {
		http.Error(w, "不正なID", http.StatusBadRequest)
		return
	}

	if err := h.teamRepo.Delete(r.Context(), teamID); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

// ListJSON チーム一覧JSON
func (h *TeamHandler) ListJSON(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamRepo.FindAll(r.Context())
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "チーム取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, teams)
}

// handleError エラーハンドリング
func (h *TeamHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("ハンドラーエラー", "error", err, "method", r.Method, "path", r.URL.Path)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// writeJSON JSONレスポンス書き込み
func (h *TeamHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSONエンコード失敗", "error", err)
	}
}
