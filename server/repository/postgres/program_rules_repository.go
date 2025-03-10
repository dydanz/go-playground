package postgres

import (
	"context"
	"database/sql"
	"go-playground/pkg/logging"
	"go-playground/server/config"
	"go-playground/server/domain"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ProgramRuleRepository struct {
	db     config.DbConnection
	logger zerolog.Logger
}

func NewProgramRuleRepository(db config.DbConnection) *ProgramRuleRepository {
	return &ProgramRuleRepository{db: db,
		logger: logging.GetLogger(),
	}
}

func (r *ProgramRuleRepository) Create(ctx context.Context, rule *domain.ProgramRule) error {
	query := `
		INSERT INTO program_rules (
			program_id, rule_name, condition_type, condition_value,
			multiplier, points_awarded, effective_from, effective_to
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.db.RW.QueryRowContext(
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
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)

	if err != nil {
		if isPgUniqueViolation(err) {
			r.logger.Error().
				Err(err).
				Msg("Failed to create program rule")
			return domain.NewResourceConflictError("program rule", "rule with this name already exists for the program")
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to create program rule")
		return domain.NewSystemError("ProgramRuleRepository.Create", err, "failed to create program rule")
	}

	return nil
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
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Error().
				Err(err).
				Msg("Failed to get program rule")
			return nil, domain.NewResourceNotFoundError("program rule", id.String(), "program rule not found")
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to get program rule")
		return nil, domain.NewSystemError("ProgramRuleRepository.GetByID", err, "failed to get program rule")
	}
	return rule, nil
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
		r.logger.Error().
			Err(err).
			Msg("Failed to query program rules")
		return nil, domain.NewSystemError("ProgramRuleRepository.GetByProgramID", err, "failed to query program rules")
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
			r.logger.Error().
				Err(err).
				Msg("Failed to scan program rule")
			return nil, domain.NewSystemError("ProgramRuleRepository.GetByProgramID", err, "failed to scan program rule")
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to iterate program rules")
		return nil, domain.NewSystemError("ProgramRuleRepository.GetByProgramID", err, "error iterating program rules")
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
			effective_to = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`
	result, err := r.db.RW.ExecContext(
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
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			r.logger.Error().
				Err(err).
				Msg("Failed to update program rule")
			return domain.NewResourceConflictError("program rule", "rule with this name already exists for the program")
		}
		r.logger.Error().
			Err(err).
			Msg("Failed to update program rule")
		return domain.NewSystemError("ProgramRuleRepository.Update", err, "failed to update program rule")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("ProgramRuleRepository.Update", err, "failed to get affected rows")
	}

	if affected == 0 {
		r.logger.Error().
			Msg("Failed to update program rule")
		return domain.NewResourceNotFoundError("program rule", rule.ID.String(), "program rule not found")
	}

	return nil
}

func (r *ProgramRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM program_rules WHERE id = $1`
	result, err := r.db.RW.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to delete program rule")
		return domain.NewSystemError("ProgramRuleRepository.Delete", err, "failed to delete program rule")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to get affected rows")
		return domain.NewSystemError("ProgramRuleRepository.Delete", err, "failed to get affected rows")
	}

	if affected == 0 {
		r.logger.Error().
			Msg("Failed to delete program rule")
		return domain.NewResourceNotFoundError("program rule", id.String(), "program rule not found")
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
		r.logger.Error().
			Err(err).
			Msg("Failed to query active program rules")
		return nil, domain.NewSystemError("ProgramRuleRepository.GetActiveRules", err, "failed to query active program rules")
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
			r.logger.Error().
				Err(err).
				Msg("Failed to scan program rule")
			return nil, domain.NewSystemError("ProgramRuleRepository.GetActiveRules", err, "failed to scan program rule")
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to iterate active program rules")
		return nil, domain.NewSystemError("ProgramRuleRepository.GetActiveRules", err, "error iterating active program rules")
	}

	return rules, nil
}
