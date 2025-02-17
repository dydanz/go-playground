package service

import (
	"context"
	"go-playground/internal/domain"
	"log"
	"time"

	"github.com/google/uuid"
)

type ProgramRulesService struct {
	programRuleRepo domain.ProgramRuleRepository
	programRepo     domain.ProgramRepository
}

func NewProgramRulesService(ruleRepo domain.ProgramRuleRepository, programRepo domain.ProgramRepository) *ProgramRulesService {
	return &ProgramRulesService{
		programRuleRepo: ruleRepo,
		programRepo:     programRepo,
	}
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

func (s *ProgramRulesService) GetProgramRulesByMerchantId(merchantID string) ([]ProgramRuleWithProgram, error) {
	mID, err := uuid.Parse(merchantID)
	if err != nil {
		log.Printf("invalid merchant ID format. error: %v", err)
		return nil, domain.NewValidationError("merchant_id", "invalid merchant ID format")
	}

	log.Printf("merchant ID: %v", mID)

	// Get all programs for the merchant
	programs, err := s.programRepo.GetByMerchantID(context.Background(), mID)
	if err != nil {
		log.Printf("failed to get merchant programs. error: %v", err)
		return nil, domain.NewSystemError("ProgramRulesService.GetProgramRulesByMerchantId", err, "failed to get merchant programs")
	}

	log.Printf("programs: %v", programs)

	var result []ProgramRuleWithProgram

	// For each program, get its rules
	for _, program := range programs {

		rules, err := s.programRuleRepo.GetByProgramID(context.Background(), program.ID)
		if err != nil {
			log.Printf("failed to get program rules. error: %v", err)
			return nil, domain.NewSystemError("ProgramRulesService.GetProgramRulesByMerchantId", err, "failed to get program rules")
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
		return []ProgramRuleWithProgram{}, nil
	}

	return result, nil
}
