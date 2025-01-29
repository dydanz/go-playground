package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
)

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(merchant *domain.Merchant) error {
	query := `
		INSERT INTO merchants (merchant_name, merchant_type)
		VALUES ($1, $2)
		RETURNING merchant_id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		merchant.Name,
		merchant.Type,
	).Scan(&merchant.ID, &merchant.CreatedAt, &merchant.UpdatedAt)
}

func (r *MerchantRepository) GetByID(id string) (*domain.Merchant, error) {
	merchant := &domain.Merchant{}
	query := `
		SELECT merchant_id, merchant_name, merchant_type, created_at, updated_at
		FROM merchants
		WHERE merchant_id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&merchant.ID,
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

func (r *MerchantRepository) GetAll() ([]*domain.Merchant, error) {
	query := `
		SELECT merchant_id, merchant_name, merchant_type, created_at, updated_at
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
		SET merchant_name = $1, merchant_type = $2
		WHERE merchant_id = $3
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		merchant.Name,
		merchant.Type,
		merchant.ID,
	).Scan(&merchant.UpdatedAt)
}

func (r *MerchantRepository) Delete(id string) error {
	query := `DELETE FROM merchants WHERE merchant_id = $1`
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