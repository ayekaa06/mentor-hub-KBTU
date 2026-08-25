-- ============================================================
-- MentorHub — Migration 000001: DOWN (rollback)
-- ============================================================

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS mentor_advisors;
DROP TABLE IF EXISTS freshman_groups;
DROP TABLE IF EXISTS mentor_groups;
DROP TABLE IF EXISTS specialties;
DROP TABLE IF EXISTS faculties;
DROP TABLE IF EXISTS academic_years;

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TYPE IF EXISTS complaint_status;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS "pgcrypto";
