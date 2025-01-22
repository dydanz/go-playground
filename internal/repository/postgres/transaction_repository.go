package postgres

import (
	"database/sql"
	"go-playground/internal/config"
	"go-playground/internal/domain"
)

type TransactionRepository struct {
	db config.DbConnection
}

func NewTransactionRepository(db config.DbConnection) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			user_id, transaction_type, points, description, status, 
			transaction_date
		) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, transaction_date, created_at, updated_at
	`
	return r.db.RW.QueryRow(
		query,
		tx.UserID,
		tx.TransactionType,
		tx.Points,
		tx.Description,
		tx.Status,
	).Scan(
		&tx.ID,
		&tx.TransactionDate,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
}

func (r *TransactionRepository) GetByID(id string) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	query := `
		SELECT id, user_id, transaction_type, points, description, 
			   status, transaction_date, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`
	err := r.db.RR.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.TransactionType,
		&tx.Points,
		&tx.Description,
		&tx.Status,
		&tx.TransactionDate,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tx, err
}

func (r *TransactionRepository) GetByUserID(userID string) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, points, description, 
			   status, transaction_date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := r.db.RR.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.TransactionType,
			&tx.Points,
			&tx.Description,
			&tx.Status,
			&tx.TransactionDate,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (r *TransactionRepository) Update(tx *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.RW.QueryRow(query, tx.Status, tx.ID).Scan(&tx.UpdatedAt)
}
