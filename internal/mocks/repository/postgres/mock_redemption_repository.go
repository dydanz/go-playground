package postgres

import (
	"go-playground/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRedemptionRepository struct {
	mock.Mock
}

func (m *MockRedemptionRepository) Create(redemption *domain.Redemption) error {
	args := m.Called(redemption)
	return args.Error(0)
}

func (m *MockRedemptionRepository) GetByID(id string) (*domain.Redemption, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) GetByUserID(userID string) ([]domain.Redemption, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) Update(redemption *domain.Redemption) error {
	args := m.Called(redemption)
	return args.Error(0)
}