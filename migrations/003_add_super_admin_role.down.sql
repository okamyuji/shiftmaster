-- super_adminロール削除
-- 注意: PostgreSQLではenum値の削除は複雑なため、このマイグレーションはロールバック不可

-- 既存のsuper_adminユーザーをadminに変更
UPDATE users SET role = 'admin' WHERE role = 'super_admin';

-- enum値の削除は別途手動で行う必要あり
-- ALTER TYPE user_role では既存の値を削除できないため
