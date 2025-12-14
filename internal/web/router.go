// Package web Webレイヤー ルーティング
package web

import (
	"context"
	"log/slog"
	"net/http"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// HealthChecker ヘルスチェックインターフェース
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// OrganizationFinder 組織検索インターフェース
type OrganizationFinder interface {
	FindByID(ctx context.Context, id sharedDomain.ID) (OrganizationInfo, error)
	FindAll(ctx context.Context) ([]OrganizationInfo, error)
}

// OrganizationInfo 組織情報
type OrganizationInfo struct {
	ID   sharedDomain.ID
	Name string
	Code string
}

// RouterDeps ルーター依存関係
type RouterDeps struct {
	Logger        *slog.Logger
	Templates     *TemplateEngine
	HealthChecker HealthChecker
	OrgFinder     OrganizationFinder
	Mux           *http.ServeMux
}

// Router HTTPルーター
type Router struct {
	mux       *http.ServeMux
	logger    *slog.Logger
	templates *TemplateEngine
	health    HealthChecker
	orgFinder OrganizationFinder
}

// NewRouter ルーター生成
func NewRouter(deps RouterDeps) *Router {
	mux := deps.Mux
	if mux == nil {
		mux = http.NewServeMux()
	}

	r := &Router{
		mux:       mux,
		logger:    deps.Logger,
		templates: deps.Templates,
		health:    deps.HealthChecker,
		orgFinder: deps.OrgFinder,
	}
	r.setupBaseRoutes()
	return r
}

// Mux ServeMux取得
func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

// ServeHTTP HTTPハンドラーインターフェース実装
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// setupBaseRoutes 基本ルート設定
func (r *Router) setupBaseRoutes() {
	// 静的ファイル
	r.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ヘルスチェック
	r.mux.HandleFunc("GET /health", r.healthHandler)

	// ダッシュボードは認証が必要なため DIコンテナで登録
}

// healthHandler ヘルスチェックハンドラー
func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	if r.health != nil {
		if err := r.health.HealthCheck(req.Context()); err != nil {
			r.logger.Error("ヘルスチェック失敗", "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// DashboardHandler ダッシュボードハンドラー
func (r *Router) DashboardHandler(w http.ResponseWriter, req *http.Request) {
	data := map[string]any{
		"Title": "ダッシュボード",
	}

	// 認証されたユーザー情報を取得
	claims := GetClaimsFromContext(req.Context())
	if claims != nil {
		data["User"] = map[string]any{
			"Email": claims.Email,
			"Role":  claims.Role,
		}

		// super_adminの場合 組織一覧を取得
		if claims.IsSuperAdmin() && r.orgFinder != nil {
			data["IsSuperAdmin"] = true
			orgs, err := r.orgFinder.FindAll(req.Context())
			if err == nil {
				data["Organizations"] = orgs

				// Cookie から選択中の組織ID取得
				selectedOrgID := r.getSelectedOrganizationID(req)

				// 組織が未選択の場合、最初の組織をデフォルトで選択
				if selectedOrgID == nil && len(orgs) > 0 {
					firstOrgID := orgs[0].ID
					selectedOrgID = &firstOrgID
				}

				if selectedOrgID != nil {
					data["SelectedOrganizationID"] = *selectedOrgID
					org, err := r.orgFinder.FindByID(req.Context(), *selectedOrgID)
					if err == nil {
						data["SelectedOrganizationName"] = org.Name
						data["OrganizationName"] = org.Name
					}
				}
			}
		} else if claims.OrganizationID != nil && r.orgFinder != nil {
			// 通常ユーザー 自分の組織名を取得
			org, err := r.orgFinder.FindByID(req.Context(), *claims.OrganizationID)
			if err == nil {
				data["OrganizationName"] = org.Name
			}
		}
	}

	if err := r.templates.Render(w, "pages/dashboard.html", data); err != nil {
		r.logger.Error("テンプレートレンダリング失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getSelectedOrganizationID Cookieから選択中の組織ID取得
func (r *Router) getSelectedOrganizationID(req *http.Request) *sharedDomain.ID {
	cookie, err := req.Cookie("selected_organization_id")
	if err != nil || cookie.Value == "" {
		return nil
	}
	id, err := sharedDomain.ParseID(cookie.Value)
	if err != nil {
		return nil
	}
	return &id
}

// SwitchOrganizationHandler 組織切り替えハンドラー super_admin用
func (r *Router) SwitchOrganizationHandler(w http.ResponseWriter, req *http.Request) {
	claims := GetClaimsFromContext(req.Context())
	if claims == nil || !claims.IsSuperAdmin() {
		http.Error(w, "権限がありません", http.StatusForbidden)
		return
	}

	orgID := req.PathValue("id")
	if orgID == "" {
		http.Error(w, "組織IDが必要です", http.StatusBadRequest)
		return
	}

	// 組織存在確認
	id, err := sharedDomain.ParseID(orgID)
	if err != nil {
		http.Error(w, "無効な組織IDです", http.StatusBadRequest)
		return
	}

	if r.orgFinder != nil {
		_, err := r.orgFinder.FindByID(req.Context(), id)
		if err != nil {
			http.Error(w, "組織が見つかりません", http.StatusNotFound)
			return
		}
	}

	// Cookieに選択組織ID保存
	http.SetCookie(w, &http.Cookie{
		Name:     "selected_organization_id",
		Value:    orgID,
		Path:     "/",
		MaxAge:   86400 * 30, // 30日間
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// リダイレクト
	referer := req.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, req, referer, http.StatusFound)
}
