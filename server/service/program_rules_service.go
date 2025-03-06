package service

import (
	"context"
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ProgramRulesService struct {
	programRuleRepo domain.ProgramRuleRepository
	programRepo     domain.ProgramRepository
	logger          zerolog.Logger
}

func NewProgramRulesService(ruleRepo domain.ProgramRuleRepository, programRepo domain.ProgramRepository) *ProgramRulesService {
	return &ProgramRulesService{
		programRuleRepo: ruleRepo,
		programRepo:     programRepo,
		logger:          logging.GetLogger(),
	}
}

func (s *ProgramRulesService) Create(req *domain.CreateProgramRuleRequest) (*domain.ProgramRule, error) {
	// Validate required fields
	if req.RuleName == "" {
		s.logger.Error().
			Msg("Rule name is required")
		return nil, domain.NewValidationError("rule_name", "rule name is required")
	}
	if req.ConditionType == "" {
		s.logger.Error().
			Msg("Condition type is required")
		return nil, domain.NewValidationError("condition_type", "condition type is required")
	}
	if req.ConditionValue == "" {
		s.logger.Error().
			Msg("Condition value is required")
		return nil, domain.NewValidationError("condition_value", "condition value is required")
	}

	rule := &domain.ProgramRule{
		ID:             uuid.New(),
		ProgramID:      req.ProgramID,
		RuleName:       req.RuleName,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
		Multiplier:     req.Multiplier,
		PointsAwarded:  req.PointsAwarded,
		EffectiveFrom:  req.EffectiveFrom,
		EffectiveTo:    req.EffectiveTo,
	}

	if err := s.programRuleRepo.Create(context.Background(), rule); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error creating program rule")
		return nil, domain.NewSystemError("ProgramRulesService.Create", err, "failed to create program rule")
	}

	return rule, nil
}

func (s *ProgramRulesService) GetByID(id string) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid rule ID format")
		return nil, domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting program rule")
		return nil, domain.NewSystemError("ProgramRulesService.GetByID", err, "failed to get program rule")
	}
	if rule == nil {
		s.logger.Error().
			Msg("Program rule not found")
		return nil, domain.NewResourceNotFoundError("program rule", id, "rule not found")
	}

	return rule, nil
}

func (s *ProgramRulesService) GetByProgramID(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid program ID format")
		return nil, domain.NewValidationError("program_id", "invalid program ID format")
	}

	rules, err := s.programRuleRepo.GetByProgramID(context.Background(), pID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting program rules")
		return nil, domain.NewSystemError("ProgramRulesService.GetByProgramID", err, "failed to get program rules")
	}
	if len(rules) == 0 {
		s.logger.Info().
			Msg("No program rules found")
		return []*domain.ProgramRule{}, nil
	}

	return rules, nil
}

func (s *ProgramRulesService) Update(id string, req *domain.UpdateProgramRuleRequest) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid rule ID format")
		return nil, domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting program rule")
		return nil, domain.NewSystemError("ProgramRulesService.Update", err, "failed to get program rule")
	}
	if rule == nil {
		s.logger.Error().
			Msg("Program rule not found")
		return nil, domain.NewResourceNotFoundError("program rule", id, "rule not found")
	}

	if req.RuleName != "" {
		rule.RuleName = req.RuleName
	}
	if req.ConditionType != "" {
		rule.ConditionType = req.ConditionType
	}
	if req.ConditionValue != "" {
		rule.ConditionValue = req.ConditionValue
	}
	if req.Multiplier != nil {
		rule.Multiplier = *req.Multiplier
	}
	if req.PointsAwarded != nil {
		rule.PointsAwarded = *req.PointsAwarded
	}
	if req.EffectiveFrom != nil {
		rule.EffectiveFrom = *req.EffectiveFrom
	}
	if req.EffectiveTo != nil {
		rule.EffectiveTo = req.EffectiveTo
	}

	if err := s.programRuleRepo.Update(context.Background(), rule); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error updating program rule")
		return nil, domain.NewSystemError("ProgramRulesService.Update", err, "failed to update program rule")
	}

	return rule, nil
}

func (s *ProgramRulesService) Delete(id string) error {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid rule ID format")
		return domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting program rule")
		return domain.NewSystemError("ProgramRulesService.Delete", err, "failed to get program rule")
	}
	if rule == nil {
		s.logger.Error().
			Msg("Program rule not found")
		return domain.NewResourceNotFoundError("program rule", id, "rule not found")
	}

	if err := s.programRuleRepo.Delete(context.Background(), ruleID); err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error deleting program rule")
		return domain.NewSystemError("ProgramRulesService.Delete", err, "failed to delete program rule")
	}

	return nil
}

func (s *ProgramRulesService) GetActiveRules(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid program ID format")
		return nil, domain.NewValidationError("program_id", "invalid program ID format")
	}

	rules, err := s.programRuleRepo.GetActiveRules(context.Background(), pID, time.Now())
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting active program rules")
		return nil, domain.NewSystemError("ProgramRulesService.GetActiveRules", err, "failed to get active program rules")
	}
	if len(rules) == 0 {
		s.logger.Info().
			Msg("No active program rules found")
		return []*domain.ProgramRule{}, nil
	}

	return rules, nil
}

type ProgramRuleWithProgram struct {
	ProgramID      uuid.UUID  `json:"program_id"`
	ProgramName    string     `json:"program_name"`
	RuleName       string     `json:"rule_name"`
	ConditionType  string     `json:"condition_type"`
	ConditionValue string     `json:"condition_value"`
	Multiplier     float64    `json:"multiplier"`
	PointsAwarded  int        `json:"points_awarded"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
}

func (s *ProgramRulesService) GetProgramRulesByMerchantId(merchantID string, page, limit int) ([]ProgramRuleWithProgram, int64, error) {
	mID, err := uuid.Parse(merchantID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Invalid merchant ID format")
		return nil, 0, domain.NewValidationError("merchant_id", "invalid merchant ID format")
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get all programs for the merchant
	programs, err := s.programRepo.GetByMerchantID(context.Background(), mID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error getting merchant programs")
		return nil, 0, domain.NewSystemError("ProgramRulesService.GetProgramRulesByMerchantId", err, "failed to get merchant programs")
	}

	var result []ProgramRuleWithProgram

	// For each program, get its rules
	for _, program := range programs {
		rules, err := s.programRuleRepo.GetByProgramID(context.Background(), program.ID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msg("Error getting program rules")
			return nil, 0, domain.NewSystemError("ProgramRulesService.GetProgramRulesByMerchantId", err, "failed to get program rules")
		}

		// Map each rule to the response format
		for _, rule := range rules {
			result = append(result, ProgramRuleWithProgram{
				ProgramID:      program.ID,
				ProgramName:    program.ProgramName,
				RuleName:       rule.RuleName,
				ConditionType:  rule.ConditionType,
				ConditionValue: rule.ConditionValue,
				Multiplier:     rule.Multiplier,
				PointsAwarded:  rule.PointsAwarded,
				EffectiveFrom:  rule.EffectiveFrom,
				EffectiveTo:    rule.EffectiveTo,
			})
		}
	}

	if len(result) == 0 {
		return []ProgramRuleWithProgram{}, 0, nil
	}

	// Calculate total count
	total := int64(len(result))

	// Apply pagination, TODO: IMPROVE PAGINATION PROPERLY
	start := offset
	end := offset + limit
	if start >= len(result) {
		return []ProgramRuleWithProgram{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], total, nil
}
