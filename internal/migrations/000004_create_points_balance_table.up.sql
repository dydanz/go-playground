CREATE TABLE IF NOT EXISTS points_balance (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_points INTEGER NOT NULL DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_balance UNIQUE (user_id)
);

-- Create an index on user_id for faster lookups
CREATE INDEX idx_points_balance_user_id ON points_balance(user_id); 