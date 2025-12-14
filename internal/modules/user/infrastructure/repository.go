// Package infrastructure ユーザーインフラストラクチャ層
package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"shiftmaster/internal/modules/user/domain"
	sharedDomain "shiftmaster/internal/shared/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// UserModel ユーザーDBモデル
type UserModel struct {
	bun.BaseModel  `bun:"table:users,alias:u"`
	ID             uuid.UUID     `bun:"id,pk,type:uuid"`
	OrganizationID uuid.NullUUID `bun:"organization_id,type:uuid"`
	Email          string        `bun:"email,notnull"`
	PasswordHash   string        `bun:"password_hash,notnull"`
	FirstName      string        `bun:"first_name,notnull"`
	LastName       string        `bun:"last_name,notnull"`
	Role           string        `bun:"role,notnull"`
	IsActive       bool          `bun:"is_active,notnull"`
	LastLoginAt    sql.NullTime  `bun:"last_login_at"`
	CreatedAt      time.Time     `bun:"created_at,notnull"`
	UpdatedAt      time.Time     `bun:"updated_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *UserModel) ToDomain() *domain.User {
	var orgID *sharedDomain.ID
	if m.OrganizationID.Valid {
		id := sharedDomain.ID(m.OrganizationID.UUID)
		orgID = &id
	}

	var lastLoginAt *time.Time
	if m.LastLoginAt.Valid {
		lastLoginAt = &m.LastLoginAt.Time
	}

	return &domain.User{
		ID:             sharedDomain.ID(m.ID),
		OrganizationID: orgID,
		Email:          m.Email,
		PasswordHash:   m.PasswordHash,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Role:           domain.UserRole(m.Role),
		IsActive:       m.IsActive,
		LastLoginAt:    lastLoginAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// UserModelFromDomain ドメインエンティティからDBモデルへ変換
func UserModelFromDomain(u *domain.User) *UserModel {
	var orgID uuid.NullUUID
	if u.OrganizationID != nil {
		orgID = uuid.NullUUID{UUID: uuid.UUID(*u.OrganizationID), Valid: true}
	}

	var lastLoginAt sql.NullTime
	if u.LastLoginAt != nil {
		lastLoginAt = sql.NullTime{Time: *u.LastLoginAt, Valid: true}
	}

	return &UserModel{
		ID:             uuid.UUID(u.ID),
		OrganizationID: orgID,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Role:           u.Role.String(),
		IsActive:       u.IsActive,
		LastLoginAt:    lastLoginAt,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// BunUserRepository Bunを使用したユーザーリポジトリ
type BunUserRepository struct {
	db *bun.DB
}

// NewBunUserRepository リポジトリ生成
func NewBunUserRepository(db *bun.DB) *BunUserRepository {
	return &BunUserRepository{db: db}
}

// FindByID IDで検索
func (r *BunUserRepository) FindByID(ctx context.Context, id sharedDomain.ID) (*domain.User, error) {
	model := new(UserModel)
	err := r.db.NewSelect().Model(model).Where("id = ?", uuid.UUID(id)).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByEmail メールアドレスで検索
func (r *BunUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	model := new(UserModel)
	err := r.db.NewSelect().Model(model).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindAll 全件取得
func (r *BunUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var models []UserModel
	err := r.db.NewSelect().Model(&models).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}
	return users, nil
}

// FindByOrganizationID 組織IDで検索
func (r *BunUserRepository) FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]domain.User, error) {
	var models []UserModel
	err := r.db.NewSelect().Model(&models).Where("organization_id = ?", uuid.UUID(organizationID)).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}
	return users, nil
}

// FindByRole ロールで検索
func (r *BunUserRepository) FindByRole(ctx context.Context, role domain.UserRole) ([]domain.User, error) {
	var models []UserModel
	err := r.db.NewSelect().Model(&models).Where("role = ?", role.String()).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}
	return users, nil
}

// Save 保存
func (r *BunUserRepository) Save(ctx context.Context, user *domain.User) error {
	model := UserModelFromDomain(user)
	_, err := r.db.NewInsert().Model(model).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	return err
}

// Delete 削除
func (r *BunUserRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*UserModel)(nil)).Where("id = ?", uuid.UUID(id)).Exec(ctx)
	return err
}

// UpdateLastLogin 最終ログイン日時更新
func (r *BunUserRepository) UpdateLastLogin(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewUpdate().Model((*UserModel)(nil)).Set("last_login_at = ?", time.Now()).Where("id = ?", uuid.UUID(id)).Exec(ctx)
	return err
}

// RefreshTokenModel リフレッシュトークンDBモデル
type RefreshTokenModel struct {
	bun.BaseModel `bun:"table:refresh_tokens,alias:rt"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	UserID        uuid.UUID `bun:"user_id,notnull,type:uuid"`
	TokenHash     string    `bun:"token_hash,notnull"`
	ExpiresAt     time.Time `bun:"expires_at,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
}

// ToDomain DBモデルからドメインエンティティへ変換
func (m *RefreshTokenModel) ToDomain() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        sharedDomain.ID(m.ID),
		UserID:    sharedDomain.ID(m.UserID),
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

// RefreshTokenModelFromDomain ドメインエンティティからDBモデルへ変換
func RefreshTokenModelFromDomain(t *domain.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		ID:        uuid.UUID(t.ID),
		UserID:    uuid.UUID(t.UserID),
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}
}

// BunRefreshTokenRepository Bunを使用したリフレッシュトークンリポジトリ
type BunRefreshTokenRepository struct {
	db *bun.DB
}

// NewBunRefreshTokenRepository リポジトリ生成
func NewBunRefreshTokenRepository(db *bun.DB) *BunRefreshTokenRepository {
	return &BunRefreshTokenRepository{db: db}
}

// FindByTokenHash トークンハッシュで検索
func (r *BunRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	model := new(RefreshTokenModel)
	err := r.db.NewSelect().Model(model).Where("token_hash = ?", tokenHash).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// FindByUserID ユーザーIDで検索
func (r *BunRefreshTokenRepository) FindByUserID(ctx context.Context, userID sharedDomain.ID) ([]domain.RefreshToken, error) {
	var models []RefreshTokenModel
	err := r.db.NewSelect().Model(&models).Where("user_id = ?", uuid.UUID(userID)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	tokens := make([]domain.RefreshToken, len(models))
	for i, m := range models {
		tokens[i] = *m.ToDomain()
	}
	return tokens, nil
}

// Save 保存
func (r *BunRefreshTokenRepository) Save(ctx context.Context, token *domain.RefreshToken) error {
	model := RefreshTokenModelFromDomain(token)
	_, err := r.db.NewInsert().Model(model).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	return err
}

// Delete 削除
func (r *BunRefreshTokenRepository) Delete(ctx context.Context, id sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*RefreshTokenModel)(nil)).Where("id = ?", uuid.UUID(id)).Exec(ctx)
	return err
}

// DeleteByUserID ユーザーIDで削除
func (r *BunRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID sharedDomain.ID) error {
	_, err := r.db.NewDelete().Model((*RefreshTokenModel)(nil)).Where("user_id = ?", uuid.UUID(userID)).Exec(ctx)
	return err
}

// DeleteExpired 期限切れトークン削除
func (r *BunRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().Model((*RefreshTokenModel)(nil)).Where("expires_at < ?", time.Now()).Exec(ctx)
	return err
}
