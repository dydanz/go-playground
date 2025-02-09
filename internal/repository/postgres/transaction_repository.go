package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go-playground/internal/config"
	"go-playground/internal/domain"
	"log"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db config.DbConnection
}

func NewTransactionRepository(db config.DbConnection) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			transaction_id, merchant_id, merchant_customers_id, program_id,
			transaction_type, transaction_amount, transaction_date,
			branch_id, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7, $8, CURRENT_TIMESTAMP)
		RETURNING transaction_date, created_at
	`
	return r.db.RW.QueryRowContext(
		ctx,
		query,
		tx.TransactionID,
		tx.MerchantID,
		tx.MerchantCustomersID,
		tx.ProgramID,
		tx.TransactionType,
		tx.TransactionAmount,
		tx.BranchID,
		tx.Status,
	).Scan(
		&tx.TransactionDate,
		&tx.CreatedAt,
	)
}

func (r *TransactionRepository) GetByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT transaction_id, merchant_id, merchant_customers_id, program_id,
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, status, created_at
		FROM transactions
		WHERE transaction_id = $1
	`
	tx := &domain.Transaction{}
	err := r.db.RR.QueryRowContext(ctx, query, transactionID).Scan(
		&tx.TransactionID,
		&tx.MerchantID,
		&tx.MerchantCustomersID,
		&tx.ProgramID,
		&tx.TransactionType,
		&tx.TransactionAmount,
		&tx.TransactionDate,
		&tx.BranchID,
		&tx.Status,
		&tx.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tx, err
}

func (r *TransactionRepository) GetByCustomerID(ctx context.Context, merchantCustomersID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, _, err := r.GetByCustomerIDWithPagination(ctx, merchantCustomersID, 0, -1)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByCustomerIDWithPagination(ctx context.Context, merchantCustomersID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM transactions WHERE merchant_customers_id = $1`
	err := r.db.RR.QueryRowContext(ctx, countQuery, merchantCustomersID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	log.Printf("function GetByCustomerIDWithPagination called with merchantCustomersID: %v, offset: %v, limit: %v", merchantCustomersID, offset, limit)

	// Build query with pagination
	query := `
		SELECT transaction_id, merchant_id, merchant_customers_id, program_id,
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, status, created_at
		FROM transactions
		WHERE merchant_customers_id = $1
		ORDER BY transaction_date DESC
	`

	var rows *sql.Rows
	if limit > 0 {
		query += ` LIMIT $2 OFFSET $3`
		rows, err = r.db.RR.QueryContext(ctx, query, merchantCustomersID, limit, offset)
	} else {
		rows, err = r.db.RR.QueryContext(ctx, query, merchantCustomersID)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.TransactionID,
			&tx.MerchantID,
			&tx.MerchantCustomersID,
			&tx.ProgramID,
			&tx.TransactionType,
			&tx.TransactionAmount,
			&tx.TransactionDate,
			&tx.BranchID,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, total, nil
}

func (r *TransactionRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Transaction, error) {
	query := `
		SELECT transaction_id, merchant_id, merchant_customers_id, program_id,
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, status, created_at
		FROM transactions
		WHERE merchant_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := r.db.RR.QueryContext(ctx, query, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.TransactionID,
			&tx.MerchantID,
			&tx.MerchantCustomersID,
			&tx.ProgramID,
			&tx.TransactionType,
			&tx.TransactionAmount,
			&tx.TransactionDate,
			&tx.BranchID,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// Notes: Table Transactions should be can not be updated/deleted.
func (r *TransactionRepository) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE transaction_id = $2`
	result, err := r.db.RW.ExecContext(ctx, query, status, transactionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}
