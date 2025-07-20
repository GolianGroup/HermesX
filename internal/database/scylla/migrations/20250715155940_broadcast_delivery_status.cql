-- Migration: broadcast_delivery_status
-- Keyspace: starnet_main

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS broadcast_delivery_status (
    broadcast_id UUID PRIMARY KEY,
    event TEXT,
    status TEXT,
    attempts INT,
    last_attempt TIMESTAMP,
);
