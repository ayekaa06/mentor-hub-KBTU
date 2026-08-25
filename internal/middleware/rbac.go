package middleware

import (
	"net/http"

	"mentorhub/internal/pkg/response"
)

// RequireRole — middleware для RBAC (Role-Based Access Control).
// Принимает список разрешённых ролей. Должен вызываться ПОСЛЕ Auth middleware.
//
// Пример:
//
//	r.Use(middleware.Auth(jwtManager))
//	r.Use(middleware.RequireRole("head"))
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	// Для быстрого O(1) поиска превращаем слайс в map
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetRole(r)
			if _, ok := allowed[userRole]; !ok {
				response.Forbidden(w, "insufficient permissions for this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
