package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/middleware"
	"go-playground/server/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authService domain.AuthService
	logger      zerolog.Logger
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logging.GetLogger(),
	}
}

// @Summary Register new user
// @Description Register a new user and send verification OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegistrationRequest true "Registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming registration request")

	var req domain.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind registration request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to register user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("User registered successfully")

	c.JSON(http.StatusCreated, user)
}

// @Summary Verify user registration
// @Description Verify user registration using OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.VerificationRequest true "Verification details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify [post]
func (h *AuthHandler) Verify(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming verification request")

	var req domain.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind verification request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	if err := h.authService.VerifyRegistration(c.Request.Context(), &req); err != nil {
		h.logger.Error().
			Err(err).
			Str("email", req.Email).
			Str("otp", req.OTP).
			Msg("Failed to verify user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Msg("User verified successfully")

	c.JSON(http.StatusOK, gin.H{"message": "verification successful"})
}

// @Summary User login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming login request")

	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind login request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	authToken, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to login user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("email", req.Email).
		Msg("User logged in successfully")

	// Set secure cookie with session token and user ID
	middleware.SetSecureCookie(c, authToken.TokenHash, authToken.UserID, authToken.UserName)

	response := domain.LoginResponse{
		Token:     authToken.TokenHash,
		ExpiresAt: authToken.ExpiresAt,
		UserID:    authToken.UserID,
		UserName:  authToken.UserName,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary User logout
// @Description Logout user and invalidate their token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security User-ID
// @Param User-ID header string true "User ID for authentication"
// @Success 200 {object} map[string]string "message: logged out successfully"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Failure 500 {object} map[string]string "error: internal server error"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming logout request")

	userID, exists := c.Get("user_id")
	if !exists {
		util.HandleError(c, domain.AuthenticationError{Message: "User unauthorized"})
		return
	}

	// First try to get token from cookie
	tokenHash, err := c.Cookie("session_token")
	if err != nil {
		// Fallback to Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			h.logger.Error().
				Msg("Token Header unauthorized")
			util.HandleError(c, domain.AuthenticationError{Message: "Token Header unauthorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.logger.Error().
				Msg("Invalid authorization format")
			util.HandleError(c, domain.AuthenticationError{Message: "Invalid authorization format"})
			return
		}

		tokenHash = parts[1]
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(string), tokenHash); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to logout user")
		util.HandleError(c, err)
		return
	}

	middleware.ClearSecureCookie(c)

	h.logger.Info().
		Str("user_id", userID.(string)).
		Msg("User logged out successfully")

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
