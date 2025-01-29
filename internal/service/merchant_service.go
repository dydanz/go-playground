package service

import (
	"go-playground/internal/domain"
)

type MerchantService struct {
	merchantRepo domain.MerchantRepository
}

func NewMerchantService(merchantRepo domain.MerchantRepository) *MerchantService {
	return &MerchantService{merchantRepo: merchantRepo}
}

func (s *MerchantService) Create(req *domain.CreateMerchantRequest) (*domain.Merchant, error) {
	merchant := &domain.Merchant{
		Name: req.Name,
		Type: req.Type,
	}

	if err := s.merchantRepo.Create(merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) GetByID(id string) (*domain.Merchant, error) {
	return s.merchantRepo.GetByID(id)
}

func (s *MerchantService) GetAll() ([]*domain.Merchant, error) {
	return s.merchantRepo.GetAll()
}

func (s *MerchantService) Update(id string, req *domain.UpdateMerchantRequest) (*domain.Merchant, error) {
	merchant, err := s.merchantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if merchant == nil {
		return nil, nil
	}

	merchant.Name = req.Name
	merchant.Type = req.Type

	if err := s.merchantRepo.Update(merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (s *MerchantService) Delete(id string) error {
	return s.merchantRepo.Delete(id)
}