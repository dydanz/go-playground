package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
)

func TestPointsService_GetLedger_Success(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()

	expectedLedger := []*domain.PointsLedger{
		{
			LedgerID:       uuid.New(),
			CustomerID:     customerID,
			ProgramID:      programID,
			PointsEarned:   100,
			PointsRedeemed: 0,
			PointsBalance:  100,
		},
	}

	mockPointsRepo.On("GetByCustomerAndProgram", mock.Anything, customerID, programID).Return(expectedLedger, nil)

	ledger, err := service.GetLedger(context.Background(), customerID, programID)
	assert.NoError(t, err)
	assert.NotNil(t, ledger)
	assert.Equal(t, expectedLedger, ledger)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_GetBalance_Success(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()
	expectedBalance := 100

	mockPointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(expectedBalance, nil)

	balance, err := service.GetBalance(context.Background(), customerID, programID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_GetBalance_Error(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()

	mockPointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(0, errors.New("database error"))

	balance, err := service.GetBalance(context.Background(), customerID, programID)
	assert.Error(t, err)
	assert.Equal(t, 0, balance)
	assert.Equal(t, "database error", err.Error())

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_EarnPoints_Success(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	mockPointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(100, nil)
	mockPointsRepo.On("Create", mock.Anything, mock.MatchedBy(func(ledger *domain.PointsLedger) bool {
		return ledger.CustomerID == customerID &&
			ledger.ProgramID == programID &&
			ledger.PointsEarned == 50 &&
			ledger.PointsBalance == 150
	})).Return(nil)

	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		return event.EventType == "points_earned" &&
			event.UserID == customerID.String() &&
			event.Details["points_earned"].(int) == 50 &&
			event.Details["new_balance"].(int) == 150
	})).Return(nil)

	err := service.EarnPoints(context.Background(), customerID, programID, 50, &transactionID)
	assert.NoError(t, err)

	mockPointsRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestPointsService_EarnPoints_InvalidPoints(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()

	err := service.EarnPoints(context.Background(), customerID, programID, 0, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPoints, err)

	mockPointsRepo.AssertNotCalled(t, "Create")
}

func TestPointsService_RedeemPoints_Success(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	mockPointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(100, nil)
	mockPointsRepo.On("Create", mock.Anything, mock.MatchedBy(func(ledger *domain.PointsLedger) bool {
		return ledger.CustomerID == customerID &&
			ledger.ProgramID == programID &&
			ledger.PointsRedeemed == 50 &&
			ledger.PointsBalance == 50
	})).Return(nil)

	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		return event.EventType == "points_redeemed" &&
			event.UserID == customerID.String() &&
			event.Details["points_redeemed"].(int) == 50 &&
			event.Details["new_balance"].(int) == 50
	})).Return(nil)

	err := service.RedeemPoints(context.Background(), customerID, programID, 50, &transactionID)
	assert.NoError(t, err)

	mockPointsRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestPointsService_RedeemPoints_InsufficientPoints(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	customerID := uuid.New()
	programID := uuid.New()

	mockPointsRepo.On("GetCurrentBalance", mock.Anything, customerID, programID).Return(50, nil)

	err := service.RedeemPoints(context.Background(), customerID, programID, 100, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientPoints, err)

	mockPointsRepo.AssertNotCalled(t, "Create")
}
