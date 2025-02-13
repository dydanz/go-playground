package service

import (
	"context"
	"database/sql"
	"go-playground/internal/domain"
	"go-playground/internal/repository/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// Mock repositories
type mockUserRepository struct {
	mock.Mock
}

type mockAuthRepository struct {
	mock.Mock
}

type mockSessionRepository struct {
	mock.Mock
}

// Implement UserRepository interface
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) UpdateTx(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *mockUserRepository) GetRandomActiveUser(ctx context.Context) (*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Implement AuthRepository interface
func (m *mockAuthRepository) CreateVerification(ctx context.Context, verification *domain.RegistrationVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

func (m *mockAuthRepository) GetVerification(ctx context.Context, userID string, otp string) (*domain.RegistrationVerification, error) {
	args := m.Called(ctx, userID, otp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *mockAuthRepository) GetLatestVerification(ctx context.Context, userID string) (*domain.RegistrationVerification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *mockAuthRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *mockAuthRepository) Commit(ctx context.Context, tx *sql.Tx) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *mockAuthRepository) MarkVerificationUsedTx(ctx context.Context, tx *sql.Tx, verificationID string) error {
	args := m.Called(ctx, tx, verificationID)
	return args.Error(0)
}

func (m *mockAuthRepository) UpdateLoginAttempts(ctx context.Context, email string, increment bool) (*domain.LoginAttempt, error) {
	args := m.Called(ctx, email, increment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginAttempt), args.Error(1)
}

func (m *mockAuthRepository) CreateToken(ctx context.Context, token *domain.AuthToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockAuthRepository) InvalidateToken(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Implement SessionRepository interface
func (m *mockSessionRepository) StoreSession(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *mockSessionRepository) GetSession(ctx context.Context, userID string) (*redis.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*redis.Session), args.Error(1)
}

func (m *mockSessionRepository) DeleteSession(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockSessionRepository) RefreshSession(ctx context.Context, userID string, newToken string, expiration time.Duration) error {
	args := m.Called(ctx, userID, newToken, expiration)
	return args.Error(0)
}

func (m *mockSessionRepository) DeleteAllSession(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// AuthServiceTestSuite defines the test suite
type AuthServiceTestSuite struct {
	suite.Suite
	userRepo    *mockUserRepository
	authRepo    *mockAuthRepository
	sessionRepo *mockSessionRepository
	authService *AuthService
}

// SetupTest is called before each test
func (s *AuthServiceTestSuite) SetupTest() {
	s.userRepo = new(mockUserRepository)
	s.authRepo = new(mockAuthRepository)
	s.sessionRepo = new(mockSessionRepository)
	s.authService = NewAuthService(s.userRepo, s.authRepo, s.sessionRepo)
}

// TestAuthServiceTestSuite runs the test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// Test cases for Register
func (s *AuthServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	req := &domain.RegistrationRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Phone:    "1234567890",
	}

	expectedUser := &domain.User{
		ID:    "user123",
		Email: req.Email,
		Name:  req.Name,
		Phone: req.Phone,
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)
	s.userRepo.On("Create", ctx, mock.AnythingOfType("*domain.CreateUserRequest")).Return(expectedUser, nil)
	s.authRepo.On("CreateVerification", ctx, mock.AnythingOfType("*domain.RegistrationVerification")).Return(nil)

	user, err := s.authService.Register(ctx, req)

	s.NoError(err)
	s.NotNil(user)
	s.Equal(expectedUser.Email, user.Email)
	s.Equal(expectedUser.Name, user.Name)
}

func (s *AuthServiceTestSuite) TestRegister_EmailExists() {
	ctx := context.Background()
	req := &domain.RegistrationRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
		Phone:    "1234567890",
	}

	existingUser := &domain.User{
		ID:    "user123",
		Email: req.Email,
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	user, err := s.authService.Register(ctx, req)

	s.Error(err)
	s.Nil(user)
	s.IsType(domain.ResourceConflictError{}, err)
}

func (s *AuthServiceTestSuite) TestRegister_InvalidInput() {
	ctx := context.Background()
	testCases := []struct {
		name    string
		req     *domain.RegistrationRequest
		wantErr bool
	}{
		{
			name: "Empty Email",
			req: &domain.RegistrationRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
				Phone:    "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Invalid Email Format",
			req: &domain.RegistrationRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
				Phone:    "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Unicode Characters",
			req: &domain.RegistrationRequest{
				Email:    "test@例子.com",
				Password: "password123",
				Name:     "测试用户",
				Phone:    "1234567890",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.userRepo.On("GetByEmail", ctx, tc.req.Email).Return(nil, nil).Maybe()
			s.userRepo.On("Create", ctx, mock.AnythingOfType("*domain.CreateUserRequest")).Return(&domain.User{}, nil).Maybe()
			s.authRepo.On("CreateVerification", ctx, mock.AnythingOfType("*domain.RegistrationVerification")).Return(nil).Maybe()

			user, err := s.authService.Register(ctx, tc.req)

			if tc.wantErr {
				s.Error(err)
				s.Nil(user)
			} else {
				s.NoError(err)
				s.NotNil(user)
			}
		})
	}
}

// Test cases for Login
func (s *AuthServiceTestSuite) TestLogin_Success() {
	ctx := context.Background()
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create a real bcrypt hash for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       "user123",
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     "Test User",
		Status:   domain.UserStatusActive,
	}

	loginAttempt := &domain.LoginAttempt{
		Email:        req.Email,
		AttemptCount: 0,
		LastAttempt:  time.Time{},
		LockedUntil:  time.Time{},
	}

	// Set up mock expectations in the correct order
	s.authRepo.On("UpdateLoginAttempts", ctx, req.Email, true).Return(loginAttempt, nil)
	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	s.authRepo.On("UpdateLoginAttempts", ctx, req.Email, false).Return(loginAttempt, nil)
	s.authRepo.On("CreateToken", ctx, mock.AnythingOfType("*domain.AuthToken")).Return(nil)
	s.sessionRepo.On("StoreSession", ctx, user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	token, err := s.authService.Login(ctx, req)

	s.NoError(err)
	s.NotNil(token)
	s.Equal(user.ID, token.UserID)
	s.Equal(user.Name, token.UserName)

	// Verify all mocks were called
	s.authRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_InvalidCredentials() {
	ctx := context.Background()
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	loginAttempt := &domain.LoginAttempt{
		Email:        req.Email,
		AttemptCount: 0,
		LastAttempt:  time.Time{},
		LockedUntil:  time.Time{},
	}

	s.authRepo.On("UpdateLoginAttempts", ctx, req.Email, true).Return(loginAttempt, nil)
	s.userRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)

	token, err := s.authService.Login(ctx, req)

	s.Error(err)
	s.Nil(token)
	s.IsType(domain.AuthenticationError{}, err)
}

func (s *AuthServiceTestSuite) TestLogin_AccountLocked() {
	ctx := context.Background()
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginAttempt := &domain.LoginAttempt{
		Email:        req.Email,
		AttemptCount: 5,
		LastAttempt:  time.Now(),
		LockedUntil:  time.Now().Add(15 * time.Minute),
	}

	s.authRepo.On("UpdateLoginAttempts", ctx, req.Email, true).Return(loginAttempt, nil)

	token, err := s.authService.Login(ctx, req)

	s.Error(err)
	s.Nil(token)
	s.IsType(domain.AuthenticationError{}, err)
}

// Test cases for Logout
func (s *AuthServiceTestSuite) TestLogout_Success() {
	ctx := context.Background()
	userID := "user123"
	tokenHash := "token123"

	s.sessionRepo.On("DeleteSession", ctx, userID).Return(nil)
	s.authRepo.On("InvalidateToken", ctx, userID).Return(nil)

	err := s.authService.Logout(ctx, userID, tokenHash)

	s.NoError(err)
}

func (s *AuthServiceTestSuite) TestLogout_Error() {
	ctx := context.Background()
	userID := "user123"
	tokenHash := "token123"

	expectedErr := domain.AuthenticationError{Message: "session not found"}
	s.sessionRepo.On("DeleteSession", ctx, userID).Return(expectedErr)

	err := s.authService.Logout(ctx, userID, tokenHash)

	s.Error(err)
}

// Test cases for VerifyRegistration
func (s *AuthServiceTestSuite) TestVerifyRegistration_Success() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	user := &domain.User{
		ID:     "user123",
		Email:  req.Email,
		Status: domain.UserStatusPending,
	}

	verification := &domain.RegistrationVerification{
		ID:        "ver123",
		UserID:    user.ID,
		OTP:       req.OTP,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockTx := &sql.Tx{}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	s.authRepo.On("GetVerification", ctx, user.ID, req.OTP).Return(verification, nil)
	s.authRepo.On("BeginTx", ctx).Return(mockTx, nil)
	s.authRepo.On("MarkVerificationUsedTx", ctx, mockTx, verification.ID).Return(nil)
	s.userRepo.On("UpdateTx", ctx, mockTx, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == user.ID && u.Status == domain.UserStatusActive
	})).Return(nil)
	s.authRepo.On("Commit", ctx, mockTx).Return(nil)

	err := s.authService.VerifyRegistration(ctx, req)

	s.NoError(err)
	s.authRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestVerifyRegistration_UserNotFound() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "nonexistent@example.com",
		OTP:   "123456",
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)

	err := s.authService.VerifyRegistration(ctx, req)

	s.Error(err)
	s.IsType(domain.ResourceNotFoundError{}, err)
}

