package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"log"
	"net/http"
	"strconv"

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
		log.Printf("error: %v", err)
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("%v error to check: \n%v\n%v", exists, userID, req.UserID)
		util.HandleError(c, domain.NewAuthenticationError("user not authenticated"))
		return
	}

	merchant, err := h.merchantService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Printf("error creating merchant: %v", err)
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
		log.Printf("error: %v", err)
		util.HandleError(c, domain.NewValidationError("id", "invalid merchant ID format"))
		return
	}

	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error: %v", err)
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Verify user has permission to update this merchant
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return
	}
	if err := h.verifyMerchantAccess(c, userID, id); err != nil {
		log.Printf("error: %v", err)
		util.HandleError(c, err)
		return
	}

	merchant, err := h.merchantService.Update(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("error: %v", err)
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
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return
	}
	if err := h.verifyMerchantAccess(c, userID, id); err != nil {
		util.HandleError(c, err)
		return
	}

	if err := h.merchantService.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err)
		return
	}

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
	log.Printf("userID: %v, merchantID: %v", userID, merchantID)
	// Get merchant by ID.
	merchant, err := h.merchantService.GetByID(c.Request.Context(), merchantID)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	if merchant.UserID != userID {
		log.Printf("error: %v", err)
		return domain.NewAuthorizationError("user does not have permission to access this merchant")
	}

	return nil
}
