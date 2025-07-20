-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    metadata JSONB NOT NULL,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    is_protected BOOLEAN NOT NULL DEFAULT FALSE,
    template VARCHAR(20),
    preferred_channel VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Enforce constraints
ALTER TABLE events ADD CONSTRAINT chk_event_name_length CHECK (
    char_length(name) >= 3 AND char_length(name) <= 100
);

ALTER TABLE events ADD CONSTRAINT chk_preferred_channel CHECK (
    preferred_channel IN ('email', 'sms', 'push')
);

-- Composite index
CREATE INDEX idx_events_preferred_channel_name ON events (preferred_channel, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_preferred_channel_name;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
