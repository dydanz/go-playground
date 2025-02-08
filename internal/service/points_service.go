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

func NewPointsService(pointsRepo domain.PointsRepository, eventRepo domain.EventLogRepository) domain.PointsServiceInterface {
	return &PointsService{
		pointsRepo: pointsRepo,
		eventRepo:  eventRepo,
	}
}

// Implementation of domain.PointsServiceInterface
func (s *PointsService) GetLedger(ctx context.Context, customerID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	return s.pointsRepo.GetByCustomerAndProgram(ctx, customerID, programID)
}

func (s *PointsService) GetBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error) {
	return s.pointsRepo.GetCurrentBalance(ctx, customerID, programID)
}

func (s *PointsService) EarnPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error {
	if points <= 0 {
		return ErrInvalidPoints
	}

	return s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:      uuid.New(),
		CustomerID:    customerID,
		ProgramID:     programID,
		PointsEarned:  points,
		TransactionID: transactionID,
		CreatedAt:     time.Now(),
	})
}

func (s *PointsService) RedeemPoints(ctx context.Context, customerID, programID uuid.UUID, points int, transactionID *uuid.UUID) error {
	if points <= 0 {
		return ErrInvalidPoints
	}

	currentBalance, err := s.pointsRepo.GetCurrentBalance(ctx, customerID, programID)
	if err != nil {
		return err
	}

	if currentBalance < points {
		return ErrInsufficientPoints
	}

	return s.pointsRepo.Create(ctx, &domain.PointsLedger{
		LedgerID:       uuid.New(),
		CustomerID:     customerID,
		ProgramID:      programID,
		PointsEarned:   0,
		PointsRedeemed: points,
		PointsBalance:  currentBalance - points,
		TransactionID:  transactionID,
	})
}

// LegacyPointsService adapts the new interface to the old one
type LegacyPointsService struct {
	service domain.PointsServiceInterface
}

func NewLegacyPointsService(service domain.PointsServiceInterface) domain.PointsService {
	return &LegacyPointsService{service: service}
}

func (s *LegacyPointsService) GetLedger(customerID string, programID string) (*domain.PointsLedger, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	progID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}

	ledgers, err := s.service.GetLedger(context.Background(), custID, progID)
	if err != nil {
		return nil, err
	}
	if len(ledgers) == 0 {
		return nil, nil
	}
	return ledgers[len(ledgers)-1], nil
}

func (s *LegacyPointsService) GetBalance(customerID string, programID string) (*domain.PointsBalance, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	progID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}

	balance, err := s.service.GetBalance(context.Background(), custID, progID)
	if err != nil {
		return nil, err
	}

	return &domain.PointsBalance{
		CustomerID: customerID,
		ProgramID:  programID,
		Balance:    balance,
	}, nil
}

func (s *LegacyPointsService) EarnPoints(req *domain.EarnPointsRequest) (*domain.PointsTransaction, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, err
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, err
	}

	txID := uuid.New()
	if err := s.service.EarnPoints(context.Background(), customerID, programID, req.Points, &txID); err != nil {
		return nil, err
	}

	return &domain.PointsTransaction{
		TransactionID: txID.String(),
		CustomerID:    req.CustomerID,
		ProgramID:     req.ProgramID,
		Points:        req.Points,
		Type:          "earn",
		CreatedAt:     time.Now(),
	}, nil
}

func (s *LegacyPointsService) RedeemPoints(req *domain.RedeemPointsRequest) (*domain.PointsTransaction, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, err
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, err
	}

	txID := uuid.New()
	if err := s.service.RedeemPoints(context.Background(), customerID, programID, req.Points, &txID); err != nil {
		return nil, err
	}

	return &domain.PointsTransaction{
		TransactionID: txID.String(),
		CustomerID:    req.CustomerID,
		ProgramID:     req.ProgramID,
		Points:        req.Points,
		Type:          "redeem",
		CreatedAt:     time.Now(),
	}, nil
}
