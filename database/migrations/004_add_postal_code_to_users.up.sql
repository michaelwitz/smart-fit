-- 004_add_postal_code_to_users.up.sql
-- Add postal_code column to USERS table after state_province

-- Add the postal_code column
ALTER TABLE USERS ADD COLUMN postal_code VARCHAR(20);
ALTER TABLE USERS ADD COLUMN locale VARCHAR(10);

-- Add comment for the new column
COMMENT ON COLUMN USERS.postal_code IS 'Postal code or ZIP code for user address';
COMMENT ON COLUMN USERS.locale IS 'Locale for user language settings (e.g., en-US, es-US)';
