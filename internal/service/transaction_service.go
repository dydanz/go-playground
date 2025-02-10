package service

import (
	"context"
	"time"

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
		MerchantID:          req.MerchantID,
		MerchantCustomersID: req.MerchantCustomersID,
		ProgramID:           req.ProgramID,
		BranchID:            req.BranchID,
		TransactionType:     req.TransactionType,
		TransactionAmount:   req.TransactionAmount,
		Status:              "pending",
		CreatedAt:           time.Now(),
	}

	createdTx, err := s.transactionRepo.Create(ctx, tx)
	if err != nil {
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
	case "redemption":
		points = int(tx.TransactionAmount * -1)
	}

	// TODO: Check if the transaction is valid for the program, branch/merchant
	// TODO: Calculate points based on the transaction type and program rules
	// TODO: Check if the customer has enough points to redeem
	// TODO: Update points balance if applicable

	if points != 0 {
		if err := s.pointsService.EarnPoints(ctx, tx.MerchantCustomersID, tx.ProgramID, points, createdTx.TransactionID); err != nil {
			return nil, err
		}
	}

	// Log the transaction event, make async
	txIDStr := tx.TransactionID.String()
	event := &domain.EventLog{
		EventType:   "transaction_created",
		UserID:      req.MerchantCustomersID.String(),
		ReferenceID: &txIDStr,
		Details: map[string]interface{}{
			"merchant_id":        tx.MerchantID,
			"transaction_type":   tx.TransactionType,
			"transaction_amount": tx.TransactionAmount,
			"points_earned":      points,
			"branch_id":          tx.BranchID,
		},
	}

	go s.eventRepo.Create(ctx, event)

	return createdTx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	txID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.transactionRepo.GetByID(ctx, txID)
}

func (s *TransactionService) GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	txs, _, err := s.transactionRepo.GetByCustomerIDWithPagination(ctx, custID, 0, -1)
	return txs, err
}

func (s *TransactionService) GetByCustomerIDWithPagination(ctx context.Context, customerID string, offset, limit int) ([]*domain.Transaction, int64, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, 0, err
	}
	return s.transactionRepo.GetByCustomerIDWithPagination(ctx, custID, offset, limit)
}

func (s *TransactionService) GetByMerchantID(ctx context.Context, merchantID string) ([]*domain.Transaction, error) {
	merchID, err := uuid.Parse(merchantID)
	if err != nil {
		return nil, err
	}
	return s.transactionRepo.GetByMerchantID(ctx, merchID)
}

func (s *TransactionService) UpdateStatus(ctx context.Context, id string, status string) error {
	txID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.transactionRepo.UpdateStatus(ctx, txID, status)
}

func (s *TransactionService) SetPointsService(pointsService domain.PointsServiceInterface) {
	s.pointsService = pointsService
}
