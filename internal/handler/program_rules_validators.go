package handler

import (
	"go-playground/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CreateProgramRuleValidator handles validation for program rule creation requests
type CreateProgramRuleValidator struct {
	*middleware.BaseValidator
}

// NewCreateProgramRuleValidator creates a new instance of CreateProgramRuleValidator
func NewCreateProgramRuleValidator() *CreateProgramRuleValidator {
	return &CreateProgramRuleValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for program rule creation
// Required fields: program_id, rule_type, points_value, is_active, start_date
// Optional fields: spend_amount, end_date, description, min_points, max_points, conditions, metadata
func (v *CreateProgramRuleValidator) GetRules() map[string][]string {
	return map[string][]string{
		"program_id":   {"required", "string"},
		"rule_type":    {"required", "string"},
		"points_value": {"required", "number"},
		"spend_amount": {"number"},
		"is_active":    {"required", "boolean"},
		"start_date":   {"required", "string"},
		"end_date":     {"string"},
		"description":  {"string"},
		"min_points":   {"number"},
		"max_points":   {"number"},
		"conditions":   {"object"},
		"metadata":     {"object"},
	}
}

// ValidateRequest validates the program rule creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *CreateProgramRuleValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// UpdateProgramRuleValidator handles validation for program rule update requests
type UpdateProgramRuleValidator struct {
	*middleware.BaseValidator
}

// NewUpdateProgramRuleValidator creates a new instance of UpdateProgramRuleValidator
func NewUpdateProgramRuleValidator() *UpdateProgramRuleValidator {
	return &UpdateProgramRuleValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for program rule updates
// Required fields: rule_type, points_value, spend_amount, is_active, start_date, description
// Optional fields: end_date
func (v *UpdateProgramRuleValidator) GetRules() map[string][]string {
	return map[string][]string{
		"rule_type":    {"required", "string"},
		"points_value": {"required", "number"},
		"spend_amount": {"required", "number"},
		"is_active":    {"required", "string"},
		"start_date":   {"required", "string"},
		"end_date":     {"string"},
		"description":  {"required", "string"},
	}
}

// ValidateRequest validates the program rule update request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UpdateProgramRuleValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
