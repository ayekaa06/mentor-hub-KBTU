package handler

import (
	"encoding/json"
	"errors"
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

// HeadHandler обрабатывает все запросы от роли Head.
type HeadHandler struct {
	headService *service.HeadService
	userService *service.UserService
}

// NewHeadHandler создаёт HeadHandler.
func NewHeadHandler(headService *service.HeadService, userService *service.UserService) *HeadHandler {
	return &HeadHandler{
		headService: headService,
		userService: userService,
	}
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

// GetDashboard godoc
// @Summary  Статистика системы для Head
// @Tags     head
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Envelope
// @Router   /head/dashboard [get]
func (h *HeadHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.headService.GetDashboard(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get dashboard stats")
		return
	}
	response.OK(w, stats)
}

// ── Users ─────────────────────────────────────────────────────────────────────

// ListUsers godoc
// @Summary  Список пользователей с фильтром по роли
// @Tags     head
// @Security BearerAuth
// @Param    role     query string false "Фильтр по роли (advisor|mentor|freshman)"
// @Param    page     query int    false "Страница (default: 1)"
// @Param    per_page query int    false "Размер страницы (default: 20)"
// @Success  200 {object} response.Envelope
// @Router   /head/users [get]
func (h *HeadHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	var roleFilter *domain.Role
	if roleStr := r.URL.Query().Get("role"); roleStr != "" {
		r := domain.Role(roleStr)
		roleFilter = &r
	}

	users, total, err := h.userService.GetAll(r.Context(), roleFilter, page, perPage)
	if err != nil {
		response.InternalError(w, "failed to fetch users")
		return
	}

	// Скрываем password_hash
	for _, u := range users {
		u.PasswordHash = ""
	}

	totalPages := (total + perPage - 1) / perPage
	response.JSONWithMeta(w, http.StatusOK, users, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateUser godoc
// @Summary  Создать пользователя (Head-только)
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateUserRequest true "Данные пользователя"
// @Success  201 {object} response.Envelope
// @Router   /head/users [post]
func (h *HeadHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			response.Conflict(w, "email already registered")
		default:
			response.InternalError(w, "failed to create user")
		}
		return
	}

	response.Created(w, user)
}

// DeactivateUser godoc
// @Summary  Деактивировать пользователя
// @Tags     head
// @Security BearerAuth
// @Param    id path string true "User UUID"
// @Success  200 {object} response.Envelope
// @Router   /head/users/{id}/deactivate [put]
func (h *HeadHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid user ID")
		return
	}

	if err := h.userService.Deactivate(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to deactivate user")
		return
	}

	response.OK(w, map[string]string{"message": "user deactivated"})
}

// DeleteUser godoc
// @Summary  Удалить пользователя
// @Tags     head
// @Security BearerAuth
// @Param    id path string true "User UUID"
// @Success  204
// @Router   /head/users/{id} [delete]
func (h *HeadHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid user ID")
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to delete user")
		return
	}

	response.NoContent(w)
}

// ── Academic Years ────────────────────────────────────────────────────────────

// CreateAcademicYear godoc
// @Summary  Создать учебный год
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateAcademicYearRequest true "Учебный год"
// @Success  201 {object} response.Envelope
// @Router   /head/academic-years [post]
func (h *HeadHandler) CreateAcademicYear(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAcademicYearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	ay, err := h.headService.CreateAcademicYear(r.Context(), &req)
	if err != nil {
		response.InternalError(w, "failed to create academic year")
		return
	}
	response.Created(w, ay)
}

// ListAcademicYears godoc
// @Summary  Список учебных лет
// @Tags     head
// @Security BearerAuth
// @Success  200 {object} response.Envelope
// @Router   /head/academic-years [get]
func (h *HeadHandler) ListAcademicYears(w http.ResponseWriter, r *http.Request) {
	years, err := h.headService.GetAcademicYears(r.Context())
	if err != nil {
		response.InternalError(w, "failed to fetch academic years")
		return
	}
	response.OK(w, years)
}

