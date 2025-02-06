package service

import (
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

	customerID := uuid.New().String()
	programID := uuid.New().String()
	txID := uuid.New()

	t.Run("GetLedger", func(t *testing.T) {
		ledger := &domain.PointsLedger{
			LedgerID:      uuid.New(),
			CustomerID:    uuid.MustParse(customerID),
			ProgramID:     uuid.MustParse(programID),
			PointsEarned:  100,
			TransactionID: &txID,
			CreatedAt:     time.Now(),
		}

		pointsRepo.On("GetByCustomerAndProgram", mock.Anything, uuid.MustParse(customerID), uuid.MustParse(programID)).Return([]*domain.PointsLedger{ledger}, nil)

		result, err := service.GetLedger(customerID, programID)
		assert.NoError(t, err)
		assert.Equal(t, ledger, result)
	})

	t.Run("GetBalance", func(t *testing.T) {
		pointsRepo.On("GetCurrentBalance", mock.Anything, uuid.MustParse(customerID), uuid.MustParse(programID)).Return(int64(100), nil)

		result, err := service.GetBalance(customerID, programID)
		assert.NoError(t, err)
		assert.Equal(t, &domain.PointsBalance{
			CustomerID: customerID,
			ProgramID:  programID,
			Balance:    100,
		}, result)
	})

	t.Run("GetBalance - Error", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)
		expectedErr := errors.New("database error")
		pointsRepo.On("GetCurrentBalance", mock.Anything, uuid.MustParse(customerID), uuid.MustParse(programID)).Return(int64(0), expectedErr)

		result, err := service.GetBalance(customerID, programID)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("EarnPoints", func(t *testing.T) {
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		req := &domain.EarnPointsRequest{
			CustomerID: customerID,
			ProgramID:  programID,
			Points:     100,
		}

		result, err := service.EarnPoints(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.TransactionID)
	})

	t.Run("EarnPoints - Invalid Points", func(t *testing.T) {
		req := &domain.EarnPointsRequest{
			CustomerID: customerID,
			ProgramID:  programID,
			Points:     0,
		}

		result, err := service.EarnPoints(req)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPoints, err)
		assert.Nil(t, result)
	})

	t.Run("RedeemPoints", func(t *testing.T) {
		pointsRepo.On("GetCurrentBalance", mock.Anything, uuid.MustParse(customerID), uuid.MustParse(programID)).Return(int64(200), nil)
		pointsRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PointsLedger")).Return(nil)

		req := &domain.RedeemPointsRequest{
			CustomerID: customerID,
			ProgramID:  programID,
			Points:     100,
		}

		result, err := service.RedeemPoints(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.TransactionID)
	})

	t.Run("RedeemPoints - Insufficient Balance", func(t *testing.T) {
		pointsRepo := new(postgres.MockPointsRepository)
		service := NewPointsService(pointsRepo, eventRepo)

		pointsRepo.On("GetCurrentBalance", mock.Anything, uuid.MustParse(customerID), uuid.MustParse(programID)).Return(int64(50), nil)

		req := &domain.RedeemPointsRequest{
			CustomerID: customerID,
			ProgramID:  programID,
			Points:     100,
		}

		result, err := service.RedeemPoints(req)
		assert.Equal(t, ErrInsufficientPoints, err)
		assert.Nil(t, result)
		pointsRepo.AssertExpectations(t)
	})
}
