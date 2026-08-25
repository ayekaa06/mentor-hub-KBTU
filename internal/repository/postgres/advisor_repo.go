package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"mentorhub/internal/domain"
)

// MentorStats — статистика по одному ментору (для дэшборда эдвайзера).
type MentorStats struct {
	Mentor       domain.User `json:"mentor"`
	StudentCount int         `json:"student_count"`
	AvgProgress  float64     `json:"avg_progress_pct"`
	OverdueTasks int         `json:"overdue_tasks"`
	LastActivity *time.Time  `json:"last_activity"`
}

// AdvisorAnalytics — агрегированная аналитика по домену эдвайзера.
type AdvisorAnalytics struct {
	TotalMentors     int     `json:"total_mentors"`
	TotalStudents    int     `json:"total_students"`
	AvgProgress      float64 `json:"avg_progress_pct"`
	OverdueTasks     int     `json:"overdue_tasks"`
	InactiveStudents int     `json:"inactive_students_7d"`
}

// AdvisorRepository содержит запросы для роли Advisor.
type AdvisorRepository struct {
	db *pgxpool.Pool
}

// NewAdvisorRepository создаёт AdvisorRepository.
func NewAdvisorRepository(db *pgxpool.Pool) *AdvisorRepository {
	return &AdvisorRepository{db: db}
}

// GetMentorsWithStats возвращает всех менторов эдвайзера со статистикой по каждому.
func (r *AdvisorRepository) GetMentorsWithStats(ctx context.Context, advisorID uuid.UUID) ([]*MentorStats, error) {
	const mentorsQ = `
		SELECT u.id, u.email, u.password_hash, u.role, u.first_name, u.last_name,
		       u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN mentor_advisors ma ON ma.mentor_id = u.id
		WHERE ma.advisor_id = $1
		ORDER BY u.last_name
	`
	rows, err := r.db.Query(ctx, mentorsQ, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentors []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		mentors = append(mentors, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*MentorStats, 0, len(mentors))
	for _, m := range mentors {
		s, err := r.getMentorStats(ctx, m)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

// getMentorStats вычисляет статистику для одного ментора.
func (r *AdvisorRepository) getMentorStats(ctx context.Context, mentor domain.User) (*MentorStats, error) {
	const q = `
		SELECT
			COUNT(DISTINCT fg.freshman_id)                                           AS student_count,
			COALESCE(
				ROUND(100.0 * COUNT(t.id) FILTER (WHERE t.status = 'approved')
				      / NULLIF(COUNT(t.id), 0), 1),
				0
			)                                                                        AS avg_progress,
			COUNT(t.id) FILTER (WHERE t.status = 'pending' AND t.due_date < NOW()) AS overdue_tasks,
			MAX(t.submitted_at)                                                      AS last_activity
		FROM mentor_groups mg
		JOIN  academic_years ay   ON ay.id = mg.academic_year_id AND ay.is_active = TRUE
		LEFT JOIN freshman_groups fg ON fg.group_id = mg.id
		LEFT JOIN tasks t            ON t.freshman_id = fg.freshman_id
		WHERE mg.mentor_id = $1
	`
	s := &MentorStats{Mentor: mentor}
	s.Mentor.PasswordHash = "" // не возвращаем хеш
	return s, r.db.QueryRow(ctx, q, mentor.ID).Scan(
		&s.StudentCount, &s.AvgProgress, &s.OverdueTasks, &s.LastActivity,
	)
}

// GetMentorStudents возвращает студентов конкретного ментора (в рамках эдвайзера).
func (r *AdvisorRepository) GetMentorStudents(ctx context.Context, advisorID, mentorID uuid.UUID) ([]*domain.User, error) {
	const q = `
		SELECT DISTINCT u.id, u.email, u.password_hash, u.role, u.first_name, u.last_name,
		       u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN freshman_groups fg ON fg.freshman_id = u.id
		JOIN mentor_groups mg   ON mg.id = fg.group_id
		JOIN mentor_advisors ma ON ma.mentor_id = mg.mentor_id
		JOIN academic_years ay  ON ay.id = mg.academic_year_id AND ay.is_active = TRUE
		WHERE mg.mentor_id = $1
		  AND ma.advisor_id = $2
		ORDER BY u.last_name
	`
	return r.scanUsers(ctx, q, mentorID, advisorID)
}

// GetInactiveStudents возвращает студентов без активности за последние N дней.
func (r *AdvisorRepository) GetInactiveStudents(ctx context.Context, advisorID uuid.UUID, days int) ([]*domain.User, error) {
	const q = `
		SELECT DISTINCT u.id, u.email, u.password_hash, u.role, u.first_name, u.last_name,
		       u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN freshman_groups fg ON fg.freshman_id = u.id
		JOIN mentor_groups mg   ON mg.id = fg.group_id
		JOIN mentor_advisors ma ON ma.mentor_id = mg.mentor_id
		JOIN academic_years ay  ON ay.id = mg.academic_year_id AND ay.is_active = TRUE
		WHERE ma.advisor_id = $1
		  AND u.is_active = TRUE
		  AND NOT EXISTS (
		      SELECT 1 FROM tasks t
		      WHERE t.freshman_id = u.id
		        AND t.submitted_at >= NOW() - ($2::int * INTERVAL '1 day')
		  )
		ORDER BY u.last_name
	`
	return r.scanUsers(ctx, q, advisorID, days)
}

// GetAdvisorAnalytics возвращает агрегированную аналитику эдвайзера.
func (r *AdvisorRepository) GetAdvisorAnalytics(ctx context.Context, advisorID uuid.UUID) (*AdvisorAnalytics, error) {
	const q = `
		SELECT
			COUNT(DISTINCT mg.mentor_id)                                              AS total_mentors,
			COUNT(DISTINCT fg.freshman_id)                                            AS total_students,
			COALESCE(
				ROUND(100.0 * COUNT(t.id) FILTER (WHERE t.status = 'approved')
				      / NULLIF(COUNT(t.id), 0), 1),
				0
			)                                                                          AS avg_progress,
			COUNT(t.id) FILTER (WHERE t.status = 'pending' AND t.due_date < NOW())   AS overdue_tasks
		FROM mentor_advisors ma
		JOIN mentor_groups mg    ON mg.mentor_id   = ma.mentor_id
		JOIN academic_years ay   ON ay.id          = mg.academic_year_id AND ay.is_active = TRUE
		LEFT JOIN freshman_groups fg ON fg.group_id = mg.id
		LEFT JOIN tasks t            ON t.freshman_id = fg.freshman_id
		WHERE ma.advisor_id = $1
	`
	var a AdvisorAnalytics
	if err := r.db.QueryRow(ctx, q, advisorID).Scan(
		&a.TotalMentors, &a.TotalStudents, &a.AvgProgress, &a.OverdueTasks,
	); err != nil {
		return nil, err
	}

	// Неактивные за 7 дней
	inactive, err := r.GetInactiveStudents(ctx, advisorID, 7)
	if err != nil {
		return nil, err
	}
	a.InactiveStudents = len(inactive)
	return &a, nil
}

// scanUsers — вспомогательный метод для чтения списка пользователей.
func (r *AdvisorRepository) scanUsers(ctx context.Context, q string, args ...any) ([]*domain.User, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.PasswordHash = ""
		users = append(users, &u)
	}
	return users, rows.Err()
}