// SetActiveYear godoc
// @Summary  Сделать учебный год активным
// @Tags     head
// @Security BearerAuth
// @Param    id path string true "AcademicYear UUID"
// @Success  200 {object} response.Envelope
// @Router   /head/academic-years/{id}/activate [put]
func (h *HeadHandler) SetActiveYear(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid academic year ID")
		return
	}
	if err := h.headService.SetActiveYear(r.Context(), id); err != nil {
		response.InternalError(w, "failed to set active year")
		return
	}
	response.OK(w, map[string]string{"message": "academic year activated"})
}

// ── Faculties & Specialties ───────────────────────────────────────────────────

// CreateFaculty godoc
// @Summary  Создать факультет
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateFacultyRequest true "Факультет"
// @Success  201 {object} response.Envelope
// @Router   /head/faculties [post]
func (h *HeadHandler) CreateFaculty(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFacultyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	f, err := h.headService.CreateFaculty(r.Context(), &req)
	if err != nil {
		response.InternalError(w, "failed to create faculty")
		return
	}
	response.Created(w, f)
}

// ListFaculties godoc
// @Summary  Список факультетов со специальностями
// @Tags     head
// @Security BearerAuth
// @Success  200 {object} response.Envelope
// @Router   /head/faculties [get]
func (h *HeadHandler) ListFaculties(w http.ResponseWriter, r *http.Request) {
	faculties, err := h.headService.GetFaculties(r.Context())
	if err != nil {
		response.InternalError(w, "failed to fetch faculties")
		return
	}
	response.OK(w, faculties)
}

// CreateSpecialty godoc
// @Summary  Создать специальность
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateSpecialtyRequest true "Специальность"
// @Success  201 {object} response.Envelope
// @Router   /head/specialties [post]
func (h *HeadHandler) CreateSpecialty(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSpecialtyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	spec, err := h.headService.CreateSpecialty(r.Context(), &req)
	if err != nil {
		response.InternalError(w, "failed to create specialty")
		return
	}
	response.Created(w, spec)
}

// ── Assignments ───────────────────────────────────────────────────────────────

// AssignMentor godoc
// @Summary  Назначить ментора к эдвайзеру
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.AssignMentorRequest true "Назначение"
// @Success  200 {object} response.Envelope
// @Router   /head/assign/mentor [post]
func (h *HeadHandler) AssignMentor(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignMentorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	mentorID, err := uuid.Parse(req.MentorID)
	if err != nil {
		response.BadRequest(w, "invalid mentor_id")
		return
	}
	advisorID, err := uuid.Parse(req.AdvisorID)
	if err != nil {
		response.BadRequest(w, "invalid advisor_id")
		return
	}

	if err := h.headService.AssignMentorToAdvisor(r.Context(), mentorID, advisorID); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.NotFound(w, "user not found")
		default:
			response.InternalError(w, err.Error())
		}
		return
	}
	response.OK(w, map[string]string{"message": "mentor assigned to advisor"})
}

// ── Task Templates ────────────────────────────────────────────────────────────

// ListTaskTemplates godoc
// @Summary  Список шаблонов задач
// @Tags     head
// @Security BearerAuth
// @Param    active_only query bool false "Только активные"
// @Success  200 {object} response.Envelope
// @Router   /head/task-templates [get]
func (h *HeadHandler) ListTaskTemplates(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"
	templates, err := h.headService.GetTaskTemplates(r.Context(), activeOnly)
	if err != nil {
		response.InternalError(w, "failed to fetch task templates")
		return
	}
	response.OK(w, templates)
}

// CreateTaskTemplate godoc
// @Summary  Создать шаблон задачи
// @Tags     head
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateTaskTemplateRequest true "Шаблон"
// @Success  201 {object} response.Envelope
// @Router   /head/task-templates [post]
func (h *HeadHandler) CreateTaskTemplate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	creatorIDStr := middleware.GetUserID(r)
	creatorID, _ := uuid.Parse(creatorIDStr)

	tmpl, err := h.headService.CreateTaskTemplate(r.Context(), &req, creatorID)
	if err != nil {
		response.InternalError(w, "failed to create task template")
		return
	}
	response.Created(w, tmpl)
}

