package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// @Summary Create transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param transaction body domain.Transaction true "Transaction details"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	var tx domain.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.transactionService.Create(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// @Summary Get transaction by ID
// @Description Get transaction details by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} domain.Transaction
// @Failure 404 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tx, err := h.transactionService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}

// @Summary Get user transactions
// @Description Get all transactions for a specific user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} domain.Transaction
// @Failure 404 {object} map[string]string
// @Router /transactions/user/{user_id} [get]
func (h *TransactionHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	transactions, err := h.transactionService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
