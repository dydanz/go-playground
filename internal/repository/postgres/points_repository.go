package postgres

import (
	"context"
	"database/sql"

	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type PointsRepository struct {
	db *sql.DB
}

func NewPointsRepository(db *sql.DB) *PointsRepository {
	return &PointsRepository{db: db}
}

// Create inserts a new points ledger entry into the database
func (r *PointsRepository) Create(ctx context.Context, ledger *domain.PointsLedger) error {
	/*
		The CTE (last_balance) fetches the points_balance of the last transaction
		for the same merchant_customers_id and program_id.

		If no previous record exists, COALESCE(..., 0) ensures we start from 0.

		The new points_balance is calculated as:
		previous_balance + points_earned - points_redeemed.

		The record is atomically inserted into the points_ledger table.
	*/
	query := `WITH last_balance AS (
			SELECT points_balance
			FROM points_ledger
			WHERE merchant_customers_id = $1
			AND program_id = $2
			ORDER BY created_at DESC
			LIMIT 1
		)
		INSERT INTO points_ledger (
			merchant_customers_id,
			program_id,
			points_earned,
			points_redeemed,
			points_balance,
			transaction_id,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			COALESCE((SELECT points_balance FROM last_balance), 0) + $3 - $4,
			$5,
			CURRENT_TIMESTAMP
		)
		RETURNING created_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		ledger.MerchantCustomersID,
		ledger.ProgramID,
		ledger.PointsEarned,
		ledger.PointsRedeemed,
		ledger.TransactionID,
	).Scan(&ledger.CreatedAt)
}

// GetByCustomerAndProgram retrieves all points ledger entries for a given customer and program
func (r *PointsRepository) GetByCustomerAndProgram(ctx context.Context, merchantCustomersID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	query := `
		SELECT ledger_id,
			   merchant_customers_id,
			   program_id,
			   points_earned,
			   points_redeemed,
			   points_balance,
			   transaction_id,
			   created_at
		FROM points_ledger
		WHERE merchant_customers_id = $1 AND program_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, merchantCustomersID, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ledgers []*domain.PointsLedger
	for rows.Next() {
		ledger := &domain.PointsLedger{}
		err := rows.Scan(
			&ledger.LedgerID,
			&ledger.MerchantCustomersID,
			&ledger.ProgramID,
			&ledger.PointsEarned,
			&ledger.PointsRedeemed,
			&ledger.PointsBalance,
			&ledger.TransactionID,
			&ledger.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ledgers = append(ledgers, ledger)
	}
	return ledgers, nil
}

// GetCurrentBalance retrieves the current points balance for a given customer and program
func (r *PointsRepository) GetCurrentBalance(ctx context.Context, merchantCustomersID, programID uuid.UUID) (int, error) {
	query := `
		SELECT points_balance
		FROM points_ledger
		WHERE merchant_customers_id = $1 AND program_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var balance int
	err := r.db.QueryRowContext(ctx, query, merchantCustomersID, programID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return balance, err
}

// GetByTransactionID retrieves a points ledger entry by its transaction ID
func (r *PointsRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*domain.PointsLedger, error) {
	query := `
		SELECT ledger_id,
			   merchant_customers_id,
			   program_id,
			   points_earned,
			   points_redeemed,
			   points_balance,
			   transaction_id,
			   created_at
		FROM points_ledger
		WHERE transaction_id = $1
	`
	ledger := &domain.PointsLedger{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&ledger.LedgerID,
		&ledger.MerchantCustomersID,
		&ledger.ProgramID,
		&ledger.PointsEarned,
		&ledger.PointsRedeemed,
		&ledger.PointsBalance,
		&ledger.TransactionID,
		&ledger.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ledger, err
}
