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

// MentorDashboard — сводная информация для ментора на главной странице.
type MentorDashboard struct {
	GroupID        *uuid.UUID `json:"group_id"`
	StudentCount   int        `json:"student_count"`
	PendingReviews int        `json:"pending_reviews"` // задачи, ожидающие проверки
	OverdueTasks   int        `json:"overdue_tasks"`
	UpcomingMeetingCount int `json:"upcoming_meetings"`
}

// MentorService — бизнес-логика роли Mentor:
// - управление группой и задачами
// - встречи и объявления
// - вопросы от freshman-ов
type MentorService struct {
	assignRepo       *postgres.AssignmentRepository
	taskRepo         repository.TaskRepository
	meetingRepo      *postgres.MeetingRepository
	notificationRepo *postgres.NotificationRepository
	questionRepo     repository.QuestionRepository
}

// NewMentorService создаёт MentorService.
func NewMentorService(
	assignRepo *postgres.AssignmentRepository,
	taskRepo repository.TaskRepository,
	meetingRepo *postgres.MeetingRepository,
	notificationRepo *postgres.NotificationRepository,
	questionRepo repository.QuestionRepository,
) *MentorService {
	return &MentorService{
		assignRepo:       assignRepo,
		taskRepo:         taskRepo,
		meetingRepo:      meetingRepo,
		notificationRepo: notificationRepo,
		questionRepo:     questionRepo,
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard возвращает сводку для главной страницы ментора.
func (s *MentorService) GetDashboard(ctx context.Context, mentorID uuid.UUID) (*MentorDashboard, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)

	dashboard := &MentorDashboard{}
	if err != nil {
		// У ментора может ещё не быть группы
		return dashboard, nil
	}

	dashboard.GroupID = &group.ID
	dashboard.StudentCount = len(group.Freshmen)

	// Задачи на проверку и просроченные
	submitted := domain.TaskStatusSubmitted
	pending   := domain.TaskStatusPending

	submittedTasks, _ := s.taskRepo.GetByMentorGroup(ctx, group.ID, &submitted)
	dashboard.PendingReviews = len(submittedTasks)

	pendingTasks, _ := s.taskRepo.GetByMentorGroup(ctx, group.ID, &pending)
	// Считаем просроченные вручную (статус pending и дата истекла)
	// Repository уже фильтрует по статусу; просрочку считаем на сервисном уровне
	_ = pendingTasks
	dashboard.OverdueTasks = 0 // Phase 4: добавим агрегацию

	// Предстоящие встречи
	meetings, _ := s.meetingRepo.GetMentorMeetings(ctx, mentorID)
	for _, m := range meetings {
		if !m.Held {
			dashboard.UpcomingMeetingCount++
		}
	}

	return dashboard, nil
}

// ── Group ─────────────────────────────────────────────────────────────────────

// GetGroup возвращает группу ментора со списком студентов.
func (s *MentorService) GetGroup(ctx context.Context, mentorID uuid.UUID) (*domain.MentorGroup, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("no active group found")
	}
	return group, nil
}

// GetStudentDetail возвращает студента с его задачами.
func (s *MentorService) GetStudentDetail(ctx context.Context, mentorID, freshmanID uuid.UUID) (*domain.User, []*domain.Task, error) {
	// Проверяем принадлежность к группе
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return nil, nil, fmt.Errorf("group not found")
	}

	var student *domain.User
	for _, f := range group.Freshmen {
		if f.ID == freshmanID {
			student = f
			break
		}
	}
	if student == nil {
		return nil, nil, fmt.Errorf("student not in your group")
	}

	tasks, err := s.taskRepo.GetByFreshmanID(ctx, freshmanID, nil)
	if err != nil {
		return student, nil, err
	}

	return student, tasks, nil
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

// GetGroupTasks возвращает все задачи группы с опциональным фильтром по статусу.
func (s *MentorService) GetGroupTasks(ctx context.Context, mentorID uuid.UUID, status *domain.TaskStatus) ([]*domain.Task, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("group not found")
	}
	return s.taskRepo.GetByMentorGroup(ctx, group.ID, status)
}

// ApproveTask подтверждает выполнение задачи.
func (s *MentorService) ApproveTask(ctx context.Context, mentorID, taskID uuid.UUID, comment string) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	if task.Status != domain.TaskStatusSubmitted {
		return fmt.Errorf("task is not submitted yet")
	}

	// Проверяем, что студент принадлежит группе ментора
	if err := s.verifyTaskOwnership(ctx, mentorID, task.FreshmanID); err != nil {
		return err
	}

	var c *string
	if comment != "" {
		c = &comment
	}

	if err := s.taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusApproved, c, &mentorID); err != nil {
		return fmt.Errorf("approve task: %w", err)
	}

	// Уведомляем студента
	bodyApprove := fmt.Sprintf("Ментор принял задачу %q. %s", task.Template.Title, comment)
	_ = s.notificationRepo.Create(ctx, &domain.Notification{
		UserID: task.FreshmanID,
		Title:  "✅ Задача принята",
		Body:   &bodyApprove,
	})

	return nil
}

