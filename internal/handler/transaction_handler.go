package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/util"
	"net/http"
	"strconv"

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

	transaction, err := h.transactionService.Create(c.Request.Context(), &req)
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

	transaction, err := h.transactionService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetCustomerTransactions godoc
// @Summary Get customer transactions
// @Description Get all transactions for a specific customer with pagination
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param customer_id path string true "Customer ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/user/{customer_id} [get]
func (h *TransactionHandler) GetByCustomerID(c *gin.Context) {
	customerID := c.Param("user_id")
	if customerID == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "customer_id",
			Message: "invalid customer ID",
		})
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	transactions, total, err := h.transactionService.GetByCustomerIDWithPagination(c.Request.Context(), customerID, offset, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)

	// Prepare response with pagination metadata
	response := gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
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

	transactions, err := h.transactionService.GetByMerchantID(c.Request.Context(), merchantID)
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

	if err := h.transactionService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}
