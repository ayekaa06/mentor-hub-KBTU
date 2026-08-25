package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/repository"
	"mentorhub/internal/repository/postgres"
)

// HeadService реализует всю бизнес-логику для роли Head:
// - управление пользователями
// - академическая структура (годы, факультеты, специальности)
// - шаблоны задач и массовое назначение
// - назначение advisor ↔ mentor
// - аналитика и жалобы
type HeadService struct {
	userRepo     repository.UserRepository
	yearRepo     repository.AcademicYearRepository
	facultyRepo  repository.FacultyRepository
	templateRepo repository.TaskTemplateRepository
	taskRepo     repository.TaskRepository
	assignRepo   *postgres.AssignmentRepository
	complaintRepo *postgres.ComplaintRepository
}

// NewHeadService создаёт HeadService.
func NewHeadService(
	userRepo repository.UserRepository,
	yearRepo repository.AcademicYearRepository,
	facultyRepo repository.FacultyRepository,
	templateRepo repository.TaskTemplateRepository,
	taskRepo repository.TaskRepository,
	assignRepo *postgres.AssignmentRepository,
	complaintRepo *postgres.ComplaintRepository,
) *HeadService {
	return &HeadService{
		userRepo:      userRepo,
		yearRepo:      yearRepo,
		facultyRepo:   facultyRepo,
		templateRepo:  templateRepo,
		taskRepo:      taskRepo,
		assignRepo:    assignRepo,
		complaintRepo: complaintRepo,
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard возвращает агрегированную статистику системы.
func (s *HeadService) GetDashboard(ctx context.Context) (*postgres.SystemStats, error) {
	return s.assignRepo.GetSystemStats(ctx)
}

// ── Academic Years ────────────────────────────────────────────────────────────

func (s *HeadService) CreateAcademicYear(ctx context.Context, req *dto.CreateAcademicYearRequest) (*domain.AcademicYear, error) {
	ay := &domain.AcademicYear{
		ID:        uuid.New(),
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  false,
	}
	if err := s.yearRepo.Create(ctx, ay); err != nil {
		return nil, fmt.Errorf("create academic year: %w", err)
	}
	return ay, nil
}

func (s *HeadService) GetAcademicYears(ctx context.Context) ([]*domain.AcademicYear, error) {
	return s.yearRepo.GetAll(ctx)
}

func (s *HeadService) SetActiveYear(ctx context.Context, yearID uuid.UUID) error {
	return s.yearRepo.SetActive(ctx, yearID)
}

// ── Faculties & Specialties ───────────────────────────────────────────────────

func (s *HeadService) CreateFaculty(ctx context.Context, req *dto.CreateFacultyRequest) (*domain.Faculty, error) {
	f := &domain.Faculty{
		ID:   uuid.New(),
		Name: req.Name,
		Code: req.Code,
	}
	if err := s.facultyRepo.CreateFaculty(ctx, f); err != nil {
		return nil, fmt.Errorf("create faculty: %w", err)
	}
	return f, nil
}

func (s *HeadService) GetFaculties(ctx context.Context) ([]*domain.Faculty, error) {
	return s.facultyRepo.GetAllFaculties(ctx)
}

func (s *HeadService) CreateSpecialty(ctx context.Context, req *dto.CreateSpecialtyRequest) (*domain.Specialty, error) {
	facultyID, err := uuid.Parse(req.FacultyID)
	if err != nil {
		return nil, fmt.Errorf("invalid faculty_id: %w", err)
	}
	spec := &domain.Specialty{
		ID:        uuid.New(),
		FacultyID: facultyID,
		Name:      req.Name,
	}
	if err := s.facultyRepo.CreateSpecialty(ctx, spec); err != nil {
		return nil, fmt.Errorf("create specialty: %w", err)
	}
	return spec, nil
}

// ── Assignments ───────────────────────────────────────────────────────────────

// AssignMentorToAdvisor привязывает ментора к эдвайзеру.
func (s *HeadService) AssignMentorToAdvisor(ctx context.Context, mentorID, advisorID uuid.UUID) error {
	// Проверяем роли
	mentor, err := s.userRepo.GetByID(ctx, mentorID)
	if err != nil {
		return ErrUserNotFound
	}
	if mentor.Role != domain.RoleMentor {
		return fmt.Errorf("user %s is not a mentor", mentorID)
	}

	advisor, err := s.userRepo.GetByID(ctx, advisorID)
	if err != nil {
		return ErrUserNotFound
	}
	if advisor.Role != domain.RoleAdvisor {
		return fmt.Errorf("user %s is not an advisor", advisorID)
	}

	return s.assignRepo.AssignMentorToAdvisor(ctx, mentorID, advisorID)
}

// CreateMentorGroup создаёт группу для ментора в текущем учебном году.
func (s *HeadService) CreateMentorGroup(ctx context.Context, mentorID uuid.UUID, specialtyID *uuid.UUID) (*domain.MentorGroup, error) {
	activeYear, err := s.yearRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("no active academic year: %w", err)
	}

	g := &domain.MentorGroup{
		ID:             uuid.New(),
		MentorID:       mentorID,
		AcademicYearID: activeYear.ID,
		SpecialtyID:    specialtyID,
	}
	if err := s.assignRepo.CreateMentorGroup(ctx, g); err != nil {
		return nil, fmt.Errorf("create mentor group: %w", err)
	}
	return g, nil
}

// AssignFreshmanToGroup добавляет первокурсника в группу ментора.
func (s *HeadService) AssignFreshmanToGroup(ctx context.Context, freshmanID, groupID uuid.UUID) error {
	freshman, err := s.userRepo.GetByID(ctx, freshmanID)
	if err != nil {
		return ErrUserNotFound
	}
	if freshman.Role != domain.RoleFreshman {
		return fmt.Errorf("user %s is not a freshman", freshmanID)
	}
	return s.assignRepo.AssignFreshmanToGroup(ctx, freshmanID, groupID)
}

// ── Task Templates ────────────────────────────────────────────────────────────

func (s *HeadService) CreateTaskTemplate(ctx context.Context, req *dto.CreateTaskTemplateRequest, createdBy uuid.UUID) (*domain.TaskTemplate, error) {
	t := &domain.TaskTemplate{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		DueDays:     req.DueDays,
		IsActive:    true,
		CreatedBy:   createdBy,
	}
	if err := s.templateRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create task template: %w", err)
	}
	return t, nil
}

