// Package repository определяет интерфейсы для доступа к данным.
// Конкретные реализации находятся в подпакете postgres/.
// Такой подход позволяет легко подменить PostgreSQL на mock в тестах.
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
)

// ── User ─────────────────────────────────────────────────────────────────────

// UserRepository — CRUD для пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context, role *domain.Role, limit, offset int) ([]*domain.User, int, error)
	Update(ctx context.Context, user *domain.User) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ── Refresh Tokens ────────────────────────────────────────────────────────────

// RefreshToken — запись refresh-токена в БД.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// RefreshTokenRepository — хранилище refresh-токенов.
// Phase 2: используется для инвалидации при logout.
type RefreshTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// ── Task Templates ────────────────────────────────────────────────────────────

// TaskTemplateRepository — управление шаблонами задач.
type TaskTemplateRepository interface {
	Create(ctx context.Context, t *domain.TaskTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TaskTemplate, error)
	GetAll(ctx context.Context, activeOnly bool) ([]*domain.TaskTemplate, error)
	Update(ctx context.Context, t *domain.TaskTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

// TaskRepository — назначенные задачи для freshman.
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetByFreshmanID(ctx context.Context, freshmanID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error)
	GetByMentorGroup(ctx context.Context, groupID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus, comment *string, reviewerID *uuid.UUID) error
	Submit(ctx context.Context, id uuid.UUID, proofURL string) error
}

// ── Notifications ─────────────────────────────────────────────────────────────

// NotificationRepository — пользовательские уведомления.
type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
}

// ── Academic Structure ────────────────────────────────────────────────────────

// AcademicYearRepository — учебные годы.
type AcademicYearRepository interface {
	Create(ctx context.Context, ay *domain.AcademicYear) error
	GetAll(ctx context.Context) ([]*domain.AcademicYear, error)
	GetActive(ctx context.Context) (*domain.AcademicYear, error)
	SetActive(ctx context.Context, id uuid.UUID) error
}

// FacultyRepository — факультеты и специальности.
type FacultyRepository interface {
	CreateFaculty(ctx context.Context, f *domain.Faculty) error
	GetAllFaculties(ctx context.Context) ([]*domain.Faculty, error)
	CreateSpecialty(ctx context.Context, s *domain.Specialty) error
	GetSpecialtiesByFaculty(ctx context.Context, facultyID uuid.UUID) ([]*domain.Specialty, error)
}

// ── Assignments ────────────────────────────────────────────────────────────────

// AssignmentRepository — назначения ментор/студент.
// Используется сервисами ментора и freshman-а.
type AssignmentRepository interface {
	// GetFreshmanGroup — группа freshman-а (в активном учебном году).
	GetFreshmanGroup(ctx context.Context, freshmanID uuid.UUID) (*domain.MentorGroup, error)
}

// ── FAQ ─────────────────────────────────────────────────────────────────────

// FAQRepository — FAQ для freshman.
type FAQRepository interface {
	GetActive(ctx context.Context) ([]*domain.FAQItem, error)
	GetAll(ctx context.Context) ([]*domain.FAQItem, error)
}

// ── Questions ────────────────────────────────────────────────────────────────

// QuestionRepository — вопросы freshman ментору.
type QuestionRepository interface {
	Create(ctx context.Context, q *domain.Question) error
	GetByFreshmanID(ctx context.Context, freshmanID uuid.UUID) ([]*domain.Question, error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) ([]*domain.Question, error)
	Answer(ctx context.Context, questionID uuid.UUID, mentorID uuid.UUID, answer string) error
}
