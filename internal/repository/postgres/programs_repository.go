package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Program struct {
	ProgramID         uuid.UUID
	MerchantID        uuid.UUID
	ProgramName       string
	PointCurrencyName string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ProgramsRepository struct {
	db *sql.DB
}

func NewProgramsRepository(db *sql.DB) *ProgramsRepository {
	return &ProgramsRepository{db: db}
}

func (r *ProgramsRepository) Create(ctx context.Context, program *Program) error {
	query := `
		INSERT INTO programs (program_id, merchant_id, program_name, point_currency_name)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		program.ProgramID,
		program.MerchantID,
		program.ProgramName,
		program.PointCurrencyName,
	).Scan(&program.CreatedAt, &program.UpdatedAt)
}

func (r *ProgramsRepository) GetByID(ctx context.Context, programID uuid.UUID) (*Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE program_id = $1`

	program := &Program{}
	err := r.db.QueryRowContext(ctx, query, programID).Scan(
		&program.ProgramID,
		&program.MerchantID,
		&program.ProgramName,
		&program.PointCurrencyName,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return program, nil
}

func (r *ProgramsRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE merchant_id = $1`

	rows, err := r.db.QueryContext(ctx, query, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []*Program
	for rows.Next() {
		program := &Program{}
		err := rows.Scan(
			&program.ProgramID,
			&program.MerchantID,
			&program.ProgramName,
			&program.PointCurrencyName,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}
	return programs, nil
}

func (r *ProgramsRepository) Update(ctx context.Context, program *Program) error {
	query := `
		UPDATE programs
		SET program_name = $1, point_currency_name = $2
		WHERE program_id = $3
		RETURNING updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		program.ProgramName,
		program.PointCurrencyName,
		program.ProgramID,
	).Scan(&program.UpdatedAt)
}

func (r *ProgramsRepository) Delete(ctx context.Context, programID uuid.UUID) error {
	query := `DELETE FROM programs WHERE program_id = $1`
	result, err := r.db.ExecContext(ctx, query, programID)
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
