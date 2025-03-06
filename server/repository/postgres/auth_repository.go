package postgres

import (
	"context"
	"database/sql"
	"go-playground/pkg/logging"
	"go-playground/server/config"
	"go-playground/server/domain"
	"log"
	"time"

	"github.com/rs/zerolog"
)

type AuthRepository struct {
	db     *sql.DB
	config *config.AuthConfig
	logger zerolog.Logger
}

func NewAuthRepository(db *sql.DB, config *config.AuthConfig) *AuthRepository {
	return &AuthRepository{
		db:     db,
		config: config,
		logger: logging.GetLogger(),
	}
}

func (r *AuthRepository) CreateVerification(ctx context.Context, verification *domain.RegistrationVerification) error {
	r.logger.Info().
		Str("user_id", verification.UserID).
		Msg("Creating new verification")

	query := `
		INSERT INTO registration_verifications (user_id, otp, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		verification.UserID,
		verification.OTP,
		verification.ExpiresAt,
	).Scan(&verification.ID, &verification.CreatedAt)

	if err != nil {
		if isPgUniqueViolation(err) {
			r.logger.Warn().
				Str("user_id", verification.UserID).
				Msg("Verification already exists for user")
			return domain.NewResourceConflictError("verification", "verification already exists for this user")
		}
		r.logger.Error().
			Err(err).
			Str("user_id", verification.UserID).
			Msg("Failed to create verification")
		return domain.NewSystemError("AuthRepository.CreateVerification", err, "failed to create verification")
	}

	r.logger.Info().
		Str("user_id", verification.UserID).
		Str("verification_id", verification.ID).
		Msg("Successfully created verification")

	return nil
}

func (r *AuthRepository) GetVerification(ctx context.Context, userID, otp string) (*domain.RegistrationVerification, error) {
	r.logger.Info().
		Str("user_id", userID).
		Msg("Fetching verification")

	verification := &domain.RegistrationVerification{}
	var usedAt sql.NullTime

	query := `
		SELECT id, user_id, otp, expires_at, created_at, used_at
		FROM registration_verifications
		WHERE user_id = $1 AND otp = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		otp,
	).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&usedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("user_id", userID).
				Msg("Verification not found")
			return nil, domain.NewResourceNotFoundError("verification", userID, "verification not found")
		}
		r.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to get verification")
		return nil, domain.NewSystemError("AuthRepository.GetVerification", err, "failed to get verification")
	}

	if usedAt.Valid {
		verification.UsedAt = usedAt.Time
	}

	r.logger.Info().
		Str("user_id", userID).
		Str("verification_id", verification.ID).
		Msg("Successfully retrieved verification")

	return verification, nil
}

func (r *AuthRepository) MarkVerificationUsed(ctx context.Context, id string) error {
	query := `
		UPDATE registration_verifications
		SET used_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND used_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("verification_id", id).
			Msg("Failed to mark verification as used")
		return domain.NewSystemError("AuthRepository.MarkVerificationUsed", err, "failed to mark verification as used")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("verification_id", id).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("AuthRepository.MarkVerificationUsed", err, "failed to get affected rows")
	}

	if affected == 0 {
		r.logger.Error().
			Str("verification_id", id).
			Msg("Verification not found or already used")
		return domain.NewResourceNotFoundError("verification", id, "verification not found or already used")
	}

	return nil
}

func (r *AuthRepository) CreateToken(ctx context.Context, token *domain.AuthToken) error {
	r.logger.Info().
		Str("user_id", token.UserID).
		Msg("Creating new auth token")

	query := `
		INSERT INTO auth_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) -- Specify the column(s) that should trigger the conflict
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		if isPgUniqueViolation(err) {
			r.logger.Warn().
				Str("user_id", token.UserID).
				Msg("Token already exists for user")
			return domain.NewResourceConflictError("auth token", "token already exists")
		}
		r.logger.Error().
			Err(err).
			Str("user_id", token.UserID).
			Msg("Failed to create auth token")
	}

	return nil
}

func (r *AuthRepository) GetTokenByHash(ctx context.Context, hash string) (*domain.AuthToken, error) {
	token := &domain.AuthToken{}
	var lastUsedAt sql.NullTime

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at
		FROM auth_tokens
		WHERE token_hash = $1
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		hash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&lastUsedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("token_hash", hash).
				Msg("Token not found")
			return nil, domain.NewResourceNotFoundError("auth token", hash, "token not found")
		}
		r.logger.Error().
			Err(err).
			Str("token_hash", hash).
			Msg("Failed to get auth token")
	}

	if lastUsedAt.Valid {
		token.LastUsedAt = lastUsedAt.Time
	}

	return token, nil
}

