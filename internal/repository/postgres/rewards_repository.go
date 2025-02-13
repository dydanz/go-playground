package postgres

import (
	"context"
	"database/sql"
	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type RewardsRepository struct {
	db *sql.DB
}

func NewRewardsRepository(db *sql.DB) *RewardsRepository {
	return &RewardsRepository{db: db}
}

func (r *RewardsRepository) Create(ctx context.Context, reward *domain.Reward) error {
	if reward.ID == uuid.Nil {
		reward.ID = uuid.New()
	}

	query := `
		INSERT INTO rewards (
			program_id, name, description, points_required,
			available_quantity, quantity, is_active,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, points_required, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		reward.ProgramID,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.AvailableQuantity,
		reward.Quantity,
		reward.IsActive,
	).Scan(
		&reward.ID,
		&reward.PointsRequired,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.NewResourceConflictError("reward", "reward with this name already exists")
		}
		return domain.NewSystemError("RewardsRepository.Create", err, "failed to create reward")
	}

	return nil
}

func (r *RewardsRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reward, error) {
	reward := &domain.Reward{}
	var availableQuantity sql.NullInt32

	query := `
		SELECT id, program_id, name, description, points_required,
			   available_quantity, quantity, is_active,
			   created_at, updated_at
		FROM rewards
		WHERE id = $1
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&reward.ID,
		&reward.ProgramID,
		&reward.Name,
		&reward.Description,
		&reward.PointsRequired,
		&availableQuantity,
		&reward.Quantity,
		&reward.IsActive,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewResourceNotFoundError("reward", id.String(), "reward not found")
		}
		return nil, domain.NewSystemError("RewardsRepository.GetByID", err, "failed to get reward")
	}

	if availableQuantity.Valid {
		q := int(availableQuantity.Int32)
		reward.AvailableQuantity = &q
	}

	return reward, nil
}

func (r *RewardsRepository) GetAll(ctx context.Context, activeOnly bool) ([]domain.Reward, error) {
	query := `
		SELECT id, program_id, name, description, points_required,
			   available_quantity, quantity, is_active,
			   created_at, updated_at
		FROM rewards
	`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY points_required ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, domain.NewSystemError("RewardsRepository.GetAll", err, "failed to query rewards")
	}
	defer rows.Close()

	var rewards []domain.Reward
	for rows.Next() {
		var reward domain.Reward
		var availableQuantity sql.NullInt32

		err := rows.Scan(
			&reward.ID,
			&reward.ProgramID,
			&reward.Name,
			&reward.Description,
			&reward.PointsRequired,
			&availableQuantity,
			&reward.Quantity,
			&reward.IsActive,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewSystemError("RewardsRepository.GetAll", err, "failed to scan reward")
		}

		if availableQuantity.Valid {
			q := int(availableQuantity.Int32)
			reward.AvailableQuantity = &q
		}

		rewards = append(rewards, reward)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewSystemError("RewardsRepository.GetAll", err, "error iterating rewards")
	}

	return rewards, nil
}

func (r *RewardsRepository) Update(ctx context.Context, reward *domain.Reward) error {
	query := `
		UPDATE rewards
		SET name = $1, description = $2, points_required = $3,
			available_quantity = $4, quantity = $5, is_active = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.AvailableQuantity,
		reward.Quantity,
		reward.IsActive,
		reward.ID,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.NewResourceConflictError("reward", "reward with this name already exists")
		}
		return domain.NewSystemError("RewardsRepository.Update", err, "failed to update reward")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("RewardsRepository.Update", err, "failed to get affected rows")
	}

	if affected == 0 {
		return domain.NewResourceNotFoundError("reward", reward.ID.String(), "reward not found")
	}

	return nil
}

func (r *RewardsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rewards WHERE id = $1`
	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return domain.NewSystemError("RewardsRepository.Delete", err, "failed to delete reward")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("RewardsRepository.Delete", err, "failed to get affected rows")
	}

	if affected == 0 {
		return domain.NewResourceNotFoundError("reward", id.String(), "reward not found")
	}

	return nil
}

func (r *RewardsRepository) GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*domain.Reward, error) {
	query := `
		SELECT id, program_id, name, description, points_required,
			   available_quantity, quantity, is_active,
			   created_at, updated_at
		FROM rewards
		WHERE program_id = $1
		ORDER BY points_required ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		programID,
	)
	if err != nil {
		return nil, domain.NewSystemError("RewardsRepository.GetByProgramID", err, "failed to query rewards")
	}
	defer rows.Close()

	var rewards []*domain.Reward
	for rows.Next() {
		reward := &domain.Reward{}
		var availableQuantity sql.NullInt32

		err := rows.Scan(
			&reward.ID,
			&reward.ProgramID,
			&reward.Name,
			&reward.Description,
			&reward.PointsRequired,
			&availableQuantity,
			&reward.Quantity,
			&reward.IsActive,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewSystemError("RewardsRepository.GetByProgramID", err, "failed to scan reward")
		}

		if availableQuantity.Valid {
			q := int(availableQuantity.Int32)
			reward.AvailableQuantity = &q
		}

		rewards = append(rewards, reward)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewSystemError("RewardsRepository.GetByProgramID", err, "error iterating rewards")
	}

	return rewards, nil
}
