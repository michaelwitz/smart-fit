-- Smart Fit - Complete Database Schema
-- This file contains the comprehensive database structure for the Smart Fit application

-- ENUM Types
CREATE TYPE sex_type AS ENUM (
    'MALE',
    'FEMALE', 
    'OTHER'
);

CREATE TYPE goal_category AS ENUM (
    'Weight',
    'Appearance', 
    'Strength',
    'Endurance'
);

CREATE TYPE food_category_type AS ENUM (
    'MEAT', 'FISH', 'GRAIN', 'VEGETABLE', 'FRUIT', 'DAIRY', 'DAIRY_ALTERNATIVE',
    'FAT', 'NIGHTSHADES', 'OIL', 'SPICE_HERB', 'SWEETENER', 'CONDIMENT', 'SNACK',
    'BEVERAGE', 'LEGUMES', 'NUTS', 'SEEDS', 'OTHER'
);

CREATE TYPE serving_unit_type AS ENUM (
    'GRAMS', 'OUNCES', 'TSP', 'TBSP', 'CUPS', 'PIECES'
);

-- Core Tables

-- Users table - user profiles with international support
CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    sex sex_type,
    phone_number VARCHAR(20),
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    city VARCHAR(100),
    state_province_code VARCHAR(50),
    country_code CHAR(2),
    postal_code VARCHAR(20),
    locale VARCHAR(10),
    timezone VARCHAR(50),
    utc_offset INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Goals table - available fitness goals