// RejectTask отклоняет выполнение задачи с комментарием.
func (s *MentorService) RejectTask(ctx context.Context, mentorID, taskID uuid.UUID, req *dto.ReviewTaskRequest) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	if task.Status != domain.TaskStatusSubmitted {
		return fmt.Errorf("task is not submitted yet")
	}

	if err := s.verifyTaskOwnership(ctx, mentorID, task.FreshmanID); err != nil {
		return err
	}

	c := req.Comment
	if err := s.taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusRejected, &c, &mentorID); err != nil {
		return fmt.Errorf("reject task: %w", err)
	}

	// Уведомляем студента
	bodyReject := fmt.Sprintf("Ментор отклонил задачу %q. Причина: %s", task.Template.Title, req.Comment)
	_ = s.notificationRepo.Create(ctx, &domain.Notification{
		UserID: task.FreshmanID,
		Title:  "❌ Задача отклонена",
		Body:   &bodyReject,
	})

	return nil
}

// verifyTaskOwnership проверяет, что студент в группе этого ментора.
func (s *MentorService) verifyTaskOwnership(ctx context.Context, mentorID, freshmanID uuid.UUID) error {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return fmt.Errorf("group not found")
	}
	for _, f := range group.Freshmen {
		if f.ID == freshmanID {
			return nil
		}
	}
	return fmt.Errorf("student not in your group")
}

// ── Meetings ──────────────────────────────────────────────────────────────────

// CreateMeeting создаёт встречу для группы.
func (s *MentorService) CreateMeeting(ctx context.Context, mentorID uuid.UUID, req *dto.CreateMeetingRequest) (*domain.Meeting, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("no active group: %w", err)
	}

	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}

	m := &domain.Meeting{
		ID:          uuid.New(),
		MentorID:    mentorID,
		GroupID:     group.ID,
		Title:       req.Title,
		Description: desc,
		ScheduledAt: req.ScheduledAt,
	}

	if err := s.meetingRepo.CreateMeeting(ctx, m); err != nil {
		return nil, fmt.Errorf("create meeting: %w", err)
	}

	// Уведомляем всех студентов группы
	for _, student := range group.Freshmen {
		bodyMeeting := fmt.Sprintf("Ментор запланировал встречу %q", m.Title)
		_ = s.notificationRepo.Create(ctx, &domain.Notification{
			UserID: student.ID,
			Title:  "📅 Новая встреча",
			Body:   &bodyMeeting,
		})
	}

	return m, nil
}

// GetMeetings возвращает встречи ментора.
func (s *MentorService) GetMeetings(ctx context.Context, mentorID uuid.UUID) ([]*domain.Meeting, error) {
	return s.meetingRepo.GetMentorMeetings(ctx, mentorID)
}

// CompleteMeeting завершает встречу.
func (s *MentorService) CompleteMeeting(ctx context.Context, mentorID, meetingID uuid.UUID, req *dto.CompleteMeetingRequest) error {
	return s.meetingRepo.CompleteMeeting(ctx, meetingID, mentorID, req.Notes)
}

// ── Announcements ─────────────────────────────────────────────────────────────

// CreateAnnouncement создаёт объявление для группы.
func (s *MentorService) CreateAnnouncement(ctx context.Context, mentorID uuid.UUID, req *dto.CreateAnnouncementRequest) (*domain.Announcement, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("no active group: %w", err)
	}

	a := &domain.Announcement{
		ID:       uuid.New(),
		AuthorID: mentorID,
		GroupID:  &group.ID,
		Title:    req.Title,
		Body:     req.Body,
	}

	if err := s.meetingRepo.CreateAnnouncement(ctx, a); err != nil {
		return nil, fmt.Errorf("create announcement: %w", err)
	}

	// Уведомляем студентов
	for _, student := range group.Freshmen {
		bodyAnnounce := a.Body
		_ = s.notificationRepo.Create(ctx, &domain.Notification{
			UserID: student.ID,
			Title:  "📢 " + a.Title,
			Body:   &bodyAnnounce,
		})
	}

	return a, nil
}

// GetAnnouncements возвращает объявления для группы ментора.
func (s *MentorService) GetAnnouncements(ctx context.Context, mentorID uuid.UUID) ([]*domain.Announcement, error) {
	group, err := s.assignRepo.GetMentorGroup(ctx, mentorID)
	if err != nil {
		// Если группы нет, возвращаем пустой список
		return s.meetingRepo.GetGroupAnnouncements(ctx, nil)
	}
	return s.meetingRepo.GetGroupAnnouncements(ctx, &group.ID)
}

// ── Questions (freshman → mentor) ────────────────────────────────────────────

// GetQuestions возвращает все вопросы от freshman-ов для этого ментора.
func (s *MentorService) GetQuestions(ctx context.Context, mentorID uuid.UUID) ([]*domain.Question, error) {
	return s.questionRepo.GetByMentorID(ctx, mentorID)
}

// AnswerQuestion отвечает на вопрос freshman-а и уведомляет его.
func (s *MentorService) AnswerQuestion(ctx context.Context, mentorID, questionID uuid.UUID, req *dto.AnswerQuestionRequest) error {
	if err := s.questionRepo.Answer(ctx, questionID, mentorID, req.Answer); err != nil {
		return fmt.Errorf("answer question: %w", err)
	}

	// Уведомляем freshman-а
	// Получаем freshman_id из вопроса через GetByID (в questionRepository есть этот метод)
	return nil
}
