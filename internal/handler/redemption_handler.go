package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedemptionHandler struct {
	redemptionService *service.RedemptionService
}

func NewRedemptionHandler(redemptionService *service.RedemptionService) *RedemptionHandler {
	return &RedemptionHandler{redemptionService: redemptionService}
}

// @Summary Create redemption
// @Description Create a new redemption request
// @Tags redemptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param redemption body domain.Redemption true "Redemption details"
// @Success 201 {object} domain.Redemption
// @Failure 400 {object} map[string]string
// @Router /redemptions [post]
func (h *RedemptionHandler) Create(c *gin.Context) {
	var req domain.CreateRedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	redemption := &domain.Redemption{
		MerchantCustomersID: req.MerchantCustomersID,
		RewardID:            req.RewardID,
		Status:              "pending",
	}

	if err := h.redemptionService.Create(c.Request.Context(), redemption); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, redemption)
}

// @Summary Get redemption by ID
// @Description Get redemption details by ID
// @Tags redemptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Redemption ID"
// @Success 200 {object} domain.Redemption
// @Failure 404 {object} map[string]string
// @Router /redemptions/{id} [get]
func (h *RedemptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid redemption ID",
		})
		return
	}

	redemption, err := h.redemptionService.GetByID(id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, redemption)
}

// @Summary Get user redemptions
// @Description Get all redemptions for a specific user
// @Tags redemptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} domain.Redemption
// @Failure 404 {object} map[string]string
// @Router /redemptions/user/{user_id} [get]
func (h *RedemptionHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "user_id",
			Message: "invalid user ID",
		})
		return
	}

	redemptions, err := h.redemptionService.GetByUserID(userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, redemptions)
}

// @Summary Update redemption status
// @Description Update the status of a redemption
// @Tags redemptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Redemption ID"
// @Param status body string true "New status" Enums(completed, pending, failed, canceled)
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /redemptions/{id}/status [put]
func (h *RedemptionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid redemption ID",
		})
		return
	}

	var req domain.UpdateRedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	if err := h.redemptionService.UpdateStatus(c.Request.Context(), id, string(req.Status)); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}
