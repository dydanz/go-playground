package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
)

type RedemptionRepository struct {
	db *sql.DB
}

func NewRedemptionRepository(db *sql.DB) *RedemptionRepository {
	return &RedemptionRepository{db: db}
}

func (r *RedemptionRepository) Create(redemption *domain.Redemption) error {
	query := `
		INSERT INTO redemptions (
			user_id, reward_id, status, redeemed_at
		) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, redeemed_at, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		redemption.UserID,
		redemption.RewardID,
		redemption.Status,
	).Scan(
		&redemption.ID,
		&redemption.RedeemedAt,
		&redemption.CreatedAt,
		&redemption.UpdatedAt,
	)
}

func (r *RedemptionRepository) GetByID(id string) (*domain.Redemption, error) {
	redemption := &domain.Redemption{}
	query := `
		SELECT id, user_id, reward_id, status, redeemed_at, 
			   created_at, updated_at
		FROM redemptions
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&redemption.ID,
		&redemption.UserID,
		&redemption.RewardID,
		&redemption.Status,
		&redemption.RedeemedAt,
		&redemption.CreatedAt,
		&redemption.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return redemption, err
}

func (r *RedemptionRepository) GetByUserID(userID string) ([]domain.Redemption, error) {
	query := `
		SELECT id, user_id, reward_id, status, redeemed_at, 
			   created_at, updated_at
		FROM redemptions
		WHERE user_id = $1
		ORDER BY redeemed_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []domain.Redemption
	for rows.Next() {
		var redemption domain.Redemption
		err := rows.Scan(
			&redemption.ID,
			&redemption.UserID,
			&redemption.RewardID,
			&redemption.Status,
			&redemption.RedeemedAt,
			&redemption.CreatedAt,
			&redemption.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, redemption)
	}
	return redemptions, nil
}

func (r *RedemptionRepository) Update(redemption *domain.Redemption) error {
	query := `
		UPDATE redemptions
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.QueryRow(query, redemption.Status, redemption.ID).
		Scan(&redemption.UpdatedAt)
}
