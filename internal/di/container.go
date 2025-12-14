// Package di 依存性注入コンテナ
package di

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"shiftmaster/internal/config"
	authApp "shiftmaster/internal/modules/auth/application"
	authInfra "shiftmaster/internal/modules/auth/infrastructure"
	authPres "shiftmaster/internal/modules/auth/presentation"
	requestApp "shiftmaster/internal/modules/request/application"
	requestDomain "shiftmaster/internal/modules/request/domain"
	requestInfra "shiftmaster/internal/modules/request/infrastructure"
	requestPres "shiftmaster/internal/modules/request/presentation"
	scheduleApp "shiftmaster/internal/modules/schedule/application"
	scheduleDomain "shiftmaster/internal/modules/schedule/domain"
	scheduleInfra "shiftmaster/internal/modules/schedule/infrastructure"
	schedulePres "shiftmaster/internal/modules/schedule/presentation"
	shiftApp "shiftmaster/internal/modules/shift/application"
	shiftDomain "shiftmaster/internal/modules/shift/domain"
	shiftInfra "shiftmaster/internal/modules/shift/infrastructure"
	shiftPres "shiftmaster/internal/modules/shift/presentation"
	staffApp "shiftmaster/internal/modules/staff/application"
	staffDomain "shiftmaster/internal/modules/staff/domain"
	staffInfra "shiftmaster/internal/modules/staff/infrastructure"
	staffPres "shiftmaster/internal/modules/staff/presentation"
	userApp "shiftmaster/internal/modules/user/application"
	userDomain "shiftmaster/internal/modules/user/domain"
	userInfra "shiftmaster/internal/modules/user/infrastructure"
	userPres "shiftmaster/internal/modules/user/presentation"
	sharedDomain "shiftmaster/internal/shared/domain"
	"shiftmaster/internal/shared/infrastructure"
	"shiftmaster/internal/web"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Container 依存性注入コンテナ
type Container struct {
	// Config アプリケーション設定
	Config *config.Config
	// Logger 構造化ロガー
	Logger *slog.Logger
	// DB データベース接続
	DB *bun.DB
	// Templates テンプレートエンジン
	Templates *web.TemplateEngine
	// Router HTTPルーター
	Router *web.Router
	// TokenService トークンサービス
	TokenService *authInfra.JWTTokenService

	// Repositories
	StaffRepo         staffDomain.StaffRepository
	TeamRepo          staffDomain.TeamRepository
	DepartmentRepo    staffDomain.DepartmentRepository
	OrganizationRepo  staffDomain.OrganizationRepository
	UserRepo          userDomain.UserRepository
	RefreshTokenRepo  userDomain.RefreshTokenRepository
	ShiftTypeRepo     shiftDomain.ShiftTypeRepository
	ScheduleRepo      scheduleDomain.ScheduleRepository
	ScheduleEntryRepo scheduleDomain.ScheduleEntryRepository
	RequestPeriodRepo requestDomain.RequestPeriodRepository
	ShiftRequestRepo  requestDomain.ShiftRequestRepository

	// UseCases
	StaffUseCase         *staffApp.StaffUseCase
	UserUseCase          *userApp.UserUseCase
	AuthUseCase          *authApp.AuthUseCase
	ShiftTypeUseCase     *shiftApp.ShiftTypeUseCase
	ScheduleUseCase      *scheduleApp.ScheduleUseCase
	RequestPeriodUseCase *requestApp.RequestPeriodUseCase
	ShiftRequestUseCase  *requestApp.ShiftRequestUseCase

	// Handlers
	StaffHandler     *staffPres.StaffHandler
	TeamHandler      *staffPres.TeamHandler
	UserHandler      *userPres.UserHandler
	AuthHandler      *authPres.AuthHandler
	ShiftTypeHandler *shiftPres.ShiftTypeHandler
	ScheduleHandler  *schedulePres.ScheduleHandler
	RequestHandler   *requestPres.RequestHandler
}

