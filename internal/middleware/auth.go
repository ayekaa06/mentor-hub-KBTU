// Package middleware содержит HTTP middleware для chi-роутера.
package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtpkg "mentorhub/internal/pkg/jwt"
	"mentorhub/internal/pkg/response"
)

// contextKey — приватный тип для ключей контекста, исключает коллизии.
type contextKey string

const (
	// ContextKeyUserID — ключ UUID текущего пользователя в контексте запроса.
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyRole — ключ роли текущего пользователя в контексте запроса.
	ContextKeyRole contextKey = "role"
)

// Auth — middleware для проверки JWT access token из заголовка Authorization.
// Формат: "Authorization: Bearer <token>"
func Auth(jwtManager *jwtpkg.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Unauthorized(w, "authorization header format: Bearer <token>")
				return
			}

			claims, err := jwtManager.Parse(parts[1])
			if err != nil {
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			// Принимаем только access token (не refresh)
			if claims.Type != "access" {
				response.Unauthorized(w, "access token required")
				return
			}

			// Пробрасываем userID и role в контекст
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, string(claims.Role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID извлекает UUID пользователя из контекста запроса.
func GetUserID(r *http.Request) string {
	v, _ := r.Context().Value(ContextKeyUserID).(string)
	return v
}

// GetRole извлекает роль пользователя из контекста запроса.
func GetRole(r *http.Request) string {
	v, _ := r.Context().Value(ContextKeyRole).(string)
	return v
}
