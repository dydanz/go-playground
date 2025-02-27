package handler

import (
	"go-playground/server/domain"
	"go-playground/server/util"
	"log"
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

// @Summary Create merchant customer
// @Description Create a new merchant customer
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Param customer body domain.CreateMerchantCustomerRequest true "Customer details"
// @Success 201 {object} domain.MerchantCustomer
// @Failure 400 {object} util.ErrorResponse
// @Failure 409 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchant-customers [post]
func (h *MerchantCustomersHandler) Create(c *gin.Context) {
	var req domain.CreateMerchantCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("MerchantCustomersHandler: Error binding request: ", err)
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Get merchant ID from context (set by auth middleware)
	if req.MerchantID == uuid.Nil {
		log.Println("MerchantCustomersHandler: Merchant ID not found in request payload")
		util.HandleError(c, domain.NewAuthenticationError("merchant not authenticated"))
		return
	}

	customer, err := h.customerService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Println("MerchantCustomersHandler: Error creating merchant customer: ", err)
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// @Summary Get merchant customer by ID
// @Description Get merchant customer details by ID
// @Tags merchant-customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.MerchantCustomer
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchant-customers/{id} [get]
func (h *MerchantCustomersHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("id", "invalid customer ID format"))
		return
	}

	// Verify merchant has access to this customer
	if err := h.verifyCustomerAccess(c, id); err != nil {
		util.HandleError(c, err)
		return
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Get merchant customers by merchant ID
// @Description Get all customers for a merchant
// @Tags merchant-customers
// @Produce json
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {array} domain.MerchantCustomer
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchant-customers/merchant/{merchant_id} [get]
func (h *MerchantCustomersHandler) GetByMerchantID(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchant_id"))
	if err != nil {
		util.HandleError(c, domain.ValidationError{Message: "invalid customer ID format"})
		return
	}

	customers, err := h.customerService.GetByMerchantID(c.Request.Context(), merchantID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	if len(customers) == 0 {
		util.EmptyResponse(c)
		return
	}

	c.JSON(http.StatusOK, customers)
}

// @Summary Update merchant customer
// @Description Update merchant customer details
// @Tags merchant-customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body domain.UpdateMerchantCustomerRequest true "Customer details"
// @Success 200 {object} domain.MerchantCustomer
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 409 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /merchant-customers/{id} [put]
func (h *MerchantCustomersHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.HandleError(c, domain.NewValidationError("id", "invalid customer ID format"))
		return
	}

	var req domain.UpdateMerchantCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.NewValidationError("request", "invalid request format"))
		return
	}

	// Verify merchant has access to this customer
	if err := h.verifyCustomerAccess(c, id); err != nil {
		util.HandleError(c, err)
		return
	}

	customer, err := h.customerService.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Helper function to verify customer access
func (h *MerchantCustomersHandler) verifyCustomerAccess(c *gin.Context, customerID uuid.UUID) error {
	merchantID, exists := c.Get("merchant_id")
	if !exists {
		return domain.NewAuthenticationError("merchant not authenticated")
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), customerID)
	if err != nil {
		return err
	}

	if customer.MerchantID != merchantID.(uuid.UUID) {
		return domain.NewAuthorizationError("merchant does not have permission to access this customer")
	}

	return nil
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
