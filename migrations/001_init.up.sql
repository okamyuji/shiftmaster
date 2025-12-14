-- 初期スキーマ作成

-- 組織テーブル
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 部門テーブル
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_departments_organization ON departments(organization_id);

-- チームテーブル
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_teams_department ON teams(department_id);

-- スキルテーブル
CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_skills_organization ON skills(organization_id);

-- スタッフテーブル
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    employee_code VARCHAR(50),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    hire_date DATE,
    employment_type VARCHAR(50) NOT NULL DEFAULT 'full_time',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_staffs_team ON staffs(team_id);
CREATE INDEX idx_staffs_active ON staffs(is_active);

-- スタッフスキル中間テーブル
CREATE TABLE IF NOT EXISTS staff_skills (
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    level INT NOT NULL DEFAULT 1,
    acquired_at DATE,
    PRIMARY KEY (staff_id, skill_id)
);

-- シフト種別テーブル
CREATE TABLE IF NOT EXISTS shift_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_minutes INT NOT NULL DEFAULT 0,
    is_night_shift BOOLEAN NOT NULL DEFAULT FALSE,
    is_holiday BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shift_types_organization ON shift_types(organization_id);

-- 勤務希望期間テーブル
CREATE TABLE IF NOT EXISTS request_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    target_year INT NOT NULL,
    target_month INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    max_requests_per_staff INT NOT NULL DEFAULT 5,
    max_requests_per_day INT NOT NULL DEFAULT 3,
    is_open BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, target_year, target_month)
);

-- 勤務希望テーブル
CREATE TABLE IF NOT EXISTS shift_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id UUID NOT NULL REFERENCES request_periods(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    target_date DATE NOT NULL,
    shift_type_id UUID REFERENCES shift_types(id) ON DELETE SET NULL,
    request_type VARCHAR(20) NOT NULL DEFAULT 'preferred',
    priority VARCHAR(20) NOT NULL DEFAULT 'optional',
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shift_requests_period ON shift_requests(period_id);
CREATE INDEX idx_shift_requests_staff ON shift_requests(staff_id);
CREATE INDEX idx_shift_requests_date ON shift_requests(target_date);

-- 勤務表テーブル
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    target_year INT NOT NULL,
    target_month INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, target_year, target_month)
);

-- 勤務表エントリテーブル
CREATE TABLE IF NOT EXISTS schedule_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    target_date DATE NOT NULL,
    shift_type_id UUID REFERENCES shift_types(id) ON DELETE SET NULL,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (schedule_id, staff_id, target_date)
);

CREATE INDEX idx_schedule_entries_schedule ON schedule_entries(schedule_id);
CREATE INDEX idx_schedule_entries_staff ON schedule_entries(staff_id);
CREATE INDEX idx_schedule_entries_date ON schedule_entries(target_date);

-- 勤務実績テーブル
CREATE TABLE IF NOT EXISTS actual_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_entry_id UUID NOT NULL REFERENCES schedule_entries(id) ON DELETE CASCADE,
    actual_start_time TIMESTAMPTZ,
    actual_end_time TIMESTAMPTZ,
    actual_break_minutes INT NOT NULL DEFAULT 0,
    overtime_minutes INT NOT NULL DEFAULT 0,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_actual_records_entry ON actual_records(schedule_entry_id);

-- 更新日時自動更新トリガー関数
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 各テーブルにトリガー適用
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON departments FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_skills_updated_at BEFORE UPDATE ON skills FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_staffs_updated_at BEFORE UPDATE ON staffs FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_shift_types_updated_at BEFORE UPDATE ON shift_types FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_request_periods_updated_at BEFORE UPDATE ON request_periods FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_shift_requests_updated_at BEFORE UPDATE ON shift_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_schedule_entries_updated_at BEFORE UPDATE ON schedule_entries FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER update_actual_records_updated_at BEFORE UPDATE ON actual_records FOR EACH ROW EXECUTE FUNCTION update_updated_at();

