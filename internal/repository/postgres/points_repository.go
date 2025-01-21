package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
)

type PointsRepository struct {
	db *sql.DB
}

func NewPointsRepository(db *sql.DB) *PointsRepository {
	return &PointsRepository{db: db}
}

func (r *PointsRepository) Create(balance *domain.PointsBalance) error {
	query := `
		INSERT INTO points_balance (user_id, total_points)
		VALUES ($1, $2)
		RETURNING id, last_updated
	`
	return r.db.QueryRow(query, balance.UserID, balance.TotalPoints).
		Scan(&balance.ID, &balance.LastUpdated)
}

func (r *PointsRepository) GetByUserID(userID string) (*domain.PointsBalance, error) {
	balance := &domain.PointsBalance{}
	query := `
		SELECT id, user_id, total_points, last_updated
		FROM points_balance
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).
		Scan(&balance.ID, &balance.UserID, &balance.TotalPoints, &balance.LastUpdated)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return balance, err
}

func (r *PointsRepository) Update(balance *domain.PointsBalance) error {
	query := `
		UPDATE points_balance
		SET total_points = $1, last_updated = CURRENT_TIMESTAMP
		WHERE user_id = $2
		RETURNING last_updated
	`
	return r.db.QueryRow(query, balance.TotalPoints, balance.UserID).
		Scan(&balance.LastUpdated)
}
