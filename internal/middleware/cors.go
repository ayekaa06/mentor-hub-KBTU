package middleware

import "net/http"

// CORS устанавливает заголовки Cross-Origin Resource Sharing.
// В production замените "*" на конкретный домен фронтенда.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count, Link")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 часа

		// Preflight request — браузер спрашивает разрешение перед реальным запросом
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
