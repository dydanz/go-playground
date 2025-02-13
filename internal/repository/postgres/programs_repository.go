package postgres

import (
	"context"
	"database/sql"
	"go-playground/internal/domain"

	"github.com/google/uuid"
)

type ProgramsRepository struct {
	db *sql.DB
}

func NewProgramsRepository(db *sql.DB) *ProgramsRepository {
	return &ProgramsRepository{db: db}
}

func (r *ProgramsRepository) Create(ctx context.Context, program *domain.Program) (*domain.Program, error) {
	merchantID := program.MerchantID.String()

	query := `
		INSERT INTO programs (merchant_id, user_id, program_name, point_currency_name)
		VALUES ($1, $2, $3, $4)
		RETURNING program_id, merchant_id, program_name, point_currency_name, created_at, updated_at`

	result := domain.Program{}
	var mID uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		query,
		merchantID,
		program.UserID,
		program.ProgramName,
		program.PointCurrencyName,
	).Scan(
		&result.ID,
		&mID,
		&result.ProgramName,
		&result.PointCurrencyName,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return nil, domain.NewResourceConflictError("program", "program with this name already exists for the merchant")
		}
		return nil, domain.NewSystemError("ProgramsRepository.Create", err, "failed to create program")
	}

	result.MerchantID = mID
	return &result, nil
}

func (r *ProgramsRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE program_id = $1`

	var program domain.Program
	var pID, mID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pID,
		&mID,
		&program.ProgramName,
		&program.PointCurrencyName,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewResourceNotFoundError("program", id.String(), "program not found")
		}
		return nil, domain.NewSystemError("ProgramsRepository.GetByID", err, "failed to get program")
	}
	program.ID = pID
	program.MerchantID = mID
	return &program, nil
}

func (r *ProgramsRepository) GetAll(ctx context.Context) ([]*domain.Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, domain.NewSystemError("ProgramsRepository.GetAll", err, "failed to query programs")
	}
	defer rows.Close()

	var programs []*domain.Program
	for rows.Next() {
		var program domain.Program
		var pID, mID uuid.UUID
		err := rows.Scan(
			&pID,
			&mID,
			&program.ProgramName,
			&program.PointCurrencyName,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewSystemError("ProgramsRepository.GetAll", err, "failed to scan program")
		}
		program.ID = pID
		program.MerchantID = mID
		programs = append(programs, &program)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewSystemError("ProgramsRepository.GetAll", err, "error iterating programs")
	}

	return programs, nil
}

func (r *ProgramsRepository) Update(ctx context.Context, program *domain.Program) error {
	query := `
		UPDATE programs
		SET program_name = $1, point_currency_name = $2
		WHERE program_id = $3
		RETURNING updated_at`

	result, err := r.db.ExecContext(
		ctx,
		query,
		program.ProgramName,
		program.PointCurrencyName,
		program.ID,
	)

	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.NewResourceConflictError("program", "program with this name already exists for the merchant")
		}
		return domain.NewSystemError("ProgramsRepository.Update", err, "failed to update program")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.NewSystemError("ProgramsRepository.Update", err, "failed to get affected rows")
	}

	if affected == 0 {
		return domain.NewResourceNotFoundError("program", program.ID.String(), "program not found")
	}

	return nil
}

func (r *ProgramsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *ProgramsRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE merchant_id = $1`

	rows, err := r.db.QueryContext(ctx, query, merchantID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramsRepository.GetByMerchantID", err, "failed to query programs")
	}
	defer rows.Close()

	var programs []*domain.Program
	for rows.Next() {
		var program domain.Program
		var pID, mID uuid.UUID
		err := rows.Scan(
			&pID,
			&mID,
			&program.ProgramName,
			&program.PointCurrencyName,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, domain.NewSystemError("ProgramsRepository.GetByMerchantID", err, "failed to scan program")
		}
		program.ID = pID
		program.MerchantID = mID
		programs = append(programs, &program)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewSystemError("ProgramsRepository.GetByMerchantID", err, "error iterating programs")
	}

	return programs, nil
}
