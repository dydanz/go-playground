package postgres

import (
	"database/sql"
	"fmt"
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

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	user := &domain.User{}
	var statusStr string
	query := `
		SELECT id, name, email, phone, password, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&statusStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Status = domain.UserStatus(statusStr)
	return user, nil
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

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	var statusStr string
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
		&statusStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user.Status = domain.UserStatus(statusStr)
	return user, err
}

func (r *UserRepository) GetAll() ([]*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
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
		var statusStr string
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Phone,
			&statusStr,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Status = domain.UserStatus(statusStr)
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
