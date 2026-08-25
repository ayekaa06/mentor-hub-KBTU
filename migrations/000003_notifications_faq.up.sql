-- ============================================================
-- MentorHub — Migration 000003: Notifications, FAQ, Questions
-- ============================================================

-- ── Notifications ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS notifications (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    body       TEXT,
    is_read    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(user_id, is_read);

-- ── FAQ ───────────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS faq_items (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    question   TEXT        NOT NULL,
    answer     TEXT        NOT NULL,
    order_num  INTEGER     NOT NULL DEFAULT 0,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_by UUID        NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_faq_items_updated_at
    BEFORE UPDATE ON faq_items
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX idx_faq_items_order ON faq_items(order_num) WHERE is_active = TRUE;

-- ── Questions (freshman → mentor) ─────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS questions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    freshman_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mentor_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body        TEXT        NOT NULL,
    answer      TEXT,
    answered_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_freshman_id ON questions(freshman_id);
CREATE INDEX idx_questions_mentor_id   ON questions(mentor_id);
CREATE INDEX idx_questions_unanswered  ON questions(mentor_id) WHERE answer IS NULL;
