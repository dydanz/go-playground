package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
	"go-playground/internal/mocks/repository/postgres"
)

// Test cases
func TestRewardsService_Create_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	programID := uuid.New()
	req := &domain.CreateRewardRequest{
		ProgramID:      programID,
		Name:           "Test Reward",
		Description:    "Test Description",
		PointsRequired: 100,
		IsActive:       true,
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Reward")).Return(nil)

	reward, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, reward)
	assert.Equal(t, req.PointsRequired, reward.PointsRequired)
	assert.Equal(t, req.IsActive, reward.IsActive)
	assert.Equal(t, DEFAULT_REDEEM_PROGRAM_ID, reward.ProgramID)
	assert.Equal(t, req.Name, reward.Name)
	assert.Equal(t, req.Description, reward.Description)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Create_InvalidPoints(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	programID := uuid.New()
	req := &domain.CreateRewardRequest{
		ProgramID:      programID,
		Name:           "Test Reward",
		Description:    "Test Description",
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

	rewardID := uuid.New()
	programID := uuid.New()
	expectedReward := &domain.Reward{
		ID:             rewardID,
		ProgramID:      programID,
		Name:           "Test Reward",
		Description:    "Test Description",
		PointsRequired: 100,
		IsActive:       true,
	}

	mockRepo.On("GetByID", rewardID).Return(expectedReward, nil)

	reward, err := service.GetByID(rewardID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedReward, reward)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	nonexistentID := uuid.New()
	mockRepo.On("GetByID", nonexistentID).Return(nil, errors.New("not found"))

	reward, err := service.GetByID(nonexistentID.String())

	assert.Error(t, err)
	assert.Nil(t, reward)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Update_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	rewardID := uuid.New()
	programID := uuid.New()
	existingReward := &domain.Reward{
		ID:             rewardID,
		ProgramID:      programID,
		Name:           "Original Name",
		Description:    "Original Description",
		PointsRequired: 100,
		IsActive:       true,
	}

	req := &domain.UpdateRewardRequest{
		Name:           "Updated Name",
		Description:    "Updated Description",
		PointsRequired: new(int),
		IsActive:       new(bool),
	}
	*req.PointsRequired = 150
	*req.IsActive = true

	mockRepo.On("GetByID", rewardID).Return(existingReward, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Reward")).Return(nil)

	reward, err := service.Update(rewardID.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, reward)
	assert.Equal(t, "Updated Name", reward.Name)
	assert.Equal(t, "Updated Description", reward.Description)
	assert.Equal(t, 150, reward.PointsRequired)
	assert.Equal(t, true, reward.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Update_InvalidPoints(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	rewardID := uuid.New()
	programID := uuid.New()
	existingReward := &domain.Reward{
		ID:             rewardID,
		ProgramID:      programID,
		Name:           "Test Reward",
		Description:    "Test Description",
		PointsRequired: 100,
		IsActive:       true,
	}

	mockRepo.On("GetByID", rewardID).Return(existingReward, nil)

	req := &domain.UpdateRewardRequest{
		PointsRequired: new(int),
		IsActive:       new(bool),
	}
	*req.PointsRequired = 0
	*req.IsActive = true

	reward, err := service.Update(rewardID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, reward)
	assert.Equal(t, "points required must be greater than 0", err.Error())
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_Delete_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	rewardID := uuid.New()
	mockRepo.On("Delete", rewardID).Return(nil)

	err := service.Delete(rewardID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_UpdateAvailability_Success(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	rewardID := uuid.New()
	programID := uuid.New()
	existingReward := &domain.Reward{
		ID:             rewardID,
		ProgramID:      programID,
		Name:           "Test Reward",
		Description:    "Test Description",
		PointsRequired: 100,
		IsActive:       false,
	}

	mockRepo.On("GetByID", rewardID).Return(existingReward, nil)
	mockRepo.On("Update", mock.MatchedBy(func(r *domain.Reward) bool {
		return r.ID == rewardID && r.IsActive == true
	})).Return(nil)

	err := service.UpdateAvailability(rewardID.String(), true)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRewardsService_UpdateAvailability_NotFound(t *testing.T) {
	mockRepo := new(postgres.MockRewardsRepository)
	service := NewRewardsService(mockRepo)

	nonexistentID := uuid.New()
	mockRepo.On("GetByID", nonexistentID).Return(nil, errors.New("not found"))

	err := service.UpdateAvailability(nonexistentID.String(), true)

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	mockRepo.AssertExpectations(t)
}
