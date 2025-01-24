package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
)

// Mock for RedemptionRepository
type MockRedemptionRepository struct {
	mock.Mock
}

func (m *MockRedemptionRepository) Create(redemption *domain.Redemption) error {
	args := m.Called(redemption)
	return args.Error(0)
}

func (m *MockRedemptionRepository) GetByID(id string) (*domain.Redemption, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) GetByUserID(userID string) ([]domain.Redemption, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Redemption), args.Error(1)
}

func (m *MockRedemptionRepository) Update(redemption *domain.Redemption) error {
	args := m.Called(redemption)
	return args.Error(0)
}

func (m *MockRewardsRepository) GetByID(id string) (*domain.Reward, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) GetAll(activeOnly bool) ([]domain.Reward, error) {
	args := m.Called(activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Reward), args.Error(1)
}

func (m *MockRewardsRepository) Create(reward *domain.Reward) error {
	args := m.Called(reward)
	return args.Error(0)
}

func (m *MockRewardsRepository) Update(reward *domain.Reward) error {
	args := m.Called(reward)
	return args.Error(0)
}

// Mock for PointsService
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

// Test cases
func (m *MockRewardsRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRedemptionService_Create_Success(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	reward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
		IsActive:       true,
	}

	balance := &domain.PointsBalance{
		UserID:      "user123",
		TotalPoints: 150,
	}

	redemption := &domain.Redemption{
		UserID:   "user123",
		RewardID: "reward123",
		Status:   "pending",
	}

	mockRewardsRepo.On("GetByID", "reward123").Return(reward, nil)
	mockPointsService.On("GetBalance", "user123").Return(balance, nil)
	mockRedemptionRepo.On("Create", redemption).Return(nil)
	mockPointsService.On("UpdateBalance", "user123", -100).Return(nil)
	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		return event.EventType == "reward_redeemed" &&
			event.UserID == "user123" &&
			event.Details["reward_id"] == "reward123" &&
			event.Details["points_used"] == 100
	})).Return(nil)

	err := service.Create(redemption)

	assert.NoError(t, err)
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
	mockRedemptionRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestRedemptionService_Create_InactiveReward(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	reward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
		IsActive:       false,
	}

	redemption := &domain.Redemption{
		UserID:   "user123",
		RewardID: "reward123",
		Status:   "pending",
	}

	mockRewardsRepo.On("GetByID", "reward123").Return(reward, nil)

	err := service.Create(redemption)

	assert.Error(t, err)
	assert.Equal(t, "reward is not available", err.Error())
	mockRewardsRepo.AssertExpectations(t)
}

func TestRedemptionService_Create_InsufficientPoints(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	reward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
		IsActive:       true,
	}

	balance := &domain.PointsBalance{
		UserID:      "user123",
		TotalPoints: 50,
	}

	redemption := &domain.Redemption{
		UserID:   "user123",
		RewardID: "reward123",
		Status:   "pending",
	}

	mockRewardsRepo.On("GetByID", "reward123").Return(reward, nil)
	mockPointsService.On("GetBalance", "user123").Return(balance, nil)

	err := service.Create(redemption)

	assert.Error(t, err)
	assert.Equal(t, "insufficient points", err.Error())
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
}

func TestRedemptionService_GetByID_Success(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	expectedRedemption := &domain.Redemption{
		ID:       "redemption123",
		UserID:   "user123",
		RewardID: "reward123",
		Status:   "completed",
	}

	mockRedemptionRepo.On("GetByID", "redemption123").Return(expectedRedemption, nil)

	redemption, err := service.GetByID("redemption123")

	assert.NoError(t, err)
	assert.Equal(t, expectedRedemption, redemption)
	mockRedemptionRepo.AssertExpectations(t)
}

func TestRedemptionService_GetByUserID_Success(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	expectedRedemptions := []domain.Redemption{
		{
			ID:       "redemption123",
			UserID:   "user123",
			RewardID: "reward123",
			Status:   "completed",
		},
		{
			ID:       "redemption124",
			UserID:   "user123",
			RewardID: "reward124",
			Status:   "pending",
		},
	}

	mockRedemptionRepo.On("GetByUserID", "user123").Return(expectedRedemptions, nil)

	redemptions, err := service.GetByUserID("user123")

	assert.NoError(t, err)
	assert.Equal(t, expectedRedemptions, redemptions)
	mockRedemptionRepo.AssertExpectations(t)
}

func TestRedemptionService_UpdateStatus_Success(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	redemption := &domain.Redemption{
		ID:       "redemption123",
		UserID:   "user123",
		RewardID: "reward123",
		Status:   "pending",
	}

	reward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
	}

	mockRedemptionRepo.On("GetByID", "redemption123").Return(redemption, nil)
	mockRewardsRepo.On("GetByID", "reward123").Return(reward, nil)
	mockPointsService.On("UpdateBalance", "user123", 100).Return(nil)
	mockRedemptionRepo.On("Update", mock.MatchedBy(func(r *domain.Redemption) bool {
		return r.ID == "redemption123" && r.Status == "canceled"
	})).Return(nil)

	err := service.UpdateStatus("redemption123", "canceled")

	assert.NoError(t, err)
	mockRedemptionRepo.AssertExpectations(t)
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
}

func TestRedemptionService_UpdateStatus_NotFound(t *testing.T) {
	mockRedemptionRepo := new(MockRedemptionRepository)
	mockRewardsRepo := new(MockRewardsRepository)
	mockEventRepo := new(MockEventLogRepository)
	mockPointsService := new(MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	mockRedemptionRepo.On("GetByID", "nonexistent").Return(nil, errors.New("not found"))

	err := service.UpdateStatus("nonexistent", "canceled")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	mockRedemptionRepo.AssertExpectations(t)
}
