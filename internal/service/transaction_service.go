package service

import (
	"context"

	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo domain.TransactionRepository
	pointsService   domain.PointsServiceInterface
	eventRepo       domain.EventLogRepository
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	pointsService domain.PointsServiceInterface,
	eventRepo domain.EventLogRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		pointsService:   pointsService,
		eventRepo:       eventRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	tx := &domain.Transaction{
		TransactionID:     uuid.New(),
		MerchantID:        req.MerchantID,
		CustomerID:        req.CustomerID,
		TransactionType:   req.TransactionType,
		TransactionAmount: req.TransactionAmount,
		BranchID:          req.BranchID,
	}

	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	// Calculate points based on transaction amount and type
	var points int
	switch tx.TransactionType {
	case "purchase":
		points = int(tx.TransactionAmount) // Example: 1 point per currency unit
	case "refund":
		points = -int(tx.TransactionAmount)
	case "bonus":
		points = int(tx.TransactionAmount * 2) // Example: Double points for bonus
	}

	// Update points balance if applicable
	if points != 0 {
		if err := s.pointsService.EarnPoints(ctx, tx.CustomerID, tx.MerchantID, points, &tx.TransactionID); err != nil {
			return nil, err
		}
	}

	// Log the transaction event
	txIDStr := tx.TransactionID.String()
	event := &domain.EventLog{
		EventType:   "transaction_created",
		UserID:      tx.CustomerID.String(),
		ReferenceID: &txIDStr,
		Details: map[string]interface{}{
			"merchant_id":        tx.MerchantID,
			"transaction_type":   tx.TransactionType,
			"transaction_amount": tx.TransactionAmount,
			"points_earned":      points,
			"branch_id":          tx.BranchID,
		},
	}
	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, transactionID)
}

func (s *TransactionService) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Transaction, error) {
	return s.transactionRepo.GetByCustomerID(ctx, customerID)
}

func (s *TransactionService) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Transaction, error) {
	return s.transactionRepo.GetByMerchantID(ctx, merchantID)
}

func (s *TransactionService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	tx, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tx.Status = status
	return s.transactionRepo.UpdateStatus(ctx, tx.TransactionID, status)
}