func (s *HeadService) GetTaskTemplates(ctx context.Context, activeOnly bool) ([]*domain.TaskTemplate, error) {
	return s.templateRepo.GetAll(ctx, activeOnly)
}

func (s *HeadService) UpdateTaskTemplate(ctx context.Context, id uuid.UUID, req *dto.UpdateTaskTemplateRequest) (*domain.TaskTemplate, error) {
	t, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("task template not found")
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.DueDays != nil {
		t.DueDays = *req.DueDays
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}

	if err := s.templateRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update task template: %w", err)
	}
	return t, nil
}

func (s *HeadService) DeleteTaskTemplate(ctx context.Context, id uuid.UUID) error {
	if _, err := s.templateRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("task template not found")
	}
	return s.templateRepo.Delete(ctx, id)
}

// AssignTasksToFreshmen массово назначает шаблон задачи первокурсникам.
func (s *HeadService) AssignTasksToFreshmen(ctx context.Context, req *dto.AssignTasksRequest, assignedBy uuid.UUID) (int, error) {
	tmpl, err := s.templateRepo.GetByID(ctx, uuid.MustParse(req.TemplateID))
	if err != nil {
		return 0, fmt.Errorf("task template not found")
	}

	assigned := 0
	for _, idStr := range req.FreshmanIDs {
		freshmanID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		dueDate := time.Now().AddDate(0, 0, tmpl.DueDays)
		if req.DueDate != nil {
			dueDate = *req.DueDate
		}

		task := &domain.Task{
			ID:         uuid.New(),
			TemplateID: tmpl.ID,
			FreshmanID: freshmanID,
			AssignedBy: assignedBy,
			Status:     domain.TaskStatusPending,
			DueDate:    dueDate,
		}

		if err := s.taskRepo.Create(ctx, task); err != nil {
			// Логируем ошибку, но продолжаем для остальных
			continue
		}
		assigned++
	}

	return assigned, nil
}

// ── Analytics ─────────────────────────────────────────────────────────────────

func (s *HeadService) GetAnalytics(ctx context.Context) (*postgres.SystemStats, error) {
	return s.assignRepo.GetSystemStats(ctx)
}

// ── Complaints ────────────────────────────────────────────────────────────────

func (s *HeadService) GetComplaints(ctx context.Context, status *domain.ComplaintStatus) ([]*domain.Complaint, error) {
	return s.complaintRepo.GetAll(ctx, status)
}

func (s *HeadService) UpdateComplaintStatus(ctx context.Context, complaintID uuid.UUID, req *dto.UpdateComplaintStatusRequest, reviewerID uuid.UUID) error {
	return s.complaintRepo.UpdateStatus(ctx, complaintID, domain.ComplaintStatus(req.Status), reviewerID)
}
