-- Migration: notification
-- Keyspace: starnet_main

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS notification (
    id UUID,
    profile_id UUID,
    has_read BOOLEAN,
    title TEXT,
    message TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (profile_id, id)
);