func (r *AuthRepository) UpdateLoginAttempts(ctx context.Context, email string, increment bool) (*domain.LoginAttempt, error) {
	var attempt domain.LoginAttempt
	var lastAttempt sql.NullTime
	r.logger.Info().
		Str("email", email).
		Msg("Updating login attempts")

	// First, try to get existing record
	query := `
		SELECT id, email, attempt_count, last_attempt_at
		FROM login_attempts
		WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&attempt.ID,
		&attempt.Email,
		&attempt.AttemptCount,
		&lastAttempt,
	)

	if err != nil && err != sql.ErrNoRows {
		r.logger.Error().
			Err(err).
			Str("email", email).
			Msg("Failed to get login attempts")
		return nil, domain.NewSystemError("AuthRepository.UpdateLoginAttempts", err, "failed to get login attempts")
	}

	if err == sql.ErrNoRows {
		// Create new record if not exists
		if !increment {
			attempt.Email = email
			attempt.AttemptCount = 0
			attempt.LastAttempt = time.Now()
			return &attempt, nil
		}

		query = `
			INSERT INTO login_attempts (email, attempt_count, last_attempt_at)
			VALUES ($1, 1, CURRENT_TIMESTAMP)
			RETURNING id, attempt_count, last_attempt_at
		`
		err = r.db.QueryRowContext(ctx, query, email).Scan(
			&attempt.ID,
			&attempt.AttemptCount,
			&attempt.LastAttempt,
		)

		if err != nil {
			if isPgUniqueViolation(err) {
				r.logger.Warn().
					Str("email", email).
					Msg("Concurrent login attempt detected")
				return nil, domain.NewResourceConflictError("login attempt", "concurrent login attempt detected")
			}
			r.logger.Error().
				Err(err).
				Str("email", email).
				Msg("Failed to create login attempt")
		}

		attempt.Email = email
		return &attempt, nil
	}

	// Update existing record
	if increment {
		query = `
			UPDATE login_attempts
			SET attempt_count = attempt_count + 1,
				last_attempt_at = CURRENT_TIMESTAMP
			WHERE email = $1
			RETURNING attempt_count, last_attempt_at
		`
	} else {
		query = `
			UPDATE login_attempts
			SET attempt_count = 0,
				last_attempt_at = CURRENT_TIMESTAMP
			WHERE email = $1
			RETURNING attempt_count, last_attempt_at
		`
	}

	err = r.db.QueryRowContext(ctx, query, email).Scan(
		&attempt.AttemptCount,
		&attempt.LastAttempt,
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("email", email).
			Msg("Failed to update login attempts")
		return nil, domain.NewSystemError("AuthRepository.UpdateLoginAttempts", err, "failed to update login attempts")
	}

	if lastAttempt.Valid {
		attempt.LastAttempt = lastAttempt.Time
	}

	return &attempt, nil
}

func (r *AuthRepository) CleanupExpiredAttempts(ctx context.Context) error {
	query := `DELETE FROM login_attempts WHERE last_attempt_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now().Add(-r.config.LoginAttemptResetPeriod))
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to cleanup expired attempts")
		return domain.NewSystemError("AuthRepository.CleanupExpiredAttempts", err, "failed to cleanup expired attempts")
	}
	return nil
}

func (r *AuthRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to begin transaction")
		return nil, domain.NewSystemError("AuthRepository.BeginTx", err, "failed to begin transaction")
	}
	return tx, nil
}

func (r *AuthRepository) Commit(ctx context.Context, tx *sql.Tx) error {
	return tx.Commit()
}

func (r *AuthRepository) MarkVerificationUsedTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `
		UPDATE registration_verifications
		SET used_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND used_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("verification_id", id).
			Msg("Failed to mark verification as used")
		return domain.NewSystemError("AuthRepository.MarkVerificationUsedTx", err, "failed to mark verification as used")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("verification_id", id).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("AuthRepository.MarkVerificationUsedTx", err, "failed to get affected rows")
	}

	if affected == 0 {
		r.logger.Error().
			Str("verification_id", id).
			Msg("Verification not found or already used")
		return domain.NewResourceNotFoundError("verification", id, "verification not found or already used")
	}

	return nil
}

func (r *AuthRepository) GetUserVerificationStatus(ctx context.Context, userID string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM registration_verifications
		WHERE user_id = $1 AND used_at IS NOT NULL
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to get verification status")
		return false, domain.NewSystemError("AuthRepository.GetUserVerificationStatus", err, "failed to get verification status")
	}
	return count > 0, nil
}

func (r *AuthRepository) CleanupExpiredVerifications(ctx context.Context) error {
	query := `DELETE FROM registration_verifications WHERE expires_at < CURRENT_TIMESTAMP AND used_at IS NULL`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to cleanup expired verifications")
		return domain.NewSystemError("AuthRepository.CleanupExpiredVerifications", err, "failed to cleanup expired verifications")
	}
	return nil
}

func (r *AuthRepository) InvalidateToken(ctx context.Context, userID string) error {
	log.Println("Invalidating token for user: ", userID)
	query := `
		UPDATE auth_tokens
		SET 
		last_used_at = CURRENT_TIMESTAMP,
		expires_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to invalidate tokens")
		return domain.NewSystemError("AuthRepository.InvalidateToken", err, "failed to invalidate tokens")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("AuthRepository.InvalidateToken", err, "failed to get affected rows")
	}

	if affected == 0 {
		r.logger.Warn().
			Str("user_id", userID).
			Msg("No active tokens found for user")
		return domain.NewValidationError("AuthRepository.InvalidateToken", "no active tokens found for user")
	}

	return nil
}

func (r *AuthRepository) GetLatestVerification(ctx context.Context, userID string) (*domain.RegistrationVerification, error) {
	verification := &domain.RegistrationVerification{}
	var usedAt sql.NullTime

	query := `
		SELECT id, user_id, otp, expires_at, created_at, used_at
		FROM registration_verifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&usedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().
				Str("user_id", userID).
				Msg("No verification found for user")
			return nil, domain.NewResourceNotFoundError("verification", userID, "no verification found for user")
		}
		r.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to get latest verification")
		return nil, domain.NewSystemError("AuthRepository.GetLatestVerification", err, "failed to get latest verification")
	}

	if usedAt.Valid {
		verification.UsedAt = usedAt.Time
	}

	return verification, nil
}
