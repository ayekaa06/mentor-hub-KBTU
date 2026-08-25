// Package domain содержит бизнес-сущности приложения.
// Структуры не зависят ни от HTTP, ни от базы данных.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role определяет роль пользователя в системе.
type Role string

const (
	RoleHead     Role = "head"
	RoleAdvisor  Role = "advisor"
	RoleMentor   Role = "mentor"
	RoleFreshman Role = "freshman"
)

// User — базовая сущность пользователя.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // никогда не отдаём в JSON
	Role         Role      `json:"role"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FullName возвращает полное имя пользователя.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// AcademicYear — учебный год.
type AcademicYear struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`       // "2025-2026"
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Faculty — факультет.
type Faculty struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Code        string       `json:"code"` // "CS", "ENG"
	CreatedAt   time.Time    `json:"created_at"`
	Specialties []*Specialty `json:"specialties,omitempty"`
}

// Specialty — специальность внутри факультета.
type Specialty struct {
	ID        uuid.UUID `json:"id"`
	FacultyID uuid.UUID `json:"faculty_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// MentorGroup — группа первокурсников, закреплённая за ментором.
type MentorGroup struct {
	ID             uuid.UUID  `json:"id"`
	MentorID       uuid.UUID  `json:"mentor_id"`
	AcademicYearID uuid.UUID  `json:"academic_year_id"`
	SpecialtyID    *uuid.UUID `json:"specialty_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`

	// Обогащённые поля (заполняются сервисом по необходимости)
	Mentor   *User   `json:"mentor,omitempty"`
	Freshmen []*User `json:"freshmen,omitempty"`
}
