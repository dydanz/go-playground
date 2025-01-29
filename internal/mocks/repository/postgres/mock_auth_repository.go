package postgres

import (
	"database/sql"
	"go-playground/internal/domain"

	"github.com/stretchr/testify/mock"
)

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

func (m *MockAuthRepository) GetLatestVerification(userID string) (*domain.RegistrationVerification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *MockAuthRepository) GetRandomActiveUser() (*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
