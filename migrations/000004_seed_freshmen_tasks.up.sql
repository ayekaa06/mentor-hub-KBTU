-- ============================================================
-- MentorHub — Migration 000004: Seed Data
-- 10 первокурсников, 1 ментор, 1 советник, 1 учебный год,
-- 4 шаблона задач, назначения задач (≥3 студента на задачу)
-- Пароль для всех тестовых пользователей: Test1234!
-- Ментор/советник: уникальные email, не конфликтуют с seed из 000001
-- ============================================================

-- ── 1. Учебный год ────────────────────────────────────────────────────────────

INSERT INTO academic_years (id, name, start_date, end_date, is_active)
VALUES (
    '10000000-0000-0000-0000-000000000001',
    '2025-2026',
    '2025-09-01',
    '2026-06-30',
    TRUE
) ON CONFLICT (id) DO NOTHING;

-- ── 2. Факультет и специальность ──────────────────────────────────────────────

INSERT INTO faculties (id, name, code)
VALUES (
    '20000000-0000-0000-0000-000000000001',
    'Школа информационных технологий',
    'IT'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO specialties (id, faculty_id, name)
VALUES (
    '30000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    'Информационные системы'
) ON CONFLICT (id) DO NOTHING;

-- ── 3. Ментор и советник (новые, уникальные email) ───────────────────────────
-- Пароль: Test1234!  (bcrypt cost 12)

INSERT INTO users (id, email, password_hash, role, first_name, last_name)
VALUES
    (
        'a0000000-0000-0000-0000-000000000001',
        'arman.seitkali@mentorhub.kz',
        '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2',
        'mentor',
        'Арман',
        'Сейткали'
    ),
    (
        'a0000000-0000-0000-0000-000000000002',
        'gulnar.bekova@mentorhub.kz',
        '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2',
        'advisor',
        'Гүлнар',
        'Бекова'
    )
ON CONFLICT (email) DO NOTHING;

-- ── 4. Группа ментора ─────────────────────────────────────────────────────────

INSERT INTO mentor_groups (id, mentor_id, academic_year_id, specialty_id)
VALUES (
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    '10000000-0000-0000-0000-000000000001',
    '30000000-0000-0000-0000-000000000001'
) ON CONFLICT (id) DO NOTHING;

-- Привязка советника к ментору
INSERT INTO mentor_advisors (mentor_id, advisor_id)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000002'
) ON CONFLICT DO NOTHING;

-- ── 5. 10 Первокурсников ──────────────────────────────────────────────────────

INSERT INTO users (id, email, password_hash, role, first_name, last_name)
VALUES
    ('c0000000-0000-0000-0000-000000000001', 'aibek.dzhaksybekov@mentorhub.kz',  '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Айбек',    'Джаксыбеков'),
    ('c0000000-0000-0000-0000-000000000002', 'dana.nurlanovna@mentorhub.kz',     '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Дана',     'Нурланова'),
    ('c0000000-0000-0000-0000-000000000003', 'yerlan.abenov@mentorhub.kz',       '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Ерлан',    'Абенов'),
    ('c0000000-0000-0000-0000-000000000004', 'ainur.seitkali@mentorhub.kz',      '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Айнур',    'Сейткали'),
    ('c0000000-0000-0000-0000-000000000005', 'timur.bekzhan@mentorhub.kz',       '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Тимур',    'Бекжан'),
    ('c0000000-0000-0000-0000-000000000006', 'zarina.ospanova@mentorhub.kz',     '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Зарина',   'Оспанова'),
    ('c0000000-0000-0000-0000-000000000007', 'aslan.tulegenov@mentorhub.kz',     '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Аслан',    'Тулегенов'),
    ('c0000000-0000-0000-0000-000000000008', 'kamila.bekova@mentorhub.kz',       '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Камила',   'Бекова'),
    ('c0000000-0000-0000-0000-000000000009', 'nurzhan.aidarov@mentorhub.kz',     '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Нуржан',   'Айдаров'),
    ('c0000000-0000-0000-0000-000000000010', 'aliya.suleimenova@mentorhub.kz',   '$2a$12$7lsagILy./BaECVla5Jho.NhPQ5xEoqd21cf9JTobaCwy6AwFYD.2', 'freshman', 'Алия',     'Сулейменова')
