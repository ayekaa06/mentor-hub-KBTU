// Package dto содержит структуры Request/Response для Auth эндпоинтов.
package dto

// LoginRequest — тело запроса POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest — тело запроса POST /auth/register.
// Регистрацию пользователей с любой ролью обычно делает Head,
// но в Phase 1 оставляем открытой для тестирования.
type RegisterRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8,max=72"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name"  validate:"required,min=1,max=100"`
	Role      string `json:"role"       validate:"required,oneof=head advisor mentor freshman"`
}

// RefreshRequest — тело запроса POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse — ответ на успешную аутентификацию.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // всегда "Bearer"
	ExpiresIn    int    `json:"expires_in"` // секунд до истечения access token
}

// MeResponse — ответ GET /auth/me.
type MeResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsActive  bool    `json:"is_active"`
}
