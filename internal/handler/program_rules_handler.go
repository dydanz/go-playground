package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"go-playground/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgramRulesHandler struct {
	programRulesService *service.ProgramRulesService
}

func NewProgramRulesHandler(service *service.ProgramRulesService) *ProgramRulesHandler {
	return &ProgramRulesHandler{
		programRulesService: service,
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
	var req domain.CreateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	rule, err := h.programRulesService.Create(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	rule, err := h.programRulesService.GetByID(id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	programID := c.Param("program_id")
	if programID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rules, err := h.programRulesService.GetByProgramID(programID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	var req domain.UpdateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	rule, err := h.programRulesService.Update(id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid rule ID",
		})
		return
	}

	if err := h.programRulesService.Delete(id); err != nil {
		util.HandleError(c, err)
		return
	}

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
	programID := c.Param("program_id")
	if programID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "program_id",
			Message: "invalid program ID",
		})
		return
	}

	rules, err := h.programRulesService.GetActiveRules(programID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

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
	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant_id is required"})
		return
	}

	var pagination domain.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	log.Printf("merchantID: %v, page: %d, limit: %d", merchantID, pagination.Page, pagination.Limit)

	rules, total, err := h.programRulesService.GetProgramRulesByMerchantId(merchantID, pagination.Page, pagination.Limit)
	if err != nil {
		// Handle different types of errors
		switch err.(type) {
		case *domain.ValidationError:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := domain.NewPaginatedResponse(rules, total, pagination.Page, pagination.Limit)
	c.JSON(http.StatusOK, response)
}
