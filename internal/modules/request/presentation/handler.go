// Package presentation 勤務希望プレゼンテーション層
package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"shiftmaster/internal/modules/request/application"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"
)

// StaffFinder スタッフ検索インターフェース
type StaffFinder interface {
	FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]StaffInfo, error)
	FindByID(ctx context.Context, id sharedDomain.ID) (*StaffInfo, error)
}

// StaffInfo スタッフ情報
type StaffInfo struct {
	ID        string
	FirstName string
	LastName  string
}

// FullName フルネーム取得
func (s *StaffInfo) FullName() string {
	return s.LastName + " " + s.FirstName
}

// RequestHandler 勤務希望HTTPハンドラー
type RequestHandler struct {
	periodUseCase  *application.RequestPeriodUseCase
	requestUseCase *application.ShiftRequestUseCase
	staffFinder    StaffFinder
	templates      *web.TemplateEngine
	logger         *slog.Logger
}

// NewRequestHandler ハンドラー生成
func NewRequestHandler(
	periodUseCase *application.RequestPeriodUseCase,
	requestUseCase *application.ShiftRequestUseCase,
	staffFinder StaffFinder,
	templates *web.TemplateEngine,
	logger *slog.Logger,
) *RequestHandler {
	return &RequestHandler{
		periodUseCase:  periodUseCase,
		requestUseCase: requestUseCase,
		staffFinder:    staffFinder,
		templates:      templates,
		logger:         logger,
	}
}

// RegisterRoutes ルート登録
func (h *RequestHandler) RegisterRoutes(mux *http.ServeMux) {
	// 受付期間管理
	mux.HandleFunc("GET /requests", h.ListPeriods)
	mux.HandleFunc("GET /requests/new", h.NewPeriod)
	mux.HandleFunc("GET /requests/{id}", h.ShowPeriod)
	mux.HandleFunc("POST /requests", h.CreatePeriod)
	mux.HandleFunc("POST /requests/{id}/open", h.OpenPeriod)
	mux.HandleFunc("POST /requests/{id}/close", h.ClosePeriod)

	// 勤務希望管理
	mux.HandleFunc("GET /requests/{period_id}/entries", h.ListRequests)
	mux.HandleFunc("GET /requests/{period_id}/entries/new", h.NewRequest)
	mux.HandleFunc("POST /requests/{period_id}/entries", h.CreateRequest)
	mux.HandleFunc("DELETE /requests/entries/{id}", h.DeleteRequest)

	// API用エンドポイント
	mux.HandleFunc("GET /api/requests", h.ListPeriodsJSON)
	mux.HandleFunc("GET /api/requests/{id}", h.ShowPeriodJSON)
	mux.HandleFunc("POST /api/requests", h.CreatePeriodJSON)
	mux.HandleFunc("GET /api/requests/{period_id}/entries", h.ListRequestsJSON)
	mux.HandleFunc("POST /api/requests/{period_id}/entries", h.CreateRequestJSON)
}

