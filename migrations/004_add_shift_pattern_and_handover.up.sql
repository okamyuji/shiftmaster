-- シフトパターン（直）テーブル
-- 1直、2直、3直などの交代勤務グループを管理

-- ローテーション種別のENUM型
CREATE TYPE rotation_type AS ENUM (
    'two_shift',           -- 二交代制
    'three_shift',         -- 三交代制
    'four_shift_three_team', -- 四直三交代制
    'custom'               -- カスタム
);

-- シフトパターンテーブル
CREATE TABLE shift_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    rotation_type rotation_type NOT NULL DEFAULT 'custom',
    color VARCHAR(7) NOT NULL DEFAULT '#6366F1',
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, code)
);

-- インデックス
CREATE INDEX idx_shift_patterns_organization_id ON shift_patterns(organization_id);
CREATE INDEX idx_shift_patterns_is_active ON shift_patterns(is_active);

-- shift_typesテーブルに新しいカラムを追加
-- shift_pattern_id: シフトパターンへの参照（オプション）
ALTER TABLE shift_types
    ADD COLUMN shift_pattern_id UUID REFERENCES shift_patterns(id) ON DELETE SET NULL;

-- handover_minutes: 申し送り時間（分）
-- 病院での引き継ぎ時間や、工場での交代準備時間
ALTER TABLE shift_types
    ADD COLUMN handover_minutes INTEGER NOT NULL DEFAULT 0;

-- インデックス
CREATE INDEX idx_shift_types_shift_pattern_id ON shift_types(shift_pattern_id);

-- コメント
COMMENT ON TABLE shift_patterns IS 'シフトパターン（直）- 1直、2直、3直などの交代勤務グループ';
COMMENT ON COLUMN shift_patterns.rotation_type IS 'ローテーション種別: two_shift=二交代, three_shift=三交代, four_shift_three_team=四直三交代, custom=カスタム';
COMMENT ON COLUMN shift_types.shift_pattern_id IS '所属するシフトパターン（直）への参照';
COMMENT ON COLUMN shift_types.handover_minutes IS '申し送り・引き継ぎ時間（分）- 勤務開始前のオーバーラップ時間';