ON CONFLICT (email) DO NOTHING;

-- Добавить всех первокурсников в группу ментора
INSERT INTO freshman_groups (freshman_id, group_id)
VALUES
    ('c0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000004', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000005', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000006', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000007', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000008', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000009', 'b0000000-0000-0000-0000-000000000001'),
    ('c0000000-0000-0000-0000-000000000010', 'b0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- ── 6. Шаблоны задач ──────────────────────────────────────────────────────────

INSERT INTO task_templates (id, title, description, due_days, created_by)
SELECT
    t.id,
    t.title,
    t.description,
    t.due_days,
    (SELECT id FROM users WHERE role = 'head' LIMIT 1)
FROM (VALUES
    (
        'd0000000-0000-0000-0000-000000000001'::uuid,
        'Подписать и сдать договор об обучении',
        'Необходимо распечатать, подписать и сдать договор об обучении в деканат.',
        14
    ),
    (
        'd0000000-0000-0000-0000-000000000002'::uuid,
        'Перезачёт дисциплины: Calculus 1',
        'Пройдите процедуру перезачёта по дисциплине "Математический анализ 1".',
        21
    ),
    (
        'd0000000-0000-0000-0000-000000000003'::uuid,
        'Перезачёт дисциплины: Physics 1',
        'Пройдите процедуру перезачёта по дисциплине "Физика 1".',
        21
    ),
    (
        'd0000000-0000-0000-0000-000000000004'::uuid,
        'Тест на определение уровня английского языка',
        'Пройдите тест на определение уровня английского языка.',
        10
    )
) AS t(id, title, description, due_days)
ON CONFLICT (id) DO NOTHING;

-- ── 7. Назначение задач первокурсникам ───────────────────────────────────────

-- Задача 1: Подписать договор — все 10 первокурсников
INSERT INTO tasks (template_id, freshman_id, assigned_by, status, due_date)
SELECT
    'd0000000-0000-0000-0000-000000000001',
    u.fid,
    'a0000000-0000-0000-0000-000000000001',
    u.st,
    NOW() + INTERVAL '14 days'
FROM (VALUES
    ('c0000000-0000-0000-0000-000000000001'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000002'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000003'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000004'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000005'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000006'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000007'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000008'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000009'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000010'::uuid, 'pending'::task_status)
) AS u(fid, st);

-- Задача 2: Перезачёт Calculus 1 — 7 студентов (1–7)
INSERT INTO tasks (template_id, freshman_id, assigned_by, status, due_date)
SELECT
    'd0000000-0000-0000-0000-000000000002',
    u.fid,
    'a0000000-0000-0000-0000-000000000001',
    u.st,
    NOW() + INTERVAL '21 days'
FROM (VALUES
    ('c0000000-0000-0000-0000-000000000001'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000002'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000003'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000004'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000005'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000006'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000007'::uuid, 'submitted'::task_status)
) AS u(fid, st);

-- Задача 3: Перезачёт Physics 1 — 7 студентов (4–10)
INSERT INTO tasks (template_id, freshman_id, assigned_by, status, due_date)
SELECT
    'd0000000-0000-0000-0000-000000000003',
    u.fid,
    'a0000000-0000-0000-0000-000000000001',
    u.st,
    NOW() + INTERVAL '21 days'
FROM (VALUES
    ('c0000000-0000-0000-0000-000000000004'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000005'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000006'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000007'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000008'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000009'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000010'::uuid, 'pending'::task_status)
) AS u(fid, st);

-- Задача 4: Тест по английскому — 5 студентов (нечётные)
INSERT INTO tasks (template_id, freshman_id, assigned_by, status, due_date)
SELECT
    'd0000000-0000-0000-0000-000000000004',
    u.fid,
    'a0000000-0000-0000-0000-000000000001',
    u.st,
    NOW() + INTERVAL '10 days'
FROM (VALUES
    ('c0000000-0000-0000-0000-000000000001'::uuid, 'submitted'::task_status),
    ('c0000000-0000-0000-0000-000000000003'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000005'::uuid, 'pending'::task_status),
    ('c0000000-0000-0000-0000-000000000007'::uuid, 'approved'::task_status),
    ('c0000000-0000-0000-0000-000000000009'::uuid, 'pending'::task_status)
) AS u(fid, st);
