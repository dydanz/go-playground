package handler

import (
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PointsHandler struct {
	pointsService *service.PointsService
}

func NewPointsHandler(pointsService *service.PointsService) *PointsHandler {
	return &PointsHandler{pointsService: pointsService}
}

// GetLedger godoc
// @Summary Get points ledger
// @Description Get points ledger entries for a customer in a program
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Param program_id path string true "Program ID"
// @Success 200 {array} domain.PointsLedger
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/ledger [get]
func (h *PointsHandler) GetLedger(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	ledger, err := h.pointsService.GetLedger(c.Request.Context(), customerID, programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ledger)
}

// GetBalance godoc
// @Summary Get points balance
// @Description Get current points balance for a customer in a program
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Param program_id path string true "Program ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/balance [get]
func (h *PointsHandler) GetBalance(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	balance, err := h.pointsService.GetBalance(c.Request.Context(), customerID, programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// EarnPoints godoc
// @Summary Earn points
// @Description Add points to a customer's balance in a program
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Param program_id path string true "Program ID"
// @Param points body EarnPointsRequest true "Points to earn"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/earn [post]
func (h *PointsHandler) EarnPoints(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	var req EarnPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transactionID *uuid.UUID
	if req.TransactionID != "" {
		id, err := uuid.Parse(req.TransactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
			return
		}
		transactionID = &id
	}

	if err := h.pointsService.EarnPoints(c.Request.Context(), customerID, programID, req.Points, transactionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points earned successfully"})
}

// RedeemPoints godoc
// @Summary Redeem points
// @Description Redeem points from a customer's balance in a program
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Param program_id path string true "Program ID"
// @Param points body RedeemPointsRequest true "Points to redeem"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/redeem [post]
func (h *PointsHandler) RedeemPoints(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	var req RedeemPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transactionID *uuid.UUID
	if req.TransactionID != "" {
		id, err := uuid.Parse(req.TransactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
			return
		}
		transactionID = &id
	}

	if err := h.pointsService.RedeemPoints(c.Request.Context(), customerID, programID, req.Points, transactionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points redeemed successfully"})
}

type EarnPointsRequest struct {
	Points        int    `json:"points" binding:"required,gt=0"`
	TransactionID string `json:"transaction_id,omitempty"`
}

type RedeemPointsRequest struct {
	Points        int    `json:"points" binding:"required,gt=0"`
	TransactionID string `json:"transaction_id,omitempty"`
}
