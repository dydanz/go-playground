DROP TABLE IF EXISTS points_balance; 

CREATE TABLE IF NOT EXISTS points_ledger (
    ledger_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_customers_id UUID NOT NULL REFERENCES merchant_customers(id),
    program_id UUID NOT NULL REFERENCES programs(program_id),
    points_earned INTEGER NOT NULL DEFAULT 0,
    points_redeemed INTEGER NOT NULL DEFAULT 0,
    points_balance INTEGER NOT NULL DEFAULT 0,
    transaction_id UUID REFERENCES transactions(transaction_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
); 

CREATE INDEX idx_points_ledger_merchant_customers_id ON points_ledger(merchant_customers_id);
CREATE INDEX idx_points_ledger_program_id ON points_ledger(program_id);
CREATE INDEX idx_points_ledger_transaction_date ON points_ledger(created_at);

