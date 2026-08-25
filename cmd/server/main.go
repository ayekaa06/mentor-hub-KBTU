// MentorHub Backend — точка входа.
// Инициализирует конфиг, подключение к БД, все сервисы и запускает HTTP-сервер
// с поддержкой graceful shutdown.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"mentorhub/internal/config"
	"mentorhub/internal/handler"
	jwtpkg "mentorhub/internal/pkg/jwt"
	"mentorhub/internal/repository/postgres"
	"mentorhub/internal/service"
)

func main() {
	// ── Logger ────────────────────────────────────────────────────────────────
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	log.Info().
		Str("port", cfg.Server.Port).
		Str("db_host", cfg.Database.Host).
		Str("db_name", cfg.Database.DBName).
		Msg("config loaded")

	// ── Database ──────────────────────────────────────────────────────────────
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse database DSN")
	}

	poolCfg.MaxConns = 20
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database pool")
	}
	defer pool.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Fatal().Err(err).Str("host", cfg.Database.Host).Msg("database ping failed")
	}
	log.Info().Msg("✅ connected to PostgreSQL")

	// ── JWT Manager ───────────────────────────────────────────────────────────
	if cfg.JWT.Secret == "" {
		log.Fatal().Msg("JWT_SECRET is required — set it in .env")
	}

	jwtManager := jwtpkg.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo         := postgres.NewUserRepository(pool)
	academicYearRepo := postgres.NewAcademicYearRepository(pool)
	facultyRepo      := postgres.NewFacultyRepository(pool)
	taskTemplateRepo := postgres.NewTaskTemplateRepository(pool)
	taskRepo         := postgres.NewTaskRepository(pool)
	assignRepo       := postgres.NewAssignmentRepository(pool)
	complaintRepo    := postgres.NewComplaintRepository(pool)
	advisorRepo      := postgres.NewAdvisorRepository(pool)
	meetingRepo      := postgres.NewMeetingRepository(pool)
	notificationRepo := postgres.NewNotificationRepository(pool)
	faqRepo          := postgres.NewFAQRepository(pool)
	questionRepo     := postgres.NewQuestionRepository(pool)

	// ── Services ──────────────────────────────────────────────────────────────
	authService := service.NewAuthService(
		userRepo,
		jwtManager,
		cfg.JWT.AccessTokenTTL,
	)

	userService := service.NewUserService(userRepo)

	headService := service.NewHeadService(
		userRepo,
		academicYearRepo,
		facultyRepo,
		taskTemplateRepo,
		taskRepo,
		assignRepo,
		complaintRepo,
	)

	advisorService := service.NewAdvisorService(
		advisorRepo,
		taskRepo,
		notificationRepo,
		complaintRepo,
	)

	mentorService := service.NewMentorService(
		assignRepo,
		taskRepo,
		meetingRepo,
		notificationRepo,
		questionRepo,
	)

	freshmanService := service.NewFreshmanService(
		assignRepo,
		taskRepo,
		meetingRepo,
		notificationRepo,
		faqRepo,
		questionRepo,
	)

	// ── Router ────────────────────────────────────────────────────────────────
	router := handler.NewRouter(
		authService,
		headService,
		userService,
		advisorService,
		mentorService,
		freshmanService,
		jwtManager,
	)

	// ── HTTP Server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("🚀 MentorHub server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server listen error")
		}
	}()

	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("shutdown signal received")

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	log.Info().Msg("✅ server stopped gracefully")
}
