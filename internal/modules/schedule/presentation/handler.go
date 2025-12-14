// Package presentation 勤務表プレゼンテーション層
package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"shiftmaster/internal/modules/schedule/application"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/web"
)

// StaffFinder スタッフ検索インターフェース
type StaffFinder interface {
	FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]StaffInfo, error)
}

// StaffInfo スタッフ情報
type StaffInfo struct {
	ID        string
	FirstName string
	LastName  string
}

// ShiftTypeFinder シフト種別検索インターフェース
type ShiftTypeFinder interface {
	FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]ShiftTypeInfo, error)
}

// ShiftTypeInfo シフト種別情報
type ShiftTypeInfo struct {
	ID   string
	Name string
	Code string
}

// ScheduleHandler 勤務表HTTPハンドラー
type ScheduleHandler struct {
	useCase         *application.ScheduleUseCase
	staffFinder     StaffFinder
	shiftTypeFinder ShiftTypeFinder
	templates       *web.TemplateEngine
	logger          *slog.Logger
}

// NewScheduleHandler ハンドラー生成
func NewScheduleHandler(
	useCase *application.ScheduleUseCase,
	staffFinder StaffFinder,
	shiftTypeFinder ShiftTypeFinder,
	templates *web.TemplateEngine,
	logger *slog.Logger,
) *ScheduleHandler {
	return &ScheduleHandler{
		useCase:         useCase,
		staffFinder:     staffFinder,
		shiftTypeFinder: shiftTypeFinder,
		templates:       templates,
		logger:          logger,
	}
}

// RegisterRoutes ルート登録
func (h *ScheduleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /schedules", h.List)
	mux.HandleFunc("GET /schedules/new", h.New)
	mux.HandleFunc("GET /schedules/{id}", h.Show)
	mux.HandleFunc("POST /schedules", h.Create)
	mux.HandleFunc("POST /schedules/{id}/publish", h.Publish)
	mux.HandleFunc("DELETE /schedules/{id}", h.Delete)

	// API用エンドポイント
	mux.HandleFunc("GET /api/schedules", h.ListJSON)
	mux.HandleFunc("GET /api/schedules/{id}", h.ShowJSON)
	mux.HandleFunc("POST /api/schedules", h.CreateJSON)
	mux.HandleFunc("POST /api/schedules/{id}/validate", h.ValidateJSON)
}

