package postgres

import (
	"context"
	"go-playground/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRedemptionRepository struct {
	mock.Mock
}

func (m *MockRedemptionRepository) Create(ctx context.Context, redemption *domain.Redemption) ([]*domain.Redemption, error) {
	args := m.Called(ctx, redemption)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Redemption, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Redemption, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) Update(ctx context.Context, redemption *domain.Redemption) error {
	args := m.Called(ctx, redemption)
	return args.Error(0)
}

func (m *MockRedemptionRepository) GetByRewardID(ctx context.Context, rewardID uuid.UUID) ([]*domain.Redemption, error) {
	args := m.Called(ctx, rewardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) GetByMerchantCustomerID(ctx context.Context, merchantCustomersID uuid.UUID) ([]*domain.Redemption, error) {
	args := m.Called(ctx, merchantCustomersID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Redemption), args.Error(1)
}
