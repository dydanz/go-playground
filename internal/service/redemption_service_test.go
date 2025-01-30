package service

import (
	"context"
	"errors"
	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
	"go-playground/internal/mocks/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedemptionService_Create_Success(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	userID := uuid.New().String()
	programID := uuid.New().String()
	rewardID := "reward123"

	reward := &domain.Reward{
		ID:             rewardID,
		PointsRequired: 100,
		IsActive:       true,
	}

	redemption := &domain.Redemption{
		UserID:    userID,
		RewardID:  rewardID,
		ProgramID: programID,
		Status:    "pending",
	}

	mockRewardsRepo.On("GetByID", rewardID).Return(reward, nil)
	mockPointsService.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(150, nil)
	mockRedemptionRepo.On("Create", redemption).Return(nil)
	mockPointsService.On("RedeemPoints", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), 100, mock.AnythingOfType("*uuid.UUID")).Return(nil)
	mockEventRepo.On("Create", mock.MatchedBy(func(event *domain.EventLog) bool {
		pointsUsed, ok := event.Details["points_used"].(int)
		return event.EventType == "reward_redeemed" &&
			event.UserID == userID &&
			event.Details["reward_id"] == rewardID &&
			ok && pointsUsed == 100 &&
			event.Details["program_id"] == programID
	})).Return(nil)

	err := service.Create(context.Background(), redemption)

	assert.NoError(t, err)
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
	mockRedemptionRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestRedemptionService_Create_InactiveReward(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	userID := uuid.New().String()
	programID := uuid.New().String()
	rewardID := "reward123"

	reward := &domain.Reward{
		ID:             rewardID,
		PointsRequired: 100,
		IsActive:       false,
	}

	redemption := &domain.Redemption{
		UserID:    userID,
		RewardID:  rewardID,
		ProgramID: programID,
		Status:    "pending",
	}

	mockRewardsRepo.On("GetByID", rewardID).Return(reward, nil)

	err := service.Create(context.Background(), redemption)

	assert.Error(t, err)
	assert.Equal(t, "reward is not available", err.Error())
	mockRewardsRepo.AssertExpectations(t)
}

func TestRedemptionService_Create_InsufficientPoints(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	userID := uuid.New().String()
	programID := uuid.New().String()
	rewardID := "reward123"

	reward := &domain.Reward{
		ID:             rewardID,
		PointsRequired: 100,
		IsActive:       true,
	}

	redemption := &domain.Redemption{
		UserID:    userID,
		RewardID:  rewardID,
		ProgramID: programID,
		Status:    "pending",
	}

	mockRewardsRepo.On("GetByID", rewardID).Return(reward, nil)
	mockPointsService.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(50, nil)

	err := service.Create(context.Background(), redemption)

	assert.Error(t, err)
	assert.Equal(t, "insufficient points", err.Error())
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
}

func TestRedemptionService_GetByID_Success(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	expectedRedemption := &domain.Redemption{
		ID:       "redemption123",
		UserID:   uuid.New().String(),
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
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	userID := uuid.New().String()
	expectedRedemptions := []domain.Redemption{
		{
			ID:       "redemption123",
			UserID:   userID,
			RewardID: "reward123",
			Status:   "completed",
		},
		{
			ID:       "redemption124",
			UserID:   userID,
			RewardID: "reward124",
			Status:   "pending",
		},
	}

	mockRedemptionRepo.On("GetByUserID", userID).Return(expectedRedemptions, nil)

	redemptions, err := service.GetByUserID(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRedemptions, redemptions)
	mockRedemptionRepo.AssertExpectations(t)
}

func TestRedemptionService_UpdateStatus_Success(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	userID := uuid.New().String()
	programID := uuid.New().String()

	redemption := &domain.Redemption{
		ID:        "redemption123",
		UserID:    userID,
		RewardID:  "reward123",
		ProgramID: programID,
		Status:    "pending",
	}

	reward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
	}

	mockRedemptionRepo.On("GetByID", "redemption123").Return(redemption, nil)
	mockRewardsRepo.On("GetByID", "reward123").Return(reward, nil)
	mockPointsService.On("EarnPoints", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"), 100, mock.AnythingOfType("*uuid.UUID")).Return(nil)
	mockRedemptionRepo.On("Update", mock.MatchedBy(func(r *domain.Redemption) bool {
		return r.ID == "redemption123" && r.Status == "canceled"
	})).Return(nil)

	err := service.UpdateStatus(context.Background(), "redemption123", "canceled")

	assert.NoError(t, err)
	mockRedemptionRepo.AssertExpectations(t)
	mockRewardsRepo.AssertExpectations(t)
	mockPointsService.AssertExpectations(t)
}

func TestRedemptionService_UpdateStatus_NotFound(t *testing.T) {
	mockRedemptionRepo := new(postgres.MockRedemptionRepository)
	mockRewardsRepo := new(postgres.MockRewardsRepository)
	mockEventRepo := new(postgres.MockEventLogRepository)
	mockPointsService := new(service.MockPointsService)

	service := NewRedemptionService(
		mockRedemptionRepo,
		mockRewardsRepo,
		mockPointsService,
		mockEventRepo,
	)

	mockRedemptionRepo.On("GetByID", "nonexistent").Return(nil, errors.New("not found"))

	err := service.UpdateStatus(context.Background(), "nonexistent", "canceled")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	mockRedemptionRepo.AssertExpectations(t)
}
