package service

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"go-playground/internal/domain"
	"go-playground/internal/repository/redis"
)

// Mock for SQL transaction
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateVerification(v *domain.RegistrationVerification) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockAuthRepository) GetVerification(userID, otp string) (*domain.RegistrationVerification, error) {
	args := m.Called(userID, otp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *MockAuthRepository) MarkVerificationUsedTx(tx *sql.Tx, verificationID string) error {
	args := m.Called(tx, verificationID)
	return args.Error(0)
}

func (m *MockAuthRepository) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockAuthRepository) CreateToken(token *domain.AuthToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthRepository) InvalidateToken(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthRepository) UpdateLoginAttempts(email string, increment bool) (*domain.LoginAttempt, error) {
	args := m.Called(email, increment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginAttempt), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	req := &domain.RegistrationRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Phone:    "1234567890",
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(nil, nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(req *domain.CreateUserRequest) bool {
		return req.Email == "test@example.com"
	})).Return(&domain.User{
		ID:    "user123",
		Email: "test@example.com",
		Name:  "Test User",
	}, nil)

	mockAuthRepo.On("CreateVerification", mock.MatchedBy(func(v *domain.RegistrationVerification) bool {
		return v.UserID == "user123" && len(v.OTP) == 6
	})).Return(nil)

	user, err := service.Register(req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	existingUser := &domain.User{
		ID:    "existing123",
		Email: "test@example.com",
	}

	mockUserRepo.On("GetByEmail", "test@example.com").Return(existingUser, nil)

	req := &domain.RegistrationRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := service.Register(req)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already exists", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       "user123",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Status:   domain.UserStatusActive,
	}

	mockAuthRepo.On("UpdateLoginAttempts", "test@example.com", true).Return(&domain.LoginAttempt{
		Email:        "test@example.com",
		AttemptCount: 0,
		LockedUntil:  time.Time{},
	}, nil)

	mockUserRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	mockAuthRepo.On("UpdateLoginAttempts", "test@example.com", false).Return(&domain.LoginAttempt{}, nil)

	mockAuthRepo.On("CreateToken", mock.MatchedBy(func(token *domain.AuthToken) bool {
		return token.UserID == "user123"
	})).Return(nil)

	mockSessionRepo.On("StoreSession", mock.Anything, "user123", mock.Anything, mock.Anything).Return(nil)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := service.Login(req)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "user123", token.UserID)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_Login_AccountLocked(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	lockedUntil := time.Now().Add(15 * time.Minute)
	mockAuthRepo.On("UpdateLoginAttempts", "test@example.com", true).Return(&domain.LoginAttempt{
		Email:        "test@example.com",
		AttemptCount: 5,
		LastAttempt:  time.Now(),
		LockedUntil:  lockedUntil,
	}, nil)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := service.Login(req)
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "account temporarily locked")

	mockAuthRepo.AssertExpectations(t)
}

/*
	func TestAuthService_VerifyRegistration_Success(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockAuthRepo := new(MockAuthRepository)
		mockSessionRepo := new(redis.MockSessionRepository)
		service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

		user := &domain.User{
			ID:     "user123",
			Email:  "test@example.com",
			Status: domain.UserStatusPending,
		}

		verification := &domain.RegistrationVerification{
			ID:        "ver123",
			UserID:    "user123",
			OTP:       "123456",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		tx := &sql.Tx{}

		mockUserRepo.On("GetByEmail", "test@example.com").Return(user, nil)
		mockAuthRepo.On("GetVerification", "user123", "123456").Return(verification, nil)

		mockAuthRepo.On("BeginTx").Return(tx, nil)
		mockAuthRepo.On("MarkVerificationUsedTx", tx, "ver123").Return(nil)
		mockAuthRepo.On("Commit").Return(tx, nil)
		mockUserRepo.On("UpdateTx", tx, mock.MatchedBy(func(u *domain.User) bool {
			return u.Status == domain.UserStatusActive
		})).Return(nil)

		mockSessionRepo.On("StoreSession", mock.Anything, "user123", mock.Anything, mock.Anything).Return(nil)

		req := &domain.VerificationRequest{
			Email: "test@example.com",
			OTP:   "123456",
		}

		err := service.VerifyRegistration(req)
		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
		mockAuthRepo.AssertExpectations(t)
		mockSessionRepo.AssertExpectations(t)

}
*/
func TestAuthService_Logout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	mockAuthRepo.On("InvalidateToken", "user123").Return(nil)

	mockSessionRepo.On("DeleteSession", mock.Anything, "user123").Return(nil)

	err := service.Logout("user123", "token123")
	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(redis.MockSessionRepository)
	service := NewAuthService(mockUserRepo, mockAuthRepo, mockSessionRepo)

	mockAuthRepo.On("UpdateLoginAttempts", "test@example.com", true).Return(&domain.LoginAttempt{
		Email:        "test@example.com",
		AttemptCount: 0,
		LastAttempt:  time.Time{},
		LockedUntil:  time.Time{},
	}, nil)

	mockUserRepo.On("GetByEmail", "test@example.com").Return(nil, nil)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	mockSessionRepo.On("StoreSession", mock.Anything, "user123", mock.Anything, mock.Anything).Return(nil)

	token, err := service.Login(req)
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, "invalid credentials", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}