func (s *AuthServiceTestSuite) TestVerifyRegistration_AlreadyVerified() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	user := &domain.User{
		ID:     "user123",
		Email:  req.Email,
		Status: domain.UserStatusActive,
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	err := s.authService.VerifyRegistration(ctx, req)

	s.Error(err)
	s.IsType(domain.ResourceConflictError{}, err)
}

func (s *AuthServiceTestSuite) TestVerifyRegistration_InvalidOTP() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	user := &domain.User{
		ID:     "user123",
		Email:  req.Email,
		Status: domain.UserStatusPending,
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	s.authRepo.On("GetVerification", ctx, user.ID, req.OTP).Return(nil, domain.ValidationError{
		Field:   "otp",
		Message: "Invalid or expired OTP",
	})

	err := s.authService.VerifyRegistration(ctx, req)

	s.Error(err)
	s.IsType(domain.ValidationError{}, err)
}

func (s *AuthServiceTestSuite) TestVerifyRegistration_ExpiredOTP() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	user := &domain.User{
		ID:     "user123",
		Email:  req.Email,
		Status: domain.UserStatusPending,
	}

	verification := &domain.RegistrationVerification{
		ID:        "ver123",
		UserID:    user.ID,
		OTP:       req.OTP,
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	s.authRepo.On("GetVerification", ctx, user.ID, req.OTP).Return(verification, domain.ValidationError{
		Field:   "otp",
		Message: "OTP has expired",
	})

	err := s.authService.VerifyRegistration(ctx, req)

	s.Error(err)
	s.IsType(domain.ValidationError{}, err)
}

