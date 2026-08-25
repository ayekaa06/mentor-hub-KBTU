-- ============================================================
-- MentorHub — Migration 000002: Tasks, Meetings, Complaints
-- Таблицы: task_templates, tasks, meetings, meeting_registrations,
--          announcements, complaints
-- ============================================================

-- ── Task Templates (создаёт Head) ────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS task_templates (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    due_days    INTEGER      NOT NULL DEFAULT 7 CHECK (due_days > 0),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by  UUID         NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_task_templates_updated_at
    BEFORE UPDATE ON task_templates
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ── Tasks (назначаются freshman) ──────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS tasks (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id  UUID        NOT NULL REFERENCES task_templates(id),
    freshman_id  UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_by  UUID        NOT NULL REFERENCES users(id),
    status       task_status NOT NULL DEFAULT 'pending',
    proof_url    TEXT,                        -- ссылка / URL файла
    comment      TEXT,                        -- комментарий при отклонении
    due_date     TIMESTAMPTZ NOT NULL,
    submitted_at TIMESTAMPTZ,
    reviewed_at  TIMESTAMPTZ,
    reviewed_by  UUID        REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_freshman_id ON tasks(freshman_id);
CREATE INDEX idx_tasks_status      ON tasks(status);
CREATE INDEX idx_tasks_due_date    ON tasks(due_date);
CREATE INDEX idx_tasks_template_id ON tasks(template_id);

-- ── Meetings ──────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS meetings (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id     UUID         NOT NULL REFERENCES mentor_groups(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    scheduled_at TIMESTAMPTZ  NOT NULL,
    held         BOOLEAN      NOT NULL DEFAULT FALSE,
    notes        TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meetings_mentor_id    ON meetings(mentor_id);
CREATE INDEX idx_meetings_group_id     ON meetings(group_id);
CREATE INDEX idx_meetings_scheduled_at ON meetings(scheduled_at);

-- ── Meeting Registrations (freshman ↔ meeting) ────────────────────────────────

CREATE TABLE IF NOT EXISTS meeting_registrations (
    freshman_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meeting_id    UUID        NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (freshman_id, meeting_id)
);

CREATE INDEX idx_meeting_registrations_meeting_id ON meeting_registrations(meeting_id);

-- ── Announcements ─────────────────────────────────────────────────────────────
-- group_id = NULL означает глобальное объявление от Head

CREATE TABLE IF NOT EXISTS announcements (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id  UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id   UUID         REFERENCES mentor_groups(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    body       TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_announcements_author_id ON announcements(author_id);
CREATE INDEX idx_announcements_group_id  ON announcements(group_id);

-- ── Complaints ────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS complaints (
    id          UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    filed_by    UUID             NOT NULL REFERENCES users(id),
    against     UUID             NOT NULL REFERENCES users(id),
    description TEXT             NOT NULL,
    status      complaint_status NOT NULL DEFAULT 'open',
    reviewed_by UUID             REFERENCES users(id),
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    CONSTRAINT ck_complaint_not_self CHECK (filed_by != against)
);

CREATE INDEX idx_complaints_status    ON complaints(status);
CREATE INDEX idx_complaints_filed_by  ON complaints(filed_by);
CREATE INDEX idx_complaints_against   ON complaints(against);
