package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
)

func TestPointsService_GetBalance_Success(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	// Add GetByID method to MockEventLogRepository before creating service
	mockEventRepo.On("GetByID", mock.Anything).Return(&domain.EventLog{}, nil)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	// Add missing GetByID method to MockEventLogRepository
	mockEventRepo.On("GetByID", mock.Anything).Return(&domain.EventLog{}, nil)

	expectedBalance := &domain.PointsBalance{
		ID:          "balance123",
		UserID:      "user123",
		TotalPoints: 100,
	}

	mockPointsRepo.On("GetByUserID", "user123").Return(expectedBalance, nil)

	balance, err := service.GetBalance("user123")
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, expectedBalance.TotalPoints, balance.TotalPoints)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_GetBalance_NotFound(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	mockPointsRepo.On("GetByUserID", "nonexistent").Return(nil, nil)

	balance, err := service.GetBalance("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, balance)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_GetBalance_Error(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	mockPointsRepo.On("GetByUserID", "user123").Return(nil, errors.New("database error"))

	balance, err := service.GetBalance("user123")
	assert.Error(t, err)
	assert.Nil(t, balance)
	assert.Equal(t, "database error", err.Error())

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_UpdateBalance_CreateNew(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	mockPointsRepo.On("GetByUserID", "user123").Return(nil, nil)
	mockPointsRepo.On("Create", mock.MatchedBy(func(balance *domain.PointsBalance) bool {
		return balance.UserID == "user123" && balance.TotalPoints == 50
	})).Return(nil)

	err := service.UpdateBalance("user123", 50)
	assert.NoError(t, err)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_UpdateBalance_UpdateExisting(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	existingBalance := &domain.PointsBalance{
		ID:          "balance123",
		UserID:      "user123",
		TotalPoints: 100,
	}

	mockPointsRepo.On("GetByUserID", "user123").Return(existingBalance, nil)
	mockPointsRepo.On("Update", mock.MatchedBy(func(balance *domain.PointsBalance) bool {
		return balance.TotalPoints == 150
	})).Return(nil)

	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		return event.EventType == "balance_update" &&
			event.UserID == "user123" &&
			event.Details["points_changed"].(int) == 50 &&
			event.Details["new_balance"].(int) == 150
	})).Return(nil)

	err := service.UpdateBalance("user123", 50)
	assert.NoError(t, err)

	mockPointsRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestPointsService_UpdateBalance_InsufficientPoints(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	existingBalance := &domain.PointsBalance{
		ID:          "balance123",
		UserID:      "user123",
		TotalPoints: 50,
	}

	mockPointsRepo.On("GetByUserID", "user123").Return(existingBalance, nil)

	err := service.UpdateBalance("user123", -100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient points balance", err.Error())

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_UpdateBalance_UpdateError(t *testing.T) {
	mockPointsRepo := new(postgres.MockPointsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	existingBalance := &domain.PointsBalance{
		ID:          "balance123",
		UserID:      "user123",
		TotalPoints: 100,
	}

	mockPointsRepo.On("GetByUserID", "user123").Return(existingBalance, nil)
	mockPointsRepo.On("Update", mock.Anything).Return(errors.New("update error"))

	err := service.UpdateBalance("user123", 50)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())

	mockPointsRepo.AssertExpectations(t)
}
