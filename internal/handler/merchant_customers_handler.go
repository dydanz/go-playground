package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchantCustomersHandler struct {
	customerService domain.MerchantCustomersService
}

func NewMerchantCustomersHandler(customerService domain.MerchantCustomersService) *MerchantCustomersHandler {
	return &MerchantCustomersHandler{customerService: customerService}
}

// Create godoc
// @Summary Create merchant customer
// @Description Create a new merchant customer
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param request body domain.CreateMerchantCustomerRequest true "Customer details"
// @Success 201 {object} domain.MerchantCustomer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers [post]
func (h *MerchantCustomersHandler) Create(c *gin.Context) {
	var req domain.CreateMerchantCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	customer, err := h.customerService.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetByID godoc
// @Summary Get merchant customer by ID
// @Description Get merchant customer details by ID
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.MerchantCustomer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers/{id} [get]
func (h *MerchantCustomersHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid customer ID format"})
		return
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetByMerchantID godoc
// @Summary Get merchant customers by merchant ID
// @Description Get all customers for a specific merchant
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {array} domain.MerchantCustomer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers/merchant/{merchant_id} [get]
func (h *MerchantCustomersHandler) GetByMerchantID(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchant_id"))
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid merchant ID format"})
		return
	}

	customers, err := h.customerService.GetByMerchantID(c.Request.Context(), merchantID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customers)
}

// Update godoc
// @Summary Update merchant customer
// @Description Update merchant customer details
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Customer ID"
// @Param request body domain.UpdateMerchantCustomerRequest true "Customer details to update"
// @Success 200 {object} domain.MerchantCustomer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers/{id} [put]
func (h *MerchantCustomersHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid customer ID format"})
		return
	}

	var req domain.UpdateMerchantCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	customer, err := h.customerService.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Delete godoc
// @Summary Delete merchant customer
// @Description Delete a merchant customer
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Customer ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers/{id} [delete]
func (h *MerchantCustomersHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid customer ID format"})
		return
	}

	if err := h.customerService.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ValidateCredentials godoc
// @Summary Validate customer credentials
// @Description Validate merchant customer login credentials
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param request body domain.CustomerLoginRequest true "Login credentials"
// @Success 200 {object} domain.MerchantCustomer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /merchant-customers/login [post]
func (h *MerchantCustomersHandler) ValidateCredentials(c *gin.Context) {
	var req domain.CustomerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	customer, err := h.customerService.ValidateCredentials(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
