package service

import (
	"context"
	"errors"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProgramNotFound = errors.New("program not found")
	ErrInvalidProgram  = errors.New("invalid program data")
)

type ProgramService struct {
	programRepo domain.ProgramRepository
}

func NewProgramService(programRepo domain.ProgramRepository) *ProgramService {
	return &ProgramService{programRepo: programRepo}
}

func (s *ProgramService) Create(req *domain.CreateProgramRequest) (*domain.Program, error) {
	program := &domain.Program{
		ID:                uuid.New().String(),
		MerchantID:        req.MerchantID,
		ProgramName:       req.ProgramName,
		PointCurrencyName: req.PointCurrencyName,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.programRepo.Create(program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) GetByID(id string) (*domain.Program, error) {
	return s.programRepo.GetByID(id)
}

func (s *ProgramService) GetAll() ([]*domain.Program, error) {
	return s.programRepo.GetAll()
}

func (s *ProgramService) Update(id string, req *domain.UpdateProgramRequest) (*domain.Program, error) {
	program, err := s.programRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ProgramName != "" {
		program.ProgramName = req.ProgramName
	}
	if req.PointCurrencyName != "" {
		program.PointCurrencyName = req.PointCurrencyName
	}
	program.UpdatedAt = time.Now()

	if err := s.programRepo.Update(program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) Delete(id string) error {
	return s.programRepo.Delete(id)
}

func (s *ProgramService) GetByMerchantID(merchantID string) ([]*domain.Program, error) {
	return s.programRepo.GetByMerchantID(merchantID)
}

func (s *ProgramService) CreateProgram(ctx context.Context, merchantID uuid.UUID, programName, pointCurrencyName string) (*domain.CreateProgramResponse, error) {
	if programName == "" || pointCurrencyName == "" {
		return nil, ErrInvalidProgram
	}

	program := &domain.Program{
		ID:                uuid.New().String(),
		MerchantID:        merchantID.String(),
		ProgramName:       programName,
		PointCurrencyName: pointCurrencyName,
	}

	if err := s.programRepo.Create(program); err != nil {
		return nil, err
	}

	response := &domain.CreateProgramResponse{
		ProgramID:         uuid.MustParse(program.ID),
		MerchantID:        merchantID,
		ProgramName:       program.ProgramName,
		PointCurrencyName: program.PointCurrencyName,
		CreatedAt:         program.CreatedAt,
		UpdatedAt:         program.UpdatedAt,
	}

	return response, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, programID uuid.UUID) (*domain.Program, error) {
	program, err := s.programRepo.GetByID(programID.String())
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, ErrProgramNotFound
	}
	return program, nil
}

func (s *ProgramService) GetMerchantPrograms(ctx context.Context, merchantID uuid.UUID) ([]*domain.Program, error) {
	return s.programRepo.GetByMerchantID(merchantID.String())
}

func (s *ProgramService) UpdateProgram(ctx context.Context, programID uuid.UUID, programName, pointCurrencyName string) (*domain.Program, error) {
	if programName == "" || pointCurrencyName == "" {
		return nil, ErrInvalidProgram
	}

	program, err := s.programRepo.GetByID(programID.String())
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, ErrProgramNotFound
	}

	program.ProgramName = programName
	program.PointCurrencyName = pointCurrencyName

	if err := s.programRepo.Update(program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, programID uuid.UUID) error {
	err := s.programRepo.Delete(programID.String())
	if err != nil {
		if errors.Is(err, ErrProgramNotFound) {
			return ErrProgramNotFound
		}
		return err
	}
	return nil
}
