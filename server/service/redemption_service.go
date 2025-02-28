package service

import (
	"context"
	"go-playground/pkg/logging"
	"go-playground/server/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RedemptionService struct {
	redemptionRepo     domain.RedemptionRepository
	rewardsRepo        domain.RewardsRepository
	pointsService      domain.PointsService
	transactionService domain.TransactionService
	eventLoggerService domain.EventLoggerService
	logger             zerolog.Logger
}

func NewRedemptionService(
	redemptionRepo domain.RedemptionRepository,
	rewardsRepo domain.RewardsRepository,
	pointsService domain.PointsService,
	transactionService domain.TransactionService,
	eventLoggerService domain.EventLoggerService,
) *RedemptionService {
	return &RedemptionService{
		redemptionRepo:     redemptionRepo,
		rewardsRepo:        rewardsRepo,
		pointsService:      pointsService,
		transactionService: transactionService,
		eventLoggerService: eventLoggerService,
		logger:             logging.GetLogger(),
	}
}

func (s *RedemptionService) Create(ctx context.Context, redemption *domain.Redemption) error {
	// Check if reward exists and is active
	reward, err := s.rewardsRepo.GetByID(ctx, redemption.RewardID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to get reward")
		return domain.NewSystemError("RedemptionService.Create", err, "failed to get reward")
	}
	if reward == nil {
		s.logger.Error().
			Str("reward_id", redemption.RewardID.String()).
			Msg("Failed to get reward")
		return domain.NewResourceNotFoundError("reward", redemption.RewardID.String(), "reward not found")
	}
	if !reward.IsActive {
		s.logger.Error().
			Str("reward_id", redemption.RewardID.String()).
			Msg("Failed to get reward")
		return domain.NewBusinessLogicError("REWARD_INACTIVE", "reward is not available")
	}

	// Parse user ID and program ID to UUID
	customerID, err := uuid.Parse(redemption.MerchantCustomersID.String())
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to parse customer ID")
		return domain.NewValidationError("customer_id", "invalid customer ID format")
	}

	// Check if user has enough points
	balance, err := s.pointsService.GetBalance(ctx, customerID, reward.ProgramID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to get points balance")
		return domain.NewSystemError("RedemptionService.Create", err, "failed to get points balance")
	}
	if balance.Balance < reward.PointsRequired {
		s.logger.Error().
			Str("reward_id", redemption.RewardID.String()).
			Msg("Failed to get reward")
		return domain.NewBusinessLogicError("INSUFFICIENT_POINTS", "insufficient points")
	}

	// Set points_used in redemption record
	redemption.PointsUsed = reward.PointsRequired

	// Create redemption record
	redemptions, err := s.redemptionRepo.Create(ctx, redemption)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to create redemption")
		return domain.NewSystemError("RedemptionService.Create", err, "failed to create redemption")
	}
	redemption = redemptions[0]

	// Deduct points by creating a redemption transaction
	transaction, err := s.transactionService.Create(ctx, &domain.CreateTransactionRequest{
		MerchantCustomersID: redemption.MerchantCustomersID,
		MerchantID:          uuid.Nil, // filled in by the transaction service
		ProgramID:           reward.ProgramID,
		TransactionType:     "redemption",
		TransactionAmount:   float64(reward.PointsRequired),
		TransactionDate:     redemption.RedemptionDate,
	})
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to create redemption transaction")
		return domain.NewSystemError("RedemptionService.Create", err, "failed to create redemption transaction")
	}

	s.logger.Info().
		Str("redemption_id", redemption.ID.String()).
		Str("paired_tx_id", transaction.TransactionID.String()).
		Msg("transaction record for redemption")

	// Log the redemption event
	go s.eventLoggerService.SaveRedemptionEvents(ctx, domain.RewardRedeemed, redemption, reward)

	return nil
}

func (s *RedemptionService) GetByID(id string) (*domain.Redemption, error) {
	redemption, err := s.redemptionRepo.GetByID(context.Background(), uuid.MustParse(id))
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to get redemption")
		return nil, domain.NewSystemError("RedemptionService.GetByID", err, "failed to get redemption")
	}
	if redemption == nil {
		s.logger.Error().
			Str("redemption_id", id).
			Msg("Failed to get redemption")
		return nil, domain.NewResourceNotFoundError("redemption", id, "redemption not found")
	}
	return redemption, nil
}

func (s *RedemptionService) GetByUserID(userID string) ([]*domain.Redemption, error) {
	redemptions, err := s.redemptionRepo.GetByUserID(context.Background(), uuid.MustParse(userID))
	if err != nil {
		return nil, domain.NewSystemError("RedemptionService.GetByUserID", err, "failed to get redemptions")
	}
	if len(redemptions) == 0 {
		return []*domain.Redemption{}, nil
	}
	return redemptions, nil
}

func (s *RedemptionService) UpdateStatus(ctx context.Context, id string, status string) error {
	redemption, err := s.redemptionRepo.GetByID(context.Background(), uuid.MustParse(id))
	if err != nil {
		return domain.NewSystemError("RedemptionService.UpdateStatus", err, "failed to get redemption")
	}
	if redemption == nil {
		return domain.NewResourceNotFoundError("redemption", id, "redemption not found")
	}

	oldStatus := redemption.Status
	redemption.Status = domain.RedemptionStatus(status)

	// If canceling a pending redemption, refund the points
	if oldStatus == "pending" && status == "canceled" {
		reward, err := s.rewardsRepo.GetByID(ctx, redemption.RewardID)
		if err != nil {
			return domain.NewSystemError("RedemptionService.UpdateStatus", err, "failed to get reward")
		}
		if reward == nil {
			return domain.NewResourceNotFoundError("reward", redemption.RewardID.String(), "reward not found")
		}

		customerID, err := uuid.Parse(redemption.MerchantCustomersID.String())
		if err != nil {
			return domain.NewValidationError("customer_id", "invalid customer ID format")
		}

		refundID := uuid.New()
		_, err = s.pointsService.EarnPoints(ctx, &domain.PointsTransaction{
			CustomerID:    customerID.String(),
			ProgramID:     reward.ProgramID.String(),
			Points:        reward.PointsRequired,
			Type:          "refund",
			TransactionID: refundID.String(),
		})
		if err != nil {
			return domain.NewSystemError("RedemptionService.UpdateStatus", err, "failed to refund points")
		}
	}

	if err := s.redemptionRepo.Update(ctx, redemption); err != nil {
		return domain.NewSystemError("RedemptionService.UpdateStatus", err, "failed to update redemption status")
	}

	return nil
}

func (s *RedemptionService) SetPointsService(pointsService domain.PointsService) {
	s.pointsService = pointsService
}