// List 勤務表一覧ページ
func (h *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	// 組織IDを取得（super_adminの場合はnil）
	orgID := h.getDefaultOrganizationID(r.Context())

	var data map[string]any

	// super_adminで組織が選択されていない場合は空リストを表示
	if orgID == "" {
		data = map[string]any{
			"Title":            "勤務表一覧",
			"Schedules":        []any{},
			"Total":            0,
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください。ヘッダーの「組織を選択」から表示する組織を選んでください。",
		}
	} else {
		result, err := h.useCase.List(r.Context(), orgID)
		if err != nil {
			h.handleError(w, r, err)
			return
		}

		data = map[string]any{
			"Title":     "勤務表一覧",
			"Schedules": result.Schedules,
			"Total":     result.Total,
		}
	}

	if isHTMXRequest(r) {
		if err := h.templates.RenderPartial(w, "pages/schedules/list.html", "schedule-list", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.Render(w, "pages/schedules/list.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// New 新規作成フォーム
func (h *ScheduleHandler) New(w http.ResponseWriter, r *http.Request) {
	// 組織IDチェック
	orgID := h.getDefaultOrganizationID(r.Context())
	if orgID == "" {
		data := map[string]any{
			"Title":            "勤務表作成",
			"NoOrgSelected":    true,
			"NoOrgSelectedMsg": "組織を選択してください。ヘッダーの「組織を選択」から表示する組織を選んでください。",
		}
		if err := h.templates.Render(w, "pages/schedules/form.html", data); err != nil {
			h.logger.Error("テンプレートレンダリング失敗", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	now := time.Now()
	// 翌月をデフォルトとする
	nextMonth := now.AddDate(0, 1, 0)

	data := map[string]any{
		"Title":          "勤務表作成",
		"DefaultYear":    nextMonth.Year(),
		"DefaultMonth":   int(nextMonth.Month()),
		"AvailableYears": []int{now.Year(), now.Year() + 1},
	}

	if err := h.templates.Render(w, "pages/schedules/form.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Show 勤務表詳細ページ
func (h *ScheduleHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	schedule, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	// 日付一覧を生成
	dates := make([]DateInfo, schedule.DaysInMonth)
	for day := 1; day <= schedule.DaysInMonth; day++ {
		date := time.Date(schedule.TargetYear, time.Month(schedule.TargetMonth), day, 0, 0, 0, 0, time.Local)
		dates[day-1] = DateInfo{
			Day:       day,
			Date:      date.Format("2006-01-02"),
			DayOfWeek: weekdayToJapanese(date.Weekday()),
			IsWeekend: date.Weekday() == time.Saturday || date.Weekday() == time.Sunday,
		}
	}

	// スタッフ一覧とスタッフごとの日別シフトマップを作成
	var staffs []StaffInfo
	staffShiftMap := make(map[string]map[string]*application.ScheduleEntryOutput) // staffID -> date -> entry

	orgID, parseErr := sharedDomain.ParseID(schedule.OrganizationID)
	if parseErr == nil && h.staffFinder != nil {
		foundStaffs, staffErr := h.staffFinder.FindActiveByOrganizationID(r.Context(), orgID)
		if staffErr == nil {
			staffs = foundStaffs
		}
	}

	// エントリをマップに変換
	for i := range schedule.Entries {
		entry := &schedule.Entries[i]
		if staffShiftMap[entry.StaffID] == nil {
			staffShiftMap[entry.StaffID] = make(map[string]*application.ScheduleEntryOutput)
		}
		staffShiftMap[entry.StaffID][entry.TargetDate] = entry
	}

	// シフト種別一覧取得
	var shiftTypes []ShiftTypeInfo
	if h.shiftTypeFinder != nil && parseErr == nil {
		foundTypes, stErr := h.shiftTypeFinder.FindByOrganizationID(r.Context(), orgID)
		if stErr == nil {
			shiftTypes = foundTypes
		}
	}

	data := map[string]any{
		"Title":         schedule.TargetPeriodLabel + " 勤務表",
		"Schedule":      schedule,
		"Dates":         dates,
		"Staffs":        staffs,
		"StaffShiftMap": staffShiftMap,
		"ShiftTypes":    shiftTypes,
	}

	if err := h.templates.Render(w, "pages/schedules/show.html", data); err != nil {
		h.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// DateInfo 日付情報
type DateInfo struct {
	Day       int
	Date      string
	DayOfWeek string
	IsWeekend bool
}

// weekdayToJapanese 曜日を日本語に変換
func weekdayToJapanese(w time.Weekday) string {
	weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}
	return weekdays[w]
}

// Create 勤務表作成
func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 組織IDチェック
	orgID := h.getDefaultOrganizationID(r.Context())
	if orgID == "" {
		http.Error(w, "組織が選択されていません", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	targetYear, _ := strconv.Atoi(r.FormValue("target_year"))
	targetMonth, _ := strconv.Atoi(r.FormValue("target_month"))

	input := &application.CreateScheduleInput{
		OrganizationID: orgID,
		TargetYear:     targetYear,
		TargetMonth:    targetMonth,
	}

	schedule, err := h.useCase.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/schedules/"+schedule.ID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/schedules/"+schedule.ID, http.StatusSeeOther)
}

// Publish 勤務表公開
func (h *ScheduleHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := h.useCase.Publish(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/schedules/"+id)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/schedules/"+id, http.StatusSeeOther)
}

// Delete 勤務表削除
func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/schedules", http.StatusSeeOther)
}

// ListJSON 勤務表一覧JSON
func (h *ScheduleHandler) ListJSON(w http.ResponseWriter, r *http.Request) {
	orgID := h.getDefaultOrganizationID(r.Context())

	result, err := h.useCase.List(r.Context(), orgID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "勤務表取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// ShowJSON 勤務表詳細JSON
func (h *ScheduleHandler) ShowJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	schedule, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == sharedDomain.ErrNotFound {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "勤務表が見つかりません"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "勤務表取得に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, schedule)
}

// CreateJSON 勤務表作成JSON
func (h *ScheduleHandler) CreateJSON(w http.ResponseWriter, r *http.Request) {
	var input application.CreateScheduleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
		return
	}

	schedule, err := h.useCase.Create(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, schedule)
}

// ValidateJSON 勤務表検証JSON
func (h *ScheduleHandler) ValidateJSON(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.useCase.Validate(r.Context(), id)
	if err != nil {
		if err == sharedDomain.ErrNotFound {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "勤務表が見つかりません"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "検証に失敗しました"})
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// CreateEntryJSON エントリ作成JSON
func (h *ScheduleHandler) CreateEntryJSON(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.PathValue("id")

	var input application.CreateEntryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの解析に失敗しました"})
		return
	}
	input.ScheduleID = scheduleID

	entry, err := h.useCase.CreateEntry(r.Context(), &input)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, entry)
}

// CreateEntry エントリ作成
func (h *ScheduleHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "フォーム解析失敗"))
		return
	}

	input := application.CreateEntryInput{
		ScheduleID:  scheduleID,
		StaffID:     r.FormValue("staff_id"),
		TargetDate:  r.FormValue("target_date"),
		ShiftTypeID: r.FormValue("shift_type_id"),
		Note:        r.FormValue("note"),
	}

	_, err := h.useCase.CreateEntry(r.Context(), &input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/schedules/"+scheduleID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/schedules/"+scheduleID, http.StatusSeeOther)
}

// handleError エラーハンドリング
func (h *ScheduleHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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
		default:
			http.Error(w, domainErr.Message, http.StatusBadRequest)
		}
		return
	}

	// その他のエラーはメッセージをそのまま返す（開発用）
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// writeJSON JSONレスポンス書き込み
func (h *ScheduleHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSONエンコード失敗", "error", err)
	}
}

// getDefaultOrganizationID デフォルト組織ID取得
func (h *ScheduleHandler) getDefaultOrganizationID(ctx interface{}) string {
	// コンテキストから認証ユーザーのOrganizationIDを取得
	if c, ok := ctx.(context.Context); ok {
		claims := web.GetClaimsFromContext(c)
		if claims != nil && claims.OrganizationID != nil {
			return claims.OrganizationID.String()
		}
	}
	return ""
}

// isHTMXRequest HTMXリクエスト判定
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
