// Package response предоставляет стандартизированный JSON-ответ для всех эндпоинтов.
// Формат: { "success": bool, "data": any, "error": string, "meta": {...} }
package response

import (
	"encoding/json"
	"net/http"
)

// Envelope — стандартная обёртка для всех ответов API.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta — метаданные пагинации.
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// JSON отправляет произвольный JSON-ответ с указанным статус-кодом.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Success: true, Data: data})
}

// JSONWithMeta отправляет ответ с пагинацией.
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Success: true, Data: data, Meta: meta})
}

// Error отправляет JSON-ответ с ошибкой.
func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Success: false, Error: msg})
}

// Shorthand helpers ──────────────────────────────────────────────────────────

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, msg)
}

func Forbidden(w http.ResponseWriter, msg string) {
	Error(w, http.StatusForbidden, msg)
}

func NotFound(w http.ResponseWriter, msg string) {
	Error(w, http.StatusNotFound, msg)
}

func Conflict(w http.ResponseWriter, msg string) {
	Error(w, http.StatusConflict, msg)
}

func InternalError(w http.ResponseWriter, msg string) {
	Error(w, http.StatusInternalServerError, msg)
}
