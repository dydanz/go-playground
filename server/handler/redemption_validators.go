package handler

import (
	"go-playground/server/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRedemptionValidator handles validation for redemption creation requests
type CreateRedemptionValidator struct {
	*middleware.BaseValidator
}

// NewCreateRedemptionValidator creates a new instance of CreateRedemptionValidator
func NewCreateRedemptionValidator() *CreateRedemptionValidator {
	return &CreateRedemptionValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for redemption creation
// Required fields: user_id, reward_id, points_amount
// Optional fields: metadata
func (v *CreateRedemptionValidator) GetRules() map[string][]string {
	return map[string][]string{
		"user_id":       {"required", "string"},
		"reward_id":     {"required", "string"},
		"points_amount": {"required", "number"},
		"metadata":      {"object"},
	}
}

// ValidateRequest validates the redemption creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *CreateRedemptionValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// UpdateRedemptionStatusValidator handles validation for redemption status update requests
type UpdateRedemptionStatusValidator struct {
	*middleware.BaseValidator
}

// NewUpdateRedemptionStatusValidator creates a new instance of UpdateRedemptionStatusValidator
func NewUpdateRedemptionStatusValidator() *UpdateRedemptionStatusValidator {
	return &UpdateRedemptionStatusValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for redemption status updates
// Required fields: status
func (v *UpdateRedemptionStatusValidator) GetRules() map[string][]string {
	return map[string][]string{
		"status": {"required", "string"},
	}
}

// ValidateRequest validates the redemption status update request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UpdateRedemptionStatusValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
