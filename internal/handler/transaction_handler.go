package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.transactionService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
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
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tx, err := h.transactionService.GetByID(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
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
	customerID, err := uuid.Parse(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	transactions, err := h.transactionService.GetByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	merchantID, err := uuid.Parse(c.Param("merchant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant ID"})
		return
	}

	transactions, err := h.transactionService.GetByMerchantID(c.Request.Context(), merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
