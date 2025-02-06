package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-playground/internal/domain"
	"go-playground/internal/repository/redis"
	"math/big"
	"time"

	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    domain.UserRepository
	authRepo    domain.AuthRepository
	sessionRepo redis.SessionRepository
}

func NewAuthService(userRepo domain.UserRepository, authRepo domain.AuthRepository, sessionRepo redis.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(req *domain.RegistrationRequest) (*domain.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  fmt.Sprintf("Error checking email: %v", err),
		}
	}
	if existingUser != nil {
		return nil, domain.ResourceConflictError{
			Resource: "user",
			Message:  "Email already exists",
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Create user
	createReq := &domain.CreateUserRequest{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
	}

	user, err := s.userRepo.Create(context.Background(), createReq)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	// Generate and store OTP
	otp := s.generateOTP()
	verification := &domain.RegistrationVerification{
		UserID:    user.ID,
		OTP:       otp,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.authRepo.CreateVerification(verification); err != nil {
		return nil, fmt.Errorf("error creating verification: %v", err)
	}

	// TODO: Send OTP via email
	log.Printf("OTP for %s: %s", user.Email, otp)

	return user, nil
}

func (s *AuthService) VerifyRegistration(req *domain.VerificationRequest) error {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return domain.ResourceNotFoundError{
			Resource: "user",
			Message:  fmt.Sprintf("Error getting user: %v", err),
		}
	}
	if user == nil {
		return domain.ResourceNotFoundError{
			Resource: "user",
			Message:  "User not found",
		}
	}

	if user.Status == domain.UserStatusActive {
		return domain.ResourceConflictError{
			Resource: "user",
			Message:  "User already verified",
		}
	}

	verification, err := s.authRepo.GetVerification(user.ID, req.OTP)
	if err != nil {
		return domain.ValidationError{
			Field:   "otp",
			Message: "Invalid or expired OTP",
		}
	}

	// Start transaction
	tx, err := s.authRepo.BeginTx()
	if err != nil {
		return fmt.Errorf("transaction error: %v", err)
	}
	defer tx.Rollback()

	// Mark verification as used
	if err := s.authRepo.MarkVerificationUsedTx(tx, verification.ID); err != nil {
		return fmt.Errorf("error marking verification used: %v", err)
	}

	// Update user status
	user.Status = domain.UserStatusActive
	if err := s.userRepo.UpdateTx(tx, user); err != nil {
		return fmt.Errorf("error updating user status: %v", err)
	}

	return tx.Commit()
}

func (s *AuthService) Login(req *domain.LoginRequest) (*domain.AuthToken, error) {
	// Check login attempts
	attempt, err := s.authRepo.UpdateLoginAttempts(req.Email, true)
	if err != nil {
		return nil, fmt.Errorf("error checking login attempts: %v", err)
	}

	if attempt.LockedUntil.After(time.Now()) {
		return nil, domain.AuthenticationError{
			Message: fmt.Sprintf("Account temporarily locked. Try again after %v", attempt.LockedUntil),
		}
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  "Error finding user",
		}
	}
	if user == nil {
		return nil, domain.AuthenticationError{
			Message: "Invalid credentials",
		}
	}

	if user.Status != domain.UserStatusActive {
		return nil, domain.AuthenticationError{
			Message: "Account not verified",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.AuthenticationError{
			Message: "Invalid credentials",
		}
	}

	// Reset login attempts on successful login
	if _, err := s.authRepo.UpdateLoginAttempts(req.Email, false); err != nil {
		return nil, fmt.Errorf("error resetting login attempts: %v", err)
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store the raw token
	authToken := &domain.AuthToken{
		UserID:    user.ID,
		TokenHash: token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.authRepo.CreateToken(authToken); err != nil {
		return nil, fmt.Errorf("error creating auth token: %v", err)
	}

	// Store session data
	if err := s.sessionRepo.StoreSession(context.Background(), user.ID, token, authToken.ExpiresAt); err != nil {
		return nil, fmt.Errorf("error storing session: %v", err)
	}

	return authToken, nil
}

func (s *AuthService) Logout(userID string, tokenHash string) error {
	if err := s.sessionRepo.DeleteSession(context.Background(), userID); err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}

	if err := s.authRepo.InvalidateToken(userID); err != nil {
		return fmt.Errorf("error invalidating token: %v", err)
	}

	return nil
}

func (s *AuthService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  fmt.Sprintf("Error finding user: %v", err),
		}
	}
	if user == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  "User not found",
		}
	}
	return user, nil
}

func (s *AuthService) GetVerificationByUserID(userID string) (*domain.RegistrationVerification, error) {
	verification, err := s.authRepo.GetLatestVerification(userID)
	if err != nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "verification",
			Message:  fmt.Sprintf("Error finding verification: %v", err),
		}
	}
	return verification, nil
}

func (s *AuthService) GetRandomActiveUser() (*domain.User, error) {
	user, err := s.userRepo.GetRandomActiveUser()
	if err != nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  fmt.Sprintf("Error finding random user: %v", err),
		}
	}
	if user == nil {
		return nil, domain.ResourceNotFoundError{
			Resource: "user",
			Message:  "No active users found",
		}
	}
	return user, nil
}

func (s *AuthService) generateOTP() string {
	const digits = "0123456789"
	result := make([]byte, 6)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "000000" // fallback OTP
		}
		result[i] = digits[num.Int64()]
	}
	return string(result)
}
