-- usersテーブルからstaff_idを削除
DROP INDEX IF EXISTS idx_users_staff_id;
ALTER TABLE users DROP COLUMN IF EXISTS staff_id;

-- スタッフ所属テーブル削除
DROP TABLE IF EXISTS staff_assignments;

-- 職位テーブル削除
DROP TABLE IF EXISTS positions;

-- 職種テーブル削除
DROP TABLE IF EXISTS job_types;
