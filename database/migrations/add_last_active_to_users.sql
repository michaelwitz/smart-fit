-- Migration: Add last_active column to USERS table
-- This column tracks user activity with 1-hour debouncing for analytics

-- Add the last_active column
ALTER TABLE USERS ADD COLUMN last_active TIMESTAMP;

-- Add comment for documentation
COMMENT ON COLUMN USERS.last_active IS 'Tracks last user activity with 1-hour debouncing for analytics';

-- Create index for efficient querying of active users
CREATE INDEX idx_users_last_active ON USERS(last_active);

-- Optional: Create a view for active users (active within last 24 hours)
CREATE VIEW active_users AS
SELECT 
    id,
    full_name,
    email,
    last_active,
    created_at,
    CASE 
        WHEN last_active IS NULL THEN 'Never active'
        WHEN last_active > (CURRENT_TIMESTAMP - INTERVAL '1 hour') THEN 'Online'
        WHEN last_active > (CURRENT_TIMESTAMP - INTERVAL '24 hours') THEN 'Active today'
        WHEN last_active > (CURRENT_TIMESTAMP - INTERVAL '7 days') THEN 'Active this week'
        ELSE 'Inactive'
    END as activity_status
FROM USERS
WHERE last_active IS NOT NULL OR created_at > (CURRENT_TIMESTAMP - INTERVAL '7 days')
ORDER BY last_active DESC NULLS LAST;
