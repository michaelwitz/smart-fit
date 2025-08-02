-- Smart Fit Girl - Complete Database Schema
-- This file shows the full database structure including users and survey system

-- ENUM Types
CREATE TYPE goal_category AS ENUM (
    'Weight',
    'Appearance', 
    'Strength',
    'Endurance'
);

-- Users table (already exists)
CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    identify_as VARCHAR(50),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country_code VARCHAR(2),
    locale VARCHAR(10),
    timezone VARCHAR(50),
    utc_offset INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Goals table - flexible goal definitions
CREATE TABLE GOALS (
    id SERIAL PRIMARY KEY,
    category goal_category NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Surveys table - immutable user surveys over time
CREATE TABLE SURVEYS (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    current_weight DECIMAL(5,2),
    target_weight DECIMAL(5,2),
    activity_level INTEGER CHECK (activity_level >= 0 AND activity_level <= 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Junction table for user survey goals (many-to-many)
-- Ensures one goal per category per survey (enforced at application level)
CREATE TABLE USER_SURVEY_GOALS (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES SURVEYS(id) ON DELETE CASCADE,
    goal_id INTEGER NOT NULL REFERENCES GOALS(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(survey_id, goal_id)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON USERS(email);
CREATE INDEX idx_surveys_user_id ON SURVEYS(user_id);
CREATE INDEX idx_surveys_created_at ON SURVEYS(created_at);
CREATE INDEX idx_user_survey_goals_survey_id ON USER_SURVEY_GOALS(survey_id);
CREATE INDEX idx_user_survey_goals_goal_id ON USER_SURVEY_GOALS(goal_id);
CREATE INDEX idx_goals_category ON GOALS(category);

-- Key Relationships:
-- 1. USERS 1:N SURVEYS (users can have multiple surveys over time)
-- 2. SURVEYS N:M GOALS (via USER_SURVEY_GOALS junction table)
-- 3. Business rule: Only one goal per category per survey (enforced in app logic)

-- Example survey data structure:
-- Survey ID 1 for User ID 5:
--   - current_weight: 170.5
--   - target_weight: 160.0
--   - activity_level: 7
--   - Goals selected:
--     * Weight: "Lose" (goal_id: 1)
--     * Appearance: "Lean" (goal_id: 5)  
--     * Strength: "Maintain" (goal_id: 8)
--     * Endurance: "Gain" (goal_id: 9)
