package handler

import (
	"errors"
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProgramRuleHandler struct {
	programRuleService *service.ProgramRuleService
}

func NewProgramRuleHandler(programRuleService *service.ProgramRuleService) *ProgramRuleHandler {
	return &ProgramRuleHandler{programRuleService: programRuleService}
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
func (h *ProgramRuleHandler) Create(c *gin.Context) {
	var req domain.CreateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.programRuleService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
func (h *ProgramRuleHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program rule ID"})
		return
	}

	rule, err := h.programRuleService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "program rule not found"})
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
func (h *ProgramRuleHandler) GetByProgramID(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	rules, err := h.programRuleService.GetByProgramID(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (h *ProgramRuleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program rule ID"})
		return
	}

	var req domain.UpdateProgramRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.programRuleService.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == errors.New("resource was not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "program rule not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
func (h *ProgramRuleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program rule ID"})
		return
	}

	err = h.programRuleService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == errors.New("resource was not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "program rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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
func (h *ProgramRuleHandler) GetActiveRules(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("program_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	rules, err := h.programRuleService.GetActiveRules(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}
