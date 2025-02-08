DROP TABLE IF EXISTS points_balance; 

CREATE TABLE IF NOT EXISTS points_ledger (
    ledger_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    program_id UUID NOT NULL,
    points_earned INTEGER NOT NULL DEFAULT 0,
    points_redeemed INTEGER NOT NULL DEFAULT 0,
    points_balance INTEGER NOT NULL DEFAULT 0,
    transaction_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES users(id),
    FOREIGN KEY (program_id) REFERENCES programs(program_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
); 