// NewContainer コンテナ生成
func NewContainer(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	// データベース接続
	db, err := newDatabase(cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// テンプレートエンジン初期化
	templates := web.NewTemplateEngine("internal/web/templates")
	if err := templates.Load(); err != nil {
		return nil, err
	}

	// JWT秘密鍵取得
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "shiftmaster-default-secret-key-change-in-production"
	}
	jwtConfig := authInfra.DefaultJWTConfig(jwtSecret)
	tokenService := authInfra.NewJWTTokenService(jwtConfig)

	// リポジトリ初期化
	staffRepo := staffInfra.NewPostgresStaffRepository(db)
	teamRepo := staffInfra.NewPostgresTeamRepository(db)
	departmentRepo := staffInfra.NewPostgresDepartmentRepository(db)
	organizationRepo := staffInfra.NewPostgresOrganizationRepository(db)
	userRepo := userInfra.NewBunUserRepository(db)
	refreshTokenRepo := userInfra.NewBunRefreshTokenRepository(db)
	shiftTypeRepo := shiftInfra.NewPostgresShiftTypeRepository(db)
	scheduleRepo := scheduleInfra.NewPostgresScheduleRepository(db)
	scheduleEntryRepo := scheduleInfra.NewPostgresScheduleEntryRepository(db)
	requestPeriodRepo := requestInfra.NewPostgresRequestPeriodRepository(db)
	shiftRequestRepo := requestInfra.NewPostgresShiftRequestRepository(db)

	// ユースケース初期化
	staffUseCase := staffApp.NewStaffUseCase(staffRepo, teamRepo, departmentRepo, logger)
	userUseCase := userApp.NewUserUseCase(userRepo, refreshTokenRepo, logger)
	authUseCase := authApp.NewAuthUseCase(userRepo, refreshTokenRepo, tokenService, logger)
	shiftTypeUseCase := shiftApp.NewShiftTypeUseCase(shiftTypeRepo, logger)
	scheduleUseCase := scheduleApp.NewScheduleUseCase(scheduleRepo, scheduleEntryRepo, shiftTypeRepo, staffRepo, nil, logger)
	requestPeriodUseCase := requestApp.NewRequestPeriodUseCase(requestPeriodRepo, shiftRequestRepo, logger)
	shiftRequestUseCase := requestApp.NewShiftRequestUseCase(shiftRequestRepo, requestPeriodRepo, logger)

	// コンテナ生成 ルーターは後で設定
	container := &Container{
		Config:               cfg,
		Logger:               logger,
		DB:                   db,
		Templates:            templates,
		TokenService:         tokenService,
		StaffRepo:            staffRepo,
		TeamRepo:             teamRepo,
		DepartmentRepo:       departmentRepo,
		OrganizationRepo:     organizationRepo,
		UserRepo:             userRepo,
		RefreshTokenRepo:     refreshTokenRepo,
		ShiftTypeRepo:        shiftTypeRepo,
		ScheduleRepo:         scheduleRepo,
		ScheduleEntryRepo:    scheduleEntryRepo,
		RequestPeriodRepo:    requestPeriodRepo,
		ShiftRequestRepo:     shiftRequestRepo,
		StaffUseCase:         staffUseCase,
		UserUseCase:          userUseCase,
		AuthUseCase:          authUseCase,
		ShiftTypeUseCase:     shiftTypeUseCase,
		ScheduleUseCase:      scheduleUseCase,
		RequestPeriodUseCase: requestPeriodUseCase,
		ShiftRequestUseCase:  shiftRequestUseCase,
	}

	// 組織ファインダーアダプター
	orgFinder := &organizationFinderAdapter{repo: organizationRepo}

	// ルーター初期化
	mux := http.NewServeMux()
	router := web.NewRouter(web.RouterDeps{
		Logger:        logger,
		Templates:     templates,
		HealthChecker: container,
		OrgFinder:     orgFinder,
		Mux:           mux,
	})
	container.Router = router

	// ハンドラー初期化
	staffHandler := staffPres.NewStaffHandler(staffUseCase, teamRepo, templates, logger)
	container.StaffHandler = staffHandler

	teamHandler := staffPres.NewTeamHandler(teamRepo, departmentRepo, templates, logger)
	container.TeamHandler = teamHandler

	shiftTypeHandler := shiftPres.NewShiftTypeHandler(shiftTypeUseCase, templates, logger)
	container.ShiftTypeHandler = shiftTypeHandler

	// スタッフ・シフト種別検索アダプター（勤務表用）
	scheduleStaffFinder := &scheduleStaffFinderAdapter{repo: staffRepo}
	shiftTypeFinder := &shiftTypeFinderAdapter{repo: shiftTypeRepo}
	scheduleHandler := schedulePres.NewScheduleHandler(scheduleUseCase, scheduleStaffFinder, shiftTypeFinder, templates, logger)
	container.ScheduleHandler = scheduleHandler

	// スタッフ検索アダプター（勤務希望用）
	staffFinder := &staffFinderAdapter{repo: staffRepo}
	requestHandler := requestPres.NewRequestHandler(requestPeriodUseCase, shiftRequestUseCase, staffFinder, templates, logger)
	container.RequestHandler = requestHandler

	userHandler := userPres.NewUserHandler(userUseCase, templates, logger)
	container.UserHandler = userHandler

	authHandler := authPres.NewAuthHandler(authUseCase, templates, logger)
	container.AuthHandler = authHandler

	// 認証ルート登録
	container.registerAuthRoutes(mux)
	container.registerAdminRoutes(mux)
	container.registerProtectedRoutes(mux)

	return container, nil
}

