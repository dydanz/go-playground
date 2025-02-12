package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go-playground/internal/domain"
	"log"
	"time"

	"github.com/google/uuid"
)

type MerchantCustomersRepository struct {
	db *sql.DB
}

func NewMerchantCustomersRepository(db *sql.DB) *MerchantCustomersRepository {
	return &MerchantCustomersRepository{db: db}
}

// Create inserts a new merchant customer into the database
func (r *MerchantCustomersRepository) Create(ctx context.Context, customer *domain.MerchantCustomer) error {
	query := `
		INSERT INTO merchant_customers (merchant_id, email, password, name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`

	log.Println("Creating merchant customer:", customer)

	err := r.db.QueryRowContext(ctx, query,
		customer.MerchantID,
		customer.Email,
		customer.Password,
		customer.Name,
		customer.Phone,
	).Scan(&customer.ID)

	if err != nil {
		log.Println("error creating merchant customer", err)
		return err
	}

	return nil
}

// GetByID retrieves a merchant customer by their ID
func (r *MerchantCustomersRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE id = $1`

	log.Println("Getting merchant customer by id:", id)

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

	if err == sql.ErrNoRows {
		log.Println("error getting merchant customer by id", err)
		return nil, errors.New("merchant customer not found")
	}
	if err != nil {
		log.Println("error getting merchant customer by id", err)
		return nil, err
	}

	return customer, nil
}

// GetByEmail retrieves a merchant customer by their email
func (r *MerchantCustomersRepository) GetByEmail(ctx context.Context, email string) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE email = $1`

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

	if err == sql.ErrNoRows {
		log.Println("error getting merchant customer by email", err)
		return nil, errors.New("merchant customer not found")
	}
	if err != nil {
		log.Println("error getting merchant customer by email", err)
		return nil, err
	}

	return customer, nil
}

// GetByPhone retrieves a merchant customer by their phone number
func (r *MerchantCustomersRepository) GetByPhone(ctx context.Context, phone string) (*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE phone = $1`

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

	if err == sql.ErrNoRows {
		log.Println("error getting merchant customer by phone", err)
		return nil, errors.New("merchant customer not found")
	}
	if err != nil {
		log.Println("error getting merchant customer by phone", err)
		return nil, err
	}

	return customer, nil
}

// GetByMerchantID retrieves all customers for a given merchant
func (r *MerchantCustomersRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.MerchantCustomer, error) {
	query := `
		SELECT id, merchant_id, email, password, name, phone, created_at, updated_at
		FROM merchant_customers
		WHERE merchant_id = $1`

	rows, err := r.db.QueryContext(ctx, query, merchantID)
	if err != nil {
		return nil, err
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
			log.Println("error getting merchant customer by merchant id", err)
			return nil, err
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		log.Println("error getting merchant customer by merchant id", err)
		return nil, err
	}

	return customers, nil
}

// Update updates an existing merchant customer
func (r *MerchantCustomersRepository) Update(ctx context.Context, customer *domain.MerchantCustomer) error {
	query := `
		UPDATE merchant_customers
		SET merchant_id = $1, email = $2, password = $3, name = $4, phone = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at`

	customer.UpdatedAt = time.Now().UTC()

	err := r.db.QueryRowContext(ctx, query,
		customer.MerchantID,
		customer.Email,
		customer.Password,
		customer.Name,
		customer.Phone,
		customer.UpdatedAt,
		customer.ID,
	).Scan(&customer.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Println("error updating merchant customer", err)
		return errors.New("merchant customer not found")
	}
	if err != nil {
		return err
	}

	return nil
}
