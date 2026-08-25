package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/middleware"
	"mentorhub/internal/pkg/response"
	"mentorhub/internal/service"
)

// MentorHandler обрабатывает HTTP-запросы для роли Mentor.
type MentorHandler struct {
	mentorService *service.MentorService
}

// NewMentorHandler создаёт MentorHandler.
func NewMentorHandler(mentorService *service.MentorService) *MentorHandler {
	return &MentorHandler{mentorService: mentorService}
}

// GetDashboard godoc
// @Summary  Дэшборд ментора
// @Tags     mentor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/dashboard [get]
func (h *MentorHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	dashboard, err := h.mentorService.GetDashboard(r.Context(), mentorID)
	if err != nil {
		response.InternalError(w, "failed to fetch mentor dashboard")
		return
	}

	response.OK(w, dashboard)
}

// GetGroup godoc
// @Summary  Группа ментора с первокурсниками
// @Tags     mentor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/group [get]
func (h *MentorHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	group, err := h.mentorService.GetGroup(r.Context(), mentorID)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}

	response.OK(w, group)
}

// GetStudentDetail godoc
// @Summary  Детали первокурсника и его задачи
// @Tags     mentor
// @Security BearerAuth
// @Param    id path string true "Freshman UUID"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/group/{id} [get]
func (h *MentorHandler) GetStudentDetail(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	freshmanID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid student ID")
		return
	}

	student, tasks, err := h.mentorService.GetStudentDetail(r.Context(), mentorID, freshmanID)
	if err != nil {
		response.Forbidden(w, err.Error())
		return
	}

	student.PasswordHash = ""
	response.OK(w, map[string]any{
		"student": student,
		"tasks":   tasks,
	})
}

// GetTasks godoc
// @Summary  Задачи группы ментора
// @Tags     mentor
// @Security BearerAuth
// @Param    status query string false "Фильтр по статусу (pending|submitted|approved|rejected)"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/tasks [get]
func (h *MentorHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	var statusFilter *domain.TaskStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := domain.TaskStatus(s)
		statusFilter = &st
	}

	tasks, err := h.mentorService.GetGroupTasks(r.Context(), mentorID, statusFilter)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, tasks)
}

// ApproveTask godoc
// @Summary  Принять выполнение задачи
// @Tags     mentor
// @Security BearerAuth
// @Param    id path string true "Task UUID"
// @Accept   json
// @Produce  json
// @Param    body body dto.ReviewTaskRequest false "Комментарий (опционально)"
// @Success  200 {object} response.Envelope
// @Router   /mentor/tasks/{id}/approve [put]
func (h *MentorHandler) ApproveTask(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid task ID")
		return
	}

	var req dto.ReviewTaskRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if err := h.mentorService.ApproveTask(r.Context(), mentorID, taskID, req.Comment); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "task approved"})
}

// RejectTask godoc
// @Summary  Отклонить выполнение задачи
// @Tags     mentor
// @Security BearerAuth
// @Param    id path string true "Task UUID"
// @Accept   json
// @Produce  json
// @Param    body body dto.ReviewTaskRequest true "Комментарий с причиной"
// @Success  200 {object} response.Envelope
// @Router   /mentor/tasks/{id}/reject [put]
func (h *MentorHandler) RejectTask(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid task ID")
		return
	}

	var req dto.ReviewTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	if err := h.mentorService.RejectTask(r.Context(), mentorID, taskID, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "task rejected"})
}

// CreateMeeting godoc
// @Summary  Запланировать встречу
// @Tags     mentor
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateMeetingRequest true "Параметры встречи"
// @Success  201 {object} response.Envelope
// @Router   /mentor/meetings [post]
func (h *MentorHandler) CreateMeeting(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	var req dto.CreateMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	meeting, err := h.mentorService.CreateMeeting(r.Context(), mentorID, &req)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, meeting)
}

// GetMeetings godoc
// @Summary  Список встреч ментора
// @Tags     mentor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/meetings [get]
func (h *MentorHandler) GetMeetings(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	meetings, err := h.mentorService.GetMeetings(r.Context(), mentorID)
	if err != nil {
		response.InternalError(w, "failed to fetch meetings")
		return
	}

	response.OK(w, meetings)
}

// CompleteMeeting godoc
// @Summary  Завершить встречу и добавить заметки
// @Tags     mentor
// @Security BearerAuth
// @Param    id path string true "Meeting UUID"
// @Accept   json
// @Produce  json
// @Param    body body dto.CompleteMeetingRequest true "Заметки"
// @Success  200 {object} response.Envelope
// @Router   /mentor/meetings/{id}/complete [put]
func (h *MentorHandler) CompleteMeeting(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	meetingID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid meeting ID")
		return
	}

	var req dto.CompleteMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	if err := h.mentorService.CompleteMeeting(r.Context(), mentorID, meetingID, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "meeting completed"})
}

// CreateAnnouncement godoc
// @Summary  Опубликовать объявление для группы
// @Tags     mentor
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateAnnouncementRequest true "Объявление"
// @Success  201 {object} response.Envelope
// @Router   /mentor/announcements [post]
func (h *MentorHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	var req dto.CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	announcement, err := h.mentorService.CreateAnnouncement(r.Context(), mentorID, &req)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, announcement)
}

// GetAnnouncements godoc
// @Summary  Список объявлений группы
// @Tags     mentor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/announcements [get]
func (h *MentorHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := middleware.GetUserID(r)
	mentorID, err := uuid.Parse(mentorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	announcements, err := h.mentorService.GetAnnouncements(r.Context(), mentorID)
	if err != nil {
		response.InternalError(w, "failed to fetch announcements")
		return
	}

	response.OK(w, announcements)
}

// ── Questions ─────────────────────────────────────────────────────────────────

// GetQuestions godoc
// @Summary  Вопросы от первокурсников
// @Tags     mentor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /mentor/questions [get]
func (h *MentorHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	mentorID, err := uuid.Parse(middleware.GetUserID(r))
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	questions, err := h.mentorService.GetQuestions(r.Context(), mentorID)
	if err != nil {
		response.InternalError(w, "failed to fetch questions")
		return
	}

	response.OK(w, questions)
}

// AnswerQuestion godoc
// @Summary  Ответить на вопрос первокурсника
// @Tags     mentor
// @Security BearerAuth
// @Param    id path string true "Question UUID"
// @Accept   json
// @Produce  json
// @Param    body body dto.AnswerQuestionRequest true "Ответ"
// @Success  200 {object} response.Envelope
// @Router   /mentor/questions/{id}/answer [put]
func (h *MentorHandler) AnswerQuestion(w http.ResponseWriter, r *http.Request) {
	mentorID, err := uuid.Parse(middleware.GetUserID(r))
	if err != nil {
		response.Unauthorized(w, "invalid mentor identity")
		return
	}

	questionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid question ID")
		return
	}

	var req dto.AnswerQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	if req.Answer == "" {
		response.BadRequest(w, "answer is required")
		return
	}

	if err := h.mentorService.AnswerQuestion(r.Context(), mentorID, questionID, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "question answered"})
}

