package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// wrappedWriter перехватывает статус-код ответа для логирования.
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Logger — middleware для структурированного логирования каждого HTTP-запроса.
// Выводит: method, path, status, latency, request_id, ip.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		// Уровень лога зависит от статус-кода
		event := log.Info()
		if ww.statusCode >= 500 {
			event = log.Error()
		} else if ww.statusCode >= 400 {
			event = log.Warn()
		}

		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("latency", latency).
			Str("ip", r.RemoteAddr).
			Str("request_id", r.Header.Get("X-Request-ID")).
			Msg("request")
	})
}
