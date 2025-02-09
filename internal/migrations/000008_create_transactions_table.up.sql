-- Drop existing foreign key constraints from points_ledger
ALTER TABLE IF EXISTS points_ledger DROP CONSTRAINT IF EXISTS points_ledger_transaction_id_fkey;

-- Drop existing transactions table
DROP TABLE IF EXISTS transactions;

-- The `transactions` table is used to record all activities associated with 
-- a `merchant_customer_id`. These records are essential for tracking and 
-- calculating whether the activities meet the `program_rules` to qualify 
-- for earning points.
CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    merchant_customers_id UUID NOT NULL REFERENCES merchant_customers(id),
    program_id UUID NOT NULL REFERENCES programs(program_id),
    transaction_type VARCHAR(50) NOT NULL, -- 'purchase', 'refund', 'bonus'
    transaction_amount DECIMAL(10,2) NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    branch_id UUID,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

-- Create indexes
CREATE INDEX idx_transactions_merchant_id ON transactions(merchant_id);
CREATE INDEX idx_transactions_merchant_customers_id ON transactions(merchant_customers_id);
CREATE INDEX idx_transactions_program_id ON transactions(program_id);
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date); 