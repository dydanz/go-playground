package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TransactionHandler struct {
	transactionService domain.TransactionService
	logger             zerolog.Logger
}

func NewTransactionHandler(transactionService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logging.GetLogger(),
	}
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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create transaction request")

	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create transaction request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	transaction, err := h.transactionService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create transaction")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("transaction_id", transaction.TransactionID.String()).
		Str("merchant_id", transaction.MerchantID.String()).
		Str("merchant_customers_id", transaction.MerchantCustomersID.String()).
		Float64("transaction_amount", transaction.TransactionAmount).
		Msg("Transaction created successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get transaction request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing transaction ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid transaction ID",
		})
		return
	}

	transaction, err := h.transactionService.GetByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("transaction_id", id).
			Msg("Failed to get transaction")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("transaction_id", transaction.TransactionID.String()).
		Str("merchant_id", transaction.MerchantID.String()).
		Str("merchant_customers_id", transaction.MerchantCustomersID.String()).
		Msg("Transaction retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get customer transactions request")

	customerID := c.Param("user_id")
	if customerID == "" {
		h.logger.Error().
			Msg("Missing customer ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "customer_id",
			Message: "invalid customer ID",
		})
		return
	}

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

	offset := (page - 1) * limit

	h.logger.Debug().
		Str("customer_id", customerID).
		Int("page", page).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Fetching customer transactions")

	transactions, total, err := h.transactionService.GetByCustomerIDWithPagination(c.Request.Context(), uuid.MustParse(customerID), offset, limit)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("customer_id", customerID).
			Msg("Failed to get customer transactions")
		util.HandleError(c, err)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	h.logger.Info().
		Str("customer_id", customerID).
		Int("transactions_count", len(transactions)).
		Int64("total_transactions", total).
		Int("page", page).
		Int("total_pages", int(totalPages)).
		Msg("Customer transactions retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get merchant transactions request")

	merchantID := c.Param("merchant_id")
	if merchantID == "" {
		h.logger.Error().
			Msg("Missing merchant ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "merchant_id",
			Message: "invalid merchant ID",
		})
		return
	}

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

	offset := (page - 1) * limit

	var transactions []*domain.Transaction
	var total int64
	var err error

	if merchantID == "null" || merchantID == "all" || merchantID == "" {
		userIDStr, _ := c.Get("user_id")
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			h.logger.Error().
				Err(err).
				Interface("user_id", userIDStr).
				Msg("Invalid user ID format")
			util.HandleError(c, err)
			return
		}
		h.logger.Debug().
			Str("user_id", userID.String()).
			Int("page", page).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Fetching user transactions")

		transactions, total, err = h.transactionService.GetByUserIDWithPagination(c.Request.Context(), userID, offset, limit)
	} else {
		h.logger.Debug().
			Str("merchant_id", merchantID).
			Int("page", page).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Fetching merchant transactions")

		transactions, total, err = h.transactionService.GetByMerchantIDWithPagination(c.Request.Context(), uuid.MustParse(merchantID), offset, limit)
	}

	if err != nil {
		h.logger.Error().
			Err(err).
			Str("merchant_id", merchantID).
			Msg("Failed to get transactions")
		util.HandleError(c, err)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	h.logger.Info().
		Str("merchant_id", merchantID).
		Int("transactions_count", len(transactions)).
		Int64("total_transactions", total).
		Int("page", page).
		Int("total_pages", int(totalPages)).
		Msg("Transactions retrieved successfully")

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

func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update transaction status request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing transaction ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid transaction ID",
		})
		return
	}

	var req domain.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update transaction status request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	if err := h.transactionService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		h.logger.Error().
			Err(err).
			Str("transaction_id", id).
			Str("status", req.Status).
			Msg("Failed to update transaction status")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("transaction_id", id).
		Str("status", req.Status).
		Msg("Transaction status updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}
