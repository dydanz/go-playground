package service

import (
	"context"
	"go-playground/server/domain"
)

type EventLoggerService struct {
	eventLogRepo domain.EventLogRepository
}

func NewEventLoggerService(eventLogRepo domain.EventLogRepository) *EventLoggerService {
	return &EventLoggerService{eventLogRepo: eventLogRepo}
}

func (s *EventLoggerService) SaveTransactionEvents(ctx context.Context, eventType domain.EventLogType, createdTx *domain.Transaction, pointsEarned int) error {
	// Log the transaction event
	event := &domain.EventLog{
		EventType:   string(domain.TransactionCreated),
		ActorID:     createdTx.MerchantCustomersID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := createdTx.TransactionID.String(); return &s }(),
		Details: map[string]interface{}{
			"transaction_id":     createdTx.TransactionID,
			"merchant_id":        createdTx.MerchantID,
			"program_id":         createdTx.ProgramID,
			"transaction_type":   createdTx.TransactionType,
			"transaction_amount": createdTx.TransactionAmount,
			"points_earned":      pointsEarned,
		},
	}

	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SaveRedemptionEvents(ctx context.Context, eventType domain.EventLogType, redemption *domain.Redemption, reward *domain.Reward) error {
	// Log the redemption event
	event := &domain.EventLog{
		EventType:   string(eventType),
		ActorID:     redemption.MerchantCustomersID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := redemption.ID.String(); return &s }(),
		Details: map[string]interface{}{
			"reward_id":     redemption.RewardID,
			"points_used":   redemption.PointsUsed,
			"redemption_id": redemption.ID,
			"program_id":    reward.ProgramID,
		},
	}

	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SaveUserUpdateEvents(ctx context.Context, eventType domain.EventLogType, user *domain.User) error {
	event := &domain.EventLog{
		EventType:   string(eventType),
		ActorID:     user.ID,
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := user.ID; return &s }(),
		Details: map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		},
	}

	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SaveMerchantUpdateEvents(ctx context.Context, eventType domain.EventLogType, merchant *domain.Merchant) error {
	event := &domain.EventLog{
		EventType:   string(eventType),
		ActorID:     merchant.ID.String(),
		ActorType:   string(domain.MerchantActorType),
		ReferenceID: func() *string { s := merchant.ID.String(); return &s }(),
		Details: map[string]interface{}{
			"user_id":     merchant.UserID,
			"merchant_id": merchant.ID,
			"name":        merchant.Name,
			"type":        merchant.Type,
			"status":      merchant.Status,
		},
	}
	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SaveProgramUpdateEvents(ctx context.Context, eventType domain.EventLogType, program *domain.Program) error {
	event := &domain.EventLog{
		EventType:   string(domain.ProgramUpdated),
		ActorID:     program.MerchantID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := program.ID.String(); return &s }(),
		Details: map[string]interface{}{
			"program_id": program.ID,
			"merchant":   program.MerchantID,
			"user_id":    program.UserID,
			"name":       program.ProgramName,
			"currency":   program.PointCurrencyName,
		},
	}
	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SaveProgramRulesEvents(ctx context.Context, eventType domain.EventLogType, program *domain.ProgramRule) error {
	event := &domain.EventLog{
		EventType:   string(eventType),
		ActorID:     "",
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := program.ID.String(); return &s }(),
		Details: map[string]interface{}{
			"id":              program.ID,
			"rulename":        program.RuleName,
			"condition_type":  program.ConditionType,
			"condition_value": program.ConditionValue,
			"multiplier":      program.Multiplier,
			"points_awarded":  program.PointsAwarded,
			"effective_from":  program.EffectiveFrom,
			"effective_to":    program.EffectiveTo,
			"program_id":      program.ProgramID,
		},
	}
	return s.eventLogRepo.Create(ctx, event)
}
func (s *EventLoggerService) SavePointUpdateEvents(ctx context.Context, eventType domain.EventLogType, ledger *domain.PointsLedger) error {
	event := &domain.EventLog{
		EventType:   string(eventType),
		ActorID:     ledger.MerchantCustomersID.String(),
		ActorType:   string(domain.MerchantUserActorType),
		ReferenceID: func() *string { s := ledger.LedgerID.String(); return &s }(),
		Details: map[string]interface{}{
			"points_id":       ledger.LedgerID,
			"customer_id":     ledger.MerchantCustomersID,
			"program_id":      ledger.ProgramID,
			"points_balance":  ledger.PointsBalance,
			"points_earned":   ledger.PointsEarned,
			"points_redeemed": ledger.PointsRedeemed,
			"transaction_id":  ledger.TransactionID,
			"created_at":      ledger.CreatedAt,
		},
	}
	return s.eventLogRepo.Create(ctx, event)
}
