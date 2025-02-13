package postgres

import (
	"context"
	"go-playground/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRewardsRepository struct {
	mock.Mock
}

func (m *MockRewardsRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reward, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) GetAll(ctx context.Context, activeOnly bool) ([]domain.Reward, error) {
	args := m.Called(ctx, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) Create(ctx context.Context, reward *domain.Reward) (*domain.Reward, error) {
	args := m.Called(ctx, reward)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) Update(ctx context.Context, reward *domain.Reward) (*domain.Reward, error) {
	args := m.Called(ctx, reward)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRewardsRepository) GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*domain.Reward, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reward), args.Error(1)
}
