package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "mentorhub/docs"
	"mentorhub/internal/middleware"
	jwtpkg "mentorhub/internal/pkg/jwt"
	"mentorhub/internal/service"
)

// NewRouter создаёт и возвращает настроенный chi-роутер.
// Все роуты разделены по ролям через RBAC middleware.
func NewRouter(
	authService *service.AuthService,
	headService *service.HeadService,
	userService *service.UserService,
	advisorService *service.AdvisorService,
	mentorService *service.MentorService,
	freshmanService *service.FreshmanService,
	jwtManager *jwtpkg.Manager,
) http.Handler {
	r := chi.NewRouter()

	// ── Глобальные middleware ─────────────────────────────────────────────────
	r.Use(chimiddleware.Recoverer)     // не падаем при панике
	r.Use(chimiddleware.RequestID)     // X-Request-Id header
	r.Use(chimiddleware.RealIP)        // X-Forwarded-For → RemoteAddr
	r.Use(middleware.CORS)             // Cross-Origin
	r.Use(middleware.Logger)           // структурированный лог

	// Обработчики
	authHandler    := NewAuthHandler(authService)
	headHandler    := NewHeadHandler(headService, userService)
	advisorHandler := NewAdvisorHandler(advisorService)
	mentorHandler  := NewMentorHandler(mentorService)
	freshmanHandler := NewFreshmanHandler(freshmanService)

	// ── Root & Health Check & Swagger (публичные) ────────────────────────────
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"mentorhub"}`))
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// ── API v1 ────────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// ── Публичные Auth роуты ─────────────────────────────────────────────
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register) // POST /api/v1/auth/register
			r.Post("/login", authHandler.Login)       // POST /api/v1/auth/login
			r.Post("/refresh", authHandler.Refresh)   // POST /api/v1/auth/refresh
		})

		// ── Защищённые роуты (требуют JWT) ───────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))

			r.Post("/auth/logout", authHandler.Logout) // POST /api/v1/auth/logout
			r.Get("/auth/me", authHandler.Me)           // GET  /api/v1/auth/me
			r.Patch("/auth/me", authHandler.UpdateProfile) // PATCH /api/v1/auth/me
			r.Get("/mentors", headHandler.ListUsers) // GET /api/v1/mentors

			// ── HEAD ── (только role=head) ────────────────────────────────────
			r.Route("/head", func(r chi.Router) {
				r.Use(middleware.RequireRole("head"))

				// Dashboard
				r.Get("/dashboard", headHandler.GetDashboard)

				// Пользователи
				r.Get("/users", headHandler.ListUsers)
				r.Post("/users", headHandler.CreateUser)
				r.Put("/users/{id}/deactivate", headHandler.DeactivateUser)
				r.Delete("/users/{id}", headHandler.DeleteUser)

				// Учебные годы
				r.Get("/academic-years", headHandler.ListAcademicYears)
				r.Post("/academic-years", headHandler.CreateAcademicYear)
				r.Put("/academic-years/{id}/activate", headHandler.SetActiveYear)

				// Факультеты и специальности
				r.Get("/faculties", headHandler.ListFaculties)
				r.Post("/faculties", headHandler.CreateFaculty)
				r.Post("/specialties", headHandler.CreateSpecialty)

				// Назначения
				r.Post("/assign/mentor", headHandler.AssignMentor)

				// Шаблоны задач
				r.Get("/task-templates", headHandler.ListTaskTemplates)
				r.Post("/task-templates", headHandler.CreateTaskTemplate)
				r.Put("/task-templates/{id}", headHandler.UpdateTaskTemplate)
				r.Delete("/task-templates/{id}", headHandler.DeleteTaskTemplate)
				r.Post("/task-templates/{id}/assign", headHandler.AssignTasks)

				// Аналитика и жалобы
				r.Get("/analytics", headHandler.GetAnalytics)
				r.Get("/complaints", headHandler.ListComplaints)
				r.Put("/complaints/{id}", headHandler.UpdateComplaintStatus)
			})

			// ── ADVISOR ── (только role=advisor) ──────────────────────────────
			r.Route("/advisor", func(r chi.Router) {
				r.Use(middleware.RequireRole("advisor"))

				r.Get("/dashboard", advisorHandler.GetDashboard)
				r.Get("/mentors", advisorHandler.GetMentors)
				r.Get("/mentors/{id}/students", advisorHandler.GetMentorStudents)
				r.Get("/mentors/{id}/tasks", advisorHandler.GetMentorTaskProgress)
				r.Get("/inactive-students", advisorHandler.GetInactiveStudents)
				r.Post("/reminders", advisorHandler.SendReminder)
				r.Post("/complaints", advisorHandler.FileComplaint)
				r.Get("/analytics", advisorHandler.GetAnalytics)
			})

			// ── MENTOR ── (только role=mentor) ────────────────────────────────
			r.Route("/mentor", func(r chi.Router) {
				r.Use(middleware.RequireRole("mentor"))

				r.Get("/dashboard", mentorHandler.GetDashboard)
				r.Get("/group", mentorHandler.GetGroup)
				r.Get("/group/{id}", mentorHandler.GetStudentDetail)
				r.Get("/tasks", mentorHandler.GetTasks)
				r.Put("/tasks/{id}/approve", mentorHandler.ApproveTask)
				r.Put("/tasks/{id}/reject", mentorHandler.RejectTask)
				r.Post("/meetings", mentorHandler.CreateMeeting)
				r.Get("/meetings", mentorHandler.GetMeetings)
				r.Put("/meetings/{id}/complete", mentorHandler.CompleteMeeting)
				r.Post("/announcements", mentorHandler.CreateAnnouncement)
				r.Get("/announcements", mentorHandler.GetAnnouncements)

				// Вопросы от freshman-ов
				r.Get("/questions", mentorHandler.GetQuestions)
				r.Put("/questions/{id}/answer", mentorHandler.AnswerQuestion)
			})

			// ── FRESHMAN ── (только role=freshman) ────────────────────────────
			r.Route("/freshman", func(r chi.Router) {
				r.Use(middleware.RequireRole("freshman"))

				// Dashboard
				r.Get("/dashboard", freshmanHandler.GetDashboard)

				// Задачи
				r.Get("/tasks", freshmanHandler.GetTasks)
				r.Get("/tasks/{id}", freshmanHandler.GetTask)
				r.Put("/tasks/{id}/submit", freshmanHandler.SubmitTask)

				// Встречи
				r.Get("/meetings", freshmanHandler.GetMeetings)

				// Объявления
				r.Get("/announcements", freshmanHandler.GetAnnouncements)

				// FAQ
				r.Get("/faq", freshmanHandler.GetFAQ)

				// Уведомления
				r.Get("/notifications", freshmanHandler.GetNotifications)
				r.Put("/notifications/read-all", freshmanHandler.MarkAllNotificationsRead)
				r.Put("/notifications/{id}/read", freshmanHandler.MarkNotificationRead)

				// Вопросы ментору
				r.Get("/questions", freshmanHandler.GetMyQuestions)
				r.Post("/questions", freshmanHandler.AskMentor)
			})
		})
	})

	return r
}


