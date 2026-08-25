package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mentorhub/internal/domain"
	"mentorhub/internal/repository"
)

// ── Task Templates ────────────────────────────────────────────────────────────

type taskTemplateRepository struct {
	db *pgxpool.Pool
}

func NewTaskTemplateRepository(db *pgxpool.Pool) repository.TaskTemplateRepository {
	return &taskTemplateRepository{db: db}
}

func (r *taskTemplateRepository) Create(ctx context.Context, t *domain.TaskTemplate) error {
	const q = `
		INSERT INTO task_templates (id, title, description, due_days, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q,
		t.ID, t.Title, t.Description, t.DueDays, t.IsActive, t.CreatedBy,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *taskTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TaskTemplate, error) {
	const q = `
		SELECT id, title, description, due_days, is_active, created_by, created_at, updated_at
		FROM task_templates WHERE id = $1
	`
	var t domain.TaskTemplate
	err := r.db.QueryRow(ctx, q, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.DueDays,
		&t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task template %s not found", id)
		}
		return nil, fmt.Errorf("get task template: %w", err)
	}
	return &t, nil
}

func (r *taskTemplateRepository) GetAll(ctx context.Context, activeOnly bool) ([]*domain.TaskTemplate, error) {
	q := `
		SELECT id, title, description, due_days, is_active, created_by, created_at, updated_at
		FROM task_templates
	`
	if activeOnly {
		q += ` WHERE is_active = TRUE`
	}
	q += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query task templates: %w", err)
	}
	defer rows.Close()

	var templates []*domain.TaskTemplate
	for rows.Next() {
		var t domain.TaskTemplate
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.DueDays,
			&t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, rows.Err()
}

func (r *taskTemplateRepository) Update(ctx context.Context, t *domain.TaskTemplate) error {
	const q = `
		UPDATE task_templates
		SET title = $2, description = $3, due_days = $4, is_active = $5
		WHERE id = $1
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, q,
		t.ID, t.Title, t.Description, t.DueDays, t.IsActive,
	).Scan(&t.UpdatedAt)
}

