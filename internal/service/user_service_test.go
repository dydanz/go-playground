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

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateTx(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetRandomActiveUser(ctx context.Context) (*domain.User, error) {
	args := m.Called(ctx)
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
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	mockUserRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, nil)

	mockUserRepo.On("Create",
		ctx,
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

	user, err := service.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	expectedUser := &domain.User{
		ID:     "123",
		Email:  "test@example.com",
		Name:   "Test User",
		Status: domain.UserStatusActive,
	}

	// Test cache hit
	mockCacheRepo.On("GetUser", "123").Return(expectedUser, nil)
	user, err := service.GetByID(ctx, "123")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Test cache miss, database hit
	mockCacheRepo.On("GetUser", "456").Return(nil, domain.ResourceNotFoundError{})
	mockUserRepo.On("GetByID", ctx, "456").Return(expectedUser, nil)
	mockCacheRepo.On("SetUser", expectedUser).Return(nil)

	user, err = service.GetByID(ctx, "456")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_GetAll(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	expectedUsers := []*domain.User{
		{
			ID:     "123",
			Email:  "test1@example.com",
			Name:   "Test User 1",
			Status: domain.UserStatusActive,
		},
		{
			ID:     "456",
			Email:  "test2@example.com",
			Name:   "Test User 2",
			Status: domain.UserStatusActive,
		},
	}

	mockUserRepo.On("GetAll", ctx).Return(expectedUsers, nil)

	users, err := service.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].ID, users[0].ID)
	assert.Equal(t, expectedUsers[1].ID, users[1].ID)
	assert.Empty(t, users[0].Password) // Ensure password is cleared

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	existingUser := &domain.User{
		ID:     "123",
		Email:  "test@example.com",
		Name:   "Old Name",
		Phone:  "1234567890",
		Status: domain.UserStatusActive,
	}

	updateReq := &domain.UpdateUserRequest{
		Name:  "New Name",
		Phone: "0987654321",
	}

	mockUserRepo.On("GetByID", ctx, "123").Return(existingUser, nil)
	mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == "123" && u.Name == "New Name" && u.Phone == "0987654321"
	})).Return(nil)

	updatedUser, err := service.Update(ctx, "123", updateReq)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "0987654321", updatedUser.Phone)
	assert.Empty(t, updatedUser.Password) // Ensure password is cleared

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockCacheRepo := new(MockCacheRepository)
	service := NewUserService(mockUserRepo, mockCacheRepo)

	existingUser := &domain.User{
		ID:     "123",
		Email:  "test@example.com",
		Status: domain.UserStatusActive,
	}

	mockUserRepo.On("GetByID", ctx, "123").Return(existingUser, nil)
	mockUserRepo.On("Delete", ctx, "123").Return(nil)

	err := service.Delete(ctx, "123")
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}
