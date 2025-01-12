package postgres

import (
	"database/sql"
	"go-cursor/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (email, password, name, phone, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, name, phone, status, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.Name,
		user.Phone,
		user.Status,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, email, password, name, phone, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, 
			phone = $2, 
			status = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Name,
		user.Phone,
		user.Status,
		user.ID,
	).Scan(
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) Delete(id int64) error {
	// Implement delete user
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetAll() ([]*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, created_at, updated_at
		FROM users
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) UpdateTx(tx *sql.Tx, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, 
			phone = $2, 
			status = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	err := tx.QueryRow(
		query,
		user.Name,
		user.Phone,
		user.Status,
		user.ID,
	).Scan(
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}
