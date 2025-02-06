package service

import (
	"context"
	"go-playground/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPointsService struct {
	mock.Mock
}

func (m *MockPointsService) GetLedger(ctx context.Context, customerID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	args := m.Called(ctx, customerID, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PointsLedger), args.Error(1)
}

func (m *MockPointsService) GetBalance(customerID, programID string) (*domain.PointsBalance, error) {
	args := m.Called(customerID, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsBalance), args.Error(1)
}

func (m *MockPointsService) EarnPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error {
	args := m.Called(ctx, customerID, programID, points, transactionID)
	return args.Error(0)
}

func (m *MockPointsService) RedeemPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error {
	args := m.Called(ctx, customerID, programID, points, transactionID)
	return args.Error(0)
}
