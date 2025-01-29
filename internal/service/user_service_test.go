package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateTx(tx *sql.Tx, user *domain.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetRandomActiveUser() (*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (m *MockCacheRepository) GetUser(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockCacheRepository) SetUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	mockUserRepo.On("GetByEmail", "test@example.com").Return(nil, nil)

	mockUserRepo.On("Create",
		mock.Anything,
		mock.MatchedBy(func(req *domain.CreateUserRequest) bool {
			return req.Email == "test@example.com" &&
				req.Name == "Test User" &&
				len(req.Password) > 0
		}),
	).Return(&domain.User{
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
	}, nil)

	req := &domain.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	user, err := service.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}
