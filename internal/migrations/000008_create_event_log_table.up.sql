-- Create event type enum
CREATE TYPE event_type AS ENUM ('transaction', 'balance_update', 'reward_redeemed');

CREATE TABLE IF NOT EXISTS event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type event_type NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reference_id UUID, -- Optional reference to related transaction/redemption
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequent query patterns
CREATE INDEX idx_event_log_user_id ON event_log(user_id);
CREATE INDEX idx_event_log_event_type ON event_log(event_type);
CREATE INDEX idx_event_log_timestamp ON event_log(event_timestamp);
CREATE INDEX idx_event_log_reference_id ON event_log(reference_id);
-- Create GIN index for JSON querying
CREATE INDEX idx_event_log_details ON event_log USING GIN (details);

-- Add foreign key constraints to ensure data integrity
ALTER TABLE event_log
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

COMMENT ON TABLE event_log IS 'Audit log for all point-related events in the system';
COMMENT ON COLUMN event_log.details IS 'JSON structure containing event-specific details';
COMMENT ON COLUMN event_log.reference_id IS 'Optional UUID reference to related transaction or redemption'; 