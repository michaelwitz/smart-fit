# Functional Specifications - Smart Fit Girl

## Purpose
This document consolidates all business logic, nutritional guidelines, and behavioral specifications for the Smart Fit Girl application.

## FOOD_CATALOG Table

### Description
The `FOOD_CATALOG` table maintains nutritional details for various foods, specifying quantities per unit of measure. This data allows users to easily track nutrient intake.

### Nutritional Information
- **Units**: All nutritional values are per one unit of `serving_units` to facilitate calculation for different serving sizes.
- **Measurement Units**:
  - `grams`
  - `ounces`
  - `fluid_ounces`
  - `cup`
  - `tsp`
  - `tbsp`
  - `scoop`
  - `item`
- **Nutrient Breakdown** (per unit):
  - `calories`
  - `protein_grams`
  - `carbs_grams`
  - `fat_grams`

### Inflammatory Response
- **is_non_inflammatory**: Indicates foods that do not cause or may reduce inflammation.

## Conventions

### Database Naming
- **Table Names**: ALL_CAPS
- **Column Names**: snake_case
- **Boolean Columns**: Prefixed with `is_`

## Database Schema

### USERS Table
Stores user profiles with international support:

- **id**: Primary key (auto-increment)
- **full_name**: User's display name (required)
- **email**: Login identifier (unique, required, immutable after creation)
- **password**: Bcrypt hashed password (required)
- **phone_number**: Optional contact number
- **identify_as**: Preferred pronouns/identity
- **city**: User's city
- **state_province**: State or province
- **postal_code**: Postal/ZIP code for international support
- **country_code**: ISO 2-letter country code
- **locale**: Language/locale preference (e.g., "en-US", "es-ES")
- **timezone**: IANA timezone identifier (e.g., "America/New_York")
- **utc_offset**: UTC offset in hours for convenience
- **created_at**: Account creation timestamp
- **updated_at**: Last profile update timestamp

### GOALS Table
Stores available fitness goals organized by categories:

- **id**: Primary key (auto-increment)
- **category**: ENUM ('Weight', 'Appearance', 'Strength', 'Endurance')
- **name**: Goal name (max 50 chars for UI display)
- **description**: Goal description (max 200 chars)
- **created_at**: Goal creation timestamp
- **updated_at**: Last goal update timestamp

### SURVEYS Table
Stores immutable user fitness surveys over time:

- **id**: Primary key (auto-increment)
- **user_id**: Foreign key to USERS table
- **current_weight**: User's current weight (DECIMAL 5,2)
- **target_weight**: User's target weight (DECIMAL 5,2)
- **activity_level**: Activity level scale 0-10 (INTEGER with CHECK constraint)
- **created_at**: Survey completion timestamp (immutable)

### USER_SURVEY_GOALS Junction Table
Links surveys to selected goals (many-to-many relationship):

- **id**: Primary key (auto-increment)
- **survey_id**: Foreign key to SURVEYS table
- **goal_id**: Foreign key to GOALS table
- **created_at**: Goal selection timestamp
- **UNIQUE constraint**: (survey_id, goal_id) - prevents duplicate goals per survey

**Business Rule**: Only one goal per category per survey (enforced at application level)

## User Authentication & Management

### Authentication Rules (MVP)
- **Login Identifier**: Email address is the unique login identifier
- **Email Immutability**: Users cannot change their email address in MVP
- **Password Management**: Users can update their password through upsert operations
- **Profile Updates**: All other profile fields (name, location, preferences) can be updated

### Authentication Flow
1. **Registration/Update**: Use upsert endpoint with email as key
2. **Login**: Email + password verification through API service
3. **Session Management**: JWT tokens issued by API service
4. **Profile Management**: Updates allowed for all fields except email

## Survey & Goals System

### Survey Workflow
1. **Survey Creation**: Users complete fitness surveys with current/target weight, activity level, and goal selections
2. **Goal Selection**: Users select multiple goals from predefined categories (one per category recommended)
3. **Data Immutability**: Surveys are immutable once created to maintain historical tracking
4. **Progress Tracking**: Multiple surveys over time allow progress analysis

### Goal Categories
- **Weight**: Goals related to weight management (loss, gain, maintenance)
- **Appearance**: Goals focused on physical appearance and body composition
- **Strength**: Goals targeting strength building and muscle development
- **Endurance**: Goals for cardiovascular fitness and stamina improvement

### Business Rules
- Surveys cannot be modified after creation (immutable for data integrity)
- Users can select multiple goals but only one per category per survey
- Goal associations are permanent once a survey is submitted
- Weight values support decimal precision for accurate tracking

## Business Logic & Guidelines
- **Nutritional Tracking**: Calculate total nutrient intake by multiplying per-unit values by serving size.
- **Goal Management**: Support setting and tracking custom fitness goals, including nutrient intake, weight management, and exercise targets.
- **User Sessions**: Maintain user authentication state through JWT tokens with appropriate expiration.