// registerAuthRoutes 認証ルート登録
func (c *Container) registerAuthRoutes(mux *http.ServeMux) {
	// ログインページ（認証済みならリダイレクト）
	mux.Handle("GET /login", web.Chain(
		http.HandlerFunc(c.AuthHandler.LoginPage),
		web.RedirectIfAuthenticated(c.TokenService, "/"),
	))

	// ログイン処理
	mux.HandleFunc("POST /login", c.AuthHandler.Login)

	// ログインAPI
	mux.HandleFunc("POST /api/auth/login", c.AuthHandler.LoginAPI)

	// トークンリフレッシュ
	mux.HandleFunc("POST /api/auth/refresh", c.AuthHandler.Refresh)

	// ログアウト
	mux.HandleFunc("POST /logout", c.AuthHandler.Logout)
	mux.HandleFunc("POST /api/auth/logout", c.AuthHandler.Logout)

	// 現在のユーザー情報
	mux.Handle("GET /api/auth/me", web.Chain(
		http.HandlerFunc(c.AuthHandler.Me),
		web.Auth(c.TokenService, c.Logger),
	))
}

// registerAdminRoutes 管理画面ルート登録
func (c *Container) registerAdminRoutes(mux *http.ServeMux) {
	// 管理者認証ミドルウェア Chain適用順序 Auth -> RequireAdmin -> handler
	adminAuth := func(h http.Handler) http.Handler {
		return web.Chain(h,
			web.Auth(c.TokenService, c.Logger),
			web.RequireAdmin(),
		)
	}

	// スーパー管理者のみ Chain適用順序 Auth -> RequireSuperAdmin -> handler
	superAdminAuth := func(h http.Handler) http.Handler {
		return web.Chain(h,
			web.Auth(c.TokenService, c.Logger),
			web.RequireSuperAdmin(),
		)
	}

	// 組織切り替え super_adminのみ
	mux.Handle("GET /admin/switch-organization/{id}", superAdminAuth(http.HandlerFunc(c.Router.SwitchOrganizationHandler)))

	// ユーザー管理
	mux.Handle("GET /admin/users", adminAuth(http.HandlerFunc(c.UserHandler.ListUsers)))
	mux.Handle("GET /admin/users/new", adminAuth(http.HandlerFunc(c.UserHandler.NewUserForm)))
	mux.Handle("POST /admin/users", adminAuth(http.HandlerFunc(c.UserHandler.CreateUser)))
	mux.Handle("GET /admin/users/{id}/edit", adminAuth(http.HandlerFunc(c.UserHandler.EditUserForm)))
	mux.Handle("PUT /admin/users/{id}", adminAuth(http.HandlerFunc(c.UserHandler.UpdateUser)))
	mux.Handle("DELETE /admin/users/{id}", adminAuth(http.HandlerFunc(c.UserHandler.DeleteUser)))
}

