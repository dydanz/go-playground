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

func (s *PointsService) GetLedger(customerID string, programID string) (*domain.PointsLedger, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	progID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}

	ledgers, err := s.pointsRepo.GetByCustomerAndProgram(context.Background(), custID, progID)
	if err != nil {
		return nil, err
	}
	if len(ledgers) == 0 {
		return nil, nil
	}
	return ledgers[len(ledgers)-1], nil
}

func (s *PointsService) GetBalance(customerID string, programID string) (*domain.PointsBalance, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	progID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}

	balance, err := s.pointsRepo.GetCurrentBalance(context.Background(), custID, progID)
	if err != nil {
		return nil, err
	}

	return &domain.PointsBalance{
		CustomerID: customerID,
		ProgramID:  programID,
		Balance:    balance,
	}, nil
}

func (s *PointsService) EarnPoints(req *domain.EarnPointsRequest) (*domain.PointsTransaction, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, err
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, err
	}

	if req.Points <= 0 {
		return nil, ErrInvalidPoints
	}

	txID := uuid.New()
	if err := s.pointsRepo.Create(context.Background(), &domain.PointsLedger{
		LedgerID:      uuid.New(),
		CustomerID:    customerID,
		ProgramID:     programID,
		PointsEarned:  req.Points,
		TransactionID: &txID,
		CreatedAt:     time.Now(),
	}); err != nil {
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

func (s *PointsService) RedeemPoints(req *domain.RedeemPointsRequest) (*domain.PointsTransaction, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, err
	}

	programID, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, err
	}

	if req.Points <= 0 {
		return nil, ErrInvalidPoints
	}

	currentBalance, err := s.pointsRepo.GetCurrentBalance(context.Background(), customerID, programID)
	if err != nil {
		return nil, err
	}

	if currentBalance < req.Points {
		return nil, ErrInsufficientPoints
	}

	txID := uuid.New()
	ledger := &domain.PointsLedger{
		LedgerID:       uuid.New(),
		CustomerID:     customerID,
		ProgramID:      programID,
		PointsEarned:   0,
		PointsRedeemed: req.Points,
		PointsBalance:  currentBalance - req.Points,
		TransactionID:  &txID,
	}

	if err := s.pointsRepo.Create(context.Background(), ledger); err != nil {
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
