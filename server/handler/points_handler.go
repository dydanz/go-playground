package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PointsHandler struct {
	pointsService domain.PointsService
	logger        zerolog.Logger
}

func NewPointsHandler(pointsService domain.PointsService) *PointsHandler {
	return &PointsHandler{
		pointsService: pointsService,
		logger:        logging.GetLogger(),
	}
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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get points ledger request")

	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		h.logger.Error().
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Missing required parameters")
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	ledger, err := h.pointsService.GetLedger(c.Request.Context(), uuid.MustParse(customerID), uuid.MustParse(programID))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Failed to get points ledger")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("customer_id", customerID).
		Str("program_id", programID).
		Int("entries_count", len(ledger)).
		Msg("Points ledger retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get points balance request")

	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		h.logger.Error().
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Missing required parameters")
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	balance, err := h.pointsService.GetBalance(c.Request.Context(), uuid.MustParse(customerID), uuid.MustParse(programID))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Failed to get points balance")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("customer_id", customerID).
		Str("program_id", programID).
		Interface("balance", balance).
		Msg("Points balance retrieved successfully")

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
// @Param points body domain.EarnPointsRequest true "Points to earn"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/earn [post]
func (h *PointsHandler) EarnPoints(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming earn points request")

	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		h.logger.Error().
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Missing required parameters")
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	var req domain.EarnPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind earn points request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	req.CustomerID = customerID
	req.ProgramID = programID

	result, err := h.pointsService.EarnPoints(c.Request.Context(), &domain.PointsTransaction{
		CustomerID:    customerID,
		ProgramID:     programID,
		Points:        req.Points,
		TransactionID: "",
	})
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("customer_id", customerID).
			Str("program_id", programID).
			Int("points", req.Points).
			Msg("Failed to earn points")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("customer_id", customerID).
		Str("program_id", programID).
		Int("points", req.Points).
		Interface("result", result).
		Msg("Points earned successfully")

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
// @Param points body domain.RedeemPointsRequest true "Points to redeem"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /points/{customer_id}/{program_id}/redeem [post]
func (h *PointsHandler) RedeemPoints(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming redeem points request")

	customerID := c.Param("customer_id")
	programID := c.Param("program_id")

	if customerID == "" || programID == "" {
		h.logger.Error().
			Str("customer_id", customerID).
			Str("program_id", programID).
			Msg("Missing required parameters")
		util.HandleError(c, domain.ValidationError{
			Message: "customer_id and program_id are required",
		})
		return
	}

	var req domain.RedeemPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind redeem points request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	req.CustomerID = customerID
	req.ProgramID = programID

	result, err := h.pointsService.RedeemPoints(c.Request.Context(), &domain.PointsTransaction{
		CustomerID:    customerID,
		ProgramID:     programID,
		Points:        req.Points,
		TransactionID: "",
	})
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("customer_id", customerID).
			Str("program_id", programID).
			Int("points", req.Points).
			Msg("Failed to redeem points")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("customer_id", customerID).
		Str("program_id", programID).
		Int("points", req.Points).
		Interface("result", result).
		Msg("Points redeemed successfully")

	c.JSON(http.StatusOK, result)
}
