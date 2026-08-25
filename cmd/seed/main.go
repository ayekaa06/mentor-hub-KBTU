// MentorHub — скрипт заполнения базы данных демо-данными.
// Создаёт пользователей всех 4 ролей, учебный год, факультет, группу, задачи, встречи, FAQ.
// Пароль для всех создаваемых пользователей: password123
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"mentorhub/internal/config"
	"mentorhub/internal/domain"
	"mentorhub/internal/pkg/hasher"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	pool, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}
	defer pool.Close()

	ctx := context.Background()
	log.Info().Msg("🌱 Starting database seeding...")

	// 1. Password Hash
	pwdHash, err := hasher.Hash("password123")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to hash default password")
	}

	// 2. Academic Year
	var academicYearID uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO academic_years (name, start_date, end_date, is_active)
		VALUES ('2025-2026', '2025-09-01', '2026-06-30', TRUE)
		ON CONFLICT DO NOTHING
		RETURNING id
	`).Scan(&academicYearID)
	if err != nil {
		_ = pool.QueryRow(ctx, `SELECT id FROM academic_years WHERE name = '2025-2026'`).Scan(&academicYearID)
	}
	log.Info().Str("id", academicYearID.String()).Msg("Academic Year 2025-2026 ready")

	// 3. Faculty & Specialty
	var facultyID, specialtyID uuid.UUID
	_ = pool.QueryRow(ctx, `
		INSERT INTO faculties (name, code) VALUES ('Факультет ИТ', 'FIT')
		ON CONFLICT DO NOTHING RETURNING id
	`).Scan(&facultyID)
	if facultyID == uuid.Nil {
		_ = pool.QueryRow(ctx, `SELECT id FROM faculties WHERE code = 'FIT'`).Scan(&facultyID)
	}

	_ = pool.QueryRow(ctx, `
		INSERT INTO specialties (faculty_id, name) VALUES ($1, 'Software Engineering')
		ON CONFLICT DO NOTHING RETURNING id
	`, facultyID).Scan(&specialtyID)
	if specialtyID == uuid.Nil {
		_ = pool.QueryRow(ctx, `SELECT id FROM specialties WHERE faculty_id = $1 LIMIT 1`, facultyID).Scan(&specialtyID)
	}

	// 4. Create Users for all 4 roles
	users := []struct {
		Email     string
		Role      domain.Role
		FirstName string
		LastName  string
	}{
		{"head@mentorhub.com", domain.RoleHead, "Алексей", "Руководителев"},
		{"advisor@mentorhub.com", domain.RoleAdvisor, "Елена", "Эдвайзерова"},
		{"mentor@mentorhub.com", domain.RoleMentor, "Иван", "Менторов"},
		{"freshman1@mentorhub.com", domain.RoleFreshman, "Аида", "Ержанова"},
		{"freshman2@mentorhub.com", domain.RoleFreshman, "Арман", "Муратов"},
		{"freshman3@mentorhub.com", domain.RoleFreshman, "Дана", "Серикова"},
		{"freshman4@mentorhub.com", domain.RoleFreshman, "Бахтияр", "Ахметов"},
		{"freshman5@mentorhub.com", domain.RoleFreshman, "Мадина", "Оспанова"},
		{"freshman6@mentorhub.com", domain.RoleFreshman, "Султан", "Болатов"},
		{"freshman7@mentorhub.com", domain.RoleFreshman, "Жания", "Касымова"},
		{"freshman8@mentorhub.com", domain.RoleFreshman, "Алихан", "Нурланов"},
		{"freshman9@mentorhub.com", domain.RoleFreshman, "Камила", "Алиева"},
		{"freshman10@mentorhub.com", domain.RoleFreshman, "Темирлан", "Кусаинов"},
	}

	userIDs := make(map[string]uuid.UUID)
	freshmanIDs := make([]uuid.UUID, 0, 10)

	for _, u := range users {
		var uid uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, role, first_name, last_name, is_active)
			VALUES ($1, $2, $3, $4, $5, TRUE)
			ON CONFLICT (email) DO UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
			RETURNING id
		`, u.Email, pwdHash, u.Role, u.FirstName, u.LastName).Scan(&uid)
		if err != nil {
			log.Error().Err(err).Str("email", u.Email).Msg("failed to seed user")
		} else {
			userIDs[u.Email] = uid
			if u.Role == domain.RoleFreshman {
				freshmanIDs = append(freshmanIDs, uid)
			}
			log.Info().Str("role", string(u.Role)).Str("email", u.Email).Msg("User ready")
		}
	}

	headID := userIDs["head@mentorhub.com"]
	advisorID := userIDs["advisor@mentorhub.com"]
	mentorID := userIDs["mentor@mentorhub.com"]

	// 5. Mentor ↔ Advisor Assignment
	if mentorID != uuid.Nil && advisorID != uuid.Nil {
		_, _ = pool.Exec(ctx, `
			INSERT INTO mentor_advisors (mentor_id, advisor_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, mentorID, advisorID)
	}

	// 6. Mentor Group & Freshmen
	var groupID uuid.UUID
	if mentorID != uuid.Nil && academicYearID != uuid.Nil {
		err = pool.QueryRow(ctx, `
			INSERT INTO mentor_groups (mentor_id, academic_year_id, specialty_id)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
			RETURNING id
		`, mentorID, academicYearID, specialtyID).Scan(&groupID)
		if err != nil {
			_ = pool.QueryRow(ctx, `SELECT id FROM mentor_groups WHERE mentor_id = $1`, mentorID).Scan(&groupID)
		}

		if groupID != uuid.Nil {
			for _, fid := range freshmanIDs {
				_, _ = pool.Exec(ctx, `
					INSERT INTO freshman_groups (freshman_id, group_id)
					VALUES ($1, $2) ON CONFLICT DO NOTHING
				`, fid, groupID)
			}
			log.Info().Int("count", len(freshmanIDs)).Str("group_id", groupID.String()).Msg("Freshmen added to mentor group")
		}
	}

	// 7. Task Templates & Assigned Tasks
	templatesData := []struct {
		Title       string
		Description string
		DueDays     int
	}{
		{
			Title:       "Зарегистрироваться на дисциплины следующего семестра",
			Description: "Выбрать обязательные и элективные курсы в портале университете и отправить скриншот формы регистрации.",
			DueDays:     14,
		},
		{
			Title:       "Подписать и сдать договор об обучении",
			Description: "Распечатать, подписать двухсторонний договор об обучении и загрузить скан-копию в формате PDF.",
			DueDays:     10,
		},
		{
			Title:       "Пройти тест по зачету дисциплин (Calculus 1)",
			Description: "Пройти аттестационный тест по математическому анализу в LMS платформе и прикрепить результат.",
			DueDays:     7,
		},
		{
			Title:       "Пройти тест по зачету дисциплин (Physics 1)",
			Description: "Пройти диагностическое тестирование по общей физике и сдать подтверждение балла.",
			DueDays:     7,
		},
		{
			Title:       "Сдать тест по определению уровня английского языка",
			Description: "Пройти онлайн-тестирование Placement Test по английскому языку для распределения по группам.",
			DueDays:     5,
		},
	}

	if headID != uuid.Nil && len(freshmanIDs) >= 10 {
		for tIdx, tData := range templatesData {
			var tid uuid.UUID
			err := pool.QueryRow(ctx, `
				INSERT INTO task_templates (title, description, due_days, is_active, created_by)
				VALUES ($1, $2, $3, TRUE, $4)
				ON CONFLICT DO NOTHING RETURNING id
			`, tData.Title, tData.Description, tData.DueDays, headID).Scan(&tid)
			if err != nil || tid == uuid.Nil {
				_ = pool.QueryRow(ctx, `SELECT id FROM task_templates WHERE title = $1`, tData.Title).Scan(&tid)
			}

			if tid != uuid.Nil {
				// Назначаем всем 10 студентам.
				// Каждое задание выполняют минимум 3-4 студента (статусы approved или submitted)!
				for fIdx, fid := range freshmanIDs {
					status := domain.TaskStatusPending
					var proofURL *string
					var comment *string

					// Логика статусов: для каждого задания студенты (fIdx % 10) имеют разные статусы.
					// По трём студентам (например, (tIdx+fIdx)%10 < 4) ставим approved / submitted!
					patternVal := (tIdx + fIdx) % 10
					if patternVal == 0 || patternVal == 1 {
						status = domain.TaskStatusApproved
						url := fmt.Sprintf("https://portal.university.edu/proofs/task_%d_user_%d.pdf", tIdx+1, fIdx+1)
						proofURL = &url
						c := "Отличная работа, задание принято!"
						comment = &c
					} else if patternVal == 2 || patternVal == 3 {
						status = domain.TaskStatusSubmitted
						url := fmt.Sprintf("https://portal.university.edu/proofs/task_%d_user_%d.png", tIdx+1, fIdx+1)
						proofURL = &url
					} else if patternVal == 4 {
						status = domain.TaskStatusRejected
						url := fmt.Sprintf("https://portal.university.edu/proofs/task_%d_user_%d.png", tIdx+1, fIdx+1)
						proofURL = &url
						c := "Нечеткая скан-копия, пересдайте пожалуйста."
						comment = &c
					} else {
						status = domain.TaskStatusPending
					}

					_, _ = pool.Exec(ctx, `
						INSERT INTO tasks (template_id, freshman_id, assigned_by, status, proof_url, comment, due_date)
						VALUES ($1, $2, $3, $4, $5, $6, NOW() + INTERVAL '7 days')
						ON CONFLICT DO NOTHING
					`, tid, fid, headID, status, proofURL, comment)
				}
				log.Info().Str("title", tData.Title).Msg("Task template and student assignments created")
			}
		}
	}

	// 8. FAQ Items
	if headID != uuid.Nil {
		_, _ = pool.Exec(ctx, `
			INSERT INTO faq_items (question, answer, order_num, is_active, created_by)
			VALUES 
			('Где получить деканатскую справку?', 'В деканате факультета ИТ (каб. 302) по будням с 10:00 до 17:00.', 1, TRUE, $1),
			('Как записаться на пересдачу теста по английскому?', 'Напишите вашему ментору в разделе "Вопросы ментору".', 2, TRUE, $1)
			ON CONFLICT DO NOTHING
		`, headID)
	}

	// 9. Meeting
	if mentorID != uuid.Nil && groupID != uuid.Nil {
		_, _ = pool.Exec(ctx, `
			INSERT INTO meetings (mentor_id, group_id, title, description, scheduled_at)
			VALUES ($1, $2, 'Консультация по выбору дисциплин и тестам', 'Обсуждение регистрации на предметы и сдачи тестов', NOW() + INTERVAL '2 days')
			ON CONFLICT DO NOTHING
		`, mentorID, groupID)
	}

	log.Info().Msg("✅ Database seeding completed successfully!")
	log.Info().Msg("🔑 10 Freshmen Accounts created (all passwords: 'password123'):")
	for i := 1; i <= 10; i++ {
		log.Info().Msgf("   - Freshman %d: freshman%d@mentorhub.com", i, i)
	}
}
