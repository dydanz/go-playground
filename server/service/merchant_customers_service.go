package service

import (
	"context"
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/util"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type MerchantCustomersService struct {
	customerRepo domain.MerchantCustomersRepository
	logger       zerolog.Logger
}

func NewMerchantCustomersService(customerRepo domain.MerchantCustomersRepository) *MerchantCustomersService {
	return &MerchantCustomersService{customerRepo: customerRepo,
		logger: logging.GetLogger(),
	}
}

func (s *MerchantCustomersService) Create(ctx context.Context, req *domain.CreateMerchantCustomerRequest) (*domain.MerchantCustomer, error) {
	var createErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.Create", func() *domain.MerchantCustomer {
		// Check if customer already exists with email or phone
		existingByEmail, _ := s.customerRepo.GetByEmail(ctx, req.Email)
		if existingByEmail != nil {
			s.logger.Error().
				Str("email", req.Email).
				Msg("Email already exists")
			createErr = domain.NewResourceConflictError("merchant customer", "email already exists")
			return nil
		}

		existingByPhone, _ := s.customerRepo.GetByPhone(ctx, req.Phone)
		if existingByPhone != nil {
			s.logger.Error().
				Str("phone", req.Phone).
				Msg("Phone already exists")
			createErr = domain.NewResourceConflictError("merchant customer", "phone already exists")
			return nil
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error hashing password")
			createErr = domain.NewSystemError("MerchantCustomersService.Create", err, "failed to hash password")
			return nil
		}

		customer := &domain.MerchantCustomer{
			ID:         uuid.New(),
			MerchantID: req.MerchantID,
			Email:      req.Email,
			Password:   string(hashedPassword),
			Name:       req.Name,
			Phone:      req.Phone,
		}

		if err := s.customerRepo.Create(ctx, customer); err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error creating merchant customer")
			createErr = domain.NewSystemError("MerchantCustomersService.Create", err, "failed to create customer")
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(createErr).
			Msg("Error creating merchant customer")
		if createErr == nil {
			createErr = domain.NewSystemError("MerchantCustomersService.Create", nil, "failed to create customer")
		}
		return nil, createErr
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByID(ctx context.Context, id uuid.UUID) (*domain.MerchantCustomer, error) {
	var getErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByID", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByID(ctx, id)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error getting merchant customer, failed to get custome")
			getErr = domain.NewSystemError("MerchantCustomersService.GetByID", err, "failed to get customer")
			return nil
		}
		if customer == nil {
			s.logger.Error().
				Str("id", id.String()).
				Msg("Customer not found")
			getErr = domain.NewResourceNotFoundError("merchant customer", id.String(), "customer not found")
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(getErr).
			Msg("Error getting merchant customer")
		return nil, getErr
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByEmail(ctx context.Context, email string) (*domain.MerchantCustomer, error) {
	var getErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByEmail", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByEmail(ctx, email)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("email", email).
				Msg("Error getting merchant customer")
			getErr = domain.NewSystemError("MerchantCustomersService.GetByEmail", err, "failed to get customer")
			return nil
		}
		if customer == nil {
			s.logger.Error().
				Str("email", email).
				Msg("Customer not found")
			getErr = domain.NewResourceNotFoundError("merchant customer", email, "customer not found")
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(getErr).
			Msg("Error getting merchant customer")
		return nil, getErr
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByPhone(ctx context.Context, phone string) (*domain.MerchantCustomer, error) {
	var getErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByPhone", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByPhone(ctx, phone)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("phone", phone).
				Msg("Error getting merchant customer")
			getErr = domain.NewSystemError("MerchantCustomersService.GetByPhone", err, "failed to get customer")
			return nil
		}
		if customer == nil {
			s.logger.Error().
				Str("phone", phone).
				Msg("Customer not found")
			getErr = domain.NewResourceNotFoundError("merchant customer", phone, "customer not found")
			return nil
		}
		return customer
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(getErr).
			Msg("Error getting merchant customer")
		return nil, getErr
	}
	return result, nil
}

