package handler

import (
	"go-playground/server/middleware"

	"github.com/gin-gonic/gin"
)

// EarnPointsValidator handles validation for point earning requests
type EarnPointsValidator struct {
	*middleware.BaseValidator
}

// NewEarnPointsValidator creates a new instance of EarnPointsValidator
func NewEarnPointsValidator() *EarnPointsValidator {
	return &EarnPointsValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for earning points
// Required fields: points
// Optional fields: transaction_id
func (v *EarnPointsValidator) GetRules() map[string][]string {
	return map[string][]string{
		"points":         {"required", "number"},
		"transaction_id": {"string"},
	}
}

// ValidateRequest validates the point earning request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *EarnPointsValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// RedeemPointsValidator handles validation for point redemption requests
type RedeemPointsValidator struct {
	*middleware.BaseValidator
}

// NewRedeemPointsValidator creates a new instance of RedeemPointsValidator
func NewRedeemPointsValidator() *RedeemPointsValidator {
	return &RedeemPointsValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for redeeming points
// Required fields: points
// Optional fields: transaction_id
func (v *RedeemPointsValidator) GetRules() map[string][]string {
	return map[string][]string{
		"points":         {"required", "number"},
		"transaction_id": {"string"},
	}
}

// ValidateRequest validates the point redemption request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *RedeemPointsValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
