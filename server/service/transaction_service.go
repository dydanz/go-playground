package service

import (
	"context"

	"go-playground/pkg/logging"
	"go-playground/server/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TransactionService struct {
	transactionRepo      domain.TransactionRepository
	pointsService        domain.PointsService
	eventLoggerService   domain.EventLoggerService
	merchantCustomerRepo domain.MerchantCustomersRepository
	logger               zerolog.Logger
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	pointsService domain.PointsService,
	eventLoggerService domain.EventLoggerService,
	merchantCustomerRepo domain.MerchantCustomersRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo:      transactionRepo,
		pointsService:        pointsService,
		eventLoggerService:   eventLoggerService,
		merchantCustomerRepo: merchantCustomerRepo,
		logger:               logging.GetLogger(),
	}
}

func (s *TransactionService) getMerchantIDByCustomerID(ctx context.Context, customerID uuid.UUID) (uuid.UUID, error) {
	customer, err := s.merchantCustomerRepo.GetByID(ctx, customerID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant customer")
		return uuid.Nil, domain.NewSystemError("TransactionService.getMerchantIDByCustomerID", err, "failed to get merchant customer")
	}
	if customer == nil {
		s.logger.Error().
			Msg("Customer not found")
		return uuid.Nil, domain.NewResourceNotFoundError("merchant customer", customerID.String(), "customer not found")
	}
	return customer.MerchantID, nil
}

func (s *TransactionService) Create(ctx context.Context, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	if req.TransactionAmount <= 0 {
		s.logger.Error().
			Msg("Transaction amount must be greater than 0")
		return nil, domain.NewValidationError("transaction_amount", "transaction amount must be greater than 0")
	}

	// Get merchant ID from customer ID
	merchantID, err := s.getMerchantIDByCustomerID(ctx, req.MerchantCustomersID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant ID")
		return nil, domain.NewSystemError("TransactionService.Create", err, "failed to get merchant ID")
	}

	transaction := &domain.Transaction{
		TransactionID:       uuid.New(),
		MerchantCustomersID: req.MerchantCustomersID,
		MerchantID:          merchantID,
		ProgramID:           req.ProgramID,
		TransactionType:     req.TransactionType,
		TransactionAmount:   req.TransactionAmount,
		TransactionDate:     req.TransactionDate,
	}

	createdTx, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error creating transaction")
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
			s.logger.Error().
				Err(err).
				Msg("Error earning points")
			return nil, domain.NewSystemError("TransactionService.Create", err, "failed to earn points")
		}
	} else if points < 0 {
		if _, err := s.pointsService.RedeemPoints(ctx, &domain.PointsTransaction{
			CustomerID:    transaction.MerchantCustomersID.String(),
			ProgramID:     transaction.ProgramID.String(),
			Points:        points,
			TransactionID: createdTx.TransactionID.String(),
		}); err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error redeeming points")
			return nil, domain.NewSystemError("TransactionService.Create", err, "failed to redeem points")
		}
	}

	// Log the transaction event
	go s.eventLoggerService.SaveTransactionEvents(ctx, domain.TransactionCreated, createdTx, points)

	return createdTx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting transaction")
		return nil, domain.NewSystemError("TransactionService.GetByID", err, "failed to get transaction")
	}
	if transaction == nil {
		s.logger.Error().
			Msg("Transaction not found")
		return nil, domain.NewResourceNotFoundError("transaction", id.String(), "transaction not found")
	}
	return transaction, nil
}

func (s *TransactionService) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting transactions")
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

func (s *TransactionService) GetByMerchantIDWithPagination(ctx context.Context, merchantID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	return s.transactionRepo.GetByMerchantIDWithPagination(ctx, merchantID, offset, limit)
}

func (s *TransactionService) GetByUserIDWithPagination(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	return s.transactionRepo.GetByUserIDWithPagination(ctx, userID, offset, limit)
}
func (s *TransactionService) UpdateStatus(ctx context.Context, id string, status string) error {
	txID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error parsing transaction ID")
		return domain.NewSystemError("TransactionService.UpdateStatus", err, "failed to parse transaction ID")
	}
	return s.transactionRepo.UpdateStatus(ctx, txID, status)
}

func (s *TransactionService) SetPointsService(pointsService domain.PointsService) {
	s.pointsService = pointsService
}