func (s *MerchantCustomersService) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.MerchantCustomer, error) {
	var getErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.GetByMerchantID", func() []*domain.MerchantCustomer {
		customers, err := s.customerRepo.GetByMerchantID(ctx, merchantID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("merchantID", merchantID.String()).
				Msg("Error getting merchant customers, failed to get customers")
			getErr = domain.NewSystemError("MerchantCustomersService.GetByMerchantID", err, "failed to get customers")
			return nil
		}
		if len(customers) == 0 {
			s.logger.Error().
				Str("merchantID", merchantID.String()).
				Msg("No customers found")
			getErr = domain.NewResourceNotFoundError("merchant customers", merchantID.String(), "no customers found")
			return nil
		}
		return customers
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(getErr).
			Msg("Error getting merchant customers")
		return nil, getErr
	}
	return result, nil
}

func (s *MerchantCustomersService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateMerchantCustomerRequest) (*domain.MerchantCustomer, error) {
	var updateErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.Update", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByID(ctx, id)
		if err != nil || customer == nil {
			s.logger.Error().
				Err(err).
				Str("id", id.String()).
				Msg("Error getting merchant customer, failed to get customer")
			updateErr = err
			return nil
		}

		// Check if email is being changed and if it's already taken
		if req.Email != "" && req.Email != customer.Email {
			existingByEmail, _ := s.customerRepo.GetByEmail(ctx, req.Email)
			if existingByEmail != nil {
				s.logger.Error().
					Str("email", req.Email).
					Msg("Email already exists")
				updateErr = domain.NewResourceConflictError("merchant customer", "email already exists")
				return nil
			}
			customer.Email = req.Email
		}

		// Check if phone is being changed and if it's already taken
		if req.Phone != "" && req.Phone != customer.Phone {
			existingByPhone, _ := s.customerRepo.GetByPhone(ctx, req.Phone)
			if existingByPhone != nil {
				s.logger.Error().
					Str("phone", req.Phone).
					Msg("Phone already exists")
				updateErr = domain.NewResourceConflictError("merchant customer", "phone already exists")
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
				s.logger.Error().
					Err(err).
					Msg("Error hashing password")
				updateErr = domain.NewSystemError("MerchantCustomersService.Update", err, "failed to hash password")
				return nil
			}
			customer.Password = string(hashedPassword)
		}

		if err := s.customerRepo.Update(ctx, customer); err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error updating merchant customer")
			updateErr = domain.NewSystemError("MerchantCustomersService.Update", err, "failed to update customer")
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		if updateErr == nil {
			updateErr = domain.NewResourceNotFoundError("merchant customer", id.String(), "customer not found")
		}
		s.logger.Error().
			Err(updateErr).
			Msg("Error updating merchant customer")
		return nil, updateErr
	}
	return result, nil
}

func (s *MerchantCustomersService) ValidateCredentials(ctx context.Context, email, password string) (*domain.MerchantCustomer, error) {
	var authErr error
	decoratedFn := util.ServiceLatencyDecorator("MerchantCustomersService.ValidateCredentials", func() *domain.MerchantCustomer {
		customer, err := s.customerRepo.GetByEmail(ctx, email)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("email", email).
				Msg("Error validating credentials, failed to get customer")
			authErr = domain.NewSystemError("MerchantCustomersService.ValidateCredentials", err, "failed to get customer")
			return nil
		}
		if customer == nil {
			s.logger.Error().
				Str("email", email).
				Msg("Invalid credentials, customer not found")
			authErr = domain.NewAuthenticationError("invalid credentials")
			return nil
		}

		err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("email", email).
				Msg("Invalid credentials, failed to compare password")
			authErr = domain.NewAuthenticationError("invalid credentials")
			return nil
		}

		return customer
	})

	result := decoratedFn()
	if result == nil {
		s.logger.Error().
			Err(authErr).
			Msg("Error validating credentials")
		return nil, authErr
	}
	return result, nil
}
