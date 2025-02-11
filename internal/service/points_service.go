package service

import (
	"context"
	"errors"
	"time"

	"go-playground/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrInsufficientPoints = errors.New("insufficient points balance")
	ErrInvalidPoints      = errors.New("invalid points amount")
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
	currentBalance, err := s.pointsRepo.GetCurrentBalance(ctx, uuid.MustParse(req.CustomerID), uuid.MustParse(req.ProgramID))
	if err != nil {
		return nil, err
	}

	if currentBalance < req.Points {
		return nil, ErrInsufficientPoints
	}

	err = s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: uuid.MustParse(req.CustomerID),
		ProgramID:           uuid.MustParse(req.ProgramID),
		PointsEarned:        0,
		PointsRedeemed:      req.Points,
		PointsBalance:       currentBalance - req.Points,
		TransactionID:       uuid.MustParse(req.TransactionID),
	})
	if err != nil {
		return nil, err
	}

	return &domain.PointsTransaction{
		TransactionID: req.TransactionID,
		CustomerID:    req.CustomerID,
		ProgramID:     req.ProgramID,
		Points:        req.Points,
		Type:          "redeem",
	}, nil
}

func (s *PointsService) GetLedger(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) ([]*domain.PointsLedger, error) {

	ledgers, err := s.pointsRepo.GetByCustomerAndProgram(context.Background(), customerID, programID)
	if err != nil {
		return nil, err
	}
	if len(ledgers) == 0 {
		return nil, nil
	}
	return ledgers, nil
}

func (s *PointsService) GetBalance(ctx context.Context, customerID uuid.UUID, programID uuid.UUID) (*domain.PointsBalance, error) {
	balance, err := s.pointsRepo.GetCurrentBalance(context.Background(), customerID, programID)
	if err != nil {
		return nil, err
	}

	return &domain.PointsBalance{
		CustomerID: customerID.String(),
		ProgramID:  programID.String(),
		Balance:    balance,
	}, nil
}

func (s *PointsService) EarnPoints(ctx context.Context, req *domain.PointsTransaction) (*domain.PointsTransaction, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, err
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, err
	}

	s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: customerID,
		ProgramID:           programID,
		PointsEarned:        req.Points,
		TransactionID:       uuid.MustParse(req.TransactionID),
		CreatedAt:           time.Now(),
	})

	return &domain.PointsTransaction{
		TransactionID: req.TransactionID,
		CustomerID:    req.CustomerID,
		ProgramID:     req.ProgramID,
		Points:        req.Points,
		Type:          "earn",
		CreatedAt:     time.Now(),
	}, nil
}
