package handler

import (
	"go-playground/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserCreateValidator handles validation for user registration requests
type UserCreateValidator struct {
	*middleware.BaseValidator
}

// NewUserCreateValidator creates a new instance of UserCreateValidator
func NewUserCreateValidator() *UserCreateValidator {
	return &UserCreateValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for user creation
// Required fields: email, password, full_name
func (v *UserCreateValidator) GetRules() map[string][]string {
	return map[string][]string{
		"email":     {"required", "string"},
		"password":  {"required", "string"},
		"full_name": {"required", "string"},
	}
}

// ValidateRequest validates the user creation request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *UserCreateValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// AuthLoginValidator handles validation for user login requests
type AuthLoginValidator struct {
	*middleware.BaseValidator
}

// NewAuthLoginValidator creates a new instance of AuthLoginValidator
func NewAuthLoginValidator() *AuthLoginValidator {
	return &AuthLoginValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for user login
// Required fields: email, password
func (v *AuthLoginValidator) GetRules() map[string][]string {
	return map[string][]string{
		"email":    {"required", "string"},
		"password": {"required", "string"},
	}
}

// ValidateRequest validates the login request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *AuthLoginValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}

// VerificationValidator handles validation for verification requests
type VerificationValidator struct {
	*middleware.BaseValidator
}

// NewVerificationValidator creates a new instance of VerificationValidator
func NewVerificationValidator() *VerificationValidator {
	return &VerificationValidator{
		BaseValidator: middleware.NewBaseValidator(),
	}
}

// GetRules returns the validation rules for verification
// Required fields: email, code
func (v *VerificationValidator) GetRules() map[string][]string {
	return map[string][]string{
		"email": {"required", "string"},
		"code":  {"required", "string"},
	}
}

// ValidateRequest validates the verification request against defined rules
// Parameters:
//   - c: Gin context containing the request
//
// Returns:
//   - error: Validation error if any
func (v *VerificationValidator) ValidateRequest(c *gin.Context) error {
	validator := middleware.NewJSONValidator(v.GetRules())
	return validator.ValidateRequest(c)
}
