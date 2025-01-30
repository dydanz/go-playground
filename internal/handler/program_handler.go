package handler

import (
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProgramHandler struct {
	programService *service.ProgramService
}

func NewProgramHandler(programService *service.ProgramService) *ProgramHandler {
	return &ProgramHandler{programService: programService}
}

// CreateProgram godoc
// @Summary Create a new loyalty program
// @Description Create a new loyalty program for a merchant
// @Tags programs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param program body CreateProgramRequest true "Program details"
// @Success 201 {object} postgres.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs [post]
func (h *ProgramHandler) Create(c *gin.Context) {
	var req CreateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchantID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant ID"})
		return
	}

	program, err := h.programService.CreateProgram(c.Request.Context(), merchantID, req.ProgramName, req.PointCurrencyName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
// @Success 200 {object} postgres.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [get]
func (h *ProgramHandler) GetByID(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	program, err := h.programService.GetProgram(c.Request.Context(), programID)
	if err != nil {
		if err == service.ErrProgramNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
// @Success 200 {array} postgres.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/merchant/{merchant_id} [get]
func (h *ProgramHandler) GetByMerchantID(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant ID"})
		return
	}

	programs, err := h.programService.GetMerchantPrograms(c.Request.Context(), merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
// @Param program body UpdateProgramRequest true "Program details"
// @Success 200 {object} postgres.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [put]
func (h *ProgramHandler) Update(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	var req UpdateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program, err := h.programService.UpdateProgram(c.Request.Context(), programID, req.ProgramName, req.PointCurrencyName)
	if err != nil {
		if err == service.ErrProgramNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid program ID"})
		return
	}

	err = h.programService.DeleteProgram(c.Request.Context(), programID)
	if err != nil {
		if err == service.ErrProgramNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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
