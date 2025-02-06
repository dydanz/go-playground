package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsService domain.PointsService
}

func NewPointsHandler(pointsService domain.PointsService) *PointsHandler {
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
	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	ledger, err := h.pointsService.GetLedger(customerID, programID)
	if err != nil {
		util.HandleError(c, err)
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
	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	balance, err := h.pointsService.GetBalance(customerID, programID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, balance)
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
	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	var req domain.EarnPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	req.CustomerID = customerID
	req.ProgramID = programID

	result, err := h.pointsService.EarnPoints(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
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
	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	var req domain.RedeemPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	req.CustomerID = customerID
	req.ProgramID = programID

	result, err := h.pointsService.RedeemPoints(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
