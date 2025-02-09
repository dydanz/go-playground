package service

import (
	"errors"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type RewardsService struct {
	rewardsRepo domain.RewardsRepository
}

func NewRewardsService(rewardsRepo domain.RewardsRepository) *RewardsService {
	return &RewardsService{
		rewardsRepo: rewardsRepo,
	}
}

func (s *RewardsService) Create(req *domain.CreateRewardRequest) (*domain.Reward, error) {
	if req.PointsRequired <= 0 {
		return nil, errors.New("points required must be greater than 0")
	}

	reward := &domain.Reward{
		Name:           req.Name,
		Description:    req.Description,
		PointsRequired: req.PointsRequired,
		IsActive:       req.IsActive,
		Quantity:       req.Quantity,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.rewardsRepo.Create(reward); err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *RewardsService) GetByID(id string) (*domain.Reward, error) {
	return s.rewardsRepo.GetByID(uuid.MustParse(id))
}

func (s *RewardsService) Update(id string, req *domain.UpdateRewardRequest) (*domain.Reward, error) {
	reward, err := s.rewardsRepo.GetByID(uuid.MustParse(id))
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		reward.Name = req.Name
	}
	if req.Description != "" {
		reward.Description = req.Description
	}
	if req.PointsRequired != nil {
		if *req.PointsRequired <= 0 {
			return nil, errors.New("points required must be greater than 0")
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

	if err := s.rewardsRepo.Update(reward); err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *RewardsService) Delete(id string) error {
	return s.rewardsRepo.Delete(uuid.MustParse(id))
}

func (s *RewardsService) UpdateAvailability(id string, available bool) error {
	reward, err := s.rewardsRepo.GetByID(uuid.MustParse(id))
	if err != nil {
		return err
	}

	reward.IsActive = available
	return s.rewardsRepo.Update(reward)
}
