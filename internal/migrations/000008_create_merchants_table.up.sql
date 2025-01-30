-- Create merchant type enum
CREATE TYPE merchant_type AS ENUM ('bank', 'e-commerce', 'repair_shop');

CREATE TABLE IF NOT EXISTS merchants (
    merchant_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_name VARCHAR(255) NOT NULL,
    merchant_type merchant_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_merchants_name ON merchants(merchant_name);
CREATE INDEX idx_merchants_type ON merchants(merchant_type);

-- Create update timestamp trigger
CREATE OR REPLACE FUNCTION update_merchants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_merchants_updated_at
    BEFORE UPDATE ON merchants
    FOR EACH ROW
    EXECUTE FUNCTION update_merchants_updated_at();