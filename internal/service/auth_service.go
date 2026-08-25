// Package service содержит бизнес-логику приложения.
// Сервисы не зависят от HTTP-слоя и напрямую от SQL — только от интерфейсов repository.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/pkg/hasher"
	jwtpkg "mentorhub/internal/pkg/jwt"
	"mentorhub/internal/repository"
)

// ── Sentinel Errors ───────────────────────────────────────────────────────────

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("account is deactivated")
	ErrEmailExists        = errors.New("email already registered")
)

// ── AuthService ───────────────────────────────────────────────────────────────

// AuthService отвечает за регистрацию, вход и управление JWT-токенами.
type AuthService struct {
	userRepo       repository.UserRepository
	jwtManager     *jwtpkg.Manager
	accessTokenTTL time.Duration
}

// NewAuthService создаёт AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *jwtpkg.Manager,
	accessTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		accessTokenTTL: accessTTL,
	}
}

// Register создаёт нового пользователя и возвращает токены.
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	// Проверяем уникальность email
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailExists
	}

	hash, err := hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         domain.Role(req.Role),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.buildTokenPair(user)
}

// Login аутентифицирует пользователя по email/password и возвращает токены.
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Не раскрываем причину — возвращаем один и тот же generic error
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if !hasher.Check(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.buildTokenPair(user)
}

// Refresh обновляет пару токенов по refresh token.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.jwtManager.Parse(refreshToken)
	if err != nil {
		return nil, jwtpkg.ErrInvalidToken
	}

	if claims.Type != "refresh" {
		return nil, jwtpkg.ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, jwtpkg.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	return s.buildTokenPair(user)
}

// GetMe возвращает данные текущего пользователя.
func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateProfile обновляет данные профиля пользователя.
func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user profile: %w", err)
	}

	return user, nil
}

// buildTokenPair генерирует access + refresh токены для пользователя.
func (s *AuthService) buildTokenPair(user *domain.User) (*dto.TokenResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTokenTTL.Seconds()),
	}, nil
}
