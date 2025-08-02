-- Test Users Seed Data
-- This file contains test users for development and testing purposes
-- Note: Passwords are hashed using bcrypt

-- Insert test users
-- Password for all test users is "password" (hashed with bcrypt)
INSERT INTO USERS (
    full_name, 
    email, 
    password, 
    phone_number, 
    identify_as, 
    city, 
    state_province, 
    postal_code, 
    country_code, 
    locale, 
    timezone, 
    utc_offset, 
    created_at, 
    updated_at
) VALUES 
-- Sophia Woytowitz - Test user for development
(
    'Sophia Woytowitz',
    'sophia.woytowitz@gmail.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt hash of "password"
    NULL,
    'female',
    'New York',
    'NY',
    NULL,
    'US',
    'en-US',
    'America/New_York',
    -5,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
-- Additional test user for variety
(
    'Alex Johnson',
    'alex.johnson@example.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt hash of "password"
    '+1-555-0123',
    'non-binary',
    'San Francisco',
    'CA',
    '94102',
    'US',
    'en-US',
    'America/Los_Angeles',
    -8,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
-- International test user
(
    'Emma Thompson',
    'emma.thompson@example.co.uk',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt hash of "password"
    '+44-20-7946-0958',
    'female',
    'London',
    'England',
    'SW1A 1AA',
    'GB',
    'en-GB',
    'Europe/London',
    0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Note: The bcrypt hash '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi' 
-- corresponds to the password "password" and can be used for testing purposes.
-- In production, users should always set strong, unique passwords.
