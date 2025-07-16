-- Migration: broadcast
-- Keyspace: starnet_main

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS broadcast (
    id UUID PRIMARY KEY,
    title TEXT,
    message TEXT,
    created_at TIMESTAMP,
);
