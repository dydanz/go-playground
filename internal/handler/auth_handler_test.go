package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-playground/internal/domain"
	"go-playground/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *domain.RegistrationRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthToken, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthToken), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, userID string, tokenHash string) error {
	args := m.Called(ctx, userID, tokenHash)
	return args.Error(0)
}

func (m *MockAuthService) VerifyRegistration(ctx context.Context, req *domain.VerificationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) GetVerificationByUserID(ctx context.Context, userID string) (*domain.RegistrationVerification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegistrationVerification), args.Error(1)
}

func (m *MockAuthService) GetRandomActiveUser(ctx context.Context) (*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// AuthHandlerTestSuite defines the test suite
type AuthHandlerTestSuite struct {
	suite.Suite
	mockAuthService *MockAuthService
	handler         *handler.AuthHandler
	router          *gin.Engine
}

// SetupTest is called before each test
func (s *AuthHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockAuthService = new(MockAuthService)
	s.handler = handler.NewAuthHandler(s.mockAuthService)
	s.router = gin.New()

	// Setup routes
	s.router.POST("/auth/register", s.handler.Register)
	s.router.POST("/auth/verify", s.handler.Verify)
	s.router.POST("/auth/login", s.handler.Login)
	s.router.POST("/auth/logout", s.handler.Logout)
}

// TestAuthHandlerTestSuite runs the test suite
func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

// Test cases for Register
func (s *AuthHandlerTestSuite) TestRegister_Success() {
	req := domain.RegistrationRequest{
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

	s.mockAuthService.On("Register", mock.Anything, &req).Return(expectedUser, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)

	var response domain.User
	s.NoError(json.Unmarshal(w.Body.Bytes(), &response))
	s.Equal(expectedUser.Email, response.Email)
	s.Equal(expectedUser.Name, response.Name)
}

func (s *AuthHandlerTestSuite) TestRegister_InvalidInput() {
	testCases := []struct {
		name    string
		request domain.RegistrationRequest
		code    int
	}{
		{
			name: "Empty Email",
			request: domain.RegistrationRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			code: http.StatusBadRequest,
		},
		{
			name: "Invalid Email Format",
			request: domain.RegistrationRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
			},
			code: http.StatusBadRequest,
		},
		{
			name: "Empty Password",
			request: domain.RegistrationRequest{
				Email:    "test@example.com",
				Password: "",
				Name:     "Test User",
			},
			code: http.StatusBadRequest,
		},
		{
			name: "Unicode Characters",
			request: domain.RegistrationRequest{
				Email:    "test@例子.com",
				Password: "password123",
				Name:     "测试用户",
				Phone:    "1234567890",
			},
			code: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.code == http.StatusCreated {
				expectedUser := &domain.User{
					ID:    "user123",
					Email: tc.request.Email,
					Name:  tc.request.Name,
					Phone: tc.request.Phone,
				}
				s.mockAuthService.On("Register", mock.Anything, &tc.request).Return(expectedUser, nil)
			}

			body, _ := json.Marshal(tc.request)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")

			s.router.ServeHTTP(w, r)
			s.Equal(tc.code, w.Code)
		})
	}
}

// Test cases for Login
func (s *AuthHandlerTestSuite) TestLogin_Success() {
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedToken := &domain.AuthToken{
		TokenHash: "token123",
		UserID:    "user123",
		UserName:  "Test User",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mockAuthService.On("Login", mock.Anything, &req).Return(expectedToken, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)

	var response domain.LoginResponse
	s.NoError(json.Unmarshal(w.Body.Bytes(), &response))
	s.Equal(expectedToken.TokenHash, response.Token)
	s.Equal(expectedToken.UserID, response.UserID)
}

func (s *AuthHandlerTestSuite) TestLogin_InvalidCredentials() {
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	s.mockAuthService.On("Login", mock.Anything, &req).Return(nil, domain.AuthenticationError{Message: "invalid credentials"})

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusUnauthorized, w.Code)
}

// Test cases for Logout
func (s *AuthHandlerTestSuite) TestLogout_Success() {
	userID := "user123"
	tokenHash := "token123"

	s.mockAuthService.On("Logout", mock.Anything, userID, tokenHash).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	r.Header.Set("Authorization", "Bearer "+tokenHash)

	// Set user_id in context
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID)
	c.Request = r

	s.handler.Logout(c)

	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthHandlerTestSuite) TestLogout_Unauthorized() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusUnauthorized, w.Code)
}

// Test cases for Verify
func (s *AuthHandlerTestSuite) TestVerify_Success() {
	req := domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	s.mockAuthService.On("VerifyRegistration", mock.Anything, &req).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthHandlerTestSuite) TestVerify_InvalidOTP() {
	req := domain.VerificationRequest{
		Email: "test@example.com",
		OTP:   "invalid",
	}

	s.mockAuthService.On("VerifyRegistration", mock.Anything, &req).Return(domain.ValidationError{Message: "invalid OTP"})

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusBadRequest, w.Code)
}
