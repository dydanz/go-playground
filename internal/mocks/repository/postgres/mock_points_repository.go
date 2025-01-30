package postgres

import (
	"context"
	"go-playground/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPointsRepository struct {
	mock.Mock
}

func (m *MockPointsRepository) Create(ctx context.Context, ledger *domain.PointsLedger) error {
	args := m.Called(ctx, ledger)
	return args.Error(0)
}

func (m *MockPointsRepository) GetByCustomerAndProgram(ctx context.Context, customerID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	args := m.Called(ctx, customerID, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PointsLedger), args.Error(1)
}

func (m *MockPointsRepository) GetCurrentBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error) {
	args := m.Called(ctx, customerID, programID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockPointsRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*domain.PointsLedger, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsLedger), args.Error(1)
}
