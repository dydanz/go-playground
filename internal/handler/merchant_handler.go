package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	merchantService domain.MerchantService
}

func NewMerchantHandler(merchantService domain.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantService: merchantService}
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
	var req domain.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		util.HandleError(c, domain.NewAuthenticationError("user not authenticated"))
		return
	}
	req.UserID = userID.(uuid.UUID)

	merchant, err := h.merchantService.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	merchant, err := h.merchantService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	merchants, err := h.merchantService.GetAll(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}

	if len(merchants) == 0 {
		util.EmptyResponse(c)
		return
	}

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Verify user has permission to update this merchant
	if err := h.verifyMerchantAccess(c, id); err != nil {
		util.HandleError(c, err)
		return
	}

	merchant, err := h.merchantService.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	// Verify user has permission to delete this merchant
	if err := h.verifyMerchantAccess(c, id); err != nil {
		util.HandleError(c, err)
		return
	}

	if err := h.merchantService.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper function to verify merchant access
func (h *MerchantHandler) verifyMerchantAccess(c *gin.Context, merchantID uuid.UUID) error {
	userID, exists := c.Get("user_id")
	if !exists {
		return domain.NewAuthenticationError("user not authenticated")
	}

	merchant, err := h.merchantService.GetByID(c.Request.Context(), merchantID)
	if err != nil {
		return err
	}

	if merchant.UserID != userID.(uuid.UUID) {
		return domain.NewAuthorizationError("user does not have permission to access this merchant")
	}

	return nil
}
