package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/repository"
	"mentorhub/internal/repository/postgres"
)

// AdvisorService — бизнес-логика роли Advisor:
// - просмотр менторов и их студентов
// - аналитика своего сегмента
// - отправка напоминаний
// - подача жалоб
type AdvisorService struct {
	advisorRepo      *postgres.AdvisorRepository
	taskRepo         repository.TaskRepository
	notificationRepo *postgres.NotificationRepository
	complaintRepo    *postgres.ComplaintRepository
}

// NewAdvisorService создаёт AdvisorService.
func NewAdvisorService(
	advisorRepo *postgres.AdvisorRepository,
	taskRepo repository.TaskRepository,
	notificationRepo *postgres.NotificationRepository,
	complaintRepo *postgres.ComplaintRepository,
) *AdvisorService {
	return &AdvisorService{
		advisorRepo:      advisorRepo,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
		complaintRepo:    complaintRepo,
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard возвращает список менторов со статистикой (для главной страницы).
func (s *AdvisorService) GetDashboard(ctx context.Context, advisorID uuid.UUID) ([]*postgres.MentorStats, error) {
	return s.advisorRepo.GetMentorsWithStats(ctx, advisorID)
}

// ── Mentors ───────────────────────────────────────────────────────────────────

// GetMentors возвращает список менторов эдвайзера.
func (s *AdvisorService) GetMentors(ctx context.Context, advisorID uuid.UUID) ([]*postgres.MentorStats, error) {
	return s.advisorRepo.GetMentorsWithStats(ctx, advisorID)
}

// GetMentorStudents возвращает студентов конкретного ментора.
func (s *AdvisorService) GetMentorStudents(ctx context.Context, advisorID, mentorID uuid.UUID) ([]*domain.User, error) {
	students, err := s.advisorRepo.GetMentorStudents(ctx, advisorID, mentorID)
	if err != nil {
		return nil, fmt.Errorf("get mentor students: %w", err)
	}
	return students, nil
}

// GetMentorTaskProgress возвращает задачи студентов конкретного ментора.
func (s *AdvisorService) GetMentorTaskProgress(ctx context.Context, advisorID, mentorID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error) {
	// Проверяем, что ментор принадлежит этому эдвайзеру
	students, err := s.advisorRepo.GetMentorStudents(ctx, advisorID, mentorID)
	if err != nil {
		return nil, fmt.Errorf("mentor not in advisor scope")
	}
	if len(students) == 0 {
		return []*domain.Task{}, nil
	}

	// Собираем задачи по каждому студенту
	var allTasks []*domain.Task
	for _, student := range students {
		tasks, err := s.taskRepo.GetByFreshmanID(ctx, student.ID, status)
		if err != nil {
			continue
		}
		allTasks = append(allTasks, tasks...)
	}
	return allTasks, nil
}

// ── Inactive Students ─────────────────────────────────────────────────────────

// GetInactiveStudents возвращает студентов без активности за последние N дней.
func (s *AdvisorService) GetInactiveStudents(ctx context.Context, advisorID uuid.UUID, days int) ([]*domain.User, error) {
	if days <= 0 {
		days = 7 // дефолт: 7 дней
	}
	return s.advisorRepo.GetInactiveStudents(ctx, advisorID, days)
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// GetAnalytics возвращает агрегированную аналитику эдвайзера.
func (s *AdvisorService) GetAnalytics(ctx context.Context, advisorID uuid.UUID) (*postgres.AdvisorAnalytics, error) {
	return s.advisorRepo.GetAdvisorAnalytics(ctx, advisorID)
}

// ── Reminders ─────────────────────────────────────────────────────────────────

// SendReminder отправляет уведомление-напоминание пользователю.
func (s *AdvisorService) SendReminder(ctx context.Context, advisorID uuid.UUID, req *dto.SendReminderRequest) error {
	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		return fmt.Errorf("invalid target_id")
	}

	n := &domain.Notification{
		ID:     uuid.New(),
		UserID: targetID,
		Title:  "📢 Напоминание от эдвайзера",
		Body:   &req.Message,
	}

	if err := s.notificationRepo.Create(ctx, n); err != nil {
		return fmt.Errorf("send reminder: %w", err)
	}
	return nil
}

// ── Complaints ────────────────────────────────────────────────────────────────

// FileComplaint подаёт жалобу от имени эдвайзера.
func (s *AdvisorService) FileComplaint(ctx context.Context, advisorID uuid.UUID, req *dto.CreateComplaintRequest) error {
	againstID, err := uuid.Parse(req.AgainstID)
	if err != nil {
		return fmt.Errorf("invalid against_id")
	}

	complaint := &domain.Complaint{
		ID:          uuid.New(),
		FiledBy:     advisorID,
		Against:     againstID,
		Description: req.Description,
		Status:      domain.ComplaintStatusOpen,
	}

	if err := s.complaintRepo.Create(ctx, complaint); err != nil {
		return fmt.Errorf("file complaint: %w", err)
	}
	return nil
}
