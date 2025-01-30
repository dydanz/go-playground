package service

import (
	"context"
	"errors"

	"go-playground/internal/repository/postgres"

	"github.com/google/uuid"
)

var (
	ErrProgramNotFound = errors.New("program not found")
	ErrInvalidProgram  = errors.New("invalid program data")
)

type ProgramService struct {
	programsRepo *postgres.ProgramsRepository
}

func NewProgramService(programsRepo *postgres.ProgramsRepository) *ProgramService {
	return &ProgramService{
		programsRepo: programsRepo,
	}
}

func (s *ProgramService) CreateProgram(ctx context.Context, merchantID uuid.UUID, programName, pointCurrencyName string) (*postgres.Program, error) {
	if programName == "" || pointCurrencyName == "" {
		return nil, ErrInvalidProgram
	}

	program := &postgres.Program{
		ProgramID:         uuid.New(),
		MerchantID:        merchantID,
		ProgramName:       programName,
		PointCurrencyName: pointCurrencyName,
	}

	if err := s.programsRepo.Create(ctx, program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, programID uuid.UUID) (*postgres.Program, error) {
	program, err := s.programsRepo.GetByID(ctx, programID)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, ErrProgramNotFound
	}
	return program, nil
}

func (s *ProgramService) GetMerchantPrograms(ctx context.Context, merchantID uuid.UUID) ([]*postgres.Program, error) {
	return s.programsRepo.GetByMerchantID(ctx, merchantID)
}

func (s *ProgramService) UpdateProgram(ctx context.Context, programID uuid.UUID, programName, pointCurrencyName string) (*postgres.Program, error) {
	if programName == "" || pointCurrencyName == "" {
		return nil, ErrInvalidProgram
	}

	program, err := s.programsRepo.GetByID(ctx, programID)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, ErrProgramNotFound
	}

	program.ProgramName = programName
	program.PointCurrencyName = pointCurrencyName

	if err := s.programsRepo.Update(ctx, program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, programID uuid.UUID) error {
	err := s.programsRepo.Delete(ctx, programID)
	if err != nil {
		if errors.Is(err, ErrProgramNotFound) {
			return ErrProgramNotFound
		}
		return err
	}
	return nil
}
