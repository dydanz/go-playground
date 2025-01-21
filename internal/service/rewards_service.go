package service

import (
	"errors"
	"go-playground/internal/domain"
)

type RewardsService struct {
	rewardsRepo domain.RewardsRepository
}

func NewRewardsService(rewardsRepo domain.RewardsRepository) *RewardsService {
	return &RewardsService{
		rewardsRepo: rewardsRepo,
	}
}

func (s *RewardsService) Create(reward *domain.Reward) error {
	if reward.PointsRequired <= 0 {
		return errors.New("points required must be greater than 0")
	}
	return s.rewardsRepo.Create(reward)
}

func (s *RewardsService) GetByID(id string) (*domain.Reward, error) {
	return s.rewardsRepo.GetByID(id)
}

func (s *RewardsService) GetAll(active bool) ([]domain.Reward, error) {
	return s.rewardsRepo.GetAll(active)
}

func (s *RewardsService) Update(reward *domain.Reward) error {
	if reward.PointsRequired <= 0 {
		return errors.New("points required must be greater than 0")
	}
	return s.rewardsRepo.Update(reward)
}

func (s *RewardsService) Delete(id string) error {
	return s.rewardsRepo.Delete(id)
}

func (s *RewardsService) UpdateAvailability(id string, available bool) error {
	reward, err := s.rewardsRepo.GetByID(id)
	if err != nil {
		return err
	}

	reward.IsActive = available
	return s.rewardsRepo.Update(reward)
}
