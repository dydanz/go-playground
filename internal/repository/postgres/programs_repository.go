package postgres

import (
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

func (r *ProgramsRepository) Create(program *domain.Program) error {
	programID := program.ID.String()
	merchantID := program.MerchantID.String()

	query := `
		INSERT INTO programs (program_id, merchant_id, program_name, point_currency_name)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(
		query,
		programID,
		merchantID,
		program.ProgramName,
		program.PointCurrencyName,
	).Scan(&program.CreatedAt, &program.UpdatedAt)
}

func (r *ProgramsRepository) GetByID(id string) (*domain.Program, error) {
	programID := id

	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE program_id = $1`

	var program domain.Program
	var pID, mID uuid.UUID
	err := r.db.QueryRow(query, programID).Scan(
		&pID,
		&mID,
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
	program.ID = pID
	program.MerchantID = mID
	return &program, nil
}

func (r *ProgramsRepository) GetAll() ([]*domain.Program, error) {
	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		program.ID = pID
		program.MerchantID = mID
		programs = append(programs, &program)
	}
	return programs, nil
}

func (r *ProgramsRepository) Update(program *domain.Program) error {
	programID := program.ID.String()

	query := `
		UPDATE programs
		SET program_name = $1, point_currency_name = $2
		WHERE program_id = $3
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		program.ProgramName,
		program.PointCurrencyName,
		programID,
	).Scan(&program.UpdatedAt)
}

func (r *ProgramsRepository) Delete(id string) error {
	programID := id

	query := `DELETE FROM programs WHERE program_id = $1`
	result, err := r.db.Exec(query, programID)
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

func (r *ProgramsRepository) GetByMerchantID(merchantID string) ([]*domain.Program, error) {
	mID := merchantID

	query := `
		SELECT program_id, merchant_id, program_name, point_currency_name, created_at, updated_at
		FROM programs
		WHERE merchant_id = $1`

	rows, err := r.db.Query(query, mID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		program.ID = pID
		program.MerchantID = mID
		programs = append(programs, &program)
	}
	return programs, nil
}
