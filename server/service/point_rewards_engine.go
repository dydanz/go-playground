package service

import (
	"fmt"
	"time"
)

type Transaction struct {
	Amount           float64
	Type             string
	Category         string
	MerchantID       string
	MerchantGroupID  string
	TransactionCount int
	MembershipTenure int // in days
}

type ProgramRule struct {
	RuleName       string
	ConditionType  string
	ConditionValue string
	Multiplier     float64
	PointsAwarded  int
	EffectiveFrom  time.Time
	EffectiveTo    *time.Time
}

type Program struct {
	ProgramID         string
	MerchantID        string
	UserID            string
	ProgramName       string
	PointCurrencyName string
	Rules             []ProgramRule
}

func evaluateRule(rule ProgramRule, tx Transaction) (bool, float64) {
	switch rule.ConditionType {
	case "program_rule_tenure":
		// Check if membership tenure meets the condition
		tenure := tx.MembershipTenure
		return float64(tenure) > parseConditionValue(rule.ConditionValue), rule.Multiplier * float64(rule.PointsAwarded)

	case "program_rule_transaction_amount":
		// Check if transaction amount meets the condition
		if tx.Amount > parseConditionValue(rule.ConditionValue) {
			if rule.PointsAwarded == 0 {
				return true, rule.Multiplier * tx.Amount
			}
			return true, rule.Multiplier * float64(rule.PointsAwarded)
		}
	case "program_rule_transaction_count":
		// Check if transaction count meets the condition
		count := tx.TransactionCount
		return float64(count) > parseConditionValue(rule.ConditionValue), rule.Multiplier * float64(rule.PointsAwarded)

	case "program_rule_transaction_type":
		// Check if transaction type matches the condition
		return tx.Type == rule.ConditionValue, rule.Multiplier * float64(rule.PointsAwarded)

	case "program_rule_transaction_category":
		// Check if transaction category matches the condition
		return tx.Category == rule.ConditionValue, rule.Multiplier * float64(rule.PointsAwarded)

	case "program_rule_transaction_merchant":
		// Check if transaction merchant matches the condition
		return tx.MerchantID == rule.ConditionValue, rule.Multiplier * float64(rule.PointsAwarded)

	case "program_rule_transaction_merchant_group":
		// Check if transaction merchant group matches the condition
		return tx.MerchantGroupID == rule.ConditionValue, rule.Multiplier * float64(rule.PointsAwarded)
	}
	return false, 0
}

// Helper function to parse condition value (e.g., "> 100")
func parseConditionValue(condition string) float64 {
	// Implement logic to parse condition strings like "> 100"
	// For simplicity, assume the condition is a numeric value
	var value float64
	fmt.Sscanf(condition, "%f", &value)
	return value
}

func calculatePoints(program Program, tx Transaction) float64 {
	totalPoints := 0.0
	basePoints := 0.0
	bonusPoints := 0.0
	now := time.Now()

	// First pass: Calculate base points from transaction amount rules
	for _, rule := range program.Rules {
		if rule.ConditionType == "program_rule_transaction_amount" {
			// Skip expired or future rules
			if now.Before(rule.EffectiveFrom) || (rule.EffectiveTo != nil && now.After(*rule.EffectiveTo)) {
				continue
			}

			matches, points := evaluateRule(rule, tx)
			if matches {
				basePoints += points
			}
		}
	}

	// Second pass: Calculate bonus points from other rules
	for _, rule := range program.Rules {
		// Skip expired or future rules
		if now.Before(rule.EffectiveFrom) || (rule.EffectiveTo != nil && now.After(*rule.EffectiveTo)) {
			continue
		}

		// Skip transaction amount rules as they were handled in first pass
		if rule.ConditionType != "program_rule_transaction_amount" {
			matches, points := evaluateRule(rule, tx)
			if matches {
				bonusPoints += points
			}
		}
	}

	totalPoints = basePoints + bonusPoints
	return totalPoints
}

/*

func main() {
    // Example transaction
    tx := Transaction{
        Amount:           150.0,
        Type:             "credit_card",
        Category:         "food",
        MerchantID:       "merchant_123",
        MerchantGroupID:  "group_456",
        TransactionCount: 12,
        MembershipTenure: 365,
    }

    // Example program with rules
    program := Program{
        ProgramID: "program_123",
        Rules: []ProgramRule{
            {
                RuleName:       "10th Transaction",
                ConditionType:  "program_rule_transaction_count",
                ConditionValue: "10",
                Multiplier:     1.0,
                PointsAwarded:  100,
                EffectiveFrom:  time.Now().AddDate(0, -1, 0), // Active since 1 month ago
                EffectiveTo:    nil,
            },
            {
                RuleName:       "Transaction above $100",
                ConditionType:  "program_rule_transaction_amount",
                ConditionValue: "100",
                Multiplier:     0.05, // 5% of transaction amount
                PointsAwarded:  0,
                EffectiveFrom:  time.Now().AddDate(0, -1, 0),
                EffectiveTo:    nil,
            },
        },
    }

    // Calculate points
    points := calculatePoints(program, tx)
    fmt.Printf("Total points awarded: %.2f\n", points)
}

*/
