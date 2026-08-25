// Package hasher предоставляет bcrypt-хеширование паролей.
package hasher

import "golang.org/x/crypto/bcrypt"

// cost — стоимость bcrypt (12 — хороший баланс безопасности и скорости).
const cost = 12

// Hash возвращает bcrypt-хеш пароля.
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Check проверяет, соответствует ли открытый пароль хешу.
// Возвращает true если совпадают.
func Check(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
