-- ============================================================
-- MentorHub — Migration 000002: DOWN
-- ============================================================

DROP TABLE IF EXISTS complaints;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS meeting_registrations;
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS tasks;
DROP TRIGGER IF EXISTS trg_task_templates_updated_at ON task_templates;
DROP TABLE IF EXISTS task_templates;
