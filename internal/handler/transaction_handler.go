package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService domain.TransactionService
}

func NewTransactionHandler(transactionService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// CreateTransaction godoc
// @Summary Create transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param transaction body domain.CreateTransactionRequest true "Transaction details"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	transaction, err := h.transactionService.Create(&req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Description Get transaction details by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid transaction ID",
		})
		return
	}

	transaction, err := h.transactionService.GetByID(id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetCustomerTransactions godoc
// @Summary Get customer transactions
// @Description Get all transactions for a specific customer
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/customer/{customer_id} [get]
func (h *TransactionHandler) GetByCustomerID(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "customer_id",
			Message: "invalid customer ID",
		})
		return
	}

	transactions, err := h.transactionService.GetByCustomerID(customerID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetMerchantTransactions godoc
// @Summary Get merchant transactions
// @Description Get all transactions for a specific merchant
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/merchant/{merchant_id} [get]
func (h *TransactionHandler) GetByMerchantID(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "merchant_id",
			Message: "invalid merchant ID",
		})
		return
	}

	transactions, err := h.transactionService.GetByMerchantID(merchantID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid transaction ID",
		})
		return
	}

	var req domain.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	if err := h.transactionService.UpdateStatus(id, req.Status); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}
