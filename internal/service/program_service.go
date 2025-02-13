package service

import (
	"context"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ProgramService struct {
	programRepo domain.ProgramRepository
}

func NewProgramService(programRepo domain.ProgramRepository) *ProgramService {
	return &ProgramService{programRepo: programRepo}
}

func (s *ProgramService) Create(ctx context.Context, req *domain.CreateProgramRequest) (*domain.Program, error) {
	// Validate required fields
	if req.ProgramName == "" {
		return nil, domain.NewValidationError("program_name", "program name is required")
	}
	if req.PointCurrencyName == "" {
		return nil, domain.NewValidationError("point_currency_name", "point currency name is required")
	}

	program := &domain.Program{
		MerchantID:        req.MerchantID,
		ProgramName:       req.ProgramName,
		PointCurrencyName: req.PointCurrencyName,
		UserID:            req.UserID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	result, err := s.programRepo.Create(ctx, program)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.Create", err, "failed to create program")
	}

	return result, nil
}

func (s *ProgramService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.GetByID", err, "failed to get program")
	}
	if program == nil {
		return nil, domain.NewResourceNotFoundError("program", id.String(), "program not found")
	}
	return program, nil
}

func (s *ProgramService) GetAll(ctx context.Context) ([]*domain.Program, error) {
	programs, err := s.programRepo.GetAll(ctx)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.GetAll", err, "failed to get programs")
	}
	if len(programs) == 0 {
		return []*domain.Program{}, nil
	}
	return programs, nil
}

func (s *ProgramService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateProgramRequest) (*domain.Program, error) {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.Update", err, "failed to get program")
	}
	if program == nil {
		return nil, domain.NewResourceNotFoundError("program", id.String(), "program not found")
	}

	if req.ProgramName != "" {
		program.ProgramName = req.ProgramName
	}
	if req.PointCurrencyName != "" {
		program.PointCurrencyName = req.PointCurrencyName
	}
	program.UpdatedAt = time.Now()

	if err := s.programRepo.Update(ctx, program); err != nil {
		return nil, domain.NewSystemError("ProgramService.Update", err, "failed to update program")
	}

	return program, nil
}

func (s *ProgramService) Delete(ctx context.Context, id uuid.UUID) error {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil {
		return domain.NewSystemError("ProgramService.Delete", err, "failed to get program")
	}
	if program == nil {
		return domain.NewResourceNotFoundError("program", id.String(), "program not found")
	}

	if err := s.programRepo.Delete(ctx, id); err != nil {
		return domain.NewSystemError("ProgramService.Delete", err, "failed to delete program")
	}
	return nil
}

func (s *ProgramService) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*domain.Program, error) {
	programs, err := s.programRepo.GetByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.GetByMerchantID", err, "failed to get programs")
	}
	if len(programs) == 0 {
		return []*domain.Program{}, nil
	}
	return programs, nil
}

func (s *ProgramService) CreateProgram(ctx context.Context, merchantID uuid.UUID, programName, pointCurrencyName string) (*domain.CreateProgramResponse, error) {
	if programName == "" {
		return nil, domain.NewValidationError("program_name", "program name is required")
	}
	if pointCurrencyName == "" {
		return nil, domain.NewValidationError("point_currency_name", "point currency name is required")
	}

	program := &domain.Program{
		MerchantID:        merchantID,
		ProgramName:       programName,
		PointCurrencyName: pointCurrencyName,
	}

	result, err := s.programRepo.Create(ctx, program)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.CreateProgram", err, "failed to create program")
	}

	response := &domain.CreateProgramResponse{
		ProgramID:         result.ID,
		MerchantID:        result.MerchantID,
		ProgramName:       result.ProgramName,
		PointCurrencyName: result.PointCurrencyName,
		CreatedAt:         result.CreatedAt,
		UpdatedAt:         result.UpdatedAt,
	}

	return response, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, programID uuid.UUID) (*domain.Program, error) {
	program, err := s.programRepo.GetByID(ctx, programID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.GetProgram", err, "failed to get program")
	}
	if program == nil {
		return nil, domain.NewResourceNotFoundError("program", programID.String(), "program not found")
	}
	return program, nil
}

func (s *ProgramService) GetMerchantPrograms(ctx context.Context, merchantID uuid.UUID) ([]*domain.Program, error) {
	programs, err := s.programRepo.GetByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.GetMerchantPrograms", err, "failed to get programs")
	}
	if len(programs) == 0 {
		return []*domain.Program{}, nil
	}
	return programs, nil
}

func (s *ProgramService) UpdateProgram(ctx context.Context, programID uuid.UUID, programName, pointCurrencyName string) (*domain.Program, error) {
	if programName == "" {
		return nil, domain.NewValidationError("program_name", "program name is required")
	}
	if pointCurrencyName == "" {
		return nil, domain.NewValidationError("point_currency_name", "point currency name is required")
	}

	program, err := s.programRepo.GetByID(ctx, programID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramService.UpdateProgram", err, "failed to get program")
	}
	if program == nil {
		return nil, domain.NewResourceNotFoundError("program", programID.String(), "program not found")
	}

	program.ProgramName = programName
	program.PointCurrencyName = pointCurrencyName

	if err := s.programRepo.Update(ctx, program); err != nil {
		return nil, domain.NewSystemError("ProgramService.UpdateProgram", err, "failed to update program")
	}

	return program, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, programID uuid.UUID) error {
	program, err := s.programRepo.GetByID(ctx, programID)
	if err != nil {
		return domain.NewSystemError("ProgramService.DeleteProgram", err, "failed to get program")
	}
	if program == nil {
		return domain.NewResourceNotFoundError("program", programID.String(), "program not found")
	}

	if err := s.programRepo.Delete(ctx, programID); err != nil {
		return domain.NewSystemError("ProgramService.DeleteProgram", err, "failed to delete program")
	}
	return nil
}
