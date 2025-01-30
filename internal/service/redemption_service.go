package service

import (
	"context"
	"errors"
	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type RedemptionService struct {
	redemptionRepo domain.RedemptionRepository
	rewardsRepo    domain.RewardsRepository
	pointsService  domain.PointsServiceInterface
	eventRepo      domain.EventLogRepository
}

func NewRedemptionService(
	redemptionRepo domain.RedemptionRepository,
	rewardsRepo domain.RewardsRepository,
	pointsService domain.PointsServiceInterface,
	eventRepo domain.EventLogRepository,
) *RedemptionService {
	return &RedemptionService{
		redemptionRepo: redemptionRepo,
		rewardsRepo:    rewardsRepo,
		pointsService:  pointsService,
		eventRepo:      eventRepo,
	}
}

func (s *RedemptionService) Create(ctx context.Context, redemption *domain.Redemption) error {
	// Check if reward exists and is active
	reward, err := s.rewardsRepo.GetByID(redemption.RewardID)
	if err != nil {
		return err
	}
	if !reward.IsActive {
		return errors.New("reward is not available")
	}

	// Parse user ID and program ID to UUID
	customerID, err := uuid.Parse(redemption.UserID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	if redemption.ProgramID == "" {
		return errors.New("program ID is required")
	}

	programID, err := uuid.Parse(redemption.ProgramID)
	if err != nil {
		return errors.New("invalid program ID format")
	}

	// Check if user has enough points
	balance, err := s.pointsService.GetBalance(ctx, customerID, programID)
	if err != nil {
		return err
	}
	if balance < reward.PointsRequired {
		return errors.New("insufficient points")
	}

	// Create redemption record
	if err := s.redemptionRepo.Create(redemption); err != nil {
		return err
	}

	// Deduct points
	redemptionID := uuid.New()
	if err := s.pointsService.RedeemPoints(ctx, customerID, programID, reward.PointsRequired, &redemptionID); err != nil {
		return err
	}

	// Set the redemption ID
	redemption.ID = redemptionID.String()

	// Log the redemption event
	event := &domain.EventLog{
		EventType:   "reward_redeemed",
		UserID:      redemption.UserID,
		ReferenceID: &redemption.ID,
		Details: map[string]interface{}{
			"reward_id":     redemption.RewardID,
			"points_used":   reward.PointsRequired,
			"redemption_id": redemption.ID,
			"program_id":    redemption.ProgramID,
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

func (s *RedemptionService) UpdateStatus(ctx context.Context, id string, status string) error {
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
		customerID, err := uuid.Parse(redemption.UserID)
		if err != nil {
			return errors.New("invalid user ID format")
		}

		programID, err := uuid.Parse(redemption.ProgramID)
		if err != nil {
			return errors.New("invalid program ID format")
		}

		refundID := uuid.New()
		if err := s.pointsService.EarnPoints(ctx, customerID, programID, reward.PointsRequired, &refundID); err != nil {
			return err
		}
	}

	return s.redemptionRepo.Update(redemption)
}
