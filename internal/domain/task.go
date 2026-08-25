package domain

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus — статус выполнения задачи.
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusSubmitted TaskStatus = "submitted"
	TaskStatusApproved  TaskStatus = "approved"
	TaskStatusRejected  TaskStatus = "rejected"
)

// TaskTemplate — шаблон задачи, создаётся Head.
// Из одного шаблона можно назначить задачи сразу всем freshman-ам.
type TaskTemplate struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDays     int       `json:"due_days"` // дней на выполнение после назначения
	IsActive    bool      `json:"is_active"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Task — конкретная задача, назначенная freshman.
type Task struct {
	ID          uuid.UUID  `json:"id"`
	TemplateID  uuid.UUID  `json:"template_id"`
	FreshmanID  uuid.UUID  `json:"freshman_id"`
	AssignedBy  uuid.UUID  `json:"assigned_by"`
	Status      TaskStatus `json:"status"`
	ProofURL    *string    `json:"proof_url,omitempty"`   // загруженное подтверждение
	Comment     *string    `json:"comment,omitempty"`     // комментарий ментора при reject
	DueDate     time.Time  `json:"due_date"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  *uuid.UUID `json:"reviewed_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`

	// Обогащённые поля
	Template *TaskTemplate `json:"template,omitempty"`
	Freshman *User         `json:"freshman,omitempty"`
}

// IsOverdue возвращает true, если дедлайн прошёл, а задача ещё не выполнена.
func (t *Task) IsOverdue() bool {
	return time.Now().After(t.DueDate) &&
		t.Status != TaskStatusApproved &&
		t.Status != TaskStatusSubmitted
}