func (r *taskTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM task_templates WHERE id = $1`, id)
	return err
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) repository.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	const q = `
		INSERT INTO tasks (id, template_id, freshman_id, assigned_by, status, due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q,
		task.ID, task.TemplateID, task.FreshmanID, task.AssignedBy, task.Status, task.DueDate,
	).Scan(&task.CreatedAt)
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	const q = `
		SELECT t.id, t.template_id, t.freshman_id, t.assigned_by, t.status,
		       t.proof_url, t.comment, t.due_date, t.submitted_at, t.reviewed_at,
		       t.reviewed_by, t.created_at,
		       tt.title, tt.description, tt.due_days
		FROM tasks t
		JOIN task_templates tt ON tt.id = t.template_id
		WHERE t.id = $1
	`
	var t domain.Task
	var tmpl domain.TaskTemplate
	err := r.db.QueryRow(ctx, q, id).Scan(
		&t.ID, &t.TemplateID, &t.FreshmanID, &t.AssignedBy, &t.Status,
		&t.ProofURL, &t.Comment, &t.DueDate, &t.SubmittedAt, &t.ReviewedAt,
		&t.ReviewedBy, &t.CreatedAt,
		&tmpl.Title, &tmpl.Description, &tmpl.DueDays,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task %s not found", id)
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	tmpl.ID = t.TemplateID
	t.Template = &tmpl
	return &t, nil
}

func (r *taskRepository) GetByFreshmanID(ctx context.Context, freshmanID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error) {
	q := `
		SELECT t.id, t.template_id, t.freshman_id, t.assigned_by, t.status,
		       t.proof_url, t.comment, t.due_date, t.submitted_at, t.reviewed_at,
		       t.reviewed_by, t.created_at,
		       tt.title, tt.description, tt.due_days
		FROM tasks t
		JOIN task_templates tt ON tt.id = t.template_id
		WHERE t.freshman_id = $1
	`
	args := []any{freshmanID}
	if status != nil {
		q += ` AND t.status = $2`
		args = append(args, *status)
	}
	q += ` ORDER BY t.due_date ASC`

	return r.queryTasks(ctx, q, args...)
}

func (r *taskRepository) GetByMentorGroup(ctx context.Context, groupID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error) {
	q := `
		SELECT t.id, t.template_id, t.freshman_id, t.assigned_by, t.status,
		       t.proof_url, t.comment, t.due_date, t.submitted_at, t.reviewed_at,
		       t.reviewed_by, t.created_at,
		       tt.title, tt.description, tt.due_days
		FROM tasks t
		JOIN task_templates tt ON tt.id = t.template_id
		JOIN freshman_groups fg ON fg.freshman_id = t.freshman_id
		WHERE fg.group_id = $1
	`
	args := []any{groupID}
	if status != nil {
		q += ` AND t.status = $2`
		args = append(args, *status)
	}
	q += ` ORDER BY t.due_date ASC`

	return r.queryTasks(ctx, q, args...)
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus, comment *string, reviewerID *uuid.UUID) error {
	const q = `
		UPDATE tasks
		SET status = $2, comment = $3, reviewed_by = $4, reviewed_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, q, id, status, comment, reviewerID)
	return err
}

func (r *taskRepository) Submit(ctx context.Context, id uuid.UUID, proofURL string) error {
	const q = `
		UPDATE tasks
		SET status = 'submitted', proof_url = $2, submitted_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`
	tag, err := r.db.Exec(ctx, q, id, proofURL)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task not found or already submitted")
	}
	return nil
}

// queryTasks — вспомогательный метод для чтения списка задач.
func (r *taskRepository) queryTasks(ctx context.Context, q string, args ...any) ([]*domain.Task, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		var tmpl domain.TaskTemplate
		if err := rows.Scan(
			&t.ID, &t.TemplateID, &t.FreshmanID, &t.AssignedBy, &t.Status,
			&t.ProofURL, &t.Comment, &t.DueDate, &t.SubmittedAt, &t.ReviewedAt,
			&t.ReviewedBy, &t.CreatedAt,
			&tmpl.Title, &tmpl.Description, &tmpl.DueDays,
		); err != nil {
			return nil, err
		}
		tmpl.ID = t.TemplateID
		t.Template = &tmpl
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

// ── Assignments: Mentor→Advisor, Freshman→Group ───────────────────────────────

// AssignmentRepository — репозиторий для назначений (mentor↔advisor, freshman↔group, analytics).
type AssignmentRepository struct {
	db *pgxpool.Pool
}

// NewAssignmentRepository создаёт AssignmentRepository.
func NewAssignmentRepository(db *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

// AssignMentorToAdvisor привязывает ментора к эдвайзеру.
func (r *AssignmentRepository) AssignMentorToAdvisor(ctx context.Context, mentorID, advisorID uuid.UUID) error {
	const q = `
		INSERT INTO mentor_advisors (mentor_id, advisor_id)
		VALUES ($1, $2)
		ON CONFLICT (mentor_id, advisor_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, q, mentorID, advisorID)
	return err
}

// UnassignMentorFromAdvisor убирает связь ментора и эдвайзера.
func (r *AssignmentRepository) UnassignMentorFromAdvisor(ctx context.Context, mentorID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM mentor_advisors WHERE mentor_id = $1`, mentorID)
	return err
}

// GetAdvisorMentors возвращает всех менторов эдвайзера с базовой статистикой.
func (r *AssignmentRepository) GetAdvisorMentors(ctx context.Context, advisorID uuid.UUID) ([]*domain.User, error) {
	const q = `
		SELECT u.id, u.email, u.password_hash, u.role, u.first_name, u.last_name,
		       u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN mentor_advisors ma ON ma.mentor_id = u.id
		WHERE ma.advisor_id = $1
		ORDER BY u.last_name
	`
	rows, err := r.db.Query(ctx, q, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentors []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		mentors = append(mentors, &u)
	}
	return mentors, rows.Err()
}

// CreateMentorGroup создаёт группу ментора.
func (r *AssignmentRepository) CreateMentorGroup(ctx context.Context, g *domain.MentorGroup) error {
	const q = `
		INSERT INTO mentor_groups (id, mentor_id, academic_year_id, specialty_id)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q, g.ID, g.MentorID, g.AcademicYearID, g.SpecialtyID).Scan(&g.CreatedAt)
}

