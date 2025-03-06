package service

import (
	"context"
	"go-playground/pkg/logging"
	"go-playground/server/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MerchantService struct {
	merchantRepo domain.MerchantRepository
	logger       zerolog.Logger
}

func NewMerchantService(merchantRepo domain.MerchantRepository) *MerchantService {
	return &MerchantService{merchantRepo: merchantRepo,
		logger: logging.GetLogger(),
	}
}

func (s *MerchantService) Create(ctx context.Context, req *domain.CreateMerchantRequest) (*domain.Merchant, error) {
	// Validate merchant type
	if !isValidMerchantType(req.Type) {
		s.logger.Error().
			Str("type", string(req.Type)).
			Msg("Invalid merchant type")
		return nil, domain.NewValidationError("type", "invalid merchant type")
	}

	// Create merchant entity
	merchant := &domain.Merchant{
		ID:     uuid.New(),
		UserID: req.UserID,
		Name:   req.Name,
		Type:   req.Type,
	}

	// Business logic validation
	if err := s.validateMerchantCreation(ctx, merchant); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error validating merchant creation")
		return nil, err
	}

	// Attempt to create the merchant
	merchant, err := s.merchantRepo.Create(ctx, merchant)
	if err != nil {
		// Repository layer will return appropriate error types
		s.logger.Error().
			Err(err).
			Msg("Error creating merchant")
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	merchant, err := s.merchantRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant")
		return nil, err // Repository layer will return appropriate error types
	}

	return merchant, nil
}

func (s *MerchantService) GetAll(ctx context.Context, userID uuid.UUID) ([]*domain.MerchantList, error) {
	merchants, err := s.merchantRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(merchants) == 0 {
		// Return empty slice instead of nil for consistent response
		return []*domain.MerchantList{}, nil
	}

	return merchants, nil
}

func (s *MerchantService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateMerchantRequest) (*domain.Merchant, error) {
	// Check if merchant exists
	merchant, err := s.merchantRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant")
		return nil, err // Repository layer will handle not found error
	}

	// Validate merchant type if provided
	if !isValidMerchantType(req.Type) {
		s.logger.Error().
			Str("type", string(req.Type)).
			Msg("Invalid merchant type")
		return nil, domain.NewValidationError("type", "invalid merchant type")
	}

	// Update fields
	merchant.Name = req.Name
	merchant.Type = req.Type

	// Business logic validation
	if err := s.validateMerchantUpdate(merchant); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error validating merchant update")
		return nil, err
	}

	// Attempt to update
	if err := s.merchantRepo.Update(ctx, merchant); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error updating merchant")
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if merchant exists
	merchant, err := s.merchantRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant")
		return err // Repository layer will handle not found error
	}

	// Business logic validation before deletion
	if err := s.validateMerchantDeletion(merchant); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error validating merchant deletion")
		return err
	}

	return s.merchantRepo.Delete(ctx, id)
}

func (s *MerchantService) GetMerchantsByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Merchant, int, error) {
	merchants, total, err := s.merchantRepo.GetMerchantsByUserID(ctx, userID, offset, limit)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchants")
		return nil, 0, domain.NewSystemError("MerchantService.GetMerchantsByUserID", err, "failed to get merchants")
	}
	if len(merchants) == 0 {
		return []*domain.Merchant{}, 0, nil
	}
	return merchants, total, nil
}

// Helper functions for business logic validation

func (s *MerchantService) validateMerchantCreation(ctx context.Context, merchant *domain.Merchant) error {
	// Example business rule: Check if user already has maximum allowed merchants
	merchants, _, err := s.merchantRepo.GetMerchantsByUserID(ctx, merchant.UserID, 1, 50)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchants")
		return domain.NewSystemError("MerchantService.validateMerchantCreation", err, "failed to check existing merchants")
	}

	if len(merchants) >= 50 { // Example maximum limit
		return domain.NewBusinessLogicError(
			"MAX_MERCHANTS_EXCEEDED",
			"user has reached the maximum number of allowed merchants",
		)
	}

	return nil
}

func (s *MerchantService) validateMerchantUpdate(merchant *domain.Merchant) error {
	// Add any business logic validation for updates
	// Example: Check if merchant name is unique within user's merchants
	return nil
}

func (s *MerchantService) validateMerchantDeletion(merchant *domain.Merchant) error {
	// Add any business logic validation for deletion
	// Example: Check if merchant has active programs or transactions
	return nil
}

func isValidMerchantType(t domain.MerchantType) bool {
	switch t {
	case domain.MerchantTypeBank,
		domain.MerchantTypeEcommerce,
		domain.MerchantTypeRepairShop:
		return true
	default:
		return false
	}
}
