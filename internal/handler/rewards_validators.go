package handler

import (
	"go-playground/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRewardValidator handles validation for reward creation requests
type CreateRewardValidator struct {
	*middleware.BaseValidator
}

// NewCreateRewardValidator creates a new instance of CreateRewardValidator
func NewCreateRewardValidator() *CreateRewardValidator {
	return &CreateRewardValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for reward creation
// Required fields: merchant_id, name, points_required, quantity
// Optional fields: description, image_url, metadata, start_date, end_date
func (v *CreateRewardValidator) GetRules() map[string][]string {
	return map[string][]string{
		"merchant_id":     {"required", "string"},
		"name":            {"required", "string"},
		"points_required": {"required", "number"},
		"quantity":        {"required", "number"},
		"description":     {"string"},
		"image_url":       {"string"},
		"metadata":        {"object"},
		"start_date":      {"string"},
		"end_date":        {"string"},
	}
}

// ValidateRequest validates the reward creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *CreateRewardValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// UpdateRewardValidator handles validation for reward update requests
type UpdateRewardValidator struct {
	*middleware.BaseValidator
}

// NewUpdateRewardValidator creates a new instance of UpdateRewardValidator
func NewUpdateRewardValidator() *UpdateRewardValidator {
	return &UpdateRewardValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for reward updates
// Required fields: name, points_required, quantity
// Optional fields: description, image_url, metadata, start_date, end_date
func (v *UpdateRewardValidator) GetRules() map[string][]string {
	return map[string][]string{
		"name":            {"required", "string"},
		"points_required": {"required", "number"},
		"quantity":        {"required", "number"},
		"description":     {"string"},
		"image_url":       {"string"},
		"metadata":        {"object"},
		"start_date":      {"string"},
		"end_date":        {"string"},
	}
}

// ValidateRequest validates the reward update request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UpdateRewardValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
