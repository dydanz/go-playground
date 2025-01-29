package postgres

import (
	"go-playground/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockPointsRepository struct {
	mock.Mock
}

func (m *MockPointsRepository) Create(balance *domain.PointsBalance) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *MockPointsRepository) GetByUserID(userID string) (*domain.PointsBalance, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsBalance), args.Error(1)
}

func (m *MockPointsRepository) Update(balance *domain.PointsBalance) error {
	args := m.Called(balance)
	return args.Error(0)
}
