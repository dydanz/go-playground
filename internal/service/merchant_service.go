package service

import (
	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type MerchantService struct {
	merchantRepo domain.MerchantRepository
}

func NewMerchantService(merchantRepo domain.MerchantRepository) *MerchantService {
	return &MerchantService{merchantRepo: merchantRepo}
}

func (s *MerchantService) Create(req *domain.CreateMerchantRequest) (*domain.Merchant, error) {
	// Validate merchant type
	if !isValidMerchantType(req.Type) {
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
	if err := s.validateMerchantCreation(merchant); err != nil {
		return nil, err
	}

	// Attempt to create the merchant
	merchant, err := s.merchantRepo.Create(merchant)
	if err != nil {
		// Repository layer will return appropriate error types
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) GetByID(id uuid.UUID) (*domain.Merchant, error) {
	merchant, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return nil, err // Repository layer will return appropriate error types
	}

	return merchant, nil
}

func (s *MerchantService) GetAll() ([]*domain.Merchant, error) {
	merchants, err := s.merchantRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(merchants) == 0 {
		// Return empty slice instead of nil for consistent response
		return []*domain.Merchant{}, nil
	}

	return merchants, nil
}

func (s *MerchantService) Update(id uuid.UUID, req *domain.UpdateMerchantRequest) (*domain.Merchant, error) {
	// Check if merchant exists
	merchant, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return nil, err // Repository layer will handle not found error
	}

	// Validate merchant type if provided
	if !isValidMerchantType(req.Type) {
		return nil, domain.NewValidationError("type", "invalid merchant type")
	}

	// Update fields
	merchant.Name = req.Name
	merchant.Type = req.Type

	// Business logic validation
	if err := s.validateMerchantUpdate(merchant); err != nil {
		return nil, err
	}

	// Attempt to update
	if err := s.merchantRepo.Update(merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) Delete(id uuid.UUID) error {
	// Check if merchant exists
	merchant, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return err // Repository layer will handle not found error
	}

	// Business logic validation before deletion
	if err := s.validateMerchantDeletion(merchant); err != nil {
		return err
	}

	return s.merchantRepo.Delete(id)
}

// Helper functions for business logic validation

func (s *MerchantService) validateMerchantCreation(merchant *domain.Merchant) error {
	// Example business rule: Check if user already has maximum allowed merchants
	merchants, err := s.merchantRepo.GetByUserID(merchant.UserID)
	if err != nil {
		return domain.NewSystemError("MerchantService.validateMerchantCreation", err, "failed to check existing merchants")
	}

	if len(merchants) >= 5 { // Example maximum limit
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