// AssignFreshmanToGroup добавляет первокурсника в группу.
func (r *AssignmentRepository) AssignFreshmanToGroup(ctx context.Context, freshmanID, groupID uuid.UUID) error {
	const q = `
		INSERT INTO freshman_groups (freshman_id, group_id)
		VALUES ($1, $2)
		ON CONFLICT (freshman_id, group_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, q, freshmanID, groupID)
	return err
}

// GetMentorGroup возвращает группу с первокурсниками.
func (r *AssignmentRepository) GetMentorGroup(ctx context.Context, mentorID uuid.UUID) (*domain.MentorGroup, error) {
	const q = `
		SELECT mg.id, mg.mentor_id, mg.academic_year_id, mg.specialty_id, mg.created_at
		FROM mentor_groups mg
		JOIN academic_years ay ON ay.id = mg.academic_year_id
		WHERE mg.mentor_id = $1 AND ay.is_active = TRUE
		LIMIT 1
	`
	var g domain.MentorGroup
	err := r.db.QueryRow(ctx, q, mentorID).Scan(
		&g.ID, &g.MentorID, &g.AcademicYearID, &g.SpecialtyID, &g.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("mentor group not found for mentor %s", mentorID)
		}
		return nil, fmt.Errorf("get mentor group: %w", err)
	}

	// Загружаем первокурсников группы
	freshmen, err := r.GetGroupFreshmen(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	g.Freshmen = freshmen
	return &g, nil
}

// GetGroupFreshmen возвращает всех первокурсников группы.
func (r *AssignmentRepository) GetGroupFreshmen(ctx context.Context, groupID uuid.UUID) ([]*domain.User, error) {
	const q = `
		SELECT u.id, u.email, u.password_hash, u.role, u.first_name, u.last_name,
		       u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN freshman_groups fg ON fg.freshman_id = u.id
		WHERE fg.group_id = $1
		ORDER BY u.last_name
	`
	rows, err := r.db.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var freshmen []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		freshmen = append(freshmen, &u)
	}
	return freshmen, rows.Err()
}

// GetFreshmanGroup возвращает группу freshman-а в активном учебном году.
func (r *AssignmentRepository) GetFreshmanGroup(ctx context.Context, freshmanID uuid.UUID) (*domain.MentorGroup, error) {
	const q = `
		SELECT mg.id, mg.mentor_id, mg.academic_year_id, mg.specialty_id, mg.created_at
		FROM mentor_groups mg
		JOIN freshman_groups fg ON fg.group_id = mg.id
		JOIN academic_years ay  ON ay.id = mg.academic_year_id
		WHERE fg.freshman_id = $1 AND ay.is_active = TRUE
		LIMIT 1
	`
	var g domain.MentorGroup
	err := r.db.QueryRow(ctx, q, freshmanID).Scan(
		&g.ID, &g.MentorID, &g.AcademicYearID, &g.SpecialtyID, &g.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no active group found for freshman %s", freshmanID)
		}
		return nil, fmt.Errorf("get freshman group: %w", err)
	}
	return &g, nil
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// SystemStats — агрегированная статистика для Head dashboard.
type SystemStats struct {
	TotalFreshmen   int `json:"total_freshmen"`
	TotalMentors    int `json:"total_mentors"`
	TotalAdvisors   int `json:"total_advisors"`
	OverdueTasks    int `json:"overdue_tasks"`
	OpenComplaints  int `json:"open_complaints"`
	AvgProgress     float64 `json:"avg_progress_pct"`
}

// GetSystemStats возвращает агрегированную статистику системы.
func (r *AssignmentRepository) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	const q = `
		SELECT
			COUNT(*) FILTER (WHERE role = 'freshman' AND is_active)  AS total_freshmen,
			COUNT(*) FILTER (WHERE role = 'mentor'   AND is_active)  AS total_mentors,
			COUNT(*) FILTER (WHERE role = 'advisor'  AND is_active)  AS total_advisors
		FROM users
	`
	var s SystemStats
	if err := r.db.QueryRow(ctx, q).Scan(
		&s.TotalFreshmen, &s.TotalMentors, &s.TotalAdvisors,
	); err != nil {
		return nil, fmt.Errorf("get system stats: %w", err)
	}

	// Просроченные задачи
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM tasks
		WHERE status IN ('pending') AND due_date < NOW()
	`).Scan(&s.OverdueTasks); err != nil {
		return nil, err
	}

	// Открытые жалобы
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM complaints WHERE status = 'open'
	`).Scan(&s.OpenComplaints); err != nil {
		return nil, err
	}

	// Средний прогресс: доля approved задач
	if err := r.db.QueryRow(ctx, `
		SELECT COALESCE(
			ROUND(100.0 * COUNT(*) FILTER (WHERE status = 'approved') / NULLIF(COUNT(*), 0), 1),
			0
		) FROM tasks
	`).Scan(&s.AvgProgress); err != nil {
		return nil, err
	}

	return &s, nil
}

// ── Complaints ────────────────────────────────────────────────────────────────

// ComplaintRepository — репозиторий жалоб.
type ComplaintRepository struct {
	db *pgxpool.Pool
}

func NewComplaintRepository(db *pgxpool.Pool) *ComplaintRepository {
	return &ComplaintRepository{db: db}
}

func (r *ComplaintRepository) Create(ctx context.Context, c *domain.Complaint) error {
	const q = `
		INSERT INTO complaints (id, filed_by, against, description, status)
		VALUES ($1, $2, $3, $4, 'open')
		RETURNING created_at
	`
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q, c.ID, c.FiledBy, c.Against, c.Description).Scan(&c.CreatedAt)
}

func (r *ComplaintRepository) GetAll(ctx context.Context, status *domain.ComplaintStatus) ([]*domain.Complaint, error) {
	q := `
		SELECT c.id, c.filed_by, c.against, c.description,
		       c.status, c.reviewed_by, c.created_at, c.resolved_at,
		       f.first_name, f.last_name,
		       a.first_name, a.last_name
		FROM complaints c
		JOIN users f ON f.id = c.filed_by
		JOIN users a ON a.id = c.against
	`
	args := []any{}
	if status != nil {
		q += ` WHERE c.status = $1`
		args = append(args, *status)
	}
	q += ` ORDER BY c.created_at DESC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query complaints: %w", err)
	}
	defer rows.Close()

	var complaints []*domain.Complaint
	for rows.Next() {
		var c domain.Complaint
		var filerFirst, filerLast, agaFirst, agaLast string
		if err := rows.Scan(
			&c.ID, &c.FiledBy, &c.Against, &c.Description,
			&c.Status, &c.ReviewedBy, &c.CreatedAt, &c.ResolvedAt,
			&filerFirst, &filerLast, &agaFirst, &agaLast,
		); err != nil {
			return nil, err
		}
		c.Filer = &domain.User{ID: c.FiledBy, FirstName: filerFirst, LastName: filerLast}
		c.Subject = &domain.User{ID: c.Against, FirstName: agaFirst, LastName: agaLast}
		complaints = append(complaints, &c)
	}
	return complaints, rows.Err()
}

func (r *ComplaintRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ComplaintStatus, reviewerID uuid.UUID) error {
	var resolvedAt *time.Time
	if status == domain.ComplaintStatusResolved || status == domain.ComplaintStatusDismissed {
		now := time.Now()
		resolvedAt = &now
	}
	const q = `
		UPDATE complaints
		SET status = $2, reviewed_by = $3, resolved_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, q, id, status, reviewerID, resolvedAt)
	return err
}
