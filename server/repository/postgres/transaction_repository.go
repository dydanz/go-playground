package postgres

import (
	"context"
	"database/sql"
	"go-playground/server/config"
	"go-playground/server/domain"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db config.DbConnection
}

func NewTransactionRepository(db config.DbConnection) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	query := `
		INSERT INTO transactions (
			merchant_id, merchant_customers_id, program_id,
			transaction_type, transaction_amount, transaction_date,
			branch_id, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		RETURNING transaction_id, transaction_date, created_at
	`
	createdTx := &domain.Transaction{}
	err := r.db.RW.QueryRowContext(
		ctx,
		query,
		tx.MerchantID,
		tx.MerchantCustomersID,
		tx.ProgramID,
		tx.TransactionType,
		tx.TransactionAmount,
		tx.TransactionDate,
		tx.BranchID,
		tx.Status,
	).Scan(
		&createdTx.TransactionID,
		&createdTx.TransactionDate,
		&createdTx.CreatedAt,
	)
	if err != nil {
		if isPgUniqueViolation(err) {
			return nil, domain.NewResourceConflictError("transaction", "duplicate transaction record")
		}
		return nil, domain.NewSystemError("TransactionRepository.Create", err, "failed to create transaction")
	}

	return createdTx, nil
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
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, domain.NewSystemError("TransactionRepository.GetByID", err, "failed to get transaction")
	}
	return tx, nil
}

func (r *TransactionRepository) GetByCustomerID(ctx context.Context, merchantCustomersID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, _, err := r.GetByCustomerIDWithPagination(ctx, merchantCustomersID, 0, -1)
	if err != nil {
		return nil, domain.NewSystemError("TransactionRepository.GetByCustomerID", err, "failed to get transactions")
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByCustomerIDWithPagination(ctx context.Context, merchantCustomersID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM transactions WHERE merchant_customers_id = $1`
	err := r.db.RR.QueryRowContext(ctx, countQuery, merchantCustomersID).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByCustomerIDWithPagination", err, "failed to get total count")
	}

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
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByCustomerIDWithPagination", err, "failed to query transactions")
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
			return nil, 0, domain.NewSystemError("TransactionRepository.GetByCustomerIDWithPagination", err, "failed to scan transaction")
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByCustomerIDWithPagination", err, "error iterating transactions")
	}

	return transactions, total, nil
}

func (r *TransactionRepository) GetByMerchantIDWithPagination(ctx context.Context, merchantID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM transactions WHERE merchant_id = $1`
	err := r.db.RR.QueryRowContext(ctx, countQuery, merchantID).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByMerchantIDWithPagination", err, "failed to get total count")
	}

	// Build query with pagination
	query := `
		SELECT transaction_id, merchant_id, merchant_customers_id, program_id,
			   transaction_type, transaction_amount, transaction_date,
			   branch_id, status, created_at
		FROM transactions
		WHERE merchant_id = $1
		ORDER BY transaction_date DESC
	`

	var rows *sql.Rows
	if limit > 0 {
		query += ` LIMIT $2 OFFSET $3`
		rows, err = r.db.RR.QueryContext(ctx, query, merchantID, limit, offset)
	} else {
		rows, err = r.db.RR.QueryContext(ctx, query, merchantID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByMerchantIDWithPagination", err, "failed to query transactions")
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
			return nil, 0, domain.NewSystemError("TransactionRepository.GetByMerchantIDWithPagination", err, "failed to scan transaction")
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByMerchantIDWithPagination", err, "error iterating transactions")
	}

	return transactions, total, nil
}

func (r *TransactionRepository) GetByUserIDWithPagination(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Transaction, int64, error) {
	// Get total count using a subquery to join merchants and transactions
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM transactions t
		INNER JOIN merchants m ON t.merchant_id = m.id
		WHERE m.user_id = $1
	`
	err := r.db.RR.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByUserIDWithPagination", err, "failed to get total count")
	}

	// Build query with pagination
	query := `
		SELECT t.transaction_id, t.merchant_id, t.merchant_customers_id, t.program_id,
			   t.transaction_type, t.transaction_amount, t.transaction_date,
			   t.branch_id, t.status, t.created_at
		FROM transactions t
		INNER JOIN merchants m ON t.merchant_id = m.id
		WHERE m.user_id = $1
		ORDER BY t.transaction_date DESC
	`

	var rows *sql.Rows
	if limit > 0 {
		query += ` LIMIT $2 OFFSET $3`
		rows, err = r.db.RR.QueryContext(ctx, query, userID, limit, offset)
	} else {
		rows, err = r.db.RR.QueryContext(ctx, query, userID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByUserIDWithPagination", err, "failed to query transactions")
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
			return nil, 0, domain.NewSystemError("TransactionRepository.GetByUserIDWithPagination", err, "failed to scan transaction")
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, domain.NewSystemError("TransactionRepository.GetByUserIDWithPagination", err, "error iterating transactions")
	}

	return transactions, total, nil
}

// Notes: Table Transactions should be can not be updated/deleted.
func (r *TransactionRepository) UpdateStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE transaction_id = $2`
	result, err := r.db.RW.ExecContext(ctx, query, status, transactionID)
	if err != nil {
		return domain.NewSystemError("TransactionRepository.UpdateStatus", err, "failed to update transaction status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("TransactionRepository.UpdateStatus", err, "failed to get affected rows")
	}

	if rowsAffected == 0 {
		return domain.NewResourceNotFoundError("transaction", transactionID.String(), "transaction not found")
	}

	return nil
}
