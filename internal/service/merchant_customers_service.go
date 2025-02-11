package service

import (
	"context"
	"errors"
	"go-playground/internal/domain"
	"go-playground/internal/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MerchantCustomersService struct {
	customerRepo domain.MerchantCustomersRepository
}

func NewMerchantCustomersService(customerRepo domain.MerchantCustomersRepository) *MerchantCustomersService {
	return &MerchantCustomersService{customerRepo: customerRepo}
}

func (s *MerchantCustomersService) Create(ctx context.Context, req *domain.CreateMerchantCustomerRequest) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.Create", func() *domain.MerchantCustomer {
		// Check if customer already exists with email or phone
		existingByEmail, _ := s.customerRepo.GetByEmail(ctx, req.Email)
		if existingByEmail != nil {
			return nil
		}

		existingByPhone, _ := s.customerRepo.GetByPhone(ctx, req.Phone)
		if existingByPhone != nil {
			return nil
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil
		}

		customer := &domain.MerchantCustomer{
			MerchantID: req.MerchantID,
			Email:      req.Email,
			Password:   string(hashedPassword),
			Name:       req.Name,
			Phone:      req.Phone,
		}

		if err := s.customerRepo.Create(ctx, customer); err != nil {
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.ResourceConflictError{
			Resource: "merchant customer",
			Message:  "customer already exists or data input is invalid",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByID(ctx context.Context, id uuid.UUID) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByID", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByID(ctx, id)
		if err != nil {
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "merchant customer",
			Message:  "merchant customer not found",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByEmail(ctx context.Context, email string) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByEmail", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "merchant customer",
			Message:  "merchant customer not found",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByPhone(ctx context.Context, phone string) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByPhone", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByPhone(ctx, phone)
		if err != nil {
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "merchant customer",
			Message:  "merchant customer not found",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByMerchantID", func() []*domain.MerchantCustomer {
		customers, err := s.customerRepo.GetByMerchantID(ctx, merchantID)
		if err != nil {
			return nil
		}
		return customers
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "merchant customers",
			Message:  "merchant customers not found",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateMerchantCustomerRequest) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.Update", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByID(ctx, id)
		if err != nil || customer == nil {
			return nil
		}

		// Check if email is being changed and if it's already taken
		if req.Email != "" && req.Email != customer.Email {
			existingByEmail, _ := s.customerRepo.GetByEmail(ctx, req.Email)
			if existingByEmail != nil {
				return nil
			}
			customer.Email = req.Email
		}

		// Check if phone is being changed and if it's already taken
		if req.Phone != "" && req.Phone != customer.Phone {
			existingByPhone, _ := s.customerRepo.GetByPhone(ctx, req.Phone)
			if existingByPhone != nil {
				return nil
			}
			customer.Phone = req.Phone
		}

		// Update other fields if provided
		if req.Name != "" {
			customer.Name = req.Name
		}

		// Update password if provided
		if req.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil
			}
			customer.Password = string(hashedPassword)
		}

		if err := s.customerRepo.Update(ctx, customer); err != nil {
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, domain.InvalidInputError{
			Message: "failed to update merchant customer",
		}
	}
	return result, nil
}

func (s *MerchantCustomersService) Delete(ctx context.Context, id uuid.UUID) error {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.Delete", func() bool {
		if err := s.customerRepo.Delete(ctx, id); err != nil {
			return false
		}
		return true
	})

	if !decoratedFn() {
		return domain.InvalidInputError{
			Message: "failed to delete merchant customer",
		}
	}
	return nil
}

func (s *MerchantCustomersService) ValidateCredentials(ctx context.Context, email, password string) (*domain.MerchantCustomer, error) {
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.ValidateCredentials", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil
		}

		err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
		if err != nil {
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		return nil, errors.New("invalid credentials")
	}
	return result, nil
}
