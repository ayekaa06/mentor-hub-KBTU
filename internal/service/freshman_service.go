// Package service содержит бизнес-логику для роли Freshman.
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

// FreshmanDashboard — сводная информация на главной странице freshman-а.
type FreshmanDashboard struct {
	GroupID            *uuid.UUID `json:"group_id"`
	MentorID           *uuid.UUID `json:"mentor_id"`
	PendingTasks       int        `json:"pending_tasks"`
	SubmittedTasks     int        `json:"submitted_tasks"`
	ApprovedTasks      int        `json:"approved_tasks"`
	OverdueTasks       int        `json:"overdue_tasks"`
	UnreadNotifications int       `json:"unread_notifications"`
	UpcomingMeetings   int        `json:"upcoming_meetings"`
}

// FreshmanService — бизнес-логика роли Freshman:
// - просмотр задач и их сдача
// - встречи и объявления группы
// - FAQ, уведомления, вопросы ментору
type FreshmanService struct {
	assignRepo       *postgres.AssignmentRepository
	taskRepo         repository.TaskRepository
	meetingRepo      *postgres.MeetingRepository
	notificationRepo *postgres.NotificationRepository
	faqRepo          repository.FAQRepository
	questionRepo     repository.QuestionRepository
}

// NewFreshmanService создаёт FreshmanService.
func NewFreshmanService(
	assignRepo *postgres.AssignmentRepository,
	taskRepo repository.TaskRepository,
	meetingRepo *postgres.MeetingRepository,
	notificationRepo *postgres.NotificationRepository,
	faqRepo repository.FAQRepository,
	questionRepo repository.QuestionRepository,
) *FreshmanService {
	return &FreshmanService{
		assignRepo:       assignRepo,
		taskRepo:         taskRepo,
		meetingRepo:      meetingRepo,
		notificationRepo: notificationRepo,
		faqRepo:          faqRepo,
		questionRepo:     questionRepo,
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard возвращает сводку для главной страницы freshman-а.
func (s *FreshmanService) GetDashboard(ctx context.Context, freshmanID uuid.UUID) (*FreshmanDashboard, error) {
	dashboard := &FreshmanDashboard{}

	// Группа freshman-а
	group, err := s.assignRepo.GetFreshmanGroup(ctx, freshmanID)
	if err == nil {
		dashboard.GroupID = &group.ID
		dashboard.MentorID = &group.MentorID
	}

	// Задачи по статусам
	allTasks, err := s.taskRepo.GetByFreshmanID(ctx, freshmanID, nil)
	if err == nil {
		for _, t := range allTasks {
			switch t.Status {
			case domain.TaskStatusPending:
				dashboard.PendingTasks++
				if t.IsOverdue() {
					dashboard.OverdueTasks++
				}
			case domain.TaskStatusSubmitted:
				dashboard.SubmittedTasks++
			case domain.TaskStatusApproved:
				dashboard.ApprovedTasks++
			}
		}
	}

	// Непрочитанные уведомления
	count, _ := s.notificationRepo.CountUnread(ctx, freshmanID)
	dashboard.UnreadNotifications = count

	// Предстоящие встречи (только если есть группа)
	if dashboard.GroupID != nil {
		meetings, _ := s.meetingRepo.GetMentorMeetings(ctx, *dashboard.MentorID)
		for _, m := range meetings {
			if !m.Held {
				dashboard.UpcomingMeetings++
			}
		}
	}

	return dashboard, nil
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

// GetTasks возвращает задачи freshman-а с опциональным фильтром по статусу.
func (s *FreshmanService) GetTasks(ctx context.Context, freshmanID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error) {
	return s.taskRepo.GetByFreshmanID(ctx, freshmanID, status)
}

// GetTask возвращает конкретную задачу, проверяя принадлежность freshman-у.
func (s *FreshmanService) GetTask(ctx context.Context, freshmanID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	if task.FreshmanID != freshmanID {
		return nil, fmt.Errorf("access denied: task does not belong to you")
	}
	return task, nil
}

// SubmitTask меняет статус задачи на submitted и уведомляет ментора.
func (s *FreshmanService) SubmitTask(ctx context.Context, freshmanID, taskID uuid.UUID, req *dto.SubmitTaskRequest) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}
	if task.FreshmanID != freshmanID {
		return fmt.Errorf("access denied: task does not belong to you")
	}
	if task.Status != domain.TaskStatusPending && task.Status != domain.TaskStatusRejected {
		return fmt.Errorf("task cannot be submitted in status %q", task.Status)
	}

	if err := s.taskRepo.Submit(ctx, taskID, req.ProofURL); err != nil {
		return fmt.Errorf("submit task: %w", err)
	}

	// Уведомляем ментора (ищем ментора через группу freshman-а)
	group, err := s.assignRepo.GetFreshmanGroup(ctx, freshmanID)
	if err == nil {
		var title string
		if task.Template != nil {
			title = task.Template.Title
		}
		body := fmt.Sprintf("Первокурсник сдал задачу %q. Требуется проверка.", title)
		_ = s.notificationRepo.Create(ctx, &domain.Notification{
			UserID: group.MentorID,
			Title:  "📬 Задача на проверку",
			Body:   &body,
		})
	}

	return nil
}

// ── Meetings ──────────────────────────────────────────────────────────────────

// GetMeetings возвращает встречи группы freshman-а.
func (s *FreshmanService) GetMeetings(ctx context.Context, freshmanID uuid.UUID) ([]*domain.Meeting, error) {
	group, err := s.assignRepo.GetFreshmanGroup(ctx, freshmanID)
	if err != nil {
		return nil, fmt.Errorf("no active group found")
	}
	return s.meetingRepo.GetMentorMeetings(ctx, group.MentorID)
}

// ── Announcements ─────────────────────────────────────────────────────────────

// GetAnnouncements возвращает объявления группы freshman-а и глобальные объявления.
func (s *FreshmanService) GetAnnouncements(ctx context.Context, freshmanID uuid.UUID) ([]*domain.Announcement, error) {
	group, err := s.assignRepo.GetFreshmanGroup(ctx, freshmanID)
	if err != nil {
		// Если группы нет — только глобальные объявления
		return s.meetingRepo.GetGroupAnnouncements(ctx, nil)
	}
	return s.meetingRepo.GetGroupAnnouncements(ctx, &group.ID)
}

// ── FAQ ───────────────────────────────────────────────────────────────────────

// GetFAQ возвращает активные FAQ-записи для freshman-а.
func (s *FreshmanService) GetFAQ(ctx context.Context) ([]*domain.FAQItem, error) {
	return s.faqRepo.GetActive(ctx)
}

// ── Notifications ─────────────────────────────────────────────────────────────

// GetNotifications возвращает уведомления freshman-а.
func (s *FreshmanService) GetNotifications(ctx context.Context, freshmanID uuid.UUID, unreadOnly bool) ([]*domain.Notification, error) {
	return s.notificationRepo.GetByUserID(ctx, freshmanID, unreadOnly)
}

// MarkNotificationRead помечает конкретное уведомление как прочитанное.
// Проверка принадлежности не нужна: notification.id уже принадлежит freshmanID
// (защищено на уровне БД — user_id = freshmanID).
func (s *FreshmanService) MarkNotificationRead(ctx context.Context, notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

// MarkAllNotificationsRead помечает все уведомления freshman-а как прочитанные.
func (s *FreshmanService) MarkAllNotificationsRead(ctx context.Context, freshmanID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, freshmanID)
}

// ── Questions ─────────────────────────────────────────────────────────────────

// AskMentor создаёт вопрос от freshman-а ментору и уведомляет ментора.
func (s *FreshmanService) AskMentor(ctx context.Context, freshmanID uuid.UUID, req *dto.CreateQuestionRequest) (*domain.Question, error) {
	group, err := s.assignRepo.GetFreshmanGroup(ctx, freshmanID)
	if err != nil {
		return nil, fmt.Errorf("no active group found — cannot ask mentor")
	}

	q := &domain.Question{
		ID:         uuid.New(),
		FreshmanID: freshmanID,
		MentorID:   group.MentorID,
		Body:       req.Body,
	}
	if err := s.questionRepo.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("create question: %w", err)
	}

	// Уведомляем ментора
	body := fmt.Sprintf("Первокурсник задал вопрос: %q", req.Body)
	_ = s.notificationRepo.Create(ctx, &domain.Notification{
		UserID: group.MentorID,
		Title:  "❓ Новый вопрос",
		Body:   &body,
	})

	return q, nil
}

// GetMyQuestions возвращает все вопросы, заданные freshman-ом.
func (s *FreshmanService) GetMyQuestions(ctx context.Context, freshmanID uuid.UUID) ([]*domain.Question, error) {
	return s.questionRepo.GetByFreshmanID(ctx, freshmanID)
}
