package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/middleware"
	"mentorhub/internal/pkg/response"
	"mentorhub/internal/service"
)

// AdvisorHandler обрабатывает все HTTP-запросы для роли Advisor.
type AdvisorHandler struct {
	advisorService *service.AdvisorService
}

// NewAdvisorHandler создаёт AdvisorHandler.
func NewAdvisorHandler(advisorService *service.AdvisorService) *AdvisorHandler {
	return &AdvisorHandler{advisorService: advisorService}
}

// GetDashboard godoc
// @Summary  Дэшборд эдвайзера
// @Tags     advisor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /advisor/dashboard [get]
func (h *AdvisorHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	stats, err := h.advisorService.GetDashboard(r.Context(), advisorID)
	if err != nil {
		response.InternalError(w, "failed to get dashboard stats")
		return
	}

	response.OK(w, stats)
}

// GetMentors godoc
// @Summary  Список менторов эдвайзера со статистикой
// @Tags     advisor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /advisor/mentors [get]
func (h *AdvisorHandler) GetMentors(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	mentors, err := h.advisorService.GetMentors(r.Context(), advisorID)
	if err != nil {
		response.InternalError(w, "failed to fetch mentors")
		return
	}

	response.OK(w, mentors)
}

// GetMentorStudents godoc
// @Summary  Список студентов конкретного ментора
// @Tags     advisor
// @Security BearerAuth
// @Param    id path string true "Mentor UUID"
// @Success  200 {object} response.Envelope
// @Router   /advisor/mentors/{id}/students [get]
func (h *AdvisorHandler) GetMentorStudents(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	mentorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid mentor ID")
		return
	}

	students, err := h.advisorService.GetMentorStudents(r.Context(), advisorID, mentorID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, students)
}

// GetMentorTaskProgress godoc
// @Summary  Прогресс по задачам студентов конкретного ментора
// @Tags     advisor
// @Security BearerAuth
// @Param    id path string true "Mentor UUID"
// @Param    status query string false "Фильтр по статусу (pending|submitted|approved|rejected)"
// @Success  200 {object} response.Envelope
// @Router   /advisor/mentors/{id}/tasks [get]
func (h *AdvisorHandler) GetMentorTaskProgress(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	mentorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid mentor ID")
		return
	}

	var statusFilter *domain.TaskStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := domain.TaskStatus(s)
		statusFilter = &st
	}

	tasks, err := h.advisorService.GetMentorTaskProgress(r.Context(), advisorID, mentorID, statusFilter)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, tasks)
}

// GetInactiveStudents godoc
// @Summary  Список неактивных студентов (без сабмитов за N дней)
// @Tags     advisor
// @Security BearerAuth
// @Param    days query int false "Количество дней неуверенности (default: 7)"
// @Success  200 {object} response.Envelope
// @Router   /advisor/inactive-students [get]
func (h *AdvisorHandler) GetInactiveStudents(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	students, err := h.advisorService.GetInactiveStudents(r.Context(), advisorID, days)
	if err != nil {
		response.InternalError(w, "failed to fetch inactive students")
		return
	}

	response.OK(w, students)
}

// SendReminder godoc
// @Summary  Отправить напоминание
// @Tags     advisor
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.SendReminderRequest true "Напоминание"
// @Success  200 {object} response.Envelope
// @Router   /advisor/reminders [post]
func (h *AdvisorHandler) SendReminder(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	var req dto.SendReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	if err := h.advisorService.SendReminder(r.Context(), advisorID, &req); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "reminder sent successfully"})
}

// FileComplaint godoc
// @Summary  Подать жалобу
// @Tags     advisor
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateComplaintRequest true "Жалоба"
// @Success  201 {object} response.Envelope
// @Router   /advisor/complaints [post]
func (h *AdvisorHandler) FileComplaint(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	var req dto.CreateComplaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	if err := h.advisorService.FileComplaint(r.Context(), advisorID, &req); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "complaint filed successfully"})
}

// GetAnalytics godoc
// @Summary  Аналитика по менторам и студентам эдвайзера
// @Tags     advisor
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /advisor/analytics [get]
func (h *AdvisorHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	advisorIDStr := middleware.GetUserID(r)
	advisorID, err := uuid.Parse(advisorIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid advisor identity")
		return
	}

	analytics, err := h.advisorService.GetAnalytics(r.Context(), advisorID)
	if err != nil {
		response.InternalError(w, "failed to get analytics")
		return
	}

	response.OK(w, analytics)
}
