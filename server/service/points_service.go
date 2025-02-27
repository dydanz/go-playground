package service

import (
	"context"
	"math"

	"go-playground/server/domain"

	"github.com/google/uuid"
)

type PointsService struct {
	pointsRepo domain.PointsRepository
	eventRepo  domain.EventLogRepository
}

func NewPointsService(pointsRepo domain.PointsRepository, eventRepo domain.EventLogRepository) *PointsService {
	return &PointsService{
		pointsRepo: pointsRepo,
		eventRepo:  eventRepo,
	}
}

func (s *PointsService) RedeemPoints(ctx context.Context, req *domain.PointsTransaction) (*domain.PointsTransaction, error) {
	absPointsRedeemed := int(math.Abs(float64(req.Points)))	

	if absPointsRedeemed <= 0 {
		return nil, domain.NewValidationError("points", "points must be greater than 0")
	}

	currentBalance, err := s.pointsRepo.GetCurrentBalance(ctx, uuid.MustParse(req.CustomerID), uuid.MustParse(req.ProgramID))
	if err != nil {
		return nil, domain.NewSystemError("PointsService.RedeemPoints", err, "failed to get current balance")
	}

	if currentBalance < absPointsRedeemed {
		return nil, domain.NewBusinessLogicError("INSUFFICIENT_POINTS", "insufficient points balance")
	}

	ledger, err := s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: uuid.MustParse(req.CustomerID),
		ProgramID:           uuid.MustParse(req.ProgramID),
		PointsEarned:        0,
		PointsRedeemed:      absPointsRedeemed,
		PointsBalance:       currentBalance - absPointsRedeemed,
		TransactionID:       uuid.MustParse(req.TransactionID),
	})
	if err != nil {
		return nil, domain.NewSystemError("PointsService.RedeemPoints", err, "failed to create points ledger entry")
	}

	return &domain.PointsTransaction{
		TransactionID: ledger.TransactionID.String(),
		CustomerID:    ledger.MerchantCustomersID.String(),
		ProgramID:     ledger.ProgramID.String(),
		Points:        ledger.PointsRedeemed,
		Type:          "redeem",
	}, nil
}

func (s *PointsService) GetLedger(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	ledgers, err := s.pointsRepo.GetByCustomerAndProgram(ctx, customerID, programID)
	if err != nil {
		return nil, domain.NewSystemError("PointsService.GetLedger", err, "failed to get points ledger")
	}
	if len(ledgers) == 0 {
		return []*domain.PointsLedger{}, nil
	}
	return ledgers, nil
}

func (s *PointsService) GetBalance(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) (*domain.PointsBalance, error) {
	balance, err := s.pointsRepo.GetCurrentBalance(ctx, customerID, programID)
	if err != nil {
		return nil, domain.NewSystemError("PointsService.GetBalance", err, "failed to get points balance")
	}

	return &domain.PointsBalance{
		CustomerID: customerID.String(),
		ProgramID:  programID.String(),
		Balance:    balance,
	}, nil
}

func (s *PointsService) EarnPoints(ctx context.Context, req *domain.PointsTransaction) (*domain.PointsTransaction, error) {
	if req.Points <= 0 {
		return nil, domain.NewValidationError("points", "points must be greater than 0")
	}

	currentBalance, err := s.pointsRepo.GetCurrentBalance(ctx, uuid.MustParse(req.CustomerID), uuid.MustParse(req.ProgramID))
	if err != nil {
		return nil, domain.NewSystemError("PointsService.EarnPoints", err, "failed to get current balance")
	}

	ledger, err := s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: uuid.MustParse(req.CustomerID),
		ProgramID:           uuid.MustParse(req.ProgramID),
		PointsEarned:        req.Points,
		PointsRedeemed:      0,
		PointsBalance:       currentBalance + req.Points,
		TransactionID:       uuid.MustParse(req.TransactionID),
	})
	if err != nil {
		return nil, domain.NewSystemError("PointsService.EarnPoints", err, "failed to create points ledger entry")
	}

	return &domain.PointsTransaction{
		TransactionID: ledger.TransactionID.String(),
		CustomerID:    ledger.MerchantCustomersID.String(),
		ProgramID:     ledger.ProgramID.String(),
		Points:        ledger.PointsEarned,
		Type:          "earn",
	}, nil
}
