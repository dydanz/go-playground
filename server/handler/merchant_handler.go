package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MerchantHandler struct {
	merchantService domain.MerchantService
	logger          zerolog.Logger
}

func NewMerchantHandler(merchantService domain.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
		logger:          logging.GetLogger(),
	}
}

// @Summary Create merchant
// @Description Create a new merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body domain.CreateMerchantRequest true "Merchant details"
// @Success 201 {object} domain.Merchant
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 409 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants [post]
func (h *MerchantHandler) Create(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create merchant request")

	var req domain.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create merchant request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().
			Interface("user_id", userID).
			Interface("req_user_id", req.UserID).
			Msg("User not authenticated")
		util.HandleError(c, domain.NewAuthenticationError("user not authenticated"))
		return
	}

	merchant, err := h.merchantService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create merchant")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("merchant_id", merchant.ID.String()).
		Str("name", merchant.Name).
		Str("user_id", merchant.UserID.String()).
		Str("type", string(merchant.Type)).
		Msg("Merchant created successfully")

	c.JSON(http.StatusCreated, merchant)
}

// @Summary Get merchant by ID
// @Description Get merchant details by ID
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} domain.Merchant
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants/{id} [get]
func (h *MerchantHandler) GetByID(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get merchant request")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", c.Param("id")).
			Msg("Invalid merchant ID format")
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	merchant, err := h.merchantService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", id.String()).
			Msg("Failed to get merchant")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("merchant_id", merchant.ID.String()).
		Str("name", merchant.Name).
		Str("user_id", merchant.UserID.String()).
		Str("type", string(merchant.Type)).
		Msg("Merchant retrieved successfully")

	c.JSON(http.StatusOK, merchant)
}

// @Summary Get all merchants
// @Description Get all merchants
// @Tags merchants
// @Produce json
// @Success 200 {array} domain.Merchant
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants [get]
func (h *MerchantHandler) GetAll(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get all merchants request")

	userIDStr, exists := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("user_id", userIDStr).
			Msg("Invalid user ID format")
		return
	}
	if !exists {
		h.logger.Error().
			Interface("user_id", userID).
			Msg("User not authenticated")
		util.HandleError(c, domain.NewAuthenticationError("user not authenticated"))
		return
	}

	merchants, err := h.merchantService.GetAll(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to get merchants")
		util.HandleError(c, err)
		return
	}

	if len(merchants) == 0 {
		h.logger.Info().
			Str("user_id", userID.String()).
			Msg("No merchants found")
		util.EmptyResponse(c)
		return
	}

	h.logger.Info().
		Str("user_id", userID.String()).
		Int("merchant_count", len(merchants)).
		Msg("Merchants retrieved successfully")

	c.JSON(http.StatusOK, merchants)
}

// @Summary Update merchant
// @Description Update merchant details
// @Tags merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Param merchant body domain.UpdateMerchantRequest true "Merchant details"
// @Success 200 {object} domain.Merchant
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants/{id} [put]
func (h *MerchantHandler) Update(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update merchant request")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", c.Param("id")).
			Msg("Invalid merchant ID format")
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update merchant request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("user_id", userIDStr).
			Msg("Invalid user ID format")
		return
	}

	if err := h.verifyMerchantAccess(c, userID, id); err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("merchant_id", id.String()).
			Msg("Failed to verify merchant access")
		util.HandleError(c, err)
		return
	}

	merchant, err := h.merchantService.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", id.String()).
			Interface("request", req).
			Msg("Failed to update merchant")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("merchant_id", merchant.ID.String()).
		Str("name", merchant.Name).
		Str("user_id", merchant.UserID.String()).
		Str("type", string(merchant.Type)).
		Msg("Merchant updated successfully")

	c.JSON(http.StatusOK, merchant)
}

// @Summary Delete merchant
// @Description Delete a merchant
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 204 "No Content"
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants/{id} [delete]
func (h *MerchantHandler) Delete(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming delete merchant request")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", c.Param("id")).
			Msg("Invalid merchant ID format")
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("user_id", userIDStr).
			Msg("Invalid user ID format")
		return
	}

	if err := h.verifyMerchantAccess(c, userID, id); err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("merchant_id", id.String()).
			Msg("Failed to verify merchant access")
		util.HandleError(c, err)
		return
	}

	if err := h.merchantService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", id.String()).
			Msg("Failed to delete merchant")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("merchant_id", id.String()).
		Msg("Merchant deleted successfully")

	c.Status(http.StatusNoContent)
}

// @Summary Get all merchants for a user
// @Description Get all Merchants for a User with pagination
// @Tags merchants
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} domain.PaginatedMerchants
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchants/user/{user_id} [get]
func (h *MerchantHandler) GetMerchantsByUserID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("user_id", "invalid user ID format"))
		return
	}

	// Get pagination parameters
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	merchants, total, err := h.merchantService.GetMerchantsByUserID(c.Request.Context(), userID, offset, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	response := domain.PaginatedMerchants{
		Merchants: merchants,
		Pagination: domain.Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			Limit:       limit,
			Total:       total,
		},
	}

	if len(merchants) == 0 {
		util.EmptyResponse(c)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to verify merchant access
func (h *MerchantHandler) verifyMerchantAccess(c *gin.Context, userID uuid.UUID, merchantID uuid.UUID) error {
	h.logger.Debug().
		Str("user_id", userID.String()).
		Str("merchant_id", merchantID.String()).
		Msg("Verifying merchant access")

	merchant, err := h.merchantService.GetByID(c.Request.Context(), merchantID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", merchantID.String()).
			Msg("Failed to get merchant for access verification")
		return err
	}

	if merchant.UserID != userID {
		h.logger.Warn().
			Str("user_id", userID.String()).
			Str("merchant_id", merchantID.String()).
			Str("merchant_owner_id", merchant.UserID.String()).
			Msg("User does not have permission to access this merchant")
		return domain.NewAuthorizationError("user does not have permission to access this merchant")
	}

	return nil
}
