package domain

import (
	"time"

	"github.com/google/uuid"
)

// ComplaintStatus — статус жалобы.
type ComplaintStatus string

const (
	ComplaintStatusOpen      ComplaintStatus = "open"
	ComplaintStatusInReview  ComplaintStatus = "in_review"
	ComplaintStatusResolved  ComplaintStatus = "resolved"
	ComplaintStatusDismissed ComplaintStatus = "dismissed"
)

// Complaint — жалоба от advisor на ментора или от freshman на кого-либо.
// Head рассматривает все жалобы.
type Complaint struct {
	ID          uuid.UUID       `json:"id"`
	FiledBy     uuid.UUID       `json:"filed_by"`
	Against     uuid.UUID       `json:"against"`
	Description string          `json:"description"`
	Status      ComplaintStatus `json:"status"`
	ReviewedBy  *uuid.UUID      `json:"reviewed_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	ResolvedAt  *time.Time      `json:"resolved_at,omitempty"`

	// Обогащённые поля
	Filer   *User `json:"filer,omitempty"`
	Subject *User `json:"subject,omitempty"`
}

// Notification — уведомление пользователю.
type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// FAQItem — вопрос-ответ для freshman.
type FAQItem struct {
	ID        uuid.UUID `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	OrderNum  int       `json:"order_num"`
	IsActive  bool      `json:"is_active"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Question — вопрос freshman-а своему ментору.
type Question struct {
	ID         uuid.UUID  `json:"id"`
	FreshmanID uuid.UUID  `json:"freshman_id"`
	MentorID   uuid.UUID  `json:"mentor_id"`
	Body       string     `json:"body"`
	Answer     *string    `json:"answer,omitempty"`
	AnsweredAt *time.Time `json:"answered_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`

	Freshman *User `json:"freshman,omitempty"`
}
