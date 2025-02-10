package postgres

import (
	"context"
	"database/sql"
	"go-playground/internal/config"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ProgramRuleRepository struct {
	db config.DbConnection
}

func NewProgramRuleRepository(db config.DbConnection) *ProgramRuleRepository {
	return &ProgramRuleRepository{db: db}
}

func (r *ProgramRuleRepository) Create(ctx context.Context, rule *domain.ProgramRule) error {
	query := `
		INSERT INTO program_rules (
			program_id, rule_name, condition_type, condition_value,
			multiplier, points_awarded, effective_from, effective_to
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	return r.db.RW.QueryRowContext(
		ctx,
		query,
		rule.ProgramID,
		rule.RuleName,
		rule.ConditionType,
		rule.ConditionValue,
		rule.Multiplier,
		rule.PointsAwarded,
		rule.EffectiveFrom,
		rule.EffectiveTo,
	).Scan(&rule.CreatedAt, &rule.UpdatedAt)
}

func (r *ProgramRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProgramRule, error) {
	query := `
		SELECT id, program_id, rule_name, condition_type, condition_value,
			   multiplier, points_awarded, effective_from, effective_to,
			   created_at, updated_at
		FROM program_rules
		WHERE id = $1
	`
	rule := &domain.ProgramRule{}
	err := r.db.RR.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.ProgramID,
		&rule.RuleName,
		&rule.ConditionType,
		&rule.ConditionValue,
		&rule.Multiplier,
		&rule.PointsAwarded,
		&rule.EffectiveFrom,
		&rule.EffectiveTo,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return rule, err
}

func (r *ProgramRuleRepository) GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*domain.ProgramRule, error) {
	query := `
		SELECT id, program_id, rule_name, condition_type, condition_value,
			   multiplier, points_awarded, effective_from, effective_to,
			   created_at, updated_at
		FROM program_rules
		WHERE program_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.RR.QueryContext(ctx, query, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.ProgramRule
	for rows.Next() {
		rule := &domain.ProgramRule{}
		err := rows.Scan(
			&rule.ID,
			&rule.ProgramID,
			&rule.RuleName,
			&rule.ConditionType,
			&rule.ConditionValue,
			&rule.Multiplier,
			&rule.PointsAwarded,
			&rule.EffectiveFrom,
			&rule.EffectiveTo,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *ProgramRuleRepository) Update(ctx context.Context, rule *domain.ProgramRule) error {
	query := `
		UPDATE program_rules
		SET rule_name = $1,
			condition_type = $2,
			condition_value = $3,
			multiplier = $4,
			points_awarded = $5,
			effective_from = $6,
			effective_to = $7
		WHERE id = $8
		RETURNING updated_at
	`
	return r.db.RW.QueryRowContext(
		ctx,
		query,
		rule.RuleName,
		rule.ConditionType,
		rule.ConditionValue,
		rule.Multiplier,
		rule.PointsAwarded,
		rule.EffectiveFrom,
		rule.EffectiveTo,
		rule.ID,
	).Scan(&rule.UpdatedAt)
}

func (r *ProgramRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM program_rules WHERE id = $1`
	result, err := r.db.RW.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ProgramRuleRepository) GetActiveRules(ctx context.Context, programID uuid.UUID, timestamp time.Time) ([]*domain.ProgramRule, error) {
	query := `
		SELECT id, program_id, rule_name, condition_type, condition_value,
			   multiplier, points_awarded, effective_from, effective_to,
			   created_at, updated_at
		FROM program_rules
		WHERE program_id = $1
		AND effective_from <= $2
		AND (effective_to IS NULL OR effective_to >= $2)
		ORDER BY created_at DESC
	`
	rows, err := r.db.RR.QueryContext(ctx, query, programID, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.ProgramRule
	for rows.Next() {
		rule := &domain.ProgramRule{}
		err := rows.Scan(
			&rule.ID,
			&rule.ProgramID,
			&rule.RuleName,
			&rule.ConditionType,
			&rule.ConditionValue,
			&rule.Multiplier,
			&rule.PointsAwarded,
			&rule.EffectiveFrom,
			&rule.EffectiveTo,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
