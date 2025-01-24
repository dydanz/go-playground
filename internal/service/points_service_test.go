package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
)

// Mock for PointsRepository
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

// Mock for EventLogRepository
type MockEventLogRepository struct {
	mock.Mock
}

func (m *MockEventLogRepository) Create(event *domain.EventLog) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventLogRepository) Update(event *domain.EventLog) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByReferenceID(referenceID string) (*domain.EventLog, error) {
	args := m.Called(referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) GetByID(id string) (*domain.EventLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *MockEventLogRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventLogRepository) GetByUserID(userID string) ([]domain.EventLog, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventLog), args.Error(1)
}

func TestPointsService_GetBalance_Success(t *testing.T) {
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
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
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	mockPointsRepo.On("GetByUserID", "nonexistent").Return(nil, nil)

	balance, err := service.GetBalance("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, balance)

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_GetBalance_Error(t *testing.T) {
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
	service := NewPointsService(mockPointsRepo, mockEventRepo)

	mockPointsRepo.On("GetByUserID", "user123").Return(nil, errors.New("database error"))

	balance, err := service.GetBalance("user123")
	assert.Error(t, err)
	assert.Nil(t, balance)
	assert.Equal(t, "database error", err.Error())

	mockPointsRepo.AssertExpectations(t)
}

func TestPointsService_UpdateBalance_CreateNew(t *testing.T) {
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
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
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
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
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
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
	mockPointsRepo := new(MockPointsRepository)
	mockEventRepo := new(MockEventLogRepository)
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
