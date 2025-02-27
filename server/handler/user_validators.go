package handler

import (
	"go-playground/server/middleware"

	"github.com/gin-gonic/gin"
)

// UpdateUserValidator handles validation for user update requests
type UpdateUserValidator struct {
	*middleware.BaseValidator
}

// NewUpdateUserValidator creates a new instance of UpdateUserValidator
func NewUpdateUserValidator() *UpdateUserValidator {
	return &UpdateUserValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for user updates
// Optional fields: full_name, email, avatar, status
func (v *UpdateUserValidator) GetRules() map[string][]string {
	return map[string][]string{
		"full_name": {"string"},
		"email":     {"string"},
		"avatar":    {"string"},
		"status":    {"string"},
	}
}

// ValidateRequest validates the user update request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UpdateUserValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
