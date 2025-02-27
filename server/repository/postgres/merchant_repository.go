package postgres

import (
	"context"
	"database/sql"
	"go-playground/server/domain"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(ctx context.Context, merchant *domain.Merchant) (*domain.Merchant, error) {
	query := `INSERT INTO merchants (user_id, merchant_name, merchant_type, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, user_id, merchant_name, merchant_type, created_at, updated_at`

	result := domain.Merchant{}
	err := r.db.QueryRowContext(ctx, query,
		merchant.UserID,
		merchant.Name,
		merchant.Type,
		merchant.CreatedAt,
		merchant.UpdatedAt,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Name,
		&result.Type,
		&result.CreatedAt,
		&result.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if isPgUniqueViolation(err) {
			return nil, domain.NewResourceConflictError("merchant", "merchant with this name already exists")
		}
		// Wrap database errors as system errors
		return nil, domain.NewSystemError("MerchantRepository.Create", err, "failed to create merchant")
	}

	return &result, nil
}

func (r *MerchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	query := `SELECT id, user_id, merchant_name, merchant_type, created_at, updated_at 
			  FROM merchants WHERE id = $1 AND status = 'active'`

	merchant := &domain.Merchant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.Name,
		&merchant.Type,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewResourceNotFoundError("merchant", id.String(), "merchant not found")
		}
		return nil, domain.NewSystemError("MerchantRepository.GetByID", err, "failed to get merchant")
	}

	return merchant, nil
}

func (r *MerchantRepository) GetMerchantsByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Merchant, int, error) {
	// First, get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM merchants WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, domain.NewSystemError("MerchantRepository.GetMerchantsByUserID", err, "failed to get total count")
	}

	// Then get paginated results
	query := `SELECT id, user_id, merchant_name, merchant_type, created_at, updated_at, status
			  FROM merchants 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		log.Printf("Error querying merchants: %v", err)
		return nil, 0, domain.NewSystemError("MerchantRepository.GetMerchantsByUserID", err, "failed to query merchants")
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
			&merchant.Status,
		)
		if err != nil {
			log.Printf("Error scanning merchant: %v", err)
			return nil, 0, domain.NewSystemError("MerchantRepository.GetMerchantsByUserID", err, "failed to scan merchant")
		}
		merchants = append(merchants, merchant)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating merchants: %v", err)
		return nil, 0, domain.NewSystemError("MerchantRepository.GetMerchantsByUserID", err, "error iterating merchants")
	}

	return merchants, total, nil
}

func (r *MerchantRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*domain.MerchantList, error) {
	query := `
		SELECT id, merchant_name
		FROM merchants
		WHERE user_id = $1 AND status = 'active'
		ORDER BY merchant_name
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []*domain.MerchantList
	for rows.Next() {
		merchant := &domain.MerchantList{}
		err := rows.Scan(
			&merchant.ID,
			&merchant.Name,
		)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}
	return merchants, nil
}

func (r *MerchantRepository) Update(ctx context.Context, merchant *domain.Merchant) error {
	query := `
		UPDATE merchants
		SET merchant_name = $1, merchant_type = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at
	`
	merchant.UpdatedAt = time.Now().UTC()

	return r.db.QueryRowContext(
		ctx,
		query,
		merchant.Name,
		merchant.Type,
		merchant.UpdatedAt,
		merchant.ID,
	).Scan(&merchant.UpdatedAt)
}

func (r *MerchantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE merchants
		SET status=$1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		"deactivated",
		id,
	)
	if err != nil {
		return domain.NewSystemError("MerchantRepository.Delete", err, "failed to delete merchant")
	}
	return nil
}

// Helper function to check for PostgreSQL unique constraint violations
func isPgUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // PostgreSQL error code for unique_violation
	}
	return false
}
