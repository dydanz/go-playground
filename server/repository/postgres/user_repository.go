package postgres

import (
	"context"
	"database/sql"
	"go-playground/server/domain"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	query := `
		INSERT INTO users (email, password, name, phone, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, name, phone, status, created_at, updated_at
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Email,
		req.Password,
		req.Name,
		req.Phone,
		domain.UserStatusPending, // Default status for new users
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
		if isPgUniqueViolation(err) {
			return nil, domain.NewResourceConflictError("user", "user with this email already exists")
		}
		return nil, domain.NewSystemError("UserRepository.Create", err, "failed to create user")
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	var statusStr string
	query := `
		SELECT id, name, email, phone, password, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&statusStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewResourceNotFoundError("user", id, "user not found")
		}
		return nil, domain.NewSystemError("UserRepository.GetByID", err, "failed to get user")
	}

	user.Status = domain.UserStatus(statusStr)
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, 
			phone = $2, 
			status = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Phone,
		user.Status,
		user.ID,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.NewResourceConflictError("user", "user with this email already exists")
		}
		return domain.NewSystemError("UserRepository.Update", err, "failed to update user")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("UserRepository.Update", err, "failed to get affected rows")
	}

	if affected == 0 {
		return domain.NewResourceNotFoundError("user", user.ID, "user not found")
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	var statusStr string
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
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
		if err == sql.ErrNoRows {
			log.Printf("No user found with the given ID.")
			return nil, nil
		}
		return nil, domain.NewSystemError("UserRepository.GetByEmail", err, "failed to get user")
	}

	user.Status = domain.UserStatus(statusStr)
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
		FROM users
	`
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user(s) found.")
			return nil, nil
		}
		return nil, domain.NewSystemError("UserRepository.GetAll", err, "failed to query users")
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
			return nil, domain.NewSystemError("UserRepository.GetAll", err, "failed to scan user")
		}
		user.Status = domain.UserStatus(statusStr)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewSystemError("UserRepository.GetAll", err, "error iterating users")
	}

	return users, nil
}

func (r *UserRepository) UpdateTx(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, 
			phone = $2, 
			status = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		user.Name,
		user.Phone,
		user.Status,
		user.ID,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.NewResourceConflictError("user", "user with this email already exists")
		}
		return domain.NewSystemError("UserRepository.UpdateTx", err, "failed to update user")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("UserRepository.UpdateTx", err, "failed to get affected rows")
	}

	if affected == 0 {
		return domain.NewResourceNotFoundError("user", user.ID, "user not found")
	}

	return nil
}

func (r *UserRepository) GetRandomActiveUser(ctx context.Context) (*domain.User, error) {
	user := &domain.User{}
	var statusStr string
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
		FROM users
		WHERE status = $1
		ORDER BY RANDOM()
		LIMIT 1
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		domain.UserStatusActive,
	).Scan(
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
		if err == sql.ErrNoRows {
			log.Printf("no active users found")
			return nil, nil
		}
		return nil, domain.NewSystemError("UserRepository.GetRandomActiveUser", err, "failed to get random active user")
	}

	user.Status = domain.UserStatus(statusStr)
	return user, nil
}
