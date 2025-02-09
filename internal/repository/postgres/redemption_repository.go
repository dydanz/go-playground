package postgres

import (
	"context"
	"database/sql"
	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type RedemptionRepository struct {
	db *sql.DB
}

func NewRedemptionRepository(db *sql.DB) *RedemptionRepository {
	return &RedemptionRepository{db: db}
}

func (r *RedemptionRepository) Create(ctx context.Context, redemption *domain.Redemption) error {
	if redemption.ID == uuid.Nil {
		redemption.ID = uuid.New()
	}

	query := `
		INSERT INTO redemptions (
			id, merchant_customers_id, reward_id, points_used,
			redemption_date, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING redemption_date, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		redemption.ID,
		redemption.MerchantCustomersID,
		redemption.RewardID,
		redemption.PointsUsed,
		redemption.Status,
	).Scan(
		&redemption.RedemptionDate,
		&redemption.CreatedAt,
		&redemption.UpdatedAt,
	)
}

func (r *RedemptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Redemption, error) {
	query := `
		SELECT id, merchant_customers_id, reward_id, points_used,
			   redemption_date, status, created_at, updated_at
		FROM redemptions
		WHERE merchant_customers_id = $1	
		ORDER BY redemption_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	redemptions := []*domain.Redemption{}
	for rows.Next() {
		redemption := &domain.Redemption{}
		err := rows.Scan(
			&redemption.ID,
			&redemption.MerchantCustomersID,
			&redemption.RewardID,
			&redemption.PointsUsed,
			&redemption.RedemptionDate,
			&redemption.Status,
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

func (r *RedemptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Redemption, error) {
	redemption := &domain.Redemption{}
	query := `
		SELECT id, merchant_customers_id, reward_id, points_used,
			   redemption_date, status, created_at, updated_at
		FROM redemptions
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&redemption.ID,
		&redemption.MerchantCustomersID,
		&redemption.RewardID,
		&redemption.PointsUsed,
		&redemption.RedemptionDate,
		&redemption.Status,
		&redemption.CreatedAt,
		&redemption.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return redemption, err
}

func (r *RedemptionRepository) GetByMerchantCustomerID(ctx context.Context, merchantCustomersID uuid.UUID) ([]*domain.Redemption, error) {
	query := `
		SELECT id, merchant_customers_id, reward_id, points_used,
			   redemption_date, status, created_at, updated_at
		FROM redemptions
		WHERE merchant_customers_id = $1
		ORDER BY redemption_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, merchantCustomersID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*domain.Redemption
	for rows.Next() {
		redemption := &domain.Redemption{}
		err := rows.Scan(
			&redemption.ID,
			&redemption.MerchantCustomersID,
			&redemption.RewardID,
			&redemption.PointsUsed,
			&redemption.RedemptionDate,
			&redemption.Status,
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

func (r *RedemptionRepository) Update(ctx context.Context, redemption *domain.Redemption) error {
	query := `
		UPDATE redemptions
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query, redemption.Status, redemption.ID).
		Scan(&redemption.UpdatedAt)
}

func (r *RedemptionRepository) GetByRewardID(ctx context.Context, rewardID uuid.UUID) ([]*domain.Redemption, error) {
	query := `
		SELECT id, merchant_customers_id, reward_id, points_used,
			   redemption_date, status, created_at, updated_at
		FROM redemptions
		WHERE reward_id = $1
		ORDER BY redemption_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, rewardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*domain.Redemption
	for rows.Next() {
		redemption := &domain.Redemption{}
		err := rows.Scan(
			&redemption.ID,
			&redemption.MerchantCustomersID,
			&redemption.RewardID,
			&redemption.PointsUsed,
			&redemption.RedemptionDate,
			&redemption.Status,
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
