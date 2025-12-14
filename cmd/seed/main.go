// Package main 初期データ投入ツール
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"shiftmaster/internal/config"
	"shiftmaster/internal/modules/user/domain"
	"shiftmaster/internal/modules/user/infrastructure"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

// OrgModel 組織DBモデル
type OrgModel struct {
	bun.BaseModel `bun:"table:organizations"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	Name          string    `bun:"name,notnull"`
	Code          string    `bun:"code"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"`
}

// DeptModel 部署DBモデル
type DeptModel struct {
	bun.BaseModel  `bun:"table:departments"`
	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code"`
	SortOrder      int       `bun:"sort_order,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

// ShiftTypeModel シフト種別DBモデル
type ShiftTypeModel struct {
	bun.BaseModel  `bun:"table:shift_types"`
	ID             uuid.UUID `bun:"id,pk,type:uuid"`
	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Code           string    `bun:"code,notnull"`
	Color          string    `bun:"color"`
	StartTime      string    `bun:"start_time,type:time"`
	EndTime        string    `bun:"end_time,type:time"`
	BreakMinutes   int       `bun:"break_minutes,notnull"`
	IsNightShift   bool      `bun:"is_night_shift,notnull"`
	IsHoliday      bool      `bun:"is_holiday,notnull"`
	SortOrder      int       `bun:"sort_order,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := config.Load()

	// データベース接続
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.URL)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("データベースクローズ失敗", "error", err)
		}
	}()

	// 接続確認
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("データベース接続確認失敗: %v", err)
	}

	ctx := context.Background()

	// 組織作成
	orgID, err := seedOrganization(ctx, db, logger)
	if err != nil {
		log.Fatalf("組織作成失敗: %v", err)
	}

	// 部署作成
	if err = seedDepartment(ctx, db, logger, orgID); err != nil {
		log.Fatalf("部署作成失敗: %v", err)
	}

	// シフト種別作成
	if err = seedShiftTypes(ctx, db, logger, orgID); err != nil {
		log.Fatalf("シフト種別作成失敗: %v", err)
	}

	// 管理者ユーザー作成
	if err = seedAdminUser(ctx, db, logger, orgID); err != nil {
		log.Fatalf("管理者ユーザー作成失敗: %v", err)
	}

	// マネージャーユーザー作成
	if err = seedManagerUser(ctx, db, logger, orgID); err != nil {
		log.Fatalf("マネージャーユーザー作成失敗: %v", err)
	}

	// 一般ユーザー作成
	if err = seedNormalUser(ctx, db, logger, orgID); err != nil {
		log.Fatalf("一般ユーザー作成失敗: %v", err)
	}

	// スーパー管理者ユーザー作成
	if err = seedSuperAdminUser(ctx, db, logger); err != nil {
		log.Fatalf("スーパー管理者ユーザー作成失敗: %v", err)
	}

	// 別組織作成（マルチテナントテスト用）
	otherOrgID, err := seedOtherOrganization(ctx, db, logger)
	if err != nil {
		log.Fatalf("別組織作成失敗: %v", err)
	}

	// 別組織の管理者ユーザー作成
	if err := seedOtherOrgAdminUser(ctx, db, logger, otherOrgID); err != nil {
		log.Fatalf("別組織管理者ユーザー作成失敗: %v", err)
	}

	logger.Info("初期データ投入完了")
}

// seedOrganization 組織作成
func seedOrganization(ctx context.Context, db *bun.DB, logger *slog.Logger) (uuid.UUID, error) {
	// 既存チェック
	var existing OrgModel
	err := db.NewSelect().Model(&existing).Where("code = ?", "DEFAULT").Scan(ctx)
	if err == nil {
		logger.Info("組織は既に存在します", "name", existing.Name, "id", existing.ID)
		return existing.ID, nil
	}
	if err != sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("組織検索失敗: %w", err)
	}

	now := time.Now()
	org := &OrgModel{
		ID:        uuid.New(),
		Name:      "サンプル病院",
		Code:      "DEFAULT",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = db.NewInsert().Model(org).Exec(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("組織保存失敗: %w", err)
	}

	logger.Info("組織作成完了", "name", org.Name, "id", org.ID)
	return org.ID, nil
}

