package handler

import (
	"go-playground/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CreateProgramValidator handles validation for program creation requests
type CreateProgramValidator struct {
	*middleware.BaseValidator
}

// NewCreateProgramValidator creates a new instance of CreateProgramValidator
func NewCreateProgramValidator() *CreateProgramValidator {
	return &CreateProgramValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for program creation
// Required fields: merchant_id, program_name, point_currency_name
func (v *CreateProgramValidator) GetRules() map[string][]string {
	return map[string][]string{
		"merchant_id":         {"required", "string"},
		"program_name":        {"required", "string"},
		"point_currency_name": {"required", "string"},
	}
}

// ValidateRequest validates the program creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *CreateProgramValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// UpdateProgramValidator handles validation for program update requests
type UpdateProgramValidator struct {
	*middleware.BaseValidator
}

// NewUpdateProgramValidator creates a new instance of UpdateProgramValidator
func NewUpdateProgramValidator() *UpdateProgramValidator {
	return &UpdateProgramValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for program updates
// Required fields: program_name, point_currency_name
func (v *UpdateProgramValidator) GetRules() map[string][]string {
	return map[string][]string{
		"program_name":        {"required", "string"},
		"point_currency_name": {"required", "string"},
	}
}

// ValidateRequest validates the program update request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UpdateProgramValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
