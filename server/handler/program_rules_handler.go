package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/service"
	"go-playground/server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ProgramRulesHandler struct {
	programRulesService *service.ProgramRulesService
	logger              zerolog.Logger
}

func NewProgramRulesHandler(service *service.ProgramRulesService) *ProgramRulesHandler {
	return &ProgramRulesHandler{
		programRulesService: service,
		logger:              logging.GetLogger(),
	}
}

// CreateProgramRule godoc
// @Summary Create program rule
// @Description Create a new program rule
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param rule body domain.CreateProgramRuleRequest true "Program rule details"
// @Success 201 {object} domain.ProgramRule
// @Failure 400 {object} map[string]string
// @Router /program-rules [post]
func (h *ProgramRulesHandler) Create(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create program rule request")

	var req domain.CreateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create program rule request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	rule, err := h.programRulesService.Create(&req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create program rule")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("rule_id", rule.ID.String()).
		Str("program_id", rule.ProgramID.String()).
		Msg("Program rule created successfully")

	c.JSON(http.StatusCreated, rule)
}

// GetProgramRule godoc
// @Summary Get program rule by ID
// @Description Get program rule details by ID
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Program Rule ID"
// @Success 200 {object} domain.ProgramRule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /program-rules/{id} [get]
func (h *ProgramRulesHandler) GetByID(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get program rule request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing rule ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	rule, err := h.programRulesService.GetByID(id)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("rule_id", id).
			Msg("Failed to get program rule")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("rule_id", rule.ID.String()).
		Str("program_id", rule.ProgramID.String()).
		Msg("Program rule retrieved successfully")

	c.JSON(http.StatusOK, rule)
}

// GetProgramRules godoc
// @Summary Get program rules by program ID
// @Description Get all program rules for a specific program
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param program_id path string true "Program ID"
// @Success 200 {array} domain.ProgramRule
// @Failure 400 {object} map[string]string
// @Router /program-rules/program/{program_id} [get]
func (h *ProgramRulesHandler) GetByProgramID(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get program rules by program request")

	programID := c.Param("program_id")
	if programID == "" {
		h.logger.Error().
			Msg("Missing program ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rules, err := h.programRulesService.GetByProgramID(programID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", programID).
			Msg("Failed to get program rules")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", programID).
		Int("rules_count", len(rules)).
		Msg("Program rules retrieved successfully")

	c.JSON(http.StatusOK, rules)
}

// UpdateProgramRule godoc
// @Summary Update program rule
// @Description Update an existing program rule
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Program Rule ID"
// @Param rule body domain.UpdateProgramRuleRequest true "Program rule details to update"
// @Success 200 {object} domain.ProgramRule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /program-rules/{id} [put]
func (h *ProgramRulesHandler) Update(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update program rule request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing rule ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	var req domain.UpdateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update program rule request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	rule, err := h.programRulesService.Update(id, &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("rule_id", id).
			Interface("request", req).
			Msg("Failed to update program rule")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("rule_id", rule.ID.String()).
		Str("program_id", rule.ProgramID.String()).
		Msg("Program rule updated successfully")

	c.JSON(http.StatusOK, rule)
}

// DeleteProgramRule godoc
// @Summary Delete program rule
// @Description Delete an existing program rule
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Program Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /program-rules/{id} [delete]
func (h *ProgramRulesHandler) Delete(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming delete program rule request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing rule ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	if err := h.programRulesService.Delete(id); err != nil {
		h.logger.Error().
			Err(err).
			Str("rule_id", id).
			Msg("Failed to delete program rule")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("rule_id", id).
		Msg("Program rule deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Program rule deleted successfully"})
}

// GetActiveProgramRules godoc
// @Summary Get active program rules
// @Description Get all active program rules for a specific program
// @Tags program-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param program_id path string true "Program ID"
// @Success 200 {array} domain.ProgramRule
// @Failure 400 {object} map[string]string
// @Router /program-rules/program/{program_id}/active [get]
func (h *ProgramRulesHandler) GetActiveRules(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get active program rules request")

	programID := c.Param("program_id")
	if programID == "" {
		h.logger.Error().
			Msg("Missing program ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rules, err := h.programRulesService.GetActiveRules(programID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("program_id", programID).
			Msg("Failed to get active program rules")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("program_id", programID).
		Int("active_rules_count", len(rules)).
		Msg("Active program rules retrieved successfully")

	c.JSON(http.StatusOK, rules)
}

// GetProgramRulesByMerchantId godoc
// @Summary Get all program rules for a merchant
// @Description Get all program rules across all programs for a specific merchant with pagination
// @Tags program-rules
// @Accept json
// @Produce json
// @Param merchant_id path string true "Merchant ID"
// @Param page query integer false "Page number (default: 1)"
// @Param limit query integer false "Items per page (default: 10, max: 100)"
// @Success 200 {object} domain.PaginatedResponse
// @Failure 400 {object} util.ErrorResponse "Invalid merchant ID format or pagination parameters"
// @Failure 500 {object} util.ErrorResponse "Internal server error"
// @Router /program-rules/by-merchant/{merchant_id} [get]
func (h *ProgramRulesHandler) GetProgramRulesByMerchantId(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get program rules by merchant request")

	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		h.logger.Error().
			Msg("Missing merchant ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant_id is required"})
		return
	}

	var pagination domain.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind pagination request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	h.logger.Debug().
		Str("merchant_id", merchantID).
		Int("page", pagination.Page).
		Int("limit", pagination.Limit).
		Msg("Fetching program rules")

	rules, total, err := h.programRulesService.GetProgramRulesByMerchantId(merchantID, pagination.Page, pagination.Limit)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", merchantID).
			Msg("Failed to get program rules")
		switch err.(type) {
		case *domain.ValidationError:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.Info().
		Str("merchant_id", merchantID).
		Int("rules_count", len(rules)).
		Int64("total_rules", total).
		Int("page", pagination.Page).
		Int("limit", pagination.Limit).
		Msg("Program rules retrieved successfully")

	response := domain.NewPaginatedResponse(rules, total, pagination.Page, pagination.Limit)
	c.JSON(http.StatusOK, response)
}
