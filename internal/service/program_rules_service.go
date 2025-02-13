package service

import (
	"context"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ProgramRulesService struct {
	programRuleRepo domain.ProgramRuleRepository
}

func NewProgramRulesService(repo domain.ProgramRuleRepository) *ProgramRulesService {
	return &ProgramRulesService{programRuleRepo: repo}
}

func (s *ProgramRulesService) Create(req *domain.CreateProgramRuleRequest) (*domain.ProgramRule, error) {
	// Validate required fields
	if req.RuleName == "" {
		return nil, domain.NewValidationError("rule_name", "rule name is required")
	}
	if req.ConditionType == "" {
		return nil, domain.NewValidationError("condition_type", "condition type is required")
	}
	if req.ConditionValue == "" {
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
		return nil, domain.NewSystemError("ProgramRulesService.Create", err, "failed to create program rule")
	}

	return rule, nil
}

func (s *ProgramRulesService) GetByID(id string) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramRulesService.GetByID", err, "failed to get program rule")
	}
	if rule == nil {
		return nil, domain.NewResourceNotFoundError("program rule", id, "rule not found")
	}

	return rule, nil
}

func (s *ProgramRulesService) GetByProgramID(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		return nil, domain.NewValidationError("program_id", "invalid program ID format")
	}

	rules, err := s.programRuleRepo.GetByProgramID(context.Background(), pID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramRulesService.GetByProgramID", err, "failed to get program rules")
	}
	if len(rules) == 0 {
		return []*domain.ProgramRule{}, nil
	}

	return rules, nil
}

func (s *ProgramRulesService) Update(id string, req *domain.UpdateProgramRuleRequest) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		return nil, domain.NewSystemError("ProgramRulesService.Update", err, "failed to get program rule")
	}
	if rule == nil {
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
		return nil, domain.NewSystemError("ProgramRulesService.Update", err, "failed to update program rule")
	}

	return rule, nil
}

func (s *ProgramRulesService) Delete(id string) error {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewValidationError("id", "invalid rule ID format")
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		return domain.NewSystemError("ProgramRulesService.Delete", err, "failed to get program rule")
	}
	if rule == nil {
		return domain.NewResourceNotFoundError("program rule", id, "rule not found")
	}

	if err := s.programRuleRepo.Delete(context.Background(), ruleID); err != nil {
		return domain.NewSystemError("ProgramRulesService.Delete", err, "failed to delete program rule")
	}

	return nil
}

func (s *ProgramRulesService) GetActiveRules(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		return nil, domain.NewValidationError("program_id", "invalid program ID format")
	}

	rules, err := s.programRuleRepo.GetActiveRules(context.Background(), pID, time.Now())
	if err != nil {
		return nil, domain.NewSystemError("ProgramRulesService.GetActiveRules", err, "failed to get active program rules")
	}
	if len(rules) == 0 {
		return []*domain.ProgramRule{}, nil
	}

	return rules, nil
}
