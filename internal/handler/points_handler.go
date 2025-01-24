package handler

import (
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsService *service.PointsService
}

func NewPointsHandler(pointsService *service.PointsService) *PointsHandler {
	return &PointsHandler{pointsService: pointsService}
}

// @Summary Get user points balance
// @Description Get points balance for a specific user
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} domain.PointsBalance
// @Failure 404 {object} map[string]string
// @Router /points/{user_id} [get]
func (h *PointsHandler) GetBalance(c *gin.Context) {
	userID := c.Param("user_id")
	balance, err := h.pointsService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

// @Summary Update points balance
// @Description Update points balance for a user
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param user_id path string true "User ID"
// @Param points body int true "Points to add/subtract"
// @Success 200 {object} domain.PointsBalance
// @Failure 400 {object} map[string]string
// @Router /points/{user_id} [put]
func (h *PointsHandler) UpdateBalance(c *gin.Context) {
	userID := c.Param("user_id")
	var points int
	if err := c.ShouldBindJSON(&points); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.pointsService.UpdateBalance(userID, points); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance updated successfully"})
}
