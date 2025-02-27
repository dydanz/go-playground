package redis

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"go-playground/server/repository/redis"
)

// MockSessionRepository is a mock implementation of the SessionRepository interface
type MockSessionRepository struct {
	mock.Mock
}

// Implement the methods of the SessionRepository interface
func (m *MockSessionRepository) StoreSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSession(ctx context.Context, userID string) (*redis.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*redis.Session), args.Error(1)
}

func (m *MockSessionRepository) DeleteSession(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) RefreshSession(ctx context.Context, userID, newToken string, expiration time.Duration) error {
	args := m.Called(ctx, userID, newToken, expiration)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteAllSession(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Add other methods as needed based on the SessionRepository interface
