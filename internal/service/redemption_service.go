package service

import (
	"errors"
	"go-playground/internal/domain"
)

type PointsServiceInterface interface {
	GetBalance(userID string) (*domain.PointsBalance, error)
	UpdateBalance(userID string, points int) error
}

type RedemptionService struct {
	redemptionRepo domain.RedemptionRepository
	rewardsRepo    domain.RewardsRepository
	pointsService  PointsServiceInterface
	eventRepo      domain.EventLogRepository
}

func NewRedemptionService(
	redemptionRepo domain.RedemptionRepository,
	rewardsRepo domain.RewardsRepository,
	pointsService PointsServiceInterface,
	eventRepo domain.EventLogRepository,
) *RedemptionService {
	return &RedemptionService{
		redemptionRepo: redemptionRepo,
		rewardsRepo:    rewardsRepo,
		pointsService:  pointsService,
		eventRepo:      eventRepo,
	}
}

func (s *RedemptionService) Create(redemption *domain.Redemption) error {
	// Check if reward exists and is active
	reward, err := s.rewardsRepo.GetByID(redemption.RewardID)
	if err != nil {
		return err
	}
	if !reward.IsActive {
		return errors.New("reward is not available")
	}

	// Check if user has enough points
	balance, err := s.pointsService.GetBalance(redemption.UserID)
	if err != nil {
		return err
	}
	if balance.TotalPoints < reward.PointsRequired {
		return errors.New("insufficient points")
	}

	// Create redemption record
	if err := s.redemptionRepo.Create(redemption); err != nil {
		return err
	}

	// Deduct points
	if err := s.pointsService.UpdateBalance(redemption.UserID, -reward.PointsRequired); err != nil {
		return err
	}

	// Log the redemption event
	event := &domain.EventLog{
		EventType:   "reward_redeemed",
		UserID:      redemption.UserID,
		ReferenceID: &redemption.ID,
		Details: map[string]interface{}{
			"reward_id":     redemption.RewardID,
			"points_used":   reward.PointsRequired,
			"redemption_id": redemption.ID,
		},
	}
	return s.eventRepo.Create(event)
}

func (s *RedemptionService) GetByID(id string) (*domain.Redemption, error) {
	return s.redemptionRepo.GetByID(id)
}

func (s *RedemptionService) GetByUserID(userID string) ([]domain.Redemption, error) {
	return s.redemptionRepo.GetByUserID(userID)
}

func (s *RedemptionService) UpdateStatus(id string, status string) error {
	redemption, err := s.redemptionRepo.GetByID(id)
	if err != nil {
		return err
	}

	oldStatus := redemption.Status
	redemption.Status = status

	// If canceling a pending redemption, refund the points
	if oldStatus == "pending" && status == "canceled" {
		reward, err := s.rewardsRepo.GetByID(redemption.RewardID)
		if err != nil {
			return err
		}
		if err := s.pointsService.UpdateBalance(redemption.UserID, reward.PointsRequired); err != nil {
			return err
		}
	}

	return s.redemptionRepo.Update(redemption)
}
