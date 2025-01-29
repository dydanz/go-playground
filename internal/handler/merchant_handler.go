package handler

import (
	"database/sql"
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantService *service.MerchantService
}

func NewMerchantHandler(merchantService *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantService: merchantService}
}

// @Summary Create merchant
// @Description Create a new merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body domain.CreateMerchantRequest true "Merchant details"
// @Success 201 {object} domain.Merchant
// @Failure 400 {object} map[string]string
// @Router /merchants [post]
func (h *MerchantHandler) Create(c *gin.Context) {
	var req domain.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := h.merchantService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, merchant)
}

// @Summary Get merchant by ID
// @Description Get merchant details by ID
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} domain.Merchant
// @Failure 404 {object} map[string]string
// @Router /merchants/{id} [get]
func (h *MerchantHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	merchant, err := h.merchantService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if merchant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// @Summary Get all merchants
// @Description Get all merchants
// @Tags merchants
// @Produce json
// @Success 200 {array} domain.Merchant
// @Router /merchants [get]
func (h *MerchantHandler) GetAll(c *gin.Context) {
	merchants, err := h.merchantService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchants)
}

// @Summary Update merchant
// @Description Update merchant details
// @Tags merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Param merchant body domain.UpdateMerchantRequest true "Merchant details"
// @Success 200 {object} domain.Merchant
// @Failure 404 {object} map[string]string
// @Router /merchants/{id} [put]
func (h *MerchantHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, err := h.merchantService.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if merchant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// @Summary Delete merchant
// @Description Delete a merchant
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /merchants/{id} [delete]
func (h *MerchantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.merchantService.Delete(id)
	if err == nil {
		c.Status(http.StatusNoContent)
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
