-- 003_add_utc_offset_to_users.up.sql
-- Add utc_offset column to USERS table

-- Add the new utc_offset column
ALTER TABLE USERS ADD COLUMN utc_offset INTEGER;

-- Add comment for the new column
COMMENT ON COLUMN USERS.utc_offset IS 'UTC offset in hours (e.g., -8 for PST, +5 for EST, 0 for UTC)';
