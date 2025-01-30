package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go-playground/internal/config"
	"go-playground/internal/domain"

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
			transaction_id, merchant_id, customer_id, 
			transaction_type, transaction_amount, transaction_date,
			branch_id
		) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6)
		RETURNING transaction_date, created_at
	`
	return r.db.RW.QueryRowContext(
		ctx,
		query,
		tx.TransactionID,
		tx.MerchantID,
		tx.CustomerID,
		tx.TransactionType,
		tx.TransactionAmount,
		tx.BranchID,
	).Scan(
		&tx.TransactionDate,
		&tx.CreatedAt,
	)
}

func (r *TransactionRepository) GetByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT transaction_id, merchant_id, customer_id, 
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, created_at
		FROM transactions
		WHERE transaction_id = $1
	`
	tx := &domain.Transaction{}
	err := r.db.RR.QueryRowContext(ctx, query, transactionID).Scan(
		&tx.TransactionID,
		&tx.MerchantID,
		&tx.CustomerID,
		&tx.TransactionType,
		&tx.TransactionAmount,
		&tx.TransactionDate,
		&tx.BranchID,
		&tx.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tx, err
}

func (r *TransactionRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Transaction, error) {
	query := `
		SELECT transaction_id, merchant_id, customer_id, 
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, created_at
		FROM transactions
		WHERE customer_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := r.db.RR.QueryContext(ctx, query, customerID)
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
			&tx.CustomerID,
			&tx.TransactionType,
			&tx.TransactionAmount,
			&tx.TransactionDate,
			&tx.BranchID,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Transaction, error) {
	query := `
		SELECT transaction_id, merchant_id, customer_id, 
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, created_at
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
			&tx.CustomerID,
			&tx.TransactionType,
			&tx.TransactionAmount,
			&tx.TransactionDate,
			&tx.BranchID,
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
