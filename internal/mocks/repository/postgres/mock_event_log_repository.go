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
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventLogRepository) Update(event *domain.EventLog) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByReferenceID(referenceID string) (*domain.EventLog, error) {
	args := m.Called(referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) GetByID(id string) (*domain.EventLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByUserID(userID string) ([]domain.EventLog, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventLog), args.Error(1)
}
