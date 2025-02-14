-- Create merchant type enum
CREATE TYPE merchant_type AS ENUM ('bank', 'e-commerce', 'repair_shop');

-- Merchant (Stores or Business Units Under a Client)
--- A merchant represents a specific brand, store, or business that operates under a client.
--- Each merchant can have multiple branches or locations.
--- Defines its own loyalty rules (e.g., point earning and redemption rules).
-- Can have multiple staff/admins managing the system.
CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    merchant_name VARCHAR(255) NOT NULL,
    merchant_type merchant_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(16) NOT NULL DEFAULT 'active'
);

-- Create indexes
CREATE INDEX idx_merchants_user ON merchants(user_id);
CREATE INDEX idx_merchants_id ON merchants(id);

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

-- End Users (Customers or Members of the Loyalty Program)
--- Registered customers who earn and redeem points.
--- Identified through unique attributes (phone number, email, loyalty card, or user ID).
--- Engages in transactions, receives promotions, and participates in reward programs.
CREATE TABLE IF NOT EXISTS merchant_customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON merchant_customers(phone); 
CREATE INDEX IF NOT EXISTS idx_users_merchant_id ON merchant_customers(merchant_id); 