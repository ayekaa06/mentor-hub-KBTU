package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"mentorhub/internal/domain"
	"mentorhub/internal/dto"
	"mentorhub/internal/pkg/hasher"
	jwtpkg "mentorhub/internal/pkg/jwt"
	"mentorhub/internal/service"
)

// mockUserRepo — in-memory реализация repository.UserRepository для тестов.
type mockUserRepo struct {
	users map[string]*domain.User
	byID  map[uuid.UUID]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*domain.User),
		byID:  make(map[uuid.UUID]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	m.users[u.Email] = u
	m.byID[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, service.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, service.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetAll(ctx context.Context, role *domain.Role, limit, offset int) ([]*domain.User, int, error) {
	var list []*domain.User
	for _, u := range m.users {
		if role == nil || u.Role == *role {
			list = append(list, u)
		}
	}
	return list, len(list), nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *domain.User) error {
	m.users[u.Email] = u
	m.byID[u.ID] = u
	return nil
}

func (m *mockUserRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	if u, ok := m.byID[id]; ok {
		u.IsActive = false
	}
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if u, ok := m.byID[id]; ok {
		delete(m.users, u.Email)
		delete(m.byID, id)
	}
	return nil
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	repo := newMockUserRepo()
	jwtMgr := jwtpkg.NewManager("secret_key_for_testing_12345", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repo, jwtMgr, 15*time.Minute)

	ctx := context.Background()

	// 1. Регистрация нового пользователя
	regReq := &dto.RegisterRequest{
		Email:     "test@mentorhub.com",
		Password:  "password123",
		Role:      "freshman",
		FirstName: "Иван",
		LastName:  "Тестов",
	}

	tokens, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Errorf("expected valid tokens, got empty")
	}

	// 2. Повторная регистрация с тем же email
	_, err = authSvc.Register(ctx, regReq)
	if err != service.ErrEmailExists {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}

	// 3. Успешный вход
	loginReq := &dto.LoginRequest{
		Email:    "test@mentorhub.com",
		Password: "password123",
	}
	loginTokens, err := authSvc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}
	if loginTokens.AccessToken == "" {
		t.Errorf("expected access token on login")
	}

	// 4. Вход с неверным паролем
	badLoginReq := &dto.LoginRequest{
		Email:    "test@mentorhub.com",
		Password: "wrongpassword",
	}
	_, err = authSvc.Login(ctx, badLoginReq)
	if err != service.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_DeactivatedUser(t *testing.T) {
	repo := newMockUserRepo()
	jwtMgr := jwtpkg.NewManager("secret_key_for_testing_12345", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repo, jwtMgr, 15*time.Minute)

	ctx := context.Background()

	pwdHash, _ := hasher.Hash("password123")
	uid := uuid.New()
	_ = repo.Create(ctx, &domain.User{
		ID:           uid,
		Email:        "inactive@mentorhub.com",
		PasswordHash: pwdHash,
		Role:         domain.RoleMentor,
		FirstName:    "Деактивирован",
		LastName:     "Тест",
		IsActive:     false,
	})

	_, err := authSvc.Login(ctx, &dto.LoginRequest{
		Email:    "inactive@mentorhub.com",
		Password: "password123",
	})
	if err != service.ErrUserInactive {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}
