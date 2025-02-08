package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	programService domain.ProgramService
}

func NewProgramHandler(programService domain.ProgramService) *ProgramHandler {
	return &ProgramHandler{programService: programService}
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
	var req domain.CreateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	program, err := h.programService.Create(&req)
	if err != nil {
		util.HandleError(c, err)
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
// @Success 200 {object} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [get]
func (h *ProgramHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	program, err := h.programService.GetByID(id)
	if err != nil {
		util.HandleError(c, err)
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
// @Success 200 {array} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/merchant/{merchant_id} [get]
func (h *ProgramHandler) GetByMerchantID(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "merchant_id",
			Message: "invalid merchant ID",
		})
		return
	}

	programs, err := h.programService.GetByMerchantID(merchantID)
	if err != nil {
		util.HandleError(c, err)
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
// @Param program body domain.UpdateProgramRequest true "Program details"
// @Success 200 {object} domain.Program
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /programs/{id} [put]
func (h *ProgramHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	var req domain.UpdateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	program, err := h.programService.Update(id, &req)
	if err != nil {
		util.HandleError(c, err)
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
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid program ID",
		})
		return
	}

	if err := h.programService.Delete(id); err != nil {
		util.HandleError(c, err)
		return
	}

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
