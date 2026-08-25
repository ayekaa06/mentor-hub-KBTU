package dto

import "time"

// ── User Management ──────────────────────────────────────────────────────────

// CreateUserRequest — создание пользователя Head-ом.
type CreateUserRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8,max=72"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Role      string `json:"role"       validate:"required,oneof=advisor mentor freshman"`
}

// UpdateProfileRequest — обновление профиля пользователем.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty"  validate:"omitempty,min=1,max=100"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// AssignAdvisorRequest — назначение advisor-а Head-ом.
type AssignAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required,uuid"`
}

// AssignMentorRequest — назначение ментора к advisor-у.
type AssignMentorRequest struct {
	MentorID  string `json:"mentor_id"  validate:"required,uuid"`
	AdvisorID string `json:"advisor_id" validate:"required,uuid"`
}

// ── Academic Structure ───────────────────────────────────────────────────────

// CreateAcademicYearRequest — создание учебного года.
type CreateAcademicYearRequest struct {
	Name      string    `json:"name"       validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date"   validate:"required"`
}

// CreateFacultyRequest — создание факультета.
type CreateFacultyRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
	Code string `json:"code" validate:"required,min=1,max=50"`
}

// CreateSpecialtyRequest — создание специальности.
type CreateSpecialtyRequest struct {
	FacultyID string `json:"faculty_id" validate:"required,uuid"`
	Name      string `json:"name"       validate:"required,min=2,max=255"`
}

// ── Pagination ───────────────────────────────────────────────────────────────

// PaginationQuery — параметры пагинации из query string.
type PaginationQuery struct {
	Page    int `schema:"page"     validate:"min=1"`
	PerPage int `schema:"per_page" validate:"min=1,max=100"`
}

func (p *PaginationQuery) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 20
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PerPage
}
