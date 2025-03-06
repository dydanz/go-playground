package postgres

import (
	"context"
	"database/sql"
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type MerchantCustomersRepository struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewMerchantCustomersRepository(db *sql.DB) *MerchantCustomersRepository {
	return &MerchantCustomersRepository{
		db:     db,
		logger: logging.GetLogger(),
	}
}

// Create inserts a new merchant customer into the database
func (r *MerchantCustomersRepository) Create(ctx context.Context, customer *domain.MerchantCustomer) error {
	query := `
		INSERT INTO merchant_customers (id, merchant_id, email, password, name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	now := time.Now().UTC()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.MerchantID,
		customer.Email,
		customer.Password,
		customer.Name,
		customer.Phone,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				r.logger.Error().
					Str("error", pqErr.Code.Name()).
					Msg("Unique violation error")
				return domain.NewResourceConflictError("merchant customer", "email or phone already exists")
			case "foreign_key_violation":
				r.logger.Error().
					Str("error", pqErr.Code.Name()).
					Msg("Foreign key violation error")
				return domain.NewValidationError("merchant_id", "invalid merchant ID")
			}
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to create merchant customer")
		return domain.NewSystemError("MerchantCustomersRepository.Create", err, "database error")
	}

	return nil
}

// GetByID retrieves a merchant customer by their ID
func (r *MerchantCustomersRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE id = $1
	`

	customer := &domain.MerchantCustomer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.MerchantID,
		&customer.Email,
		&customer.Password,
		&customer.Name,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("id", id.String()).
				Msg("No merchant customer found")
			return nil, nil
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to get merchant customer by ID")
		return nil, domain.NewSystemError("MerchantCustomersRepository.GetByID", err, "database error")
	}

	return customer, nil
}

// GetByEmail retrieves a merchant customer by their email
func (r *MerchantCustomersRepository) GetByEmail(ctx context.Context, email string) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE email = $1
	`

	customer := &domain.MerchantCustomer{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&customer.ID,
		&customer.MerchantID,
		&customer.Email,
		&customer.Password,
		&customer.Name,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("email", email).
				Msg("No merchant customer found")
			return nil, nil
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to get merchant customer by email")
		return nil, domain.NewSystemError("MerchantCustomersRepository.GetByEmail", err, "database error")
	}

	return customer, nil
}

// GetByPhone retrieves a merchant customer by their phone number
func (r *MerchantCustomersRepository) GetByPhone(ctx context.Context, phone string) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE phone = $1
	`

	customer := &domain.MerchantCustomer{}
	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&customer.ID,
		&customer.MerchantID,
		&customer.Email,
		&customer.Password,
		&customer.Name,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("phone", phone).
				Msg("No merchant customer found")
			return nil, nil
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to get merchant customer by phone")
		return nil, domain.NewSystemError("MerchantCustomersRepository.GetByPhone", err, "database error")
	}

	return customer, nil
}

// GetByMerchantID retrieves all customers for a given merchant
func (r *MerchantCustomersRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE merchant_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, merchantID)
	if err != nil {
		return nil, domain.NewSystemError("MerchantCustomersRepository.GetByMerchantID", err, "database error")
	}
	defer rows.Close()

	var customers []*domain.MerchantCustomer
	for rows.Next() {
		customer := &domain.MerchantCustomer{}
		err := rows.Scan(
			&customer.ID,
			&customer.MerchantID,
			&customer.Email,
			&customer.Password,
			&customer.Name,
			&customer.Phone,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().
				Err(err).
				Msg("Failed to get merchant customers by merchant ID")
			return nil, domain.NewSystemError("MerchantCustomersRepository.GetByMerchantID", err, "error scanning row")
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to get merchant customers by merchant ID")
		return nil, domain.NewSystemError("MerchantCustomersRepository.GetByMerchantID", err, "error iterating rows")
	}

	return customers, nil
}

// Update updates an existing merchant customer
func (r *MerchantCustomersRepository) Update(ctx context.Context, customer *domain.MerchantCustomer) error {
	query := `
		UPDATE merchant_customers
		SET email = $1, password = $2, name = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		customer.Email,
		customer.Password,
		customer.Name,
		customer.Phone,
		customer.ID,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				r.logger.Error().
					Str("error", pqErr.Code.Name()).
					Msg("Unique violation error")
				return domain.NewResourceConflictError("merchant customer", "email or phone already exists")
			}
		}
		return domain.NewSystemError("MerchantCustomersRepository.Update", err, "database error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("MerchantCustomersRepository.Update", err, "error getting affected rows")
	}

	if rows == 0 {
		r.logger.Warn().
			Str("id", customer.ID.String()).
			Msg("No merchant customer found")
		return domain.NewResourceNotFoundError("merchant customer", customer.ID.String(), "customer not found")
	}

	return nil
}