// seedDepartment 部署作成
func seedDepartment(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	// 既存チェック
	var existing DeptModel
	err := db.NewSelect().Model(&existing).Where("organization_id = ? AND code = ?", orgID, "NURSING").Scan(ctx)
	if err == nil {
		logger.Info("部署は既に存在します", "name", existing.Name, "id", existing.ID)
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("部署検索失敗: %w", err)
	}

	now := time.Now()
	dept := &DeptModel{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           "看護部",
		Code:           "NURSING",
		SortOrder:      1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = db.NewInsert().Model(dept).Exec(ctx)
	if err != nil {
		return fmt.Errorf("部署保存失敗: %w", err)
	}

	logger.Info("部署作成完了", "name", dept.Name, "id", dept.ID)
	return nil
}

// seedShiftTypes シフト種別作成
func seedShiftTypes(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	shiftTypes := []struct {
		Name         string
		Code         string
		Color        string
		StartTime    string
		EndTime      string
		BreakMinutes int
		IsNightShift bool
		IsHoliday    bool
		SortOrder    int
	}{
		{"日勤", "D", "#4CAF50", "08:30:00", "17:30:00", 60, false, false, 1},
		{"夜勤", "N", "#3F51B5", "17:00:00", "09:00:00", 120, true, false, 2},
		{"早番", "E", "#FF9800", "06:30:00", "15:30:00", 60, false, false, 3},
		{"遅番", "L", "#9C27B0", "12:00:00", "21:00:00", 60, false, false, 4},
		{"公休", "O", "#9E9E9E", "00:00:00", "00:00:00", 0, false, true, 5},
		{"有給", "Y", "#2196F3", "00:00:00", "00:00:00", 0, false, true, 6},
	}

	now := time.Now()
	for _, st := range shiftTypes {
		// 既存チェック
		var existing ShiftTypeModel
		err := db.NewSelect().Model(&existing).Where("organization_id = ? AND code = ?", orgID, st.Code).Scan(ctx)
		if err == nil {
			logger.Info("シフト種別は既に存在します", "name", existing.Name, "code", existing.Code)
			continue
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("シフト種別検索失敗: %w", err)
		}

		model := &ShiftTypeModel{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           st.Name,
			Code:           st.Code,
			Color:          st.Color,
			StartTime:      st.StartTime,
			EndTime:        st.EndTime,
			BreakMinutes:   st.BreakMinutes,
			IsNightShift:   st.IsNightShift,
			IsHoliday:      st.IsHoliday,
			SortOrder:      st.SortOrder,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		_, err = db.NewInsert().Model(model).Exec(ctx)
		if err != nil {
			return fmt.Errorf("シフト種別保存失敗: %w", err)
		}

		logger.Info("シフト種別作成完了", "name", st.Name, "code", st.Code)
	}

	return nil
}

// seedAdminUser 管理者ユーザー作成
func seedAdminUser(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	userRepo := infrastructure.NewBunUserRepository(db)

	// 既存チェック
	existing, err := userRepo.FindByEmail(ctx, "admin@example.com")
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if existing != nil {
		logger.Info("管理者ユーザーは既に存在します", "email", "admin@example.com")
		return nil
	}

	// パスワードハッシュ生成
	password := "Password123$!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("パスワードハッシュ化失敗: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:             sharedDomain.NewID(),
		OrganizationID: &orgID,
		Email:          "admin@example.com",
		PasswordHash:   string(hashedPassword),
		FirstName:      "管理者",
		LastName:       "システム",
		Role:           domain.RoleAdmin,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("ユーザー保存失敗: %w", err)
	}

	logger.Info("管理者ユーザー作成完了",
		"email", user.Email,
		"password", password,
		"role", user.Role,
	)

	return nil
}

// seedSuperAdminUser スーパー管理者ユーザー作成
func seedSuperAdminUser(ctx context.Context, db *bun.DB, logger *slog.Logger) error {
	userRepo := infrastructure.NewBunUserRepository(db)

	// 既存チェック
	existing, err := userRepo.FindByEmail(ctx, "superadmin@example.com")
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if existing != nil {
		logger.Info("スーパー管理者ユーザーは既に存在します", "email", "superadmin@example.com")
		return nil
	}

	// パスワードハッシュ生成
	password := "SuperAdmin123$!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("パスワードハッシュ化失敗: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:             sharedDomain.NewID(),
		OrganizationID: nil, // super_adminは組織に属さない
		Email:          "superadmin@example.com",
		PasswordHash:   string(hashedPassword),
		FirstName:      "管理者",
		LastName:       "スーパー",
		Role:           domain.RoleSuperAdmin,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("ユーザー保存失敗: %w", err)
	}

	logger.Info("スーパー管理者ユーザー作成完了",
		"email", user.Email,
		"password", password,
		"role", user.Role,
	)

	return nil
}