CREATE TABLE GOALS (
    id SERIAL PRIMARY KEY,
    category goal_category NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Goals table - direct goal assignment (many-to-many)
CREATE TABLE USER_GOALS (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    goal_id INTEGER NOT NULL REFERENCES GOALS(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, goal_id)
);

-- Food Catalog table - comprehensive food database
CREATE TABLE FOOD_CATALOG (
    id SERIAL PRIMARY KEY,
    food_name VARCHAR(255) NOT NULL,
    category food_category_type NOT NULL,
    serving_units serving_unit_type NOT NULL,
    calories DECIMAL(8,2) NOT NULL,
    protein_grams DECIMAL(6,2) NOT NULL,
    carbs_grams DECIMAL(6,2) NOT NULL,
    fat_grams DECIMAL(6,2) NOT NULL,
    is_non_inflammatory BOOLEAN DEFAULT false,
    is_probiotic BOOLEAN DEFAULT false,
    is_prebiotic BOOLEAN DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Food User Likes table - user food preferences (many-to-many)
CREATE TABLE FOOD_USER_LIKES (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    food_id INTEGER NOT NULL REFERENCES FOOD_CATALOG(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, food_id)
);

-- Meals table - meal definitions with nutritional totals
CREATE TABLE MEALS (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_calories DECIMAL(8,2),
    total_protein DECIMAL(6,2),
    total_carbs DECIMAL(6,2),
    total_fat DECIMAL(6,2),
    prep_time INTEGER, -- in minutes
    prep_instructions TEXT, -- HTML or Markdown formatted cooking instructions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Meals table - user meal consumption tracking
CREATE TABLE USER_MEALS (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    meal_id INTEGER NOT NULL REFERENCES MEALS(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    meal_number INTEGER NOT NULL CHECK (meal_number >= 1 AND meal_number <= 6), -- 1-6 for meal ordering throughout the day
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, meal_id, date, meal_number)
);

-- Meal Ingredients table - meal composition with quantities
CREATE TABLE MEAL_INGREDIENTS (
    id SERIAL PRIMARY KEY,
    meal_id INTEGER NOT NULL REFERENCES MEALS(id) ON DELETE CASCADE,
    food_id INTEGER NOT NULL REFERENCES FOOD_CATALOG(id) ON DELETE CASCADE,
    quantity DECIMAL(8,3) NOT NULL,
    unit serving_unit_type NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_users_email ON USERS(email);
CREATE INDEX idx_users_username ON USERS(username);
CREATE INDEX idx_users_full_name ON USERS(full_name);
CREATE INDEX idx_users_country_code ON USERS(country_code);
CREATE INDEX idx_users_created_at ON USERS(created_at);
CREATE INDEX idx_goals_category ON GOALS(category);
CREATE INDEX idx_user_goals_user_id ON USER_GOALS(user_id);
CREATE INDEX idx_user_goals_goal_id ON USER_GOALS(goal_id);
CREATE INDEX idx_food_catalog_category ON FOOD_CATALOG(category);
CREATE INDEX idx_food_catalog_serving_units ON FOOD_CATALOG(serving_units);
CREATE INDEX idx_food_user_likes_user_id ON FOOD_USER_LIKES(user_id);
CREATE INDEX idx_food_user_likes_food_id ON FOOD_USER_LIKES(food_id);
CREATE INDEX idx_meals_name ON MEALS(name);
CREATE INDEX idx_user_meals_user_id ON USER_MEALS(user_id);
CREATE INDEX idx_user_meals_date ON USER_MEALS(date);
CREATE INDEX idx_user_meals_meal_id ON USER_MEALS(meal_id);
CREATE INDEX idx_meal_ingredients_meal_id ON MEAL_INGREDIENTS(meal_id);
CREATE INDEX idx_meal_ingredients_food_id ON MEAL_INGREDIENTS(food_id);

-- Add comments for documentation
COMMENT ON TABLE USERS IS 'User accounts for authentication and profile management';
COMMENT ON COLUMN USERS.email IS 'Email address used as login identifier (unique, immutable after creation)';
COMMENT ON COLUMN USERS.username IS 'Optional username for display purposes';
COMMENT ON COLUMN USERS.password_hash IS 'Bcrypt hashed password for authentication';
COMMENT ON COLUMN USERS.full_name IS 'Users full name for display purposes';
COMMENT ON COLUMN USERS.sex IS 'Biological sex (MALE, FEMALE, OTHER)';
COMMENT ON COLUMN USERS.phone_number IS 'Optional phone number for contact';
COMMENT ON COLUMN USERS.city IS 'City of residence';
COMMENT ON COLUMN USERS.state_province_code IS 'State or province of residence';
COMMENT ON COLUMN USERS.postal_code IS 'Postal code or ZIP code for user address';
COMMENT ON COLUMN USERS.country_code IS '2-letter country code (ISO 3166-1 alpha-2)';
COMMENT ON COLUMN USERS.locale IS 'Locale for user language settings (e.g., en-US, es-US)';
COMMENT ON COLUMN USERS.timezone IS 'User timezone (IANA timezone format, e.g., America/New_York)';
COMMENT ON COLUMN USERS.utc_offset IS 'UTC offset in hours (e.g., -8 for PST, +5 for EST, 0 for UTC)';

COMMENT ON TABLE GOALS IS 'Available fitness goals organized by categories';
COMMENT ON COLUMN GOALS.category IS 'Goal category (Weight, Appearance, Strength, Endurance)';
COMMENT ON COLUMN GOALS.name IS 'Goal name (max 50 chars for UI display)';
COMMENT ON COLUMN GOALS.description IS 'Goal description (max 200 chars)';

COMMENT ON TABLE USER_GOALS IS 'Junction table linking users to their selected goals';
COMMENT ON COLUMN USER_GOALS.user_id IS 'Foreign key to USERS table';
COMMENT ON COLUMN USER_GOALS.goal_id IS 'Foreign key to GOALS table';

COMMENT ON TABLE FOOD_CATALOG IS 'Comprehensive food database with nutritional information and health properties';
COMMENT ON COLUMN FOOD_CATALOG.category IS 'Food category from enum (MEAT, FISH, GRAIN, etc.)';
COMMENT ON COLUMN FOOD_CATALOG.serving_units IS 'Unit of measurement from enum (GRAMS, OUNCES, etc.)';
COMMENT ON COLUMN FOOD_CATALOG.is_non_inflammatory IS 'Boolean flag indicating anti-inflammatory properties';
COMMENT ON COLUMN FOOD_CATALOG.is_probiotic IS 'Boolean flag indicating probiotic content';
COMMENT ON COLUMN FOOD_CATALOG.is_prebiotic IS 'Boolean flag indicating prebiotic content';

COMMENT ON TABLE FOOD_USER_LIKES IS 'Junction table tracking user food preferences';
COMMENT ON COLUMN FOOD_USER_LIKES.user_id IS 'Foreign key to USERS table';
COMMENT ON COLUMN FOOD_USER_LIKES.food_id IS 'Foreign key to FOOD_CATALOG table';

COMMENT ON TABLE MEALS IS 'Stores meal definitions with nutritional totals and preparation instructions';
COMMENT ON COLUMN MEALS.prep_time IS 'Preparation time in minutes';
COMMENT ON COLUMN MEALS.prep_instructions IS 'HTML or Markdown formatted cooking instructions';

COMMENT ON TABLE USER_MEALS IS 'Tracks user meal consumption by date and meal number';
COMMENT ON COLUMN USER_MEALS.meal_number IS 'Meal order (1-6 for meal ordering throughout the day)';

COMMENT ON TABLE MEAL_INGREDIENTS IS 'Junction table defining meal composition with quantities';
COMMENT ON COLUMN MEAL_INGREDIENTS.quantity IS 'Amount of food item';
COMMENT ON COLUMN MEAL_INGREDIENTS.unit IS 'Unit of measurement from serving_units enum';

-- Key Relationships:
-- 1. USERS 1:N USER_GOALS (users can have multiple goals)
-- 2. USERS 1:N FOOD_USER_LIKES (users can like multiple foods)
-- 3. USERS 1:N USER_MEALS (users can consume multiple meals)
-- 4. MEALS 1:N MEAL_INGREDIENTS (meals can have multiple ingredients)
-- 5. FOOD_CATALOG 1:N MEAL_INGREDIENTS (foods can be used in multiple meals)
-- 6. GOALS 1:N USER_GOALS (goals can be assigned to multiple users)
