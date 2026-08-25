package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/pkg/hasher"
	"mentorhub/internal/repository"
)

// UserService управляет пользователями (Head-уровень операций).
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService создаёт UserService.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetAll возвращает пагинированный список пользователей с опциональным фильтром по роли.
func (s *UserService) GetAll(ctx context.Context, role *domain.Role, page, perPage int) ([]*domain.User, int, error) {
	offset := (page - 1) * perPage
	return s.userRepo.GetAll(ctx, role, perPage, offset)
}

// GetByID возвращает пользователя по ID.
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Create создаёт нового пользователя (вызывается Head-ом).
func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*domain.User, error) {
	// Проверяем уникальность email
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailExists
	}

	hash, err := hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         domain.Role(req.Role),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user.PasswordHash = "" // не возвращаем хеш
	return user, nil
}

// Deactivate деактивирует пользователя.
func (s *UserService) Deactivate(ctx context.Context, id uuid.UUID) error {
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return ErrUserNotFound
	}
	return s.userRepo.Deactivate(ctx, id)
}

// Delete удаляет пользователя из системы (только Head).
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(ctx, id)
}
