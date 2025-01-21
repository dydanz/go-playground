package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
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
// @Param redemption body domain.Redemption true "Redemption details"
// @Success 201 {object} domain.Redemption
// @Failure 400 {object} map[string]string
// @Router /redemptions [post]
func (h *RedemptionHandler) Create(c *gin.Context) {
	var redemption domain.Redemption
	if err := c.ShouldBindJSON(&redemption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.redemptionService.Create(&redemption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Param id path string true "Redemption ID"
// @Success 200 {object} domain.Redemption
// @Failure 404 {object} map[string]string
// @Router /redemptions/{id} [get]
func (h *RedemptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	redemption, err := h.redemptionService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
// @Param user_id path string true "User ID"
// @Success 200 {array} domain.Redemption
// @Failure 404 {object} map[string]string
// @Router /redemptions/user/{user_id} [get]
func (h *RedemptionHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	redemptions, err := h.redemptionService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
// @Param id path string true "Redemption ID"
// @Param status body string true "New status" Enums(completed, pending, failed, canceled)
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /redemptions/{id}/status [put]
func (h *RedemptionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var status string
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.redemptionService.UpdateStatus(id, status); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
