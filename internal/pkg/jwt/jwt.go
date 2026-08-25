// Package jwt управляет созданием и валидацией JWT-токенов.
// Алгоритм: HMAC-SHA256 (HS256).
// Access token: короткоживущий (15 мин по умолчанию).
// Refresh token: долгоживущий (7 дней по умолчанию).
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"mentorhub/internal/domain"
)

var (
	// ErrInvalidToken — токен не прошёл валидацию.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken — токен истёк.
	ErrExpiredToken = errors.New("token expired")
)

// Claims — кастомные claims JWT.
type Claims struct {
	UserID string      `json:"uid"`
	Role   domain.Role `json:"role"`
	Type   string      `json:"typ"` // "access" | "refresh"
	jwt.RegisteredClaims
}

// Manager управляет JWT токенами.
type Manager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewManager создаёт новый Manager.
func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// GenerateAccessToken создаёт короткоживущий access token.
func (m *Manager) GenerateAccessToken(userID uuid.UUID, role domain.Role) (string, error) {
	return m.generate(userID.String(), role, "access", m.accessTokenTTL)
}

// GenerateRefreshToken создаёт долгоживущий refresh token.
func (m *Manager) GenerateRefreshToken(userID uuid.UUID, role domain.Role) (string, error) {
	return m.generate(userID.String(), role, "refresh", m.refreshTokenTTL)
}

// Parse разбирает и валидирует токен. Возвращает Claims при успехе.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// generate — внутренний метод генерации токена.
func (m *Manager) generate(userID string, role domain.Role, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}
