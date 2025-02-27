package handler

import (
	"go-playground/server/domain"
	"go-playground/server/middleware"
	"go-playground/server/util"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
	var req domain.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	var req domain.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	if err := h.authService.VerifyRegistration(c.Request.Context(), &req); err != nil {
		util.HandleError(c, err)
		return
	}

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
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	authToken, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		log.Println("error on AuthHandler.Login: ", err)
		util.HandleError(c, err)
		return
	}

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
			util.HandleError(c, domain.AuthenticationError{Message: "Token Header unauthorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.HandleError(c, domain.AuthenticationError{Message: "Invalid authorization format"})
			return
		}

		tokenHash = parts[1]
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(string), tokenHash); err != nil {
		util.HandleError(c, err)
		return
	}

	middleware.ClearSecureCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
