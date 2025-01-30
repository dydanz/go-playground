package service

import (
	"context"
	"errors"

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

	currentBalance, err := s.pointsRepo.GetCurrentBalance(ctx, customerID, programID)
	if err != nil {
		return err
	}

	ledger := &domain.PointsLedger{
		LedgerID:       uuid.New(),
		CustomerID:     customerID,
		ProgramID:      programID,
		PointsEarned:   points,
		PointsRedeemed: 0,
		PointsBalance:  currentBalance + points,
		TransactionID:  transactionID,
	}

	if err := s.pointsRepo.Create(ctx, ledger); err != nil {
		return err
	}

	// Log the points earned event
	event := &domain.EventLog{
		EventType: "points_earned",
		UserID:    customerID.String(),
		Details: map[string]interface{}{
			"program_id":     programID.String(),
			"points_earned":  points,
			"new_balance":    ledger.PointsBalance,
			"transaction_id": transactionID,
		},
	}
	return s.eventRepo.Create(event)
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

	ledger := &domain.PointsLedger{
		LedgerID:       uuid.New(),
		CustomerID:     customerID,
		ProgramID:      programID,
		PointsEarned:   0,
		PointsRedeemed: points,
		PointsBalance:  currentBalance - points,
		TransactionID:  transactionID,
	}

	if err := s.pointsRepo.Create(ctx, ledger); err != nil {
		return err
	}

	// Log the points redeemed event
	event := &domain.EventLog{
		EventType: "points_redeemed",
		UserID:    customerID.String(),
		Details: map[string]interface{}{
			"program_id":      programID.String(),
			"points_redeemed": points,
			"new_balance":     ledger.PointsBalance,
			"transaction_id":  transactionID,
		},
	}
	return s.eventRepo.Create(event)
}
