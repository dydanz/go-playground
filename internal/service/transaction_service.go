package service

import (
	"context"
	"log"

	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo      domain.TransactionRepository
	pointsService        domain.PointsService
	eventRepo            domain.EventLogRepository
	merchantCustomerRepo domain.MerchantCustomersRepository
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	pointsService domain.PointsService,
	eventRepo domain.EventLogRepository,
	merchantCustomerRepo domain.MerchantCustomersRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo:      transactionRepo,
		pointsService:        pointsService,
		eventRepo:            eventRepo,
		merchantCustomerRepo: merchantCustomerRepo,
	}
}

func (s *TransactionService) getMerchantIDByCustomerID(ctx context.Context, customerID uuid.UUID) (uuid.UUID, error) {
	customer, err := s.merchantCustomerRepo.GetByID(ctx, customerID)
	if err != nil {
		return uuid.Nil, domain.NewSystemError("TransactionService.getMerchantIDByCustomerID", err, "failed to get merchant customer")
	}
	if customer == nil {
		return uuid.Nil, domain.NewResourceNotFoundError("merchant customer", customerID.String(), "customer not found")
	}
	return customer.MerchantID, nil
}

func (s *TransactionService) Create(ctx context.Context, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	if req.TransactionAmount <= 0 {
		log.Println("TransactionService: Transaction amount must be greater than 0")
		return nil, domain.NewValidationError("transaction_amount", "transaction amount must be greater than 0")
	}

	// Get merchant ID from customer ID
	merchantID, err := s.getMerchantIDByCustomerID(ctx, req.MerchantCustomersID)
	if err != nil {
		log.Println("TransactionService: Error getting merchant ID: ", err)
		return nil, err
	}

	transaction := &domain.Transaction{
		TransactionID:       uuid.New(),
		MerchantCustomersID: req.MerchantCustomersID,
		MerchantID:          merchantID,
		ProgramID:           req.ProgramID,
		TransactionType:     req.TransactionType,
		TransactionAmount:   req.TransactionAmount,
	}

	createdTx, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		log.Println("TransactionService: Error creating transaction: ", err)
		return nil, domain.NewSystemError("TransactionService.Create", err, "failed to create transaction")
	}

	// Calculate points based on transaction amount and type
	var points int
	switch transaction.TransactionType {
	case "purchase":
		points = int(transaction.TransactionAmount) // Example: 1 point per currency unit
	case "refund":
		points = -int(transaction.TransactionAmount)
	case "bonus":
		points = int(transaction.TransactionAmount * 2) // Example: Double points for bonus
	case "redemption":
		points = int(transaction.TransactionAmount * -1)
	}

	// TODO: Check if the transaction is valid for the program, branch/merchant
	// TODO: Calculate points based on the transaction type and program rules
	// TODO: Check if the customer has enough points to redeem
	// TODO: Update points balance if applicable

	if points > 0 {
		if _, err := s.pointsService.EarnPoints(ctx, &domain.PointsTransaction{
			CustomerID:    transaction.MerchantCustomersID.String(),
			ProgramID:     transaction.ProgramID.String(),
			Points:        points,
			TransactionID: createdTx.TransactionID.String(),
		}); err != nil {
			log.Println("TransactionService: Error earning points: ", err)
			return nil, err
		}
	} else if points < 0 {
		if _, err := s.pointsService.RedeemPoints(ctx, &domain.PointsTransaction{
			CustomerID:    transaction.MerchantCustomersID.String(),
			ProgramID:     transaction.ProgramID.String(),
			Points:        points,
			TransactionID: createdTx.TransactionID.String(),
		}); err != nil {
			log.Println("TransactionService: Error redeeming points: ", err)
			return nil, err
		}
	}

	// Log the transaction event
	event := &domain.EventLog{
		EventType:   string(domain.TransactionCreated),
		ActorID:     req.MerchantCustomersID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := createdTx.TransactionID.String(); return &s }(),
		Details: map[string]interface{}{
			"transaction_id":     createdTx.TransactionID,
			"merchant_id":        merchantID,
			"program_id":         req.ProgramID,
			"transaction_type":   req.TransactionType,
			"transaction_amount": req.TransactionAmount,
			"points_earned":      points,
		},
	}

	go s.eventRepo.Create(ctx, event)

	return createdTx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewSystemError("TransactionService.GetByID", err, "failed to get transaction")
	}
	if transaction == nil {
		return nil, domain.NewResourceNotFoundError("transaction", id.String(), "transaction not found")
	}
	return transaction, nil
}

func (s *TransactionService) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, domain.NewSystemError("TransactionService.GetByCustomerID", err, "failed to get transactions")
	}
	if len(transactions) == 0 {
		return []*domain.Transaction{}, nil
	}
	return transactions, nil
}

func (s *TransactionService) GetByCustomerIDWithPagination(ctx context.Context, customerID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	return s.transactionRepo.GetByCustomerIDWithPagination(ctx, customerID, offset, limit)
}

func (s *TransactionService) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, domain.NewSystemError("TransactionService.GetByMerchantID", err, "failed to get transactions")
	}
	if len(transactions) == 0 {
		return []*domain.Transaction{}, nil
	}
	return transactions, nil
}

func (s *TransactionService) UpdateStatus(ctx context.Context, id string, status string) error {
	txID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.transactionRepo.UpdateStatus(ctx, txID, status)
}

func (s *TransactionService) SetPointsService(pointsService domain.PointsService) {
	s.pointsService = pointsService
}
