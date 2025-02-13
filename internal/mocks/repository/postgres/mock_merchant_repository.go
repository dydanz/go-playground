package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
)

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) Create(ctx context.Context, merchant *domain.Merchant) (*domain.Merchant, error) {
	args := m.Called(ctx, merchant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) GetAll(ctx context.Context) ([]*domain.Merchant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) Update(ctx context.Context, merchant *domain.Merchant) error {
	args := m.Called(ctx, merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMerchantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Merchant, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Merchant), args.Error(1)
}
