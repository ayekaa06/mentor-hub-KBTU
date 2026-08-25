// Package handler содержит HTTP-обработчики и роутер приложения.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"mentorhub/internal/dto"
	"mentorhub/internal/middleware"
	"mentorhub/internal/pkg/response"
	"mentorhub/internal/service"
)

// AuthHandler обрабатывает Auth-запросы: register, login, refresh, me, logout.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создаёт AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary     Регистрация пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body dto.RegisterRequest true "Данные пользователя"
// @Success     201 {object} response.Envelope{data=dto.TokenResponse}
// @Failure     400 {object} response.Envelope
// @Failure     409 {object} response.Envelope
// @Router      /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	tokens, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			response.Conflict(w, "email already registered")
		default:
			response.InternalError(w, "registration failed")
		}
		return
	}

	response.Created(w, tokens)
}

// Login godoc
// @Summary     Вход в систему
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body dto.LoginRequest true "Учётные данные"
// @Success     200 {object} response.Envelope{data=dto.TokenResponse}
// @Failure     401 {object} response.Envelope
// @Router      /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	tokens, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Unauthorized(w, "invalid email or password")
		case errors.Is(err, service.ErrUserInactive):
			response.Forbidden(w, "account is deactivated, contact administrator")
		default:
			response.InternalError(w, "login failed")
		}
		return
	}

	response.OK(w, tokens)
}

// Refresh godoc
// @Summary     Обновление пары токенов
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body dto.RefreshRequest true "Refresh token"
// @Success     200 {object} response.Envelope{data=dto.TokenResponse}
// @Failure     401 {object} response.Envelope
// @Router      /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(w, "invalid or expired refresh token")
		return
	}

	response.OK(w, tokens)
}

// Me godoc
// @Summary     Получить данные текущего пользователя
// @Tags        auth
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} response.Envelope{data=dto.MeResponse}
// @Failure     401 {object} response.Envelope
// @Router      /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid user identity in token")
		return
	}

	user, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		response.NotFound(w, "user not found")
		return
	}

	response.OK(w, dto.MeResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
	})
}

// UpdateProfile godoc
// @Summary     Обновить данные текущего пользователя
// @Tags        auth
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body dto.UpdateProfileRequest true "Новые данные профиля"
// @Success     200 {object} response.Envelope{data=dto.MeResponse}
// @Failure     400 {object} response.Envelope
// @Failure     401 {object} response.Envelope
// @Router      /auth/me [patch]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(w, "invalid user identity in token")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		response.InternalError(w, "failed to update profile")
		return
	}

	response.OK(w, dto.MeResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
	})
}

// Logout godoc
// @Summary     Выход из системы
// @Tags        auth
// @Security    BearerAuth
// @Success     200 {object} response.Envelope
// @Router      /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Phase 1: stateless JWT — клиент удаляет токены сам.
	// Phase 2: добавим инвалидацию refresh token в БД.
	response.OK(w, map[string]string{"message": "logged out successfully"})
}
