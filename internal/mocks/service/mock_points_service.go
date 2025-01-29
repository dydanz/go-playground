package service

import (
	"go-playground/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockPointsService struct {
	mock.Mock
}

func (m *MockPointsService) GetBalance(userID string) (*domain.PointsBalance, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsBalance), args.Error(1)
}

func (m *MockPointsService) UpdateBalance(userID string, points int) error {
	args := m.Called(userID, points)
	return args.Error(0)
}
