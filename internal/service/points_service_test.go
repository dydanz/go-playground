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
			TransactionID:       txID,
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
		assert.Equal(t, &domain.PointsBalance{
			CustomerID: customerID.String(),
			ProgramID:  programID.String(),
			Balance:    100,
		}, result)
	})

	t.Run("GetBalance - Error", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)
		expectedErr := errors.New("database error")
		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(0, expectedErr)

		result, err := service.GetBalance(ctx, customerID, programID)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("EarnPoints", func(t *testing.T) {
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		req := &domain.PointsTransaction{
			CustomerID:    customerID.String(),
			ProgramID:     programID.String(),
			Points:        100,
			Type:          "earn",
			TransactionID: txID.String(),
		}

		result, err := service.EarnPoints(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Points, result.Points)
		assert.Equal(t, req.CustomerID, result.CustomerID)
		assert.Equal(t, req.ProgramID, result.ProgramID)
		assert.Equal(t, "earn", result.Type)
	})

	t.Run("RedeemPoints", func(t *testing.T) {
		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(200, nil)
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		req := &domain.PointsTransaction{
			CustomerID:    customerID.String(),
			ProgramID:     programID.String(),
			Points:        100,
			Type:          "redeem",
			TransactionID: txID.String(),
		}

		result, err := service.RedeemPoints(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Points, result.Points)
		assert.Equal(t, req.CustomerID, result.CustomerID)
		assert.Equal(t, req.ProgramID, result.ProgramID)
		assert.Equal(t, "redeem", result.Type)
	})

	t.Run("RedeemPoints - Insufficient Balance", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)

		pointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(50, nil)

		req := &domain.PointsTransaction{
			CustomerID:    customerID.String(),
			ProgramID:     programID.String(),
			Points:        100,
			Type:          "redeem",
			TransactionID: txID.String(),
		}

		result, err := service.RedeemPoints(ctx, req)
		assert.Equal(t, ErrInsufficientPoints, err)
		assert.Nil(t, result)
		pointsRepo.AssertExpectations(t)
	})
}
