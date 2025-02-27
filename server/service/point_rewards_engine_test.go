package service

import (
	"testing"
	"time"
)

func TestCalculatePoints(t *testing.T) {
	// Helper function to create a time pointer
	timePtr := func(t time.Time) *time.Time {
		return &t
	}

	// Base time for testing
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		program  Program
		tx       Transaction
		expected float64
	}{
		{
			name: "Transaction-Based Points - Minimum Spend",
			program: Program{
				ProgramID: "prog1",
				Rules: []ProgramRule{
					{
						RuleName:       "Spend $100 get 10 points",
						ConditionType:  "program_rule_transaction_amount",
						ConditionValue: "100",
						Multiplier:     1.0,
						PointsAwarded:  10,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount: 150.0,
			},
			expected: 10.0,
		},
		{
			name: "Transaction-Based Points - Category Multiplier",
			program: Program{
				ProgramID: "prog2",
				Rules: []ProgramRule{
					{
						RuleName:       "2x points on dining",
						ConditionType:  "program_rule_transaction_category",
						ConditionValue: "dining",
						Multiplier:     2.0,
						PointsAwarded:  100,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:   100.0,
				Category: "dining",
			},
			expected: 200.0,
		},
		{
			name: "Frequency-Based Points - Transaction Count",
			program: Program{
				ProgramID: "prog3",
				Rules: []ProgramRule{
					{
						RuleName:       "Bonus for 5+ transactions",
						ConditionType:  "program_rule_transaction_count",
						ConditionValue: "5",
						Multiplier:     1.0,
						PointsAwarded:  500,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:           100.0,
				TransactionCount: 6,
			},
			expected: 500.0,
		},
		{
			name: "Loyalty Milestone - Membership Tenure",
			program: Program{
				ProgramID: "prog4",
				Rules: []ProgramRule{
					{
						RuleName:       "1 year membership bonus",
						ConditionType:  "program_rule_tenure",
						ConditionValue: "365",
						Multiplier:     1.0,
						PointsAwarded:  1000,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:           100.0,
				MembershipTenure: 366,
			},
			expected: 1000.0,
		},
		{
			name: "Category-Based Points - Merchant Group",
			program: Program{
				ProgramID: "prog5",
				Rules: []ProgramRule{
					{
						RuleName:       "Partner merchant group bonus",
						ConditionType:  "program_rule_transaction_merchant_group",
						ConditionValue: "premium_partners",
						Multiplier:     3.0,
						PointsAwarded:  100,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:          100.0,
				MerchantGroupID: "premium_partners",
			},
			expected: 300.0,
		},
		{
			name: "Multiple Rules Combined",
			program: Program{
				ProgramID: "prog6",
				Rules: []ProgramRule{
					{
						RuleName:       "Base points per amount",
						ConditionType:  "program_rule_transaction_amount",
						ConditionValue: "50",
						Multiplier:     0.1, // 0.1 points per dollar
						PointsAwarded:  0,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
					{
						RuleName:       "Category bonus",
						ConditionType:  "program_rule_transaction_category",
						ConditionValue: "electronics",
						Multiplier:     2.0,
						PointsAwarded:  50,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:   100.0,
				Category: "electronics",
			},
			expected: 110.0, // 10 points from amount (100 * 0.1) + 100 points from category bonus (50 * 2)
		},
		// {
		// 	name: "Tiered Spending - Multiple Tiers",
		// 	program: Program{
		// 		ProgramID: "prog7",
		// 		Rules: []ProgramRule{
		// 			{
		// 				RuleName:       "Tier 1: First $100",
		// 				ConditionType:  "program_rule_transaction_amount",
		// 				ConditionValue: "0",
		// 				Multiplier:     1.0, // 1 point per dollar
		// 				PointsAwarded:  0,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 			{
		// 				RuleName:       "Tier 2: $101-$200",
		// 				ConditionType:  "program_rule_transaction_amount",
		// 				ConditionValue: "100",
		// 				Multiplier:     2.0, // 2 points per dollar
		// 				PointsAwarded:  0,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 			{
		// 				RuleName:       "Tier 3: Above $200",
		// 				ConditionType:  "program_rule_transaction_amount",
		// 				ConditionValue: "200",
		// 				Multiplier:     3.0, // 3 points per dollar
		// 				PointsAwarded:  0,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 		},
		// 	},
		// 	tx: Transaction{
		// 		Amount: 250.0,
		// 	},
		// 	expected: 450.0, // (100 * 1) + (100 * 2) + (50 * 3)
		// },
		{
			name: "Promotional Points - Limited Time Offer",
			program: Program{
				ProgramID: "prog8",
				Rules: []ProgramRule{
					{
						RuleName:       "Holiday Season Double Points",
						ConditionType:  "program_rule_transaction_amount",
						ConditionValue: "0",
						Multiplier:     2.0,
						PointsAwarded:  0,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount: 100.0,
			},
			expected: 200.0,
		},
		// {
		// 	name: "Complex Combination - Multiple Rules",
		// 	program: Program{
		// 		ProgramID: "prog9",
		// 		Rules: []ProgramRule{
		// 			{
		// 				RuleName:       "Base Points",
		// 				ConditionType:  "program_rule_transaction_amount",
		// 				ConditionValue: "0",
		// 				Multiplier:     1.0,
		// 				PointsAwarded:  0,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 			{
		// 				RuleName:       "Premium Merchant Bonus",
		// 				ConditionType:  "program_rule_transaction_merchant_group",
		// 				ConditionValue: "premium",
		// 				Multiplier:     2.0,
		// 				PointsAwarded:  50,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 			{
		// 				RuleName:       "Loyal Customer Bonus",
		// 				ConditionType:  "program_rule_tenure",
		// 				ConditionValue: "365",
		// 				Multiplier:     1.5,
		// 				PointsAwarded:  100,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 			{
		// 				RuleName:       "Frequent Shopper Bonus",
		// 				ConditionType:  "program_rule_transaction_count",
		// 				ConditionValue: "10",
		// 				Multiplier:     1.0,
		// 				PointsAwarded:  200,
		// 				EffectiveFrom:  yesterday,
		// 				EffectiveTo:    timePtr(tomorrow),
		// 			},
		// 		},
		// 	},
		// 	tx: Transaction{
		// 		Amount:           300.0,
		// 		MerchantGroupID:  "premium",
		// 		MembershipTenure: 400,
		// 		TransactionCount: 15,
		// 	},
		// 	expected: 800.0, // 300 (base) + 100 (premium) + 150 (loyal) + 200 (frequent)
		// },
		{
			name: "Category Specific with Minimum Spend",
			program: Program{
				ProgramID: "prog10",
				Rules: []ProgramRule{
					{
						RuleName:       "Electronics Category Base",
						ConditionType:  "program_rule_transaction_category",
						ConditionValue: "electronics",
						Multiplier:     1.0,
						PointsAwarded:  50,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
					{
						RuleName:       "High Value Purchase Bonus",
						ConditionType:  "program_rule_transaction_amount",
						ConditionValue: "500",
						Multiplier:     2.0,
						PointsAwarded:  100,
						EffectiveFrom:  yesterday,
						EffectiveTo:    timePtr(tomorrow),
					},
				},
			},
			tx: Transaction{
				Amount:   600.0,
				Category: "electronics",
			},
			expected: 250.0, // 50 (category) + 200 (high value)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePoints(tt.program, tt.tx)
			if got != tt.expected {
				t.Errorf("calculatePoints() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEvaluateRule(t *testing.T) {
	tests := []struct {
		name        string
		rule        ProgramRule
		tx          Transaction
		wantMatches bool
		wantPoints  float64
	}{
		{
			name: "Transaction Amount Rule - Above Threshold",
			rule: ProgramRule{
				ConditionType:  "program_rule_transaction_amount",
				ConditionValue: "100",
				Multiplier:     0.1,
				PointsAwarded:  0,
			},
			tx: Transaction{
				Amount: 150.0,
			},
			wantMatches: true,
			wantPoints:  15.0, // 150 * 0.1
		},
		{
			name: "Transaction Amount Rule - Below Threshold",
			rule: ProgramRule{
				ConditionType:  "program_rule_transaction_amount",
				ConditionValue: "100",
				Multiplier:     0.1,
				PointsAwarded:  0,
			},
			tx: Transaction{
				Amount: 50.0,
			},
			wantMatches: false,
			wantPoints:  0,
		},
		{
			name: "Transaction Type Rule - Matching",
			rule: ProgramRule{
				ConditionType:  "program_rule_transaction_type",
				ConditionValue: "credit_card",
				Multiplier:     2.0,
				PointsAwarded:  50,
			},
			tx: Transaction{
				Type: "credit_card",
			},
			wantMatches: true,
			wantPoints:  100.0, // 50 * 2
		},
		{
			name: "Membership Tenure Rule - Exceeding",
			rule: ProgramRule{
				ConditionType:  "program_rule_tenure",
				ConditionValue: "365",
				Multiplier:     1.0,
				PointsAwarded:  1000,
			},
			tx: Transaction{
				MembershipTenure: 400,
			},
			wantMatches: true,
			wantPoints:  1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, gotPoints := evaluateRule(tt.rule, tt.tx)
			if gotMatches != tt.wantMatches {
				t.Errorf("evaluateRule() matches = %v, want %v", gotMatches, tt.wantMatches)
			}
			if gotPoints != tt.wantPoints {
				t.Errorf("evaluateRule() points = %v, want %v", gotPoints, tt.wantPoints)
			}
		})
	}
}
