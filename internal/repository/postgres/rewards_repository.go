package postgres

import (
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

func (r *RewardsRepository) Create(reward *domain.Reward) error {
	if reward.ID == uuid.Nil {
		reward.ID = uuid.New()
	}

	query := `
		INSERT INTO rewards (
			id, program_id, name, description, points_required,
			available_quantity, quantity, is_active,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		reward.ID,
		reward.ProgramID,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.AvailableQuantity,
		reward.Quantity,
		reward.IsActive,
	).Scan(
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)
}

func (r *RewardsRepository) GetByID(id uuid.UUID) (*domain.Reward, error) {
	reward := &domain.Reward{}
	var availableQuantity sql.NullInt32

	query := `
		SELECT id, program_id, name, description, points_required,
			   available_quantity, quantity, is_active,
			   created_at, updated_at
		FROM rewards
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
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

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if availableQuantity.Valid {
		q := int(availableQuantity.Int32)
		reward.AvailableQuantity = &q
	}

	return reward, nil
}

func (r *RewardsRepository) GetAll(activeOnly bool) ([]domain.Reward, error) {
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

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if availableQuantity.Valid {
			q := int(availableQuantity.Int32)
			reward.AvailableQuantity = &q
		}

		rewards = append(rewards, reward)
	}

	return rewards, nil
}

func (r *RewardsRepository) Update(reward *domain.Reward) error {
	query := `
		UPDATE rewards
		SET name = $1, description = $2, points_required = $3,
			available_quantity = $4, quantity = $5, is_active = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.AvailableQuantity,
		reward.Quantity,
		reward.IsActive,
		reward.ID,
	).Scan(&reward.UpdatedAt)
}

func (r *RewardsRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM rewards WHERE id = $1`
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

func (r *RewardsRepository) GetByProgramID(programID uuid.UUID) ([]*domain.Reward, error) {
	query := `
		SELECT id, program_id, name, description, points_required,
			   available_quantity, quantity, is_active,
			   created_at, updated_at
		FROM rewards
		WHERE program_id = $1
		ORDER BY points_required ASC
	`

	rows, err := r.db.Query(query, programID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if availableQuantity.Valid {
			q := int(availableQuantity.Int32)
			reward.AvailableQuantity = &q
		}

		rewards = append(rewards, reward)
	}

	return rewards, nil
}