// seedManagerUser マネージャーユーザー作成
func seedManagerUser(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	userRepo := infrastructure.NewBunUserRepository(db)

	// 既存チェック
	existing, err := userRepo.FindByEmail(ctx, "manager@example.com")
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if existing != nil {
		logger.Info("マネージャーユーザーは既に存在します", "email", "manager@example.com")
		return nil
	}

	// パスワードハッシュ生成
	password := "Manager123$!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("パスワードハッシュ化失敗: %w", err)
	}

	now := time.Now()
	id := sharedDomain.NewID()
	user := &domain.User{
		ID:             id,
		OrganizationID: &orgID,
		Email:          "manager@example.com",
		PasswordHash:   string(hashedPassword),
		FirstName:      "太郎",
		LastName:       "マネージャー",
		Role:           domain.RoleManager,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("ユーザー保存失敗: %w", err)
	}

	logger.Info("マネージャーユーザー作成完了",
		"email", user.Email,
		"password", password,
		"role", user.Role,
	)

	return nil
}

// seedNormalUser 一般ユーザー作成
func seedNormalUser(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	userRepo := infrastructure.NewBunUserRepository(db)

	// 既存チェック
	existing, err := userRepo.FindByEmail(ctx, "user@example.com")
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if existing != nil {
		logger.Info("一般ユーザーは既に存在します", "email", "user@example.com")
		return nil
	}

	// パスワードハッシュ生成
	password := "User123$!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("パスワードハッシュ化失敗: %w", err)
	}

	now := time.Now()
	id := sharedDomain.NewID()
	user := &domain.User{
		ID:             id,
		OrganizationID: &orgID,
		Email:          "user@example.com",
		PasswordHash:   string(hashedPassword),
		FirstName:      "花子",
		LastName:       "一般",
		Role:           domain.RoleUser,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("ユーザー保存失敗: %w", err)
	}

	logger.Info("一般ユーザー作成完了",
		"email", user.Email,
		"password", password,
		"role", user.Role,
	)

	return nil
}

// seedOtherOrganization 別組織作成（マルチテナントテスト用）
func seedOtherOrganization(ctx context.Context, db *bun.DB, logger *slog.Logger) (uuid.UUID, error) {
	// 既存チェック
	var existing OrgModel
	err := db.NewSelect().Model(&existing).Where("code = ?", "OTHER").Scan(ctx)
	if err == nil {
		logger.Info("別組織は既に存在します", "name", existing.Name, "id", existing.ID)
		return existing.ID, nil
	}
	if err != sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("組織検索失敗: %w", err)
	}

	now := time.Now()
	org := &OrgModel{
		ID:        uuid.New(),
		Name:      "テスト医療センター",
		Code:      "OTHER",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = db.NewInsert().Model(org).Exec(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("組織保存失敗: %w", err)
	}

	logger.Info("別組織作成完了", "name", org.Name, "id", org.ID)
	return org.ID, nil
}

// seedOtherOrgAdminUser 別組織管理者ユーザー作成
func seedOtherOrgAdminUser(ctx context.Context, db *bun.DB, logger *slog.Logger, orgID uuid.UUID) error {
	userRepo := infrastructure.NewBunUserRepository(db)

	// 既存チェック
	existing, err := userRepo.FindByEmail(ctx, "otheradmin@example.com")
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if existing != nil {
		logger.Info("別組織管理者ユーザーは既に存在します", "email", "otheradmin@example.com")
		return nil
	}

	// パスワードハッシュ生成
	password := "OtherAdmin123$!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("パスワードハッシュ化失敗: %w", err)
	}

	now := time.Now()
	id := sharedDomain.NewID()
	user := &domain.User{
		ID:             id,
		OrganizationID: &orgID,
		Email:          "otheradmin@example.com",
		PasswordHash:   string(hashedPassword),
		FirstName:      "次郎",
		LastName:       "他院管理者",
		Role:           domain.RoleAdmin,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("ユーザー保存失敗: %w", err)
	}

	logger.Info("別組織管理者ユーザー作成完了",
		"email", user.Email,
		"password", password,
		"role", user.Role,
		"organization_id", orgID,
	)

	return nil
}
