package service

import (
	"go-playground/internal/domain"
)

type TransactionService struct {
	transactionRepo domain.TransactionRepository
	pointsService   *PointsService
	eventRepo       domain.EventLogRepository
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	pointsService *PointsService,
	eventRepo domain.EventLogRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		pointsService:   pointsService,
		eventRepo:       eventRepo,
	}
}

func (s *TransactionService) Create(tx *domain.Transaction) error {
	if err := s.transactionRepo.Create(tx); err != nil {
		return err
	}

	// Update points balance
	points := tx.Points
	if tx.TransactionType == "redeem" {
		points = -points
	}

	if err := s.pointsService.UpdateBalance(tx.UserID, points); err != nil {
		return err
	}

	// Log the transaction event
	event := &domain.EventLog{
		EventType:   "transaction",
		UserID:      tx.UserID,
		ReferenceID: tx.ID,
		Details: map[string]interface{}{
			"type":        tx.TransactionType,
			"points":      tx.Points,
			"description": tx.Description,
		},
	}
	return s.eventRepo.Create(event)
}

func (s *TransactionService) GetByID(id string) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

func (s *TransactionService) GetByUserID(userID string) ([]domain.Transaction, error) {
	return s.transactionRepo.GetByUserID(userID)
}

func (s *TransactionService) UpdateStatus(id string, status string) error {
	tx, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return err
	}

	tx.Status = status
	return s.transactionRepo.Update(tx)
}