// registerProtectedRoutes 認証必須ルート登録
func (c *Container) registerProtectedRoutes(mux *http.ServeMux) {
	// 認証ミドルウェア
	auth := func(h http.Handler) http.Handler {
		return web.Chain(h, web.Auth(c.TokenService, c.Logger))
	}

	// ダッシュボード（認証必須）
	mux.Handle("GET /{$}", auth(http.HandlerFunc(c.Router.DashboardHandler)))

	// スタッフ管理
	mux.Handle("GET /staffs", auth(http.HandlerFunc(c.StaffHandler.List)))
	mux.Handle("GET /staffs/new", auth(http.HandlerFunc(c.StaffHandler.New)))
	mux.Handle("POST /staffs", auth(http.HandlerFunc(c.StaffHandler.Create)))
	mux.Handle("GET /staffs/{id}", auth(http.HandlerFunc(c.StaffHandler.Show)))
	mux.Handle("GET /staffs/{id}/edit", auth(http.HandlerFunc(c.StaffHandler.Edit)))
	mux.Handle("PUT /staffs/{id}", auth(http.HandlerFunc(c.StaffHandler.Update)))
	mux.Handle("DELETE /staffs/{id}", auth(http.HandlerFunc(c.StaffHandler.Delete)))

	// チーム管理
	mux.Handle("GET /teams", auth(http.HandlerFunc(c.TeamHandler.List)))
	mux.Handle("GET /teams/new", auth(http.HandlerFunc(c.TeamHandler.New)))
	mux.Handle("POST /teams", auth(http.HandlerFunc(c.TeamHandler.Create)))
	mux.Handle("GET /teams/{id}/edit", auth(http.HandlerFunc(c.TeamHandler.Edit)))
	mux.Handle("PUT /teams/{id}", auth(http.HandlerFunc(c.TeamHandler.Update)))
	mux.Handle("DELETE /teams/{id}", auth(http.HandlerFunc(c.TeamHandler.Delete)))

	// シフト種別管理
	mux.Handle("GET /shifts", auth(http.HandlerFunc(c.ShiftTypeHandler.List)))
	mux.Handle("GET /shifts/new", auth(http.HandlerFunc(c.ShiftTypeHandler.New)))
	mux.Handle("POST /shifts", auth(http.HandlerFunc(c.ShiftTypeHandler.Create)))
	mux.Handle("GET /shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.Show)))
	mux.Handle("GET /shifts/{id}/edit", auth(http.HandlerFunc(c.ShiftTypeHandler.Edit)))
	mux.Handle("PUT /shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.Update)))
	mux.Handle("DELETE /shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.Delete)))

	// 勤務表管理
	mux.Handle("GET /schedules", auth(http.HandlerFunc(c.ScheduleHandler.List)))
	mux.Handle("GET /schedules/new", auth(http.HandlerFunc(c.ScheduleHandler.New)))
	mux.Handle("POST /schedules", auth(http.HandlerFunc(c.ScheduleHandler.Create)))
	mux.Handle("GET /schedules/{id}", auth(http.HandlerFunc(c.ScheduleHandler.Show)))
	mux.Handle("POST /schedules/{id}/entries", auth(http.HandlerFunc(c.ScheduleHandler.CreateEntry)))
	mux.Handle("POST /schedules/{id}/publish", auth(http.HandlerFunc(c.ScheduleHandler.Publish)))
	mux.Handle("DELETE /schedules/{id}", auth(http.HandlerFunc(c.ScheduleHandler.Delete)))

	// 勤務希望管理
	mux.Handle("GET /requests", auth(http.HandlerFunc(c.RequestHandler.ListPeriods)))
	mux.Handle("GET /requests/new", auth(http.HandlerFunc(c.RequestHandler.NewPeriod)))
	mux.Handle("POST /requests", auth(http.HandlerFunc(c.RequestHandler.CreatePeriod)))
	mux.Handle("GET /requests/{id}", auth(http.HandlerFunc(c.RequestHandler.ShowPeriod)))
	mux.Handle("POST /requests/{id}/open", auth(http.HandlerFunc(c.RequestHandler.OpenPeriod)))
	mux.Handle("POST /requests/{id}/close", auth(http.HandlerFunc(c.RequestHandler.ClosePeriod)))
	mux.Handle("GET /requests/{period_id}/entries", auth(http.HandlerFunc(c.RequestHandler.ListRequests)))
	mux.Handle("GET /requests/{period_id}/entries/new", auth(http.HandlerFunc(c.RequestHandler.NewRequest)))
	mux.Handle("POST /requests/{period_id}/entries", auth(http.HandlerFunc(c.RequestHandler.CreateRequest)))
	mux.Handle("DELETE /requests/entries/{id}", auth(http.HandlerFunc(c.RequestHandler.DeleteRequest)))

	// API シフト種別
	mux.Handle("GET /api/shifts", auth(http.HandlerFunc(c.ShiftTypeHandler.ListJSON)))
	mux.Handle("GET /api/shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.ShowJSON)))
	mux.Handle("POST /api/shifts", auth(http.HandlerFunc(c.ShiftTypeHandler.CreateJSON)))
	mux.Handle("PUT /api/shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.UpdateJSON)))
	mux.Handle("DELETE /api/shifts/{id}", auth(http.HandlerFunc(c.ShiftTypeHandler.DeleteJSON)))
}

// Close リソース解放
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// newDatabase Bun ORMデータベース接続生成
func newDatabase(databaseURL string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseURL)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// 接続確認
	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

// HealthCheck データベースヘルスチェック
func (c *Container) HealthCheck(ctx context.Context) error {
	return c.DB.PingContext(ctx)
}

// Transaction トランザクション実行ヘルパー
func (c *Container) Transaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return infrastructure.RunInTransaction(ctx, c.DB, fn)
}

// organizationFinderAdapter 組織検索アダプター
type organizationFinderAdapter struct {
	repo staffDomain.OrganizationRepository
}

// FindByID 組織をIDで検索
func (a *organizationFinderAdapter) FindByID(ctx context.Context, id sharedDomain.ID) (web.OrganizationInfo, error) {
	org, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return web.OrganizationInfo{}, err
	}
	if org == nil {
		return web.OrganizationInfo{}, nil
	}
	return web.OrganizationInfo{
		ID:   org.ID,
		Name: org.Name,
		Code: org.Code,
	}, nil
}

// FindAll 全組織取得
func (a *organizationFinderAdapter) FindAll(ctx context.Context) ([]web.OrganizationInfo, error) {
	orgs, err := a.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]web.OrganizationInfo, len(orgs))
	for i, org := range orgs {
		result[i] = web.OrganizationInfo{
			ID:   org.ID,
			Name: org.Name,
			Code: org.Code,
		}
	}
	return result, nil
}

// staffFinderAdapter スタッフ検索アダプター
type staffFinderAdapter struct {
	repo staffDomain.StaffRepository
}

// FindActiveByOrganizationID 組織IDで有効スタッフを検索
func (a *staffFinderAdapter) FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]requestPres.StaffInfo, error) {
	staffs, err := a.repo.FindActiveByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	result := make([]requestPres.StaffInfo, len(staffs))
	for i, s := range staffs {
		result[i] = requestPres.StaffInfo{
			ID:        s.ID.String(),
			FirstName: s.FirstName,
			LastName:  s.LastName,
		}
	}
	return result, nil
}

