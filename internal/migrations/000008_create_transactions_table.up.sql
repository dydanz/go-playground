-- Drop existing foreign key constraints from points_ledger
ALTER TABLE IF EXISTS points_ledger DROP CONSTRAINT IF EXISTS points_ledger_transaction_id_fkey;

-- Drop existing transactions table
DROP TABLE IF EXISTS transactions;

-- Create new transactions table
CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    program_id UUID NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'purchase', 'refund', 'bonus'
    transaction_amount DECIMAL(10,2) NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    branch_id UUID,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (merchant_id) REFERENCES merchants(merchant_id),
    FOREIGN KEY (customer_id) REFERENCES users(id),
    FOREIGN KEY (program_id) REFERENCES programs(program_id)
    -- TODO: Add branch_table
    -- FOREIGN KEY (branch_id) REFERENCES merchant_branches(branch_id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_transactions_merchant_id ON transactions(merchant_id);
CREATE INDEX idx_transactions_customer_id ON transactions(customer_id);
CREATE INDEX idx_transactions_program_id ON transactions(program_id);
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date); 