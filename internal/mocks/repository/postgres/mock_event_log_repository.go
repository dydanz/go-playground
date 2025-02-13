package postgres

import (
	"context"
	"go-playground/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockEventLogRepository struct {
	mock.Mock
}

func (m *MockEventLogRepository) Create(ctx context.Context, event *domain.EventLog) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventLogRepository) Update(ctx context.Context, event *domain.EventLog) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByReferenceID(ctx context.Context, referenceID string) (*domain.EventLog, error) {
	args := m.Called(ctx, referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) GetByID(ctx context.Context, id string) (*domain.EventLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByUserID(ctx context.Context, userID string) ([]domain.EventLog, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventLog), args.Error(1)
}
