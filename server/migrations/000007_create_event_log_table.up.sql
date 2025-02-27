-- Create event type enum
CREATE TYPE event_type AS ENUM (
    'transaction_created',  -- Transaction created
    'program_id_created', -- Program ID created
    'program_rule_created', -- Program rule created
    'points_earned',        -- Points earned
    'points_redeemed',      -- Points redeemed
    'points_balance_updated', -- Points balance updated
    'reward_redeemed'       -- Reward redeemed
);

CREATE TYPE actor_type AS ENUM (
    'client',
    'merchant',
    'merchant_user',
    'superadmin'
);

CREATE TABLE IF NOT EXISTS event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type event_type NOT NULL,
    -- Use a Polymorphic Association instead of a foreign key
    actor_id UUID NOT NULL,
    actor_type actor_type NOT NULL,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reference_id UUID, -- Optional reference to related transaction/redemption
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequent query patterns
CREATE INDEX idx_event_log_user_id ON event_log(actor_id);
CREATE INDEX idx_event_log_event_type ON event_log(event_type);
CREATE INDEX idx_event_log_timestamp ON event_log(event_timestamp);
CREATE INDEX idx_event_log_reference_id ON event_log(reference_id);
-- Create GIN index for JSON querying
CREATE INDEX idx_event_log_details ON event_log USING GIN (details);

COMMENT ON TABLE event_log IS 'Audit log for all point-related events in the system';
COMMENT ON COLUMN event_log.details IS 'JSON structure containing event-specific details';
COMMENT ON COLUMN event_log.reference_id IS 'Optional UUID reference to related transaction or redemption'; 