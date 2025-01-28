package handler

import (
	"go-playground/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	authService domain.AuthService
}

func NewTestHandler(authService domain.AuthService) *TestHandler {
	return &TestHandler{authService: authService}
}

// @Summary Get verification code for testing
// @Description Get OTP verification code by email (for testing purposes only)
// @Tags test
// @Accept json
// @Produce json
// @Param email query string true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/test/get-verification/code [get]
func (h *TestHandler) GetVerificationCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	verification, err := h.authService.GetVerificationByUserID(user.ID)
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