// ListPeriods 受付期間一覧ページ
func (h *RequestHandler) ListPeriods(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationIDFromContext(r.Context())

	var periods []application.RequestPeriodOutput
	var err error

	if orgID != "" {
		periods, err = h.periodUseCase.ListPeriodsByOrganization(r.Context(), orgID)
	} else {
		// super_admin で組織未選択時は空配列
		periods = []application.RequestPeriodOutput{}
	}

	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":   "勤務希望受付期間一覧",
		"Periods": periods,
	}

	if isHTMXRequest(r) {
		if err := h.templates.RenderPartial(w, "pages/request_periods/list.html", "period-list", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.Render(w, "pages/request_periods/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewPeriod 受付期間新規作成フォーム
func (h *RequestHandler) NewPeriod(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	data := map[string]any{
		"Title":          "受付期間作成",
		"CurrentYear":    now.Year(),
		"DefaultYear":    nextMonth.Year(),
		"DefaultMonth":   int(nextMonth.Month()),
		"AvailableYears": []int{now.Year(), now.Year() + 1},
	}

	if err := h.templates.Render(w, "pages/request_periods/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ShowPeriod 受付期間詳細ページ
func (h *RequestHandler) ShowPeriod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := h.getOrganizationIDFromContext(r.Context())

	period, err := h.periodUseCase.GetPeriodByIDWithOrg(r.Context(), id, orgID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	requests, err := h.requestUseCase.ListByPeriod(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":    period.TargetPeriodLabel + " 勤務希望",
		"Period":   period,
		"Requests": requests.Requests,
		"Total":    requests.Total,
	}

	if err := h.templates.Render(w, "pages/request_periods/show.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreatePeriod 受付期間作成
func (h *RequestHandler) CreatePeriod(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	targetYear, _ := strconv.Atoi(r.FormValue("target_year"))
	targetMonth, _ := strconv.Atoi(r.FormValue("target_month"))
	maxPerStaff, _ := strconv.Atoi(r.FormValue("max_requests_per_staff"))
	maxPerDay, _ := strconv.Atoi(r.FormValue("max_requests_per_day"))

	input := &application.CreateRequestPeriodInput{
		OrganizationID:      h.getDefaultOrganizationID(r.Context()),
		TargetYear:          targetYear,
		TargetMonth:         targetMonth,
		StartDate:           r.FormValue("start_date"),
		EndDate:             r.FormValue("end_date"),
		MaxRequestsPerStaff: maxPerStaff,
		MaxRequestsPerDay:   maxPerDay,
	}

	period, err := h.periodUseCase.CreatePeriod(r.Context(), input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/requests/"+period.ID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/requests/"+period.ID, http.StatusSeeOther)
}

// OpenPeriod 受付開始
func (h *RequestHandler) OpenPeriod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := h.periodUseCase.OpenPeriod(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/requests/"+id)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/requests/"+id, http.StatusSeeOther)
}

// ClosePeriod 受付終了
func (h *RequestHandler) ClosePeriod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := h.periodUseCase.ClosePeriod(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/requests/"+id)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/requests/"+id, http.StatusSeeOther)
}

// ListRequests 勤務希望一覧
func (h *RequestHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	periodID := r.PathValue("period_id")

	period, err := h.periodUseCase.GetPeriodByID(r.Context(), periodID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	requests, err := h.requestUseCase.ListByPeriod(r.Context(), periodID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	// スタッフ名を設定
	if h.staffFinder != nil {
		for i := range requests.Requests {
			req := &requests.Requests[i]
			staffID, parseErr := sharedDomain.ParseID(req.StaffID)
			if parseErr == nil {
				staff, staffErr := h.staffFinder.FindByID(r.Context(), staffID)
				if staffErr == nil && staff != nil {
					req.StaffName = staff.FullName()
				}
			}
		}
	}

	data := map[string]any{
		"Title":    period.TargetPeriodLabel + " 勤務希望一覧",
		"Period":   period,
		"Requests": requests.Requests,
		"Total":    requests.Total,
	}

	if isHTMXRequest(r) {
		if err := h.templates.RenderPartial(w, "pages/shift_requests/list.html", "request-list", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.Render(w, "pages/shift_requests/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewRequest 勤務希望新規作成フォーム
func (h *RequestHandler) NewRequest(w http.ResponseWriter, r *http.Request) {
	periodID := r.PathValue("period_id")

	period, err := h.periodUseCase.GetPeriodByID(r.Context(), periodID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	data := map[string]any{
		"Title":  "勤務希望登録",
		"Period": period,
	}

	// クレームから組織IDを取得してスタッフ一覧を取得
	claims := web.GetClaimsFromContext(r.Context())
	if claims != nil && h.staffFinder != nil {
		var orgID sharedDomain.ID
		if claims.OrganizationID != nil {
			orgID = *claims.OrganizationID
		} else if period.OrganizationID != "" {
			// super_adminの場合、Periodの組織IDを使用
			parsedID, parseErr := sharedDomain.ParseID(period.OrganizationID)
			if parseErr == nil {
				orgID = parsedID
			}
		}

		var emptyID sharedDomain.ID
		if orgID != emptyID {
			staffs, staffErr := h.staffFinder.FindActiveByOrganizationID(r.Context(), orgID)
			if staffErr == nil {
				data["Staffs"] = staffs
			} else {
				h.logger.Warn("スタッフ一覧取得失敗", "error", staffErr)
			}
		}
	}

	if err := h.templates.Render(w, "pages/shift_requests/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateRequest 勤務希望作成
func (h *RequestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	periodID := r.PathValue("period_id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	input := &application.CreateShiftRequestInput{
		PeriodID:    periodID,
		StaffID:     r.FormValue("staff_id"),
		TargetDate:  r.FormValue("target_date"),
		ShiftTypeID: r.FormValue("shift_type_id"),
		RequestType: r.FormValue("request_type"),
		Priority:    r.FormValue("priority"),
		Comment:     r.FormValue("comment"),
	}

	_, err := h.requestUseCase.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/requests/"+periodID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/requests/"+periodID, http.StatusSeeOther)
}

// DeleteRequest 勤務希望削除
func (h *RequestHandler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// 削除前にリクエストを取得してリダイレクト先を特定
	request, err := h.requestUseCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := h.requestUseCase.Delete(r.Context(), id); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/requests/"+request.PeriodID, http.StatusSeeOther)
}

// ListPeriodsJSON 受付期間一覧JSON
func (h *RequestHandler) ListPeriodsJSON(w http.ResponseWriter, r *http.Request) {
	orgID := h.getOrganizationIDFromContext(r.Context())

	var periods []application.RequestPeriodOutput
	var err error

	if orgID != "" {
		periods, err = h.periodUseCase.ListPeriodsByOrganization(r.Context(), orgID)
	} else {
		// super_admin で組織未選択時は空配列
		periods = []application.RequestPeriodOutput{}
	}

	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "受付期間取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, periods)
}

// ShowPeriodJSON 受付期間詳細JSON
func (h *RequestHandler) ShowPeriodJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	period, err := h.periodUseCase.GetPeriodByID(r.Context(), id)
	if err != nil {
		if err == sharedDomain.ErrNotFound {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "受付期間が見つかりません"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "受付期間取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, period)
}

// CreatePeriodJSON 受付期間作成JSON
func (h *RequestHandler) CreatePeriodJSON(w http.ResponseWriter, r *http.Request) {
	var input application.CreateRequestPeriodInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
		return
	}

	period, err := h.periodUseCase.CreatePeriod(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, period)
}

// ListRequestsJSON 勤務希望一覧JSON
func (h *RequestHandler) ListRequestsJSON(w http.ResponseWriter, r *http.Request) {
	periodID := r.PathValue("period_id")

	requests, err := h.requestUseCase.ListByPeriod(r.Context(), periodID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "勤務希望取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, requests)
}

// CreateRequestJSON 勤務希望作成JSON
func (h *RequestHandler) CreateRequestJSON(w http.ResponseWriter, r *http.Request) {
	periodID := r.PathValue("period_id")

	var input application.CreateShiftRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
		return
	}
	input.PeriodID = periodID

	request, err := h.requestUseCase.Create(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, request)
}

// handleError エラーハンドリング
func (h *RequestHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Warn("ハンドラーエラー", "error", err, "method", r.Method, "path", r.URL.Path)

	// ドメインエラーの場合は適切なステータスコードを返す
	var domainErr *sharedDomain.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case sharedDomain.ErrCodeNotFound:
			http.Error(w, domainErr.Message, http.StatusNotFound)
		case sharedDomain.ErrCodeValidation, sharedDomain.ErrCodeInvalidInput:
			http.Error(w, domainErr.Message, http.StatusBadRequest)
		case sharedDomain.ErrCodeConflict:
			http.Error(w, domainErr.Message, http.StatusConflict)
		case sharedDomain.ErrCodeUnauthorized:
			http.Error(w, domainErr.Message, http.StatusUnauthorized)
		case sharedDomain.ErrCodeForbidden:
			http.Error(w, domainErr.Message, http.StatusForbidden)
		default:
			http.Error(w, domainErr.Message, http.StatusBadRequest)
		}
		return
	}

	// その他のエラーはメッセージをそのまま返す（開発用）
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// writeJSON JSONレスポンス書き込み
func (h *RequestHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSONエンコード失敗", "error", err)
	}
}

// getDefaultOrganizationID デフォルト組織ID取得
func (h *RequestHandler) getDefaultOrganizationID(ctx interface{}) string {
	// コンテキストから認証ユーザーのOrganizationIDを取得
	if c, ok := ctx.(context.Context); ok {
		claims := web.GetClaimsFromContext(c)
		if claims != nil && claims.OrganizationID != nil {
			return claims.OrganizationID.String()
		}
	}
	return ""
}

// getOrganizationIDFromContext コンテキストから組織ID取得
func (h *RequestHandler) getOrganizationIDFromContext(ctx context.Context) string {
	claims := web.GetClaimsFromContext(ctx)
	if claims != nil && claims.OrganizationID != nil {
		return claims.OrganizationID.String()
	}
	return ""
}

// isHTMXRequest HTMXリクエスト判定
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