// FindByID スタッフIDで検索
func (a *staffFinderAdapter) FindByID(ctx context.Context, id sharedDomain.ID) (*requestPres.StaffInfo, error) {
	staff, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, nil
	}
	return &requestPres.StaffInfo{
		ID:        staff.ID.String(),
		FirstName: staff.FirstName,
		LastName:  staff.LastName,
	}, nil
}

// scheduleStaffFinderAdapter スタッフ検索アダプター（勤務表用）
type scheduleStaffFinderAdapter struct {
	repo staffDomain.StaffRepository
}

// FindActiveByOrganizationID 組織IDで有効スタッフを検索
func (a *scheduleStaffFinderAdapter) FindActiveByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]schedulePres.StaffInfo, error) {
	staffs, err := a.repo.FindActiveByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	result := make([]schedulePres.StaffInfo, len(staffs))
	for i, s := range staffs {
		result[i] = schedulePres.StaffInfo{
			ID:        s.ID.String(),
			FirstName: s.FirstName,
			LastName:  s.LastName,
		}
	}
	return result, nil
}

// shiftTypeFinderAdapter シフト種別検索アダプター
type shiftTypeFinderAdapter struct {
	repo shiftDomain.ShiftTypeRepository
}

// FindByOrganizationID 組織IDでシフト種別を検索
func (a *shiftTypeFinderAdapter) FindByOrganizationID(ctx context.Context, orgID sharedDomain.ID) ([]schedulePres.ShiftTypeInfo, error) {
	shiftTypes, err := a.repo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	result := make([]schedulePres.ShiftTypeInfo, len(shiftTypes))
	for i, st := range shiftTypes {
		result[i] = schedulePres.ShiftTypeInfo{
			ID:   st.ID.String(),
			Name: st.Name,
			Code: st.Code,
		}
	}
	return result, nil
}
