package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
)

type RewardsRepository struct {
	db *sql.DB
}

func NewRewardsRepository(db *sql.DB) *RewardsRepository {
	return &RewardsRepository{db: db}
}

func (r *RewardsRepository) Create(reward *domain.Reward) error {
	query := `
		INSERT INTO rewards (
			name, description, points_required, is_active, quantity
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.IsActive,
		reward.Quantity,
	).Scan(
		&reward.ID,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)
}

func (r *RewardsRepository) GetByID(id string) (*domain.Reward, error) {
	reward := &domain.Reward{}
	var quantity sql.NullInt32

	query := `
		SELECT id, name, description, points_required, is_active, 
			   quantity, created_at, updated_at
		FROM rewards
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&reward.ID,
		&reward.Name,
		&reward.Description,
		&reward.PointsRequired,
		&reward.IsActive,
		&quantity,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if quantity.Valid {
		q := int(quantity.Int32)
		reward.Quantity = &q
	}

	return reward, nil
}

func (r *RewardsRepository) GetAll(activeOnly bool) ([]domain.Reward, error) {
	query := `
		SELECT id, name, description, points_required, is_active, 
			   quantity, created_at, updated_at
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
		var quantity sql.NullInt32

		err := rows.Scan(
			&reward.ID,
			&reward.Name,
			&reward.Description,
			&reward.PointsRequired,
			&reward.IsActive,
			&quantity,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if quantity.Valid {
			q := int(quantity.Int32)
			reward.Quantity = &q
		}

		rewards = append(rewards, reward)
	}

	return rewards, nil
}

func (r *RewardsRepository) Update(reward *domain.Reward) error {
	query := `
		UPDATE rewards
		SET name = $1, description = $2, points_required = $3,
			is_active = $4, quantity = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		reward.Name,
		reward.Description,
		reward.PointsRequired,
		reward.IsActive,
		reward.Quantity,
		reward.ID,
	).Scan(&reward.UpdatedAt)
}

func (r *RewardsRepository) Delete(id string) error {
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
