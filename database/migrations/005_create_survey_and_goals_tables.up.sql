-- Create ENUM for goal categories
CREATE TYPE goal_category AS ENUM (
    'Weight',
    'Appearance', 
    'Strength',
    'Endurance'
);

-- Create GOALS table
CREATE TABLE GOALS (
    id SERIAL PRIMARY KEY,
    category goal_category NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create SURVEYS table (immutable, no updated_at)
CREATE TABLE SURVEYS (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    current_weight DECIMAL(5,2),
    target_weight DECIMAL(5,2),
    activity_level INTEGER CHECK (activity_level >= 0 AND activity_level <= 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create USER_SURVEY_GOALS junction table (many-to-many)
CREATE TABLE USER_SURVEY_GOALS (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES SURVEYS(id) ON DELETE CASCADE,
    goal_id INTEGER NOT NULL REFERENCES GOALS(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(survey_id, goal_id)
);

-- Create indexes for better performance
CREATE INDEX idx_surveys_user_id ON SURVEYS(user_id);
CREATE INDEX idx_surveys_created_at ON SURVEYS(created_at);
CREATE INDEX idx_user_survey_goals_survey_id ON USER_SURVEY_GOALS(survey_id);
CREATE INDEX idx_user_survey_goals_goal_id ON USER_SURVEY_GOALS(goal_id);
CREATE INDEX idx_goals_category ON GOALS(category);

-- Note: To add new categories in the future, use:
-- ALTER TYPE goal_category ADD VALUE 'NewCategory';
-- This can be done in new migration files as needed
