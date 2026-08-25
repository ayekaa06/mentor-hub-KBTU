package dto

import "time"

// ── Task Templates ───────────────────────────────────────────────────────────

// CreateTaskTemplateRequest — создание шаблона задачи (только Head).
type CreateTaskTemplateRequest struct {
	Title       string `json:"title"       validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=2000"`
	DueDays     int    `json:"due_days"    validate:"required,min=1,max=365"`
}

// UpdateTaskTemplateRequest — обновление шаблона.
type UpdateTaskTemplateRequest struct {
	Title       *string `json:"title,omitempty"       validate:"omitempty,min=3,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	DueDays     *int    `json:"due_days,omitempty"    validate:"omitempty,min=1,max=365"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// AssignTasksRequest — массовое назначение шаблона freshman-ам.
type AssignTasksRequest struct {
	TemplateID  string     `json:"template_id"   validate:"required,uuid"`
	FreshmanIDs []string   `json:"freshman_ids"  validate:"required,min=1,dive,uuid"`
	DueDate     *time.Time `json:"due_date,omitempty"` // если nil — now + due_days
}

// ── Task Actions ─────────────────────────────────────────────────────────────

// SubmitTaskRequest — freshman загружает подтверждение выполнения.
type SubmitTaskRequest struct {
	ProofURL string `json:"proof_url" validate:"required,url"`
}

// ReviewTaskRequest — ментор подтверждает или отклоняет задачу.
type ReviewTaskRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment" validate:"max=1000"`
}

// ── Meetings ─────────────────────────────────────────────────────────────────

// CreateMeetingRequest — создание встречи ментором.
type CreateMeetingRequest struct {
	Title       string    `json:"title"       validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	ScheduledAt time.Time `json:"scheduled_at" validate:"required"`
}

// CompleteMeetingRequest — завершение встречи ментором.
type CompleteMeetingRequest struct {
	Notes string `json:"notes" validate:"max=2000"`
}

// ── Announcements ─────────────────────────────────────────────────────────────

// CreateAnnouncementRequest — создание объявления.
type CreateAnnouncementRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Body    string `json:"body"  validate:"required,min=1"`
	GroupID string `json:"group_id,omitempty" validate:"omitempty,uuid"`
}

// ── Complaints ────────────────────────────────────────────────────────────────

// CreateComplaintRequest — подача жалобы.
type CreateComplaintRequest struct {
	AgainstID   string `json:"against_id"  validate:"required,uuid"`
	Description string `json:"description" validate:"required,min=10,max=2000"`
}

// UpdateComplaintStatusRequest — обновление статуса жалобы Head-ом.
type UpdateComplaintStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=in_review resolved dismissed"`
}

