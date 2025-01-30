package service

import (
	"context"
	"errors"
	"go-playground/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ProgramRuleService struct {
	programRuleRepo domain.ProgramRuleRepository
	eventRepo       domain.EventLogRepository
}

func NewProgramRuleService(
	programRuleRepo domain.ProgramRuleRepository,
	eventRepo domain.EventLogRepository,
) *ProgramRuleService {
	return &ProgramRuleService{
		programRuleRepo: programRuleRepo,
		eventRepo:       eventRepo,
	}
}

func (s *ProgramRuleService) Create(ctx context.Context, req *domain.CreateProgramRuleRequest) (*domain.ProgramRule, error) {
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

	if err := s.programRuleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	// Log event
	ruleIDStr := rule.ID.String()
	event := &domain.EventLog{
		EventType:   "program_rule_created",
		ReferenceID: &ruleIDStr,
		Details: map[string]interface{}{
			"program_id":      rule.ProgramID,
			"rule_name":       rule.RuleName,
			"condition_type":  rule.ConditionType,
			"condition_value": rule.ConditionValue,
			"multiplier":      rule.Multiplier,
			"points_awarded":  rule.PointsAwarded,
			"effective_from":  rule.EffectiveFrom,
			"effective_to":    rule.EffectiveTo,
		},
	}
	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *ProgramRuleService) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProgramRule, error) {
	return s.programRuleRepo.GetByID(ctx, id)
}

func (s *ProgramRuleService) GetByProgramID(ctx context.Context, programID uuid.UUID) ([]*domain.ProgramRule, error) {
	return s.programRuleRepo.GetByProgramID(ctx, programID)
}

func (s *ProgramRuleService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateProgramRuleRequest) (*domain.ProgramRule, error) {
	rule, err := s.programRuleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("resources not found")
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

	if err := s.programRuleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}

	// Log event
	ruleIDStr := rule.ID.String()
	event := &domain.EventLog{
		EventType:   "program_rule_updated",
		ReferenceID: &ruleIDStr,
		Details: map[string]interface{}{
			"program_id":      rule.ProgramID,
			"rule_name":       rule.RuleName,
			"condition_type":  rule.ConditionType,
			"condition_value": rule.ConditionValue,
			"multiplier":      rule.Multiplier,
			"points_awarded":  rule.PointsAwarded,
			"effective_from":  rule.EffectiveFrom,
			"effective_to":    rule.EffectiveTo,
		},
	}
	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *ProgramRuleService) Delete(ctx context.Context, id uuid.UUID) error {
	rule, err := s.programRuleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rule == nil {
		return errors.New("resources not found")
	}

	if err := s.programRuleRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log event
	ruleIDStr := rule.ID.String()
	event := &domain.EventLog{
		EventType:   "program_rule_deleted",
		ReferenceID: &ruleIDStr,
		Details: map[string]interface{}{
			"program_id": rule.ProgramID,
			"rule_name":  rule.RuleName,
		},
	}
	return s.eventRepo.Create(event)
}

func (s *ProgramRuleService) GetActiveRules(ctx context.Context, programID uuid.UUID) ([]*domain.ProgramRule, error) {
	return s.programRuleRepo.GetActiveRules(ctx, programID, time.Now())
}
