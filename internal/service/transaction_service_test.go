package service

import (
	"context"
	"errors"
	"testing"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
	"go-playground/internal/mocks/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_Create_Success(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	branchID := uuid.New()
	req := &domain.CreateTransactionRequest{
		MerchantID:        uuid.New(),
		CustomerID:        uuid.New(),
		TransactionType:   "purchase",
		TransactionAmount: 100,
		BranchID:          &branchID,
	}

	mockTransactionRepo.On("Create", mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.MerchantID == req.MerchantID &&
			tx.CustomerID == req.CustomerID &&
			tx.TransactionType == req.TransactionType &&
			tx.TransactionAmount == req.TransactionAmount &&
			tx.BranchID == req.BranchID
	})).Return(nil)

	mockPointsService.On("EarnPoints", mock.Anything, req.CustomerID, req.MerchantID, 100, mock.AnythingOfType("*uuid.UUID")).Return(nil)

	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		return event.EventType == "transaction_created" &&
			event.UserID == req.CustomerID.String() &&
			event.Details["merchant_id"] == req.MerchantID &&
			event.Details["transaction_type"] == req.TransactionType &&
			event.Details["transaction_amount"] == req.TransactionAmount &&
			event.Details["points_earned"] == 100 &&
			event.Details["branch_id"] == req.BranchID
	})).Return(nil)

	tx, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, req.MerchantID, tx.MerchantID)
	assert.Equal(t, req.CustomerID, tx.CustomerID)
	assert.Equal(t, req.TransactionType, tx.TransactionType)
	assert.Equal(t, req.TransactionAmount, tx.TransactionAmount)
	assert.Equal(t, req.BranchID, tx.BranchID)

	mockTransactionRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestTransactionService_Create_RepositoryError(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	branchID := uuid.New()
	req := &domain.CreateTransactionRequest{
		MerchantID:        uuid.New(),
		CustomerID:        uuid.New(),
		TransactionType:   "purchase",
		TransactionAmount: 100,
		BranchID:          &branchID,
	}

	mockTransactionRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))

	tx, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.Equal(t, "database error", err.Error())

	mockTransactionRepo.AssertExpectations(t)
	mockPointsService.AssertNotCalled(t, "EarnPoints")
	mockEventRepo.AssertNotCalled(t, "Create")
}

func TestTransactionService_GetByID_Success(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	txID := uuid.New()
	expectedTx := &domain.Transaction{
		TransactionID: txID,
		MerchantID:    uuid.New(),
		CustomerID:    uuid.New(),
	}

	mockTransactionRepo.On("GetByID", mock.Anything, txID).Return(expectedTx, nil)

	tx, err := service.GetByID(context.Background(), txID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, tx)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_GetByID_NotFound(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	txID := uuid.New()
	mockTransactionRepo.On("GetByID", mock.Anything, txID).Return(nil, errors.New("not found"))

	tx, err := service.GetByID(context.Background(), txID)

	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.Equal(t, "not found", err.Error())

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_GetByCustomerID_Success(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	customerID := uuid.New()
	expectedTxs := []*domain.Transaction{
		{TransactionID: uuid.New(), CustomerID: customerID},
		{TransactionID: uuid.New(), CustomerID: customerID},
	}

	mockTransactionRepo.On("GetByCustomerID", mock.Anything, customerID).Return(expectedTxs, nil)

	txs, err := service.GetByCustomerID(context.Background(), customerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTxs, txs)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_GetByMerchantID_Success(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	merchantID := uuid.New()
	expectedTxs := []*domain.Transaction{
		{TransactionID: uuid.New(), MerchantID: merchantID},
		{TransactionID: uuid.New(), MerchantID: merchantID},
	}

	mockTransactionRepo.On("GetByMerchantID", mock.Anything, merchantID).Return(expectedTxs, nil)

	txs, err := service.GetByMerchantID(context.Background(), merchantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTxs, txs)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateStatus_Success(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	txID := uuid.New()
	existingTx := &domain.Transaction{TransactionID: txID, Status: "pending"}
	newStatus := "completed"

	mockTransactionRepo.On("GetByID", mock.Anything, txID).Return(existingTx, nil)
	mockTransactionRepo.On("UpdateStatus", mock.Anything, txID, newStatus).Return(nil)

	err := service.UpdateStatus(context.Background(), txID, newStatus)

	assert.NoError(t, err)
	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateStatus_NotFound(t *testing.T) {
	mockTransactionRepo := new(postgres.MockTransactionRepository)
	mockPointsService := new(service.MockPointsService)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewTransactionService(mockTransactionRepo, mockPointsService, mockEventRepo)

	txID := uuid.New()
	mockTransactionRepo.On("GetByID", mock.Anything, txID).Return(nil, errors.New("not found"))

	err := service.UpdateStatus(context.Background(), txID, "completed")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	mockTransactionRepo.AssertExpectations(t)
}
