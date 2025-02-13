package postgres

import (
	"context"
	"database/sql"
	"go-playground/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateVerification(ctx context.Context, v *domain.RegistrationVerification) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockAuthRepository) GetVerification(ctx context.Context, userID, otp string) (*domain.RegistrationVerification, error) {
	args := m.Called(ctx, userID, otp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *MockAuthRepository) MarkVerificationUsedTx(ctx context.Context, tx *sql.Tx, verificationID string) error {
	args := m.Called(ctx, tx, verificationID)
	return args.Error(0)
}

func (m *MockAuthRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockAuthRepository) CreateToken(ctx context.Context, token *domain.AuthToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepository) InvalidateToken(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) UpdateLoginAttempts(ctx context.Context, email string, increment bool) (*domain.LoginAttempt, error) {
	args := m.Called(ctx, email, increment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginAttempt), args.Error(1)
}

func (m *MockAuthRepository) GetLatestVerification(ctx context.Context, userID string) (*domain.RegistrationVerification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *MockAuthRepository) GetUserVerificationStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) CleanupExpiredVerifications(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) CleanupExpiredAttempts(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) GetTokenByHash(ctx context.Context, hash string) (*domain.AuthToken, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthToken), args.Error(1)
}
