
-- each program rule type has a condition_type and condition_value
-- one program_rule_type will be implement one rewards calculation function
CREATE TYPE program_rule_type AS ENUM (
    'program_rule_tenure', -- Tenure of the program
    'program_rule_transaction_amount', -- Transaction amount
    'program_rule_transaction_count', -- Transaction count
    'program_rule_transaction_type', -- Transaction type, ie using credit card, debit card, etc
    'program_rule_transaction_date', -- Transaction date, ie transaction on or after/before certain date
    'program_rule_transaction_category', -- Transaction category, ie food, travel, etc
    'program_rule_transaction_merchant', -- Transaction merchant, ie transaction at merchantid
    'program_rule_transaction_merchant_group' -- Transaction merchant group, ie merchant group id
);

CREATE TABLE IF NOT EXISTS program_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    program_id UUID NOT NULL REFERENCES programs(program_id),
    rule_name VARCHAR(255) NOT NULL,
    condition_type VARCHAR(50) NOT NULL,
    condition_value TEXT NOT NULL,
    multiplier DECIMAL(10,2) NOT NULL DEFAULT 1.0,
    points_awarded INTEGER NOT NULL DEFAULT 0,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_effective_dates CHECK (effective_to IS NULL OR effective_from <= effective_to)
);

CREATE INDEX idx_program_rules_program_id ON program_rules(program_id);
CREATE INDEX idx_program_rules_effective_dates ON program_rules(effective_from, effective_to);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_program_rules_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_program_rules_updated_at
    BEFORE UPDATE ON program_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_program_rules_updated_at();