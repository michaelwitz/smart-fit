-- 002_create_users_table.up.sql
-- Migration to create USERS table for authentication and user management

-- Create USERS table
CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    identify_as VARCHAR(50),
    city VARCHAR(100),
    state_province VARCHAR(100),
    country_code CHAR(2),
    timezone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON USERS(email);
CREATE INDEX idx_users_full_name ON USERS(full_name);
CREATE INDEX idx_users_country_code ON USERS(country_code);
CREATE INDEX idx_users_created_at ON USERS(created_at);

-- Add comments for documentation
COMMENT ON TABLE USERS IS 'User accounts for authentication and profile management';
COMMENT ON COLUMN USERS.full_name IS 'Users full name for display purposes';
COMMENT ON COLUMN USERS.email IS 'Email address used as login identifier (unique)';
COMMENT ON COLUMN USERS.password IS 'Hashed password for authentication';
COMMENT ON COLUMN USERS.phone_number IS 'Optional phone number for contact';
COMMENT ON COLUMN USERS.identify_as IS 'Gender identity or preferred identification';
COMMENT ON COLUMN USERS.city IS 'City of residence';
COMMENT ON COLUMN USERS.state_province IS 'State or province of residence';
COMMENT ON COLUMN USERS.country_code IS '2-letter country code (ISO 3166-1 alpha-2)';
COMMENT ON COLUMN USERS.timezone IS 'User timezone (IANA timezone format, e.g., America/New_York)';
