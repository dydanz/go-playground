package postgres

import (
	"database/sql"
	"fmt"
	"go-cursor/internal/config"
	"go-cursor/internal/domain"
	"log"
)

type AuthRepository struct {
	db     *sql.DB
	config *config.AuthConfig
}

func NewAuthRepository(db *sql.DB, config *config.AuthConfig) *AuthRepository {
	repo := &AuthRepository{
		db:     db,
		config: config,
	}

	// Ensure table exists
	if err := repo.ensureLoginAttemptsTable(); err != nil {
		log.Printf("Failed to create login_attempts table: %v", err)
	}

	return repo
}

func (r *AuthRepository) ensureLoginAttemptsTable() error {
	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS login_attempts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) NOT NULL,
			attempt_count INT NOT NULL DEFAULT 1,
			last_attempt_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			locked_until TIMESTAMP WITH TIME ZONE
	)
	`
	if _, err := r.db.Exec(createTable); err != nil {
		return err
	}

	// Add unique constraint if not exists
	addConstraint := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (
			SELECT 1 
			FROM pg_constraint 
			WHERE conname = 'login_attempts_email_key'
		) THEN
			ALTER TABLE login_attempts ADD CONSTRAINT login_attempts_email_key UNIQUE(email);
		END IF;
	END $$;
	`
	_, err := r.db.Exec(addConstraint)
	return err
}

func (r *AuthRepository) CreateVerification(verification *domain.RegistrationVerification) error {
	query := `
		INSERT INTO registration_verifications (user_id, otp, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, verification.UserID, verification.OTP, verification.ExpiresAt)
	return err
}

func (r *AuthRepository) GetVerification(userID, otp string) (*domain.RegistrationVerification, error) {
	verification := &domain.RegistrationVerification{}
	var usedAt sql.NullTime

	// Get the verification record
	query := `
		SELECT id, user_id, otp, expires_at, created_at, used_at
		FROM registration_verifications
			WHERE user_id = $1 
			AND otp = $2 
			AND expires_at > NOW()
			ORDER BY created_at DESC
			LIMIT 1
	`

	err := r.db.QueryRow(query, userID, otp).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&usedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired OTP")
	}
	if err != nil {
		return nil, err
	}

	// Check if OTP is already used
	if usedAt.Valid {
		return nil, fmt.Errorf("OTP already used")
	}

	// Set the used_at field if valid
	if usedAt.Valid {
		verification.UsedAt = usedAt.Time
	}

	return verification, nil
}

func (r *AuthRepository) MarkVerificationUsed(id string) error {
	query := `
			UPDATE registration_verifications
			SET used_at = NOW()
			WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AuthRepository) CreateToken(token *domain.AuthToken) error {
	query := `
		INSERT INTO auth_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = CURRENT_TIMESTAMP,
			last_used_at = NULL
		RETURNING id
	`
	return r.db.QueryRow(query, token.UserID, token.TokenHash, token.ExpiresAt).Scan(&token.ID)
}

func (r *AuthRepository) GetTokenByHash(hash string) (*domain.AuthToken, error) {
	token := &domain.AuthToken{}
	var lastUsedAt sql.NullTime

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at
		FROM auth_tokens
		WHERE token_hash = $1 
		AND expires_at > NOW()
		LIMIT 1
	`
	err := r.db.QueryRow(query, hash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&lastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	if lastUsedAt.Valid {
		token.LastUsedAt = lastUsedAt.Time
	}

	// Update last_used_at
	updateQuery := `
		UPDATE auth_tokens 
		SET last_used_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err = r.db.Exec(updateQuery, token.ID)
	if err != nil {
		log.Printf("Failed to update last_used_at: %v", err)
	}

	return token, nil
}

func (r *AuthRepository) UpdateLoginAttempts(email string, increment bool) (*domain.LoginAttempt, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	attempt := &domain.LoginAttempt{}
	if increment {
		var lockedUntil sql.NullTime

		query := `
			INSERT INTO login_attempts (email, attempt_count, last_attempt_at)
			VALUES ($1, 1, NOW())
				ON CONFLICT (email) DO UPDATE
				SET 
					attempt_count = CASE
						WHEN login_attempts.last_attempt_at < NOW() - $2::interval THEN 1
						ELSE login_attempts.attempt_count + 1
					END,
					last_attempt_at = NOW(),
					locked_until = CASE
						WHEN login_attempts.last_attempt_at < NOW() - $2::interval THEN NULL
						WHEN login_attempts.attempt_count + 1 >= $3 THEN NOW() + $4::interval
						ELSE NULL
					END
				RETURNING id, email, attempt_count, last_attempt_at, locked_until
		`
		err = tx.QueryRow(
			query,
			email,
			fmt.Sprintf("%d seconds", int(r.config.LoginAttemptResetPeriod.Seconds())),
			r.config.MaxLoginAttempts,
			fmt.Sprintf("%d seconds", int(r.config.LockDuration.Seconds())),
		).Scan(
			&attempt.ID,
			&attempt.Email,
			&attempt.AttemptCount,
			&attempt.LastAttempt,
			&lockedUntil,
		)

		if lockedUntil.Valid {
			attempt.LockedUntil = lockedUntil.Time
		}
	} else {
		query := `DELETE FROM login_attempts WHERE email = $1`
		_, err = tx.Exec(query, email)
	}

	if err != nil {
		return nil, err
	}

	return attempt, tx.Commit()
}

func (r *AuthRepository) CleanupExpiredAttempts() error {
	query := `
		DELETE FROM login_attempts 
		WHERE last_attempt_at < NOW() - $1::interval
		OR (locked_until IS NOT NULL AND locked_until < NOW())
	`
	_, err := r.db.Exec(query, fmt.Sprintf("%d seconds", int(r.config.LoginAttemptResetPeriod.Seconds())))
	return err
}

func (r *AuthRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *AuthRepository) MarkVerificationUsedTx(tx *sql.Tx, id string) error {
	query := `
		UPDATE registration_verifications
		SET used_at = NOW()
		WHERE id = $1
	`
	_, err := tx.Exec(query, id)
	return err
}

func (r *AuthRepository) GetUserVerificationStatus(userID string) (bool, error) {
	var status string
	err := r.db.QueryRow(`
		SELECT status 
		FROM users 
		WHERE id = $1
	`, userID).Scan(&status)

	if err != nil {
		return false, err
	}

	return status == string(domain.UserStatusActive), nil
}

func (r *AuthRepository) CleanupExpiredVerifications() error {
	query := `
		DELETE FROM registration_verifications 
		WHERE expires_at < NOW() 
		OR used_at IS NOT NULL
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *AuthRepository) InvalidateToken(userID string) error {
	query := `
		UPDATE auth_tokens 
		SET expires_at = '1970-01-01 00:00:01'
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}
