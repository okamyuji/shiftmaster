-- super_adminロール追加

-- 既存のenum型に新しい値を追加
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'super_admin' BEFORE 'admin';

-- super_admin用にorganization_idをNULL許容に変更済み（users作成時に設定済み）
