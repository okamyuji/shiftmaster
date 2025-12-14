-- 職種テーブル
CREATE TABLE job_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#6B7280',
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_job_types_organization_id ON job_types(organization_id);
CREATE INDEX idx_job_types_is_active ON job_types(is_active);

-- 職位テーブル
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_positions_organization_id ON positions(organization_id);
CREATE INDEX idx_positions_is_active ON positions(is_active);
CREATE INDEX idx_positions_level ON positions(level);

-- スタッフ所属テーブル（多対多）
CREATE TABLE staff_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    job_type_id UUID REFERENCES job_types(id) ON DELETE SET NULL,
    position_id UUID REFERENCES positions(id) ON DELETE SET NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_at_least_one_assignment CHECK (
        team_id IS NOT NULL OR job_type_id IS NOT NULL OR position_id IS NOT NULL
    )
);

CREATE INDEX idx_staff_assignments_staff_id ON staff_assignments(staff_id);
CREATE INDEX idx_staff_assignments_team_id ON staff_assignments(team_id);
CREATE INDEX idx_staff_assignments_job_type_id ON staff_assignments(job_type_id);
CREATE INDEX idx_staff_assignments_position_id ON staff_assignments(position_id);
CREATE INDEX idx_staff_assignments_is_primary ON staff_assignments(is_primary);
CREATE INDEX idx_staff_assignments_date_range ON staff_assignments(start_date, end_date);

-- usersテーブルにstaff_idを追加（ユーザーとスタッフの関連）
ALTER TABLE users ADD COLUMN staff_id UUID REFERENCES staffs(id) ON DELETE SET NULL;
CREATE INDEX idx_users_staff_id ON users(staff_id);

-- 初期職種データ（サンプル病院用）
INSERT INTO job_types (organization_id, name, code, description, color, sort_order)
SELECT
    id,
    '看護師',
    'NS',
    '正看護師',
    '#3B82F6',
    1
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

INSERT INTO job_types (organization_id, name, code, description, color, sort_order)
SELECT
    id,
    '准看護師',
    'ANS',
    '准看護師',
    '#60A5FA',
    2
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

INSERT INTO job_types (organization_id, name, code, description, color, sort_order)
SELECT
    id,
    '介護士',
    'CW',
    '介護福祉士',
    '#10B981',
    3
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

-- 初期職位データ（サンプル病院用）
INSERT INTO positions (organization_id, name, code, description, level, sort_order)
SELECT
    id,
    '師長',
    'DIR',
    '看護師長',
    1,
    1
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

INSERT INTO positions (organization_id, name, code, description, level, sort_order)
SELECT
    id,
    '主任',
    'MGR',
    '看護主任',
    2,
    2
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

INSERT INTO positions (organization_id, name, code, description, level, sort_order)
SELECT
    id,
    'リーダー',
    'LDR',
    'チームリーダー',
    3,
    3
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;

INSERT INTO positions (organization_id, name, code, description, level, sort_order)
SELECT
    id,
    '一般',
    'STF',
    '一般スタッフ',
    4,
    4
FROM organizations WHERE name = 'サンプル病院'
ON CONFLICT DO NOTHING;
