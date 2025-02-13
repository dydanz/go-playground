package handler

import (
	"go-playground/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InternalLoadTestHandler struct {
	authService domain.AuthService
}

func NewInternalLoadTestHandler(authService domain.AuthService) *InternalLoadTestHandler {
	return &InternalLoadTestHandler{authService: authService}
}

// @Summary Get verification code for testing
// @Description Get OTP verification code by email (for testing purposes only)
// @Tags internal-load-testing
// @Accept json
// @Produce json
// @Param email query string true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/test/get-verification/code [get]
func (h *InternalLoadTestHandler) GetVerificationCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := h.authService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	verification, err := h.authService.GetVerificationByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if verification == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "verification code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": verification.OTP})
}

// @Summary Get random verified user
// @Description Get a random verified user's credentials (for testing purposes only)
// @Tags internal-load-testing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/test/random-user [get]
func (h *InternalLoadTestHandler) GetRandomVerifiedUser(c *gin.Context) {
	user, err := h.authService.GetRandomActiveUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active users found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":    user.Email,
		"password": user.Password, // Return default test password or handle as needed
	})
}
