-- ユーザーテーブル削除

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;

