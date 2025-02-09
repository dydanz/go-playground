package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPointsService(t *testing.T) {
	pointsRepo := new(postgres.MockPointsRepository)
	eventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(pointsRepo, eventRepo)
	ctx := context.Background()

	customerID := uuid.New()
	programID := uuid.New()
	txID := uuid.New()

	t.Run("GetLedger", func(t *testing.T) {
		ledger := &domain.PointsLedger{
			LedgerID:            uuid.New(),
			MerchantCustomersID: customerID,
			ProgramID:           programID,
			PointsEarned:        100,
			TransactionID:       &txID,
			CreatedAt:           time.Now(),
		}

		pointsRepo.On("GetByCustomerAndProgram", mock.Anything, customerID, programID).Return([]*domain.PointsLedger{ledger}, nil)

		result, err := service.GetLedger(ctx, customerID, programID)
		assert.NoError(t, err)
		assert.Equal(t, []*domain.PointsLedger{ledger}, result)
	})

	t.Run("GetBalance", func(t *testing.T) {
		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(100, nil)

		result, err := service.GetBalance(ctx, customerID, programID)
		assert.NoError(t, err)
		assert.Equal(t, 100, result)
	})

	t.Run("GetBalance - Error", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)
		expectedErr := errors.New("database error")
		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(0, expectedErr)

		result, err := service.GetBalance(ctx, customerID, programID)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, 0, result)
	})

	t.Run("EarnPoints", func(t *testing.T) {
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		err := service.EarnPoints(ctx, customerID, programID, 100, &txID)
		assert.NoError(t, err)
	})

	t.Run("EarnPoints - Invalid Points", func(t *testing.T) {
		err := service.EarnPoints(ctx, customerID, programID, 0, &txID)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPoints, err)
	})

	t.Run("RedeemPoints", func(t *testing.T) {
		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(200, nil)
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		err := service.RedeemPoints(ctx, customerID, programID, 100, &txID)
		assert.NoError(t, err)
	})

	t.Run("RedeemPoints - Insufficient Balance", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)

		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(50, nil)

		err := service.RedeemPoints(ctx, customerID, programID, 100, &txID)
		assert.Equal(t, ErrInsufficientPoints, err)
		pointsRepo.AssertExpectations(t)
	})
}
