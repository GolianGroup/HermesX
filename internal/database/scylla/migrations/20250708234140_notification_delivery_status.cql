-- Migration: notification_delivery_status
-- Keyspace: starnet_main

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS notification_delivery_status (
    notification_id UUID,
    profile_id UUID,
    channel TEXT,
    status TEXT,
    attempts INT,
    last_attempt TIMESTAMP,
    PRIMARY KEY (profile_id, notification_id)
) WITH default_time_to_live = 86400; -- expires in 1 day
