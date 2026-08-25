// Package postgres содержит PostgreSQL реализации repository-интерфейсов.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mentorhub/internal/domain"
	"mentorhub/internal/repository"
)

// userRepository реализует repository.UserRepository через pgxpool.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создаёт новый userRepository.
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

// Create добавляет нового пользователя и заполняет CreatedAt/UpdatedAt из БД.
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (id, email, password_hash, role, first_name, last_name, avatar_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, q,
		user.ID, user.Email, user.PasswordHash, user.Role,
		user.FirstName, user.LastName, user.AvatarURL, user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// GetByID возвращает пользователя по UUID.
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, role, first_name, last_name,
		       avatar_url, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s not found", id)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

// GetByEmail возвращает пользователя по email.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, role, first_name, last_name,
		       avatar_url, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

// GetAll возвращает список пользователей с фильтром по роли и пагинацией.
func (r *userRepository) GetAll(ctx context.Context, role *domain.Role, limit, offset int) ([]*domain.User, int, error) {
	// Считаем total
	countQ := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []any{}
	idx := 1

	if role != nil {
		countQ += fmt.Sprintf(" AND role = $%d", idx)
		args = append(args, *role)
		idx++
	}

	var total int
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Получаем данные
	dataQ := fmt.Sprintf(`
		SELECT id, email, password_hash, role, first_name, last_name,
		       avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE 1=1 %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`,
		func() string {
			if role != nil {
				return fmt.Sprintf("AND role = $%d", idx-1)
			}
			return ""
		}(),
		idx, idx+1,
	)

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &u.AvatarURL, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}

	return users, total, rows.Err()
}

// Update обновляет профиль пользователя (имя, аватар).
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users
		SET first_name = $2, last_name = $3, avatar_url = $4
		WHERE id = $1
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, q,
		user.ID, user.FirstName, user.LastName, user.AvatarURL,
	).Scan(&user.UpdatedAt)
}

// Deactivate деактивирует пользователя (не удаляет из БД).
func (r *userRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET is_active = FALSE WHERE id = $1`, id)
	return err
}

// Delete удаляет пользователя (CASCADE удалит связанные записи).
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
