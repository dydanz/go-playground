package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(merchant *domain.Merchant) error {
	query := `
		INSERT INTO merchants (user_id, merchant_name, merchant_type, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	now := time.Now().UTC()
	merchant.CreatedAt = now
	merchant.UpdatedAt = now

	return r.db.QueryRow(
		query,
		merchant.UserID,
		merchant.Name,
		merchant.Type,
	).Scan(&merchant.ID, &merchant.CreatedAt, &merchant.UpdatedAt)
}

func (r *MerchantRepository) GetByID(id uuid.UUID) (*domain.Merchant, error) {
	merchant := &domain.Merchant{}
	query := `
		SELECT id, user_id, merchant_name, merchant_type, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.Name,
		&merchant.Type,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return merchant, err
}

func (r *MerchantRepository) GetByUserID(userID uuid.UUID) ([]*domain.Merchant, error) {
	query := `
		SELECT id, user_id, merchant_name, merchant_type, created_at, updated_at
		FROM merchants
		WHERE user_id = $1
		ORDER BY merchant_name
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []*domain.Merchant
	for rows.Next() {
		merchant := &domain.Merchant{}
		err := rows.Scan(
			&merchant.ID,
			&merchant.UserID,
			&merchant.Name,
			&merchant.Type,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}
	return merchants, nil
}

func (r *MerchantRepository) GetAll() ([]*domain.Merchant, error) {
	query := `
		SELECT id, user_id, merchant_name, merchant_type, created_at, updated_at
		FROM merchants
		ORDER BY merchant_name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []*domain.Merchant
	for rows.Next() {
		merchant := &domain.Merchant{}
		err := rows.Scan(
			&merchant.ID,
			&merchant.UserID,
			&merchant.Name,
			&merchant.Type,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}
	return merchants, nil
}

func (r *MerchantRepository) Update(merchant *domain.Merchant) error {
	query := `
		UPDATE merchants
		SET merchant_name = $1, merchant_type = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at
	`
	merchant.UpdatedAt = time.Now().UTC()

	return r.db.QueryRow(
		query,
		merchant.Name,
		merchant.Type,
		merchant.UpdatedAt,
		merchant.ID,
	).Scan(&merchant.UpdatedAt)
}

func (r *MerchantRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM merchants WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
