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
		return nil, err
	}

	return rule, nil
}

func (s *ProgramRulesService) GetByID(id string) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.programRuleRepo.GetByID(context.Background(), ruleID)
}

func (s *ProgramRulesService) GetByProgramID(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}
	return s.programRuleRepo.GetByProgramID(context.Background(), pID)
}

func (s *ProgramRulesService) Update(id string, req *domain.UpdateProgramRuleRequest) (*domain.ProgramRule, error) {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	rule, err := s.programRuleRepo.GetByID(context.Background(), ruleID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return rule, nil
}

func (s *ProgramRulesService) Delete(id string) error {
	ruleID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.programRuleRepo.Delete(context.Background(), ruleID)
}

func (s *ProgramRulesService) GetActiveRules(programID string) ([]*domain.ProgramRule, error) {
	pID, err := uuid.Parse(programID)
	if err != nil {
		return nil, err
	}
	return s.programRuleRepo.GetActiveRules(context.Background(), pID, time.Now())
}
