-- Create redemption status enum
CREATE TYPE redemption_status AS ENUM ('completed', 'pending', 'failed');

-- Redemptions refer to the process of redeeming points for rewards.
-- Disperse points to rewards.
CREATE TABLE IF NOT EXISTS redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_customers_id UUID NOT NULL REFERENCES merchant_customers(id),
    reward_id UUID NOT NULL REFERENCES rewards(id) ,
    points_used INTEGER NOT NULL CHECK (points_used > 0),
    redemption_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status redemption_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequent query patterns
CREATE INDEX idx_redemptions_merchant_customers_id ON redemptions(merchant_customers_id);
CREATE INDEX idx_redemptions_reward_id ON redemptions(reward_id);
CREATE INDEX idx_redemptions_status ON redemptions(status);
CREATE INDEX idx_redemptions_date ON redemptions(redemption_date);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_redemptions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create a trigger to automatically update updated_at
CREATE TRIGGER update_redemptions_updated_at
    BEFORE UPDATE ON redemptions
    FOR EACH ROW
    EXECUTE FUNCTION update_redemptions_updated_at(); 