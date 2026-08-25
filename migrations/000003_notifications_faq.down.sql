-- ============================================================
-- MentorHub — Migration 000003: DOWN
-- ============================================================

DROP TABLE IF EXISTS questions;
DROP TRIGGER IF EXISTS trg_faq_items_updated_at ON faq_items;
DROP TABLE IF EXISTS faq_items;
DROP TABLE IF EXISTS notifications;
