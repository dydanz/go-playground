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

type ProgramHandler struct {
	programService domain.ProgramService
	logger         zerolog.Logger
}

func NewProgramHandler(programService domain.ProgramService) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
		logger:         logging.GetLogger(),
	}
}

// CreateProgram godoc
// @Summary Create a new loyalty program
// @Description Create a new loyalty program for a merchant
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param program body domain.CreateProgramRequest true "Program details"
// @Success 201 {object} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs [post]
func (h *ProgramHandler) Create(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create program request")

	var req domain.CreateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create program request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	program, err := h.programService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create program")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", program.ID.String()).
		Str("merchant_id", program.MerchantID.String()).
		Msg("Program created successfully")

	c.JSON(http.StatusCreated, program)
}

// GetProgram godoc
// @Summary Get a loyalty program by ID
// @Description Get details of a specific loyalty program
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Program ID"
// @Success 200 {object} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [get]
func (h *ProgramHandler) GetByID(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get program request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing program ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	program, err := h.programService.GetByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", id).
			Msg("Failed to get program")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", program.ID.String()).
		Str("merchant_id", program.MerchantID.String()).
		Msg("Program retrieved successfully")

	c.JSON(http.StatusOK, program)
}

// GetMerchantPrograms godoc
// @Summary Get all programs for a merchant
// @Description Get all loyalty programs for a specific merchant
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {array} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/merchant/{merchant_id} [get]
func (h *ProgramHandler) GetByMerchantID(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get programs by merchant request")

	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		h.logger.Error().
			Msg("Missing merchant ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "merchant_id",
			Message: "invalid merchant ID",
		})
		return
	}

	programs, err := h.programService.GetByMerchantID(c.Request.Context(), uuid.MustParse(merchantID))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", merchantID).
			Msg("Failed to get programs")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("merchant_id", merchantID).
		Int("programs_count", len(programs)).
		Msg("Programs retrieved successfully")

	c.JSON(http.StatusOK, programs)
}

// UpdateProgram godoc
// @Summary Update a loyalty program
// @Description Update details of a specific loyalty program
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Program ID"
// @Param program body domain.UpdateProgramRequest true "Program details"
// @Success 200 {object} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [put]
func (h *ProgramHandler) Update(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update program request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing program ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	var req domain.UpdateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update program request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	program, err := h.programService.Update(c.Request.Context(), uuid.MustParse(id), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", id).
			Interface("request", req).
			Msg("Failed to update program")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", program.ID.String()).
		Str("merchant_id", program.MerchantID.String()).
		Msg("Program updated successfully")

	c.JSON(http.StatusOK, program)
}

// DeleteProgram godoc
// @Summary Delete a loyalty program
// @Description Delete a specific loyalty program
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Program ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [delete]
func (h *ProgramHandler) Delete(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming delete program request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing program ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	if err := h.programService.Delete(c.Request.Context(), uuid.MustParse(id)); err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", id).
			Msg("Failed to delete program")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", id).
		Msg("Program deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Program deleted successfully"})
}

type CreateProgramRequest struct {
	MerchantID        string `json:"merchant_id" binding:"required"`
	ProgramName       string `json:"program_name" binding:"required"`
	PointCurrencyName string `json:"point_currency_name" binding:"required"`
}

type UpdateProgramRequest struct {
	ProgramName       string `json:"program_name" binding:"required"`
	PointCurrencyName string `json:"point_currency_name" binding:"required"`
}
