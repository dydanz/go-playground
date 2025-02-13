package postgres

import (
	"context"
	"go-playground/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	args := m.Called(ctx, transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByCustomerID(ctx context.Context, merchantCustomersID uuid.UUID) ([]*domain.Transaction, error) {
	args := m.Called(ctx, merchantCustomersID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByCustomerIDWithPagination(ctx context.Context, merchantCustomersID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	args := m.Called(ctx, merchantCustomersID, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Transaction, error) {
	args := m.Called(ctx, merchantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	args := m.Called(ctx, transactionID, status)
	return args.Error(0)
}
