package service

import (
	"context"
	"fmt"
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  domain.UserRepository
	cacheRepo domain.CacheRepository
}

func NewUserService(userRepo domain.UserRepository, cacheRepo domain.CacheRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *UserService) Create(req *domain.CreateUserRequest) (*domain.User, error) {
	decoratedFn := util.ServiceLatencyDecorator("UserService.Create", func() *domain.User {
		// Check if email already exists
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err != nil || existingUser != nil {
			return nil
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil
		}

		user := &domain.CreateUserRequest{
			Email:    req.Email,
			Password: string(hashedPassword),
			Name:     req.Name,
			Phone:    req.Phone,
		}

		createdUser, err := s.userRepo.Create(context.Background(), user)
		if err != nil {
			return nil
		}

		// Clear password before returning
		createdUser.Password = ""
		return createdUser
	})

	result := decoratedFn()
	if result == nil {
		return nil, fmt.Errorf("failed to create user")
	}
	return result, nil
}

func (s *UserService) GetByID(id string) (*domain.User, error) {
	decoratedFn := util.ServiceLatencyDecorator("UserService.GetByID", func() *domain.User {
		// Try to get from cache first
		if user, err := s.cacheRepo.GetUser(id); err == nil && user != nil {
			user.Status = user.Status
			return user
		}

		// If not in cache, get from database
		user, err := s.userRepo.GetByID(id)
		if err != nil {
			return nil
		}

		// Convert status to string representation
		user.Status = user.Status

		// Store in cache for future requests
		if err := s.cacheRepo.SetUser(user); err != nil {
			log.Printf("Failed to cache user: %v", err)
		}

		return user
	})

	result := decoratedFn()
	if result == nil {
		return nil, fmt.Errorf("user not found")
	}
	return result, nil
}

func (s *UserService) GetAll() ([]domain.User, error) {
	decoratedFn := util.ServiceLatencyDecorator("UserService.GetAll", func() []domain.User {
		usersPtr, err := s.userRepo.GetAll()
		if err != nil {
			return nil
		}

		users := make([]domain.User, len(usersPtr))
		for i, u := range usersPtr {
			users[i] = *u
			users[i].Password = ""
		}

		return users
	})

	result := decoratedFn()
	if result == nil {
		return nil, fmt.Errorf("failed to get users")
	}
	return result, nil
}

func (s *UserService) Update(id string, req *domain.UpdateUserRequest) (*domain.User, error) {
	decoratedFn := util.ServiceLatencyDecorator("UserService.Update", func() *domain.User {
		// Get existing user
		user, err := s.userRepo.GetByID(id)
		if err != nil || user == nil {
			return nil
		}

		// Update fields
		user.Name = req.Name
		user.Phone = req.Phone
		user.UpdatedAt = time.Now()

		if err := s.userRepo.Update(user); err != nil {
			return nil
		}

		// Clear password before returning
		user.Password = ""
		return user
	})

	result := decoratedFn()
	if result == nil {
		return nil, fmt.Errorf("failed to update user")
	}
	return result, nil
}

func (s *UserService) Delete(id string) error {
	decoratedFn := util.ServiceLatencyDecorator("UserService.Delete", func() bool {
		// Get existing user
		user, err := s.userRepo.GetByID(id)
		if err != nil || user == nil {
			return false
		}

		if err := s.userRepo.Delete(id); err != nil {
			return false
		}
		return true
	})

	if !decoratedFn() {
		return fmt.Errorf("failed to delete user")
	}
	return nil
}
