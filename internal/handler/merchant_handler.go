package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantService domain.MerchantService
}

func NewMerchantHandler(merchantService domain.MerchantService) *MerchantHandler {
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
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	merchant, err := h.merchantService.Create(&req)
	if err != nil {
		util.HandleError(c, err)
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
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid merchant ID",
		})
		return
	}

	merchant, err := h.merchantService.GetByID(id)
	if err != nil {
		util.HandleError(c, err)
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
		util.HandleError(c, err)
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
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid merchant ID",
		})
		return
	}

	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	merchant, err := h.merchantService.Update(id, &req)
	if err != nil {
		util.HandleError(c, err)
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
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid merchant ID",
		})
		return
	}

	if err := h.merchantService.Delete(id); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Merchant deleted successfully"})
}
