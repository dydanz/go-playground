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

func (r *PointsRepository) Create(ctx context.Context, ledger *domain.PointsLedger) error {
	/*
		The CTE (last_balance) fetches the points_balance of the last transaction
		for the same customer_id and program_id.

		If no previous record exists, COALESCE(..., 0) ensures we start from 0.

		The new points_balance is calculated as:
		previous_balance + points_earned - points_redeemed.

		The record is atomically inserted into the points_ledger table.
	*/
	query := `
		INSERT INTO points_ledger (
			ledger_id, customer_id, program_id, 
			points_earned, points_redeemed, points_balance, 
			transaction_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	query = `WITH last_balance AS (
			SELECT points_balance
			FROM points_ledger
			WHERE customer_id = $1
			AND program_id = $2
			ORDER BY created_at DESC
			LIMIT 1
		)
		INSERT INTO points_ledger (customer_id, program_id, points_earned, points_redeemed, points_balance, transaction_id)
		VALUES (
			$1, -- customer_id
			$2, -- program_id
			$3, -- points_earned
			$4, -- points_redeemed
			COALESCE((SELECT points_balance FROM last_balance), 0) + $3 - $4, -- Calculate new balance
			$5  -- transaction_id
		)
		RETURNING created_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		ledger.CustomerID,
		ledger.ProgramID,
		ledger.PointsEarned,
		ledger.PointsRedeemed,
		ledger.TransactionID,
	).Scan(&ledger.CreatedAt)
}

func (r *PointsRepository) GetByCustomerAndProgram(ctx context.Context, customerID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	query := `
		SELECT ledger_id, customer_id, program_id, 
			   points_earned, points_redeemed, points_balance, 
			   transaction_id, created_at
		FROM points_ledger
		WHERE customer_id = $1 AND program_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ledgers []*domain.PointsLedger
	for rows.Next() {
		ledger := &domain.PointsLedger{}
		err := rows.Scan(
			&ledger.LedgerID,
			&ledger.CustomerID,
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

func (r *PointsRepository) GetCurrentBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error) {
	query := `
		SELECT points_balance
		FROM points_ledger
		WHERE customer_id = $1 AND program_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var balance int
	err := r.db.QueryRowContext(ctx, query, customerID, programID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return balance, err
}

func (r *PointsRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*domain.PointsLedger, error) {
	query := `
		SELECT ledger_id, customer_id, program_id, 
			   points_earned, points_redeemed, points_balance, 
			   transaction_id, created_at
		FROM points_ledger
		WHERE transaction_id = $1
	`
	ledger := &domain.PointsLedger{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&ledger.LedgerID,
		&ledger.CustomerID,
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
