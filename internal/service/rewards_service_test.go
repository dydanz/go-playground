package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
)

// Test cases
func TestRewardsService_Create_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	req := &domain.CreateRewardRequest{
		PointsRequired: 100,
		IsActive:       true,
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Reward")).Return(nil)

	reward, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, reward)
	assert.Equal(t, req.PointsRequired, reward.PointsRequired)
	assert.Equal(t, req.IsActive, reward.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Create_InvalidPoints(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	req := &domain.CreateRewardRequest{
		PointsRequired: 0,
		IsActive:       true,
	}

	reward, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, reward)
	assert.Equal(t, "points required must be greater than 0", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestRewardsService_GetByID_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	expectedReward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
		IsActive:       true,
	}

	mockRepo.On("GetByID", "reward123").Return(expectedReward, nil)

	reward, err := service.GetByID("reward123")

	assert.NoError(t, err)
	assert.Equal(t, expectedReward, reward)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	mockRepo.On("GetByID", "nonexistent").Return(nil, errors.New("not found"))

	reward, err := service.GetByID("nonexistent")

	assert.Error(t, err)
	assert.Nil(t, reward)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_GetAll_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	expectedRewards := []domain.Reward{
		{
			ID:             "reward123",
			PointsRequired: 100,
			IsActive:       true,
		},
		{
			ID:             "reward124",
			PointsRequired: 200,
			IsActive:       true,
		},
	}

	mockRepo.On("GetAll", true).Return(expectedRewards, nil)

	rewards, err := service.GetAll(true)

	assert.NoError(t, err)
	assert.Equal(t, expectedRewards, rewards)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Update_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	req := &domain.UpdateRewardRequest{
		PointsRequired: new(int),
		IsActive:       new(bool),
	}
	*req.PointsRequired = 150
	*req.IsActive = true

	mockRepo.On("GetByID", "reward123").Return(&domain.Reward{ID: "reward123"}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Reward")).Return(nil)

	reward, err := service.Update("reward123", req)

	assert.NoError(t, err)
	assert.NotNil(t, reward)
	assert.Equal(t, req.PointsRequired, reward.PointsRequired)
	assert.Equal(t, req.IsActive, reward.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Update_InvalidPoints(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	req := &domain.UpdateRewardRequest{
		PointsRequired: new(int),
		IsActive:       new(bool),
	}
	*req.PointsRequired = 0
	*req.IsActive = true

	reward, err := service.Update("reward123", req)

	assert.Error(t, err)
	assert.Nil(t, reward)
	assert.Equal(t, "points required must be greater than 0", err.Error())
	mockRepo.AssertNotCalled(t, "Update")
}

func TestRewardsService_Delete_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	mockRepo.On("Delete", "reward123").Return(nil)

	err := service.Delete("reward123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_UpdateAvailability_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	existingReward := &domain.Reward{
		ID:             "reward123",
		PointsRequired: 100,
		IsActive:       false,
	}

	mockRepo.On("GetByID", "reward123").Return(existingReward, nil)
	mockRepo.On("Update", mock.MatchedBy(func(r *domain.Reward) bool {
		return r.ID == "reward123" && r.IsActive == true
	})).Return(nil)

	err := service.UpdateAvailability("reward123", true)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_UpdateAvailability_NotFound(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	mockRepo.On("GetByID", "nonexistent").Return(nil, errors.New("not found"))

	err := service.UpdateAvailability("nonexistent", true)

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	mockRepo.AssertExpectations(t)
}
