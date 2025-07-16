-- Migration: broadcast_read
-- Keyspace: starnet_main

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS broadcast_read (
    profile_id UUID,
    broadcast_id UUID,
    read_at TIMESTAMP,
    PRIMARY KEY (profile_id, broadcast_id)
);
