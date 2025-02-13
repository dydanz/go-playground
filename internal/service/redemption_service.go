package service

import (
	"context"
	"go-playground/internal/domain"
	"log"

	"github.com/google/uuid"
)

type RedemptionService struct {
	redemptionRepo     domain.RedemptionRepository
	rewardsRepo        domain.RewardsRepository
	pointsService      domain.PointsService
	transactionService domain.TransactionService
	eventRepo          domain.EventLogRepository
}

func NewRedemptionService(
	redemptionRepo domain.RedemptionRepository,
	rewardsRepo domain.RewardsRepository,
	pointsService domain.PointsService,
	transactionService domain.TransactionService,
	eventRepo domain.EventLogRepository,
) *RedemptionService {
	return &RedemptionService{
		redemptionRepo:     redemptionRepo,
		rewardsRepo:        rewardsRepo,
		pointsService:      pointsService,
		transactionService: transactionService,
		eventRepo:          eventRepo,
	}
}

func (s *RedemptionService) Create(ctx context.Context, redemption *domain.Redemption) error {
	// Check if reward exists and is active
	reward, err := s.rewardsRepo.GetByID(ctx, redemption.RewardID)
	if err != nil {
		return domain.NewSystemError("RedemptionService.Create", err, "failed to get reward")
	}
	if reward == nil {
		return domain.NewResourceNotFoundError("reward", redemption.RewardID.String(), "reward not found")
	}
	if !reward.IsActive {
		return domain.NewBusinessLogicError("REWARD_INACTIVE", "reward is not available")
	}

	// Parse user ID and program ID to UUID
	customerID, err := uuid.Parse(redemption.MerchantCustomersID.String())
	if err != nil {
		return domain.NewValidationError("customer_id", "invalid customer ID format")
	}

	// Check if user has enough points
	balance, err := s.pointsService.GetBalance(ctx, customerID, reward.ProgramID)
	if err != nil {
		return domain.NewSystemError("RedemptionService.Create", err, "failed to get points balance")
	}
	if balance.Balance < reward.PointsRequired {
		return domain.NewBusinessLogicError("INSUFFICIENT_POINTS", "insufficient points")
	}

	// Set points_used in redemption record
	redemption.PointsUsed = reward.PointsRequired

	// Create redemption record
	redemptions, err := s.redemptionRepo.Create(ctx, redemption)
	if err != nil {
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
	})
	if err != nil {
		log.Printf("error creating redemption transaction for redemption id: %s, error: %v", redemption.ID, err)
		return domain.NewSystemError("RedemptionService.Create", err, "failed to create redemption transaction")
	}

	log.Println("transaction record for redemption id: ", redemption.ID, "paired tx-id: ", transaction.TransactionID)

	// Log the redemption event
	event := &domain.EventLog{
		EventType:   string(domain.RewardRedeemed),
		ActorID:     redemption.MerchantCustomersID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := redemption.ID.String(); return &s }(),
		Details: map[string]interface{}{
			"reward_id":     redemption.RewardID,
			"points_used":   reward.PointsRequired,
			"redemption_id": redemption.ID,
			"program_id":    reward.ProgramID,
		},
	}

	go s.eventRepo.Create(ctx, event)

	return nil
}

func (s *RedemptionService) GetByID(id string) (*domain.Redemption, error) {
	redemption, err := s.redemptionRepo.GetByID(context.Background(), uuid.MustParse(id))
	if err != nil {
		return nil, domain.NewSystemError("RedemptionService.GetByID", err, "failed to get redemption")
	}
	if redemption == nil {
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
