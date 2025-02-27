package handler

import (
	"go-playground/server/middleware"

	"github.com/gin-gonic/gin"
)

// CreateTransactionValidator handles validation for transaction creation requests
type CreateTransactionValidator struct {
	*middleware.BaseValidator
}

// NewCreateTransactionValidator creates a new instance of CreateTransactionValidator
func NewCreateTransactionValidator() *CreateTransactionValidator {
	return &CreateTransactionValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for transaction creation
// Required fields: customer_id, merchant_id, amount
// Optional fields: description, currency, metadata, tags
func (v *CreateTransactionValidator) GetRules() map[string][]string {
	return map[string][]string{
		"customer_id": {"required", "string"},
		"merchant_id": {"required", "string"},
		"amount":      {"required", "number"},
		"description": {"string"},
		"currency":    {"string"},
		"metadata":    {"object"},
		"tags":        {"array"},
	}
}

// ValidateRequest validates the transaction creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *CreateTransactionValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
