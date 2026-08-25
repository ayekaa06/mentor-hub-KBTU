package domain

import (
	"time"

	"github.com/google/uuid"
)

// Meeting — встреча ментора с группой.
type Meeting struct {
	ID          uuid.UUID `json:"id"`
	MentorID    uuid.UUID `json:"mentor_id"`
	GroupID     uuid.UUID `json:"group_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Held        bool      `json:"held"`
	Notes       *string   `json:"notes,omitempty"` // заметки после встречи
	CreatedAt   time.Time `json:"created_at"`

	// Обогащённые поля
	Mentor    *User   `json:"mentor,omitempty"`
	Attendees []*User `json:"attendees,omitempty"`
}

// Announcement — объявление ментора группе или Head всей системе.
type Announcement struct {
	ID        uuid.UUID  `json:"id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	GroupID   *uuid.UUID `json:"group_id,omitempty"` // nil = глобальное от Head
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`

	Author *User `json:"author,omitempty"`
}
