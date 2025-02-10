CREATE TABLE IF NOT EXISTS points_balance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_customers_id UUID NOT NULL REFERENCES merchant_customers(id),
    total_points INTEGER NOT NULL DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_balance UNIQUE (merchant_customers_id)
);

-- Create an index on user_id for faster lookups
CREATE INDEX idx_points_balance_user_id ON points_balance(merchant_customers_id); 
