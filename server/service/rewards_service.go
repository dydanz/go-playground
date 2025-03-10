package service

import (
	"context"
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RewardsService struct {
	rewardsRepo domain.RewardsRepository
	logger      zerolog.Logger
}

func NewRewardsService(rewardsRepo domain.RewardsRepository) *RewardsService {
	return &RewardsService{
		rewardsRepo: rewardsRepo,
		logger:      logging.GetLogger(),
	}
}

func (s *RewardsService) Create(ctx context.Context, req *domain.CreateRewardRequest) (*domain.Reward, error) {
	if req.Name == "" {
		s.logger.Error().
			Msg("Reward name is required")
		return nil, domain.NewValidationError("name", "reward name is required")
	}
	if req.PointsRequired <= 0 {
		s.logger.Error().
			Msg("Points required must be greater than 0")
		return nil, domain.NewValidationError("points_required", "points required must be greater than 0")
	}

	reward := &domain.Reward{
		Name:           req.Name,
		ProgramID:      req.ProgramID,
		Description:    req.Description,
		PointsRequired: req.PointsRequired,
		IsActive:       req.IsActive,
		Quantity:       req.Quantity,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := s.rewardsRepo.Create(ctx, reward)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error creating reward")
		return nil, domain.NewSystemError("RewardsService.Create", err, "failed to create reward")
	}

	return result, nil
}

func (s *RewardsService) GetByID(ctx context.Context, id string) (*domain.Reward, error) {
	reward, err := s.rewardsRepo.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting reward")
		return nil, domain.NewSystemError("RewardsService.GetByID", err, "failed to get reward")
	}
	if reward == nil {
		s.logger.Error().
			Msg("Reward not found")
		return nil, domain.NewResourceNotFoundError("reward", id, "reward not found")
	}
	return reward, nil
}

func (s *RewardsService) GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*domain.Reward, error) {
	rewards, err := s.rewardsRepo.GetByProgramID(ctx, programID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting rewards")
		return nil, domain.NewSystemError("RewardsService.GetByProgramID", err, "failed to get rewards")
	}
	return rewards, nil
}

func (s *RewardsService) Update(ctx context.Context, id string, req *domain.UpdateRewardRequest) (*domain.Reward, error) {
	reward, err := s.rewardsRepo.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting reward")
		return nil, domain.NewSystemError("RewardsService.Update", err, "failed to get reward")
	}
	if reward == nil {
		s.logger.Error().
			Msg("Reward not found")
		return nil, domain.NewResourceNotFoundError("reward", id, "reward not found")
	}

	if req.Name != "" {
		reward.Name = req.Name
	}
	if req.Description != "" {
		reward.Description = req.Description
	}
	if req.PointsRequired != nil {
		if *req.PointsRequired <= 0 {
			return nil, domain.NewValidationError("points_required", "points required must be greater than 0")
		}
		reward.PointsRequired = *req.PointsRequired
	}
	if req.IsActive != nil {
		reward.IsActive = *req.IsActive
	}
	if req.Quantity != nil {
		reward.Quantity = *req.Quantity
	}
	reward.UpdatedAt = time.Now()

	result, err := s.rewardsRepo.Update(ctx, reward)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error updating reward")
		return nil, domain.NewSystemError("RewardsService.Update", err, "failed to update reward")
	}

	return result, nil
}

func (s *RewardsService) Delete(ctx context.Context, id string) error {
	reward, err := s.rewardsRepo.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting reward")
		return domain.NewSystemError("RewardsService.Delete", err, "failed to get reward")
	}
	if reward == nil {
		s.logger.Error().
			Msg("Reward not found")
		return domain.NewResourceNotFoundError("reward", id, "reward not found")
	}

	if err := s.rewardsRepo.Delete(ctx, uuid.MustParse(id)); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error deleting reward")
		return domain.NewSystemError("RewardsService.Delete", err, "failed to delete reward")
	}
	return nil
}

func (s *RewardsService) UpdateAvailability(ctx context.Context, id string, available bool) (*domain.Reward, error) {
	reward, err := s.rewardsRepo.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting reward")
		return nil, domain.NewSystemError("RewardsService.UpdateAvailability", err, "failed to get reward")
	}

	reward.IsActive = available
	result, err := s.rewardsRepo.Update(ctx, reward)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error updating reward availability")
		return nil, domain.NewSystemError("RewardsService.UpdateAvailability", err, "failed to update reward availability")
	}

	return result, nil
}