// Test cases for GetUserByEmail
func (s *AuthServiceTestSuite) TestGetUserByEmail_Success() {
	ctx := context.Background()
	email := "test@example.com"

	expectedUser := &domain.User{
		ID:    "user123",
		Email: email,
		Name:  "Test User",
	}

	s.userRepo.On("GetByEmail", ctx, email).Return(expectedUser, nil)

	user, err := s.authService.GetUserByEmail(ctx, email)

	s.NoError(err)
	s.NotNil(user)
	s.Equal(expectedUser.ID, user.ID)
	s.Equal(expectedUser.Email, user.Email)
}

func (s *AuthServiceTestSuite) TestGetUserByEmail_NotFound() {
	ctx := context.Background()
	email := "nonexistent@example.com"

	s.userRepo.On("GetByEmail", ctx, email).Return(nil, nil)

	user, err := s.authService.GetUserByEmail(ctx, email)

	s.Error(err)
	s.Nil(user)
	s.IsType(domain.ResourceNotFoundError{}, err)
}

// Test cases for GetRandomActiveUser
func (s *AuthServiceTestSuite) TestGetRandomActiveUser_Success() {
	ctx := context.Background()

	expectedUser := &domain.User{
		ID:     "user123",
		Email:  "test@example.com",
		Name:   "Test User",
		Status: domain.UserStatusActive,
	}

	s.userRepo.On("GetRandomActiveUser", ctx).Return(expectedUser, nil)

	user, err := s.authService.GetRandomActiveUser(ctx)

	s.NoError(err)
	s.NotNil(user)
	s.Equal(expectedUser.ID, user.ID)
	s.Equal(domain.UserStatusActive, user.Status)
}

func (s *AuthServiceTestSuite) TestGetRandomActiveUser_NoActiveUsers() {
	ctx := context.Background()

	s.userRepo.On("GetRandomActiveUser", ctx).Return(nil, nil)

	user, err := s.authService.GetRandomActiveUser(ctx)

	s.Error(err)
	s.Nil(user)
	s.IsType(domain.ResourceNotFoundError{}, err)
}

// Test cases for edge cases and error handling
func (s *AuthServiceTestSuite) TestLogin_DatabaseError() {
	ctx := context.Background()
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedErr := &domain.SystemError{
		Op:      "UpdateLoginAttempts",
		Message: "database connection error",
	}

	s.authRepo.On("UpdateLoginAttempts", ctx, req.Email, true).Return(nil, expectedErr)

	token, err := s.authService.Login(ctx, req)

	s.Error(err)
	s.Nil(token)
	s.IsType(domain.AuthenticationError{}, err)
	s.Contains(err.Error(), expectedErr.Message)
}

func (s *AuthServiceTestSuite) TestRegister_DatabaseError() {
	ctx := context.Background()
	req := &domain.RegistrationRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Phone:    "1234567890",
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(nil, domain.SystemError{
		Op:      "GetByEmail",
		Message: "database connection error",
	})

	user, err := s.authService.Register(ctx, req)

	s.Error(err)
	s.Nil(user)
	s.IsType(domain.ResourceNotFoundError{}, err)
}

func (s *AuthServiceTestSuite) TestVerifyRegistration_TransactionError() {
	ctx := context.Background()
	req := &domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	user := &domain.User{
		ID:     "user123",
		Email:  req.Email,
		Status: domain.UserStatusPending,
	}

	verification := &domain.RegistrationVerification{
		ID:        "ver123",
		UserID:    user.ID,
		OTP:       req.OTP,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	s.userRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	s.authRepo.On("GetVerification", ctx, user.ID, req.OTP).Return(verification, nil)
	s.authRepo.On("BeginTx", ctx).Return(nil, domain.SystemError{
		Op:      "BeginTx",
		Message: "transaction error",
	})

	err := s.authService.VerifyRegistration(ctx, req)

	s.Error(err)
	s.IsType(domain.SystemError{}, err)
}
