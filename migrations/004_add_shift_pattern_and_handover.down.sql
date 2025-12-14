-- shift_typesから追加カラムを削除
DROP INDEX IF EXISTS idx_shift_types_shift_pattern_id;
ALTER TABLE shift_types DROP COLUMN IF EXISTS handover_minutes;
ALTER TABLE shift_types DROP COLUMN IF EXISTS shift_pattern_id;

-- shift_patternsテーブル削除
DROP TABLE IF EXISTS shift_patterns;

-- rotation_type ENUM削除
DROP TYPE IF EXISTS rotation_type;
