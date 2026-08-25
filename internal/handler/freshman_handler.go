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

// FreshmanHandler обрабатывает HTTP-запросы для роли Freshman.
type FreshmanHandler struct {
	freshmanService *service.FreshmanService
}

// NewFreshmanHandler создаёт FreshmanHandler.
func NewFreshmanHandler(freshmanService *service.FreshmanService) *FreshmanHandler {
	return &FreshmanHandler{freshmanService: freshmanService}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (h *FreshmanHandler) freshmanID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(middleware.GetUserID(r))
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard godoc
// @Summary  Дэшборд freshman-а
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/dashboard [get]
func (h *FreshmanHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	dashboard, err := h.freshmanService.GetDashboard(r.Context(), freshmanID)
	if err != nil {
		response.InternalError(w, "failed to fetch dashboard")
		return
	}

	response.OK(w, dashboard)
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

// GetTasks godoc
// @Summary  Мои задачи
// @Tags     freshman
// @Security BearerAuth
// @Param    status query string false "Фильтр по статусу (pending|submitted|approved|rejected)"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/tasks [get]
func (h *FreshmanHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	var statusFilter *domain.TaskStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := domain.TaskStatus(s)
		statusFilter = &st
	}

	tasks, err := h.freshmanService.GetTasks(r.Context(), freshmanID, statusFilter)
	if err != nil {
		response.InternalError(w, "failed to fetch tasks")
		return
	}

	response.OK(w, tasks)
}

// GetTask godoc
// @Summary  Детали задачи
// @Tags     freshman
// @Security BearerAuth
// @Param    id path string true "Task UUID"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/tasks/{id} [get]
func (h *FreshmanHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid task ID")
		return
	}

	task, err := h.freshmanService.GetTask(r.Context(), freshmanID, taskID)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}

	response.OK(w, task)
}

// SubmitTask godoc
// @Summary  Сдать задачу
// @Tags     freshman
// @Security BearerAuth
// @Param    id path string true "Task UUID"
// @Accept   json
// @Produce  json
// @Param    body body dto.SubmitTaskRequest true "Ссылка на подтверждение"
// @Success  200 {object} response.Envelope
// @Router   /freshman/tasks/{id}/submit [put]
func (h *FreshmanHandler) SubmitTask(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid task ID")
		return
	}

	var req dto.SubmitTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	if req.ProofURL == "" {
		response.BadRequest(w, "proof_url is required")
		return
	}

	if err := h.freshmanService.SubmitTask(r.Context(), freshmanID, taskID, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "task submitted"})
}

// ── Meetings ──────────────────────────────────────────────────────────────────

// GetMeetings godoc
// @Summary  Встречи группы
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/meetings [get]
func (h *FreshmanHandler) GetMeetings(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	meetings, err := h.freshmanService.GetMeetings(r.Context(), freshmanID)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}

	response.OK(w, meetings)
}

// ── Announcements ─────────────────────────────────────────────────────────────

// GetAnnouncements godoc
// @Summary  Объявления группы и глобальные объявления
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/announcements [get]
func (h *FreshmanHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	announcements, err := h.freshmanService.GetAnnouncements(r.Context(), freshmanID)
	if err != nil {
		response.InternalError(w, "failed to fetch announcements")
		return
	}

	response.OK(w, announcements)
}

// ── FAQ ───────────────────────────────────────────────────────────────────────

// GetFAQ godoc
// @Summary  FAQ для первокурсников
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/faq [get]
func (h *FreshmanHandler) GetFAQ(w http.ResponseWriter, r *http.Request) {
	faq, err := h.freshmanService.GetFAQ(r.Context())
	if err != nil {
		response.InternalError(w, "failed to fetch FAQ")
		return
	}

	response.OK(w, faq)
}

// ── Notifications ─────────────────────────────────────────────────────────────

// GetNotifications godoc
// @Summary  Уведомления freshman-а
// @Tags     freshman
// @Security BearerAuth
// @Param    unread query bool false "Только непрочитанные"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/notifications [get]
func (h *FreshmanHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	unreadOnly := r.URL.Query().Get("unread") == "true"

	notifications, err := h.freshmanService.GetNotifications(r.Context(), freshmanID, unreadOnly)
	if err != nil {
		response.InternalError(w, "failed to fetch notifications")
		return
	}

	response.OK(w, notifications)
}

// MarkNotificationRead godoc
// @Summary  Пометить уведомление как прочитанное
// @Tags     freshman
// @Security BearerAuth
// @Param    id path string true "Notification UUID"
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/notifications/{id}/read [put]
func (h *FreshmanHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	notificationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid notification ID")
		return
	}

	if err := h.freshmanService.MarkNotificationRead(r.Context(), notificationID); err != nil {
		response.InternalError(w, "failed to mark notification as read")
		return
	}

	response.OK(w, map[string]string{"message": "notification marked as read"})
}

// MarkAllNotificationsRead godoc
// @Summary  Пометить все уведомления как прочитанные
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/notifications/read-all [put]
func (h *FreshmanHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	if err := h.freshmanService.MarkAllNotificationsRead(r.Context(), freshmanID); err != nil {
		response.InternalError(w, "failed to mark all notifications as read")
		return
	}

	response.OK(w, map[string]string{"message": "all notifications marked as read"})
}

// ── Questions ─────────────────────────────────────────────────────────────────

// AskMentor godoc
// @Summary  Задать вопрос ментору
// @Tags     freshman
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateQuestionRequest true "Вопрос"
// @Success  201 {object} response.Envelope
// @Router   /freshman/questions [post]
func (h *FreshmanHandler) AskMentor(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	var req dto.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	if req.Body == "" {
		response.BadRequest(w, "body is required")
		return
	}

	question, err := h.freshmanService.AskMentor(r.Context(), freshmanID, &req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, question)
}

// GetMyQuestions godoc
// @Summary  Мои вопросы ментору
// @Tags     freshman
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /freshman/questions [get]
func (h *FreshmanHandler) GetMyQuestions(w http.ResponseWriter, r *http.Request) {
	freshmanID, err := h.freshmanID(r)
	if err != nil {
		response.Unauthorized(w, "invalid identity")
		return
	}

	questions, err := h.freshmanService.GetMyQuestions(r.Context(), freshmanID)
	if err != nil {
		response.InternalError(w, "failed to fetch questions")
		return
	}

	response.OK(w, questions)
}
