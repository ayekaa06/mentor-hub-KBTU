-- ============================================================
-- MentorHub — Migration 000001: Core Schema
-- Таблицы: users, academic_years, faculties, specialties,
--          mentor_groups, freshman_groups, mentor_advisors,
--          refresh_tokens
-- ============================================================

-- UUID генерация через pgcrypto
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── Enum Types ────────────────────────────────────────────────────────────────

CREATE TYPE user_role AS ENUM ('head', 'advisor', 'mentor', 'freshman');
CREATE TYPE task_status AS ENUM ('pending', 'submitted', 'approved', 'rejected');
CREATE TYPE complaint_status AS ENUM ('open', 'in_review', 'resolved', 'dismissed');

-- ── Trigger: auto-update updated_at ──────────────────────────────────────────

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ── Users ─────────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          user_role   NOT NULL,
    first_name    VARCHAR(100) NOT NULL DEFAULT '',
    last_name     VARCHAR(100) NOT NULL DEFAULT '',
    avatar_url    TEXT,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX idx_users_email    ON users(email);
CREATE INDEX idx_users_role     ON users(role);
CREATE INDEX idx_users_active   ON users(is_active);

-- ── Academic Years ────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS academic_years (
    id         UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(50) NOT NULL,          -- "2025-2026"
    start_date DATE        NOT NULL,
    end_date   DATE        NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_year_dates CHECK (end_date > start_date)
);

-- ── Faculties ─────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS faculties (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(50)  UNIQUE NOT NULL,  -- "CS", "ENG", "MATH"
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ── Specialties ───────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS specialties (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    faculty_id UUID        NOT NULL REFERENCES faculties(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_specialties_faculty_id ON specialties(faculty_id);

-- ── Mentor Groups ─────────────────────────────────────────────────────────────
-- Одна группа = один ментор + один учебный год

CREATE TABLE IF NOT EXISTS mentor_groups (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    academic_year_id UUID        NOT NULL REFERENCES academic_years(id),
    specialty_id     UUID        REFERENCES specialties(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mentor_groups_mentor_id ON mentor_groups(mentor_id);
CREATE INDEX idx_mentor_groups_year_id   ON mentor_groups(academic_year_id);

-- ── Freshman → Group (M2M) ────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS freshman_groups (
    freshman_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id    UUID        NOT NULL REFERENCES mentor_groups(id) ON DELETE CASCADE,
    joined_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (freshman_id, group_id)
);

CREATE INDEX idx_freshman_groups_group_id     ON freshman_groups(group_id);
CREATE INDEX idx_freshman_groups_freshman_id  ON freshman_groups(freshman_id);

-- ── Mentor → Advisor (M2M) ───────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS mentor_advisors (
    mentor_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    advisor_id  UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (mentor_id, advisor_id)
);

CREATE INDEX idx_mentor_advisors_advisor_id ON mentor_advisors(advisor_id);
CREATE INDEX idx_mentor_advisors_mentor_id  ON mentor_advisors(mentor_id);

-- ── Refresh Tokens ────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- ── Seed: Head пользователь ───────────────────────────────────────────────────
-- Пароль: Admin1234!  (bcrypt cost 12)
INSERT INTO users (email, password_hash, role, first_name, last_name)
VALUES (
    'head@mentorhub.kz',
    '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2',
    'head',
    'System',
    'Head'
) ON CONFLICT (email) DO NOTHING;