// UpdateTaskTemplate godoc
// @Summary  Обновить шаблон задачи
// @Tags     head
// @Security BearerAuth
// @Param    id   path string                       true "Template UUID"
// @Param    body body dto.UpdateTaskTemplateRequest true "Обновление"
// @Success  200 {object} response.Envelope
// @Router   /head/task-templates/{id} [put]
func (h *HeadHandler) UpdateTaskTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}

	var req dto.UpdateTaskTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	tmpl, err := h.headService.UpdateTaskTemplate(r.Context(), id, &req)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.OK(w, tmpl)
}

// DeleteTaskTemplate godoc
// @Summary  Удалить шаблон задачи
// @Tags     head
// @Security BearerAuth
// @Param    id path string true "Template UUID"
// @Success  204
// @Router   /head/task-templates/{id} [delete]
func (h *HeadHandler) DeleteTaskTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}
	if err := h.headService.DeleteTaskTemplate(r.Context(), id); err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.NoContent(w)
}

// AssignTasks godoc
// @Summary  Массово назначить шаблон задачи первокурсникам
// @Tags     head
// @Security BearerAuth
// @Param    id   path string              true "Template UUID"
// @Param    body body dto.AssignTasksRequest true "Список freshman IDs"
// @Success  200 {object} response.Envelope
// @Router   /head/task-templates/{id}/assign [post]
func (h *HeadHandler) AssignTasks(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")

	var req dto.AssignTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	req.TemplateID = templateID // берём из URL

	assignerIDStr := middleware.GetUserID(r)
	assignerID, _ := uuid.Parse(assignerIDStr)

	count, err := h.headService.AssignTasksToFreshmen(r.Context(), &req, assignerID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, map[string]int{"assigned_count": count})
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// GetAnalytics godoc
// @Summary  Аналитика системы
// @Tags     head
// @Security BearerAuth
// @Success  200 {object} response.Envelope
// @Router   /head/analytics [get]
func (h *HeadHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.headService.GetAnalytics(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get analytics")
		return
	}
	response.OK(w, stats)
}

// ── Complaints ────────────────────────────────────────────────────────────────

// ListComplaints godoc
// @Summary  Список жалоб
// @Tags     head
// @Security BearerAuth
// @Param    status query string false "Фильтр: open|in_review|resolved|dismissed"
// @Success  200 {object} response.Envelope
// @Router   /head/complaints [get]
func (h *HeadHandler) ListComplaints(w http.ResponseWriter, r *http.Request) {
	var statusFilter *domain.ComplaintStatus
	if s := r.URL.Query().Get("status"); s != "" {
		cs := domain.ComplaintStatus(s)
		statusFilter = &cs
	}

	complaints, err := h.headService.GetComplaints(r.Context(), statusFilter)
	if err != nil {
		response.InternalError(w, "failed to fetch complaints")
		return
	}
	response.OK(w, complaints)
}

// UpdateComplaintStatus godoc
// @Summary  Обновить статус жалобы
// @Tags     head
// @Security BearerAuth
// @Param    id   path string                           true "Complaint UUID"
// @Param    body body dto.UpdateComplaintStatusRequest true "Новый статус"
// @Success  200 {object} response.Envelope
// @Router   /head/complaints/{id} [put]
func (h *HeadHandler) UpdateComplaintStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid complaint ID")
		return
	}

	var req dto.UpdateComplaintStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	reviewerIDStr := middleware.GetUserID(r)
	reviewerID, _ := uuid.Parse(reviewerIDStr)

	if err := h.headService.UpdateComplaintStatus(r.Context(), id, &req, reviewerID); err != nil {
		response.InternalError(w, "failed to update complaint")
		return
	}
	response.OK(w, map[string]string{"message": "complaint updated"})
}
