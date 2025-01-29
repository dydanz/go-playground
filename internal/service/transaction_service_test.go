package service

import (
	"errors"
	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(id string) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(userID string) ([]domain.Transaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func TestTransactionService_Create(t *testing.T) {
	tests := []struct {
		name           string
		tx             *domain.Transaction
		setupMocks     func(*MockTransactionRepository, *PointsService, *postgres.MockEventLogRepository)
		expectedError  error
		transactionErr error
		pointsErr      error
		eventErr       error
	}{
		{
			name: "successful earn transaction",
			tx: &domain.Transaction{
				ID:              "tx1",
				UserID:          "user1",
				TransactionType: "earn",
				Points:          100,
				Description:     "Test earn",
				Status:          "completed",
			},
			setupMocks: func(tr *MockTransactionRepository, ps *PointsService, er *postgres.MockEventLogRepository) {
				mockPointsRepo := ps.pointsRepo.(*postgres.MockPointsRepository)
				mockPointsRepo.On("GetByUserID", "user1").Return(&domain.PointsBalance{UserID: "user1", TotalPoints: 200}, nil)
				mockPointsRepo.On("Update", mock.AnythingOfType("*domain.PointsBalance")).Return(nil)
				tr.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)
				er.On("Create", mock.AnythingOfType("*domain.EventLog")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "successful redeem transaction",
			tx: &domain.Transaction{
				ID:              "tx2",
				UserID:          "user1",
				TransactionType: "redeem",
				Points:          50,
				Description:     "Test redeem",
				Status:          "completed",
			},
			setupMocks: func(tr *MockTransactionRepository, ps *PointsService, er *postgres.MockEventLogRepository) {
				mockPointsRepo := ps.pointsRepo.(*postgres.MockPointsRepository)
				mockPointsRepo.On("GetByUserID", "user1").Return(&domain.PointsBalance{UserID: "user1", TotalPoints: 200}, nil)
				mockPointsRepo.On("Update", mock.AnythingOfType("*domain.PointsBalance")).Return(nil)
				tr.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)
				er.On("Create", mock.AnythingOfType("*domain.EventLog")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "transaction creation fails",
			tx: &domain.Transaction{
				ID:              "tx3",
				UserID:          "user1",
				TransactionType: "earn",
				Points:          100,
			},
			setupMocks: func(tr *MockTransactionRepository, ps *PointsService, er *postgres.MockEventLogRepository) {
				tr.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockTxRepo := new(MockTransactionRepository)
			mockEventRepo := new(postgres.MockEventLogRepository)
			mockPointsRepo := new(postgres.MockPointsRepository)
			pointsService := NewPointsService(mockPointsRepo, mockEventRepo)

			// Setup mock expectations
			tt.setupMocks(mockTxRepo, pointsService, mockEventRepo)

			// Create service
			service := NewTransactionService(mockTxRepo, pointsService, mockEventRepo)

			// Execute test
			err := service.Create(tt.tx)

			// Assert results
			assert.Equal(t, tt.expectedError, err)
			mockTxRepo.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
		})
	}
}

func TestTransactionService_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTx    *domain.Transaction
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   "tx1",
			expectedTx: &domain.Transaction{
				ID:              "tx1",
				UserID:          "user1",
				TransactionType: "earn",
				Points:          100,
			},
			expectedError: nil,
		},
		{
			name:          "transaction not found",
			id:            "nonexistent",
			expectedTx:    nil,
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockTxRepo := new(MockTransactionRepository)
			mockTxRepo.On("GetByID", tt.id).Return(tt.expectedTx, tt.expectedError)

			// Create service
			service := NewTransactionService(mockTxRepo, nil, nil)

			// Execute test
			tx, err := service.GetByID(tt.id)

			// Assert results
			assert.Equal(t, tt.expectedTx, tx)
			assert.Equal(t, tt.expectedError, err)
			mockTxRepo.AssertExpectations(t)
		})
	}
}

func TestTransactionService_GetByUserID(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		expectedTxs   []domain.Transaction
		expectedError error
	}{
		{
			name:   "successful retrieval",
			userID: "user1",
			expectedTxs: []domain.Transaction{
				{
					ID:              "tx1",
					UserID:          "user1",
					TransactionType: "earn",
					Points:          100,
				},
				{
					ID:              "tx2",
					UserID:          "user1",
					TransactionType: "redeem",
					Points:          50,
				},
			},
			expectedError: nil,
		},
		{
			name:          "no transactions found",
			userID:        "user2",
			expectedTxs:   []domain.Transaction{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockTxRepo := new(MockTransactionRepository)
			mockTxRepo.On("GetByUserID", tt.userID).Return(tt.expectedTxs, tt.expectedError)

			// Create service
			service := NewTransactionService(mockTxRepo, nil, nil)

			// Execute test
			txs, err := service.GetByUserID(tt.userID)

			// Assert results
			assert.Equal(t, tt.expectedTxs, txs)
			assert.Equal(t, tt.expectedError, err)
			mockTxRepo.AssertExpectations(t)
		})
	}
}

func TestTransactionService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		newStatus     string
		setupMocks    func(*MockTransactionRepository)
		expectedError error
	}{
		{
			name:      "successful status update",
			id:        "tx1",
			newStatus: "completed",
			setupMocks: func(tr *MockTransactionRepository) {
				tr.On("GetByID", "tx1").Return(&domain.Transaction{
					ID:     "tx1",
					Status: "pending",
				}, nil)
				tr.On("Update", mock.AnythingOfType("*domain.Transaction")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "transaction not found",
			id:        "nonexistent",
			newStatus: "completed",
			setupMocks: func(tr *MockTransactionRepository) {
				tr.On("GetByID", "nonexistent").Return(nil, errors.New("not found"))
			},
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockTxRepo := new(MockTransactionRepository)
			tt.setupMocks(mockTxRepo)

			// Create service
			service := NewTransactionService(mockTxRepo, nil, nil)

			// Execute test
			err := service.UpdateStatus(tt.id, tt.newStatus)

			// Assert results
			assert.Equal(t, tt.expectedError, err)
			mockTxRepo.AssertExpectations(t)
		})
	}
}
