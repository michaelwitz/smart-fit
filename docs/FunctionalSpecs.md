# Functional Specifications - Smart Fit Girl

## Purpose

This document consolidates all business logic, nutritional guidelines, and behavioral specifications for the Smart Fit Girl application.

## Table of Contents

1. [Database Schema & Conventions](#database-schema--conventions)
2. [User Account Management](#user-account-management)
3. [User Fitness & Goals System](#user-fitness--goals-system)
4. [User Activity Tracking](#user-activity-tracking)
5. [Nutritional Tracking System](#nutritional-tracking-system)
6. [Backend Automation & Analytics](#backend-automation--analytics)
7. [Business Logic & Guidelines](#business-logic--guidelines)

---

## Database Schema & Conventions

### Database Naming Conventions

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
- **sex**: Biological sex (MALE, FEMALE, OTHER)

- **city**: User's city
- **state_province**: State or province
- **postal_code**: Postal/ZIP code for international support
- **country_code**: ISO 2-letter country code
- **locale**: Language/locale preference (e.g., "en-US", "es-ES")
- **timezone**: IANA timezone identifier (e.g., "America/New_York")
- **utc_offset**: UTC offset in hours for convenience
- **last_active**: Last activity timestamp with 1-hour debouncing (nullable)
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

## User Account Management

### Account Creation & Registration

- **Registration Method**: Users create accounts by providing email, password, full name, and optional profile information
- **Email Uniqueness**: Each email address can only be associated with one account
- **Password Requirements**: Passwords are hashed using bcrypt before storage
- **Profile Setup**: Users can provide optional demographic and location information during registration
- **Account Activation**: Accounts are immediately active upon creation (no email verification required for MVP)

### Authentication & Login

- **Login Identifier**: Email address is the unique login identifier
- **Password Verification**: Bcrypt hash comparison for secure authentication
- **Session Management**: JWT tokens issued upon successful authentication
- **Token Expiration**: JWT tokens expire after 1 week for user convenience while maintaining security
- **Login Tracking**: Successful logins update the user's `last_active` timestamp

### Password Management

- **Password Updates**: Users can change their passwords through profile update operations
- **Password Reset Flow**: Secure password reset via email token system
  1. **Request Reset**: User provides email address to request password reset
  2. **Token Generation**: System generates secure random token with 1-hour expiration
  3. **Email Delivery**: Reset token sent to user's email address via SendGrid
  4. **Token Validation**: User provides token and new password to complete reset
  5. **Password Update**: New password is hashed and stored, invalidating the reset token
- **Token Security**: Reset tokens are single-use and automatically expire after 1 hour
- **Email Integration**: Password reset emails sent through SendGrid service

### Profile Management

- **Editable Fields**: All profile fields except email can be updated
- **Email Immutability**: Email addresses cannot be changed in MVP (business rule)
- **International Support**: Profile supports international addresses, locales, and timezones
- **Optional Fields**: Most profile fields are optional except email, password, and full name
- **Update Tracking**: Profile updates modify the `updated_at` timestamp

### Account Security

- **Password Hashing**: All passwords stored as bcrypt hashes (never plain text)
- **Token-Based Authentication**: JWT tokens for stateless session management (1-week expiration)
- **Secure Password Reset**: Time-limited, single-use tokens for password recovery
- **Activity Logging**: Login activities tracked for analytics (with privacy considerations)

## User Fitness & Goals System

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

## User Activity Tracking

### Activity Monitoring

- **Purpose**: Track user engagement and identify active users for analytics and engagement strategies
- **Implementation**: `last_active` timestamp field in USERS table with 1-hour debouncing
- **Debounce Logic**: Activity updates only occur if the last recorded activity is more than 1 hour old or null
- **Triggers**: Activity is tracked on:
  - Successful user authentication (login)
  - Survey creation and completion
  - Other authenticated user interactions with the API

### Activity Status Categories

- **Online**: Active within the last 1 hour
- **Active Today**: Active within the last 24 hours
- **Active This Week**: Active within the last 7 days
- **Inactive**: No activity beyond 7 days or never active

### Database Implementation

- **Column**: `last_active TIMESTAMP` (nullable)
- **Index**: Created on `last_active` for efficient querying
- **View**: `active_users` view provides categorized activity status
- **Privacy**: Activity tracking respects user privacy - only timestamps, no detailed action logging

## Nutritional Tracking System

### FOOD_CATALOG Table

Stores nutritional details for various foods, specifying quantities per unit of measure:

- **id**: Primary key (auto-increment)
- **name**: Food item name (required)
- **serving_units**: Unit of measurement for nutritional values
- **calories**: Calories per unit
- **protein_grams**: Protein content per unit
- **carbs_grams**: Carbohydrate content per unit
- **fat_grams**: Fat content per unit
- **is_non_inflammatory**: Boolean flag indicating anti-inflammatory properties
- **created_at**: Food catalog entry creation timestamp
- **updated_at**: Last update timestamp

### Nutritional Information Standards

- **Units**: All nutritional values are per one unit of `serving_units` to facilitate calculation for different serving sizes
- **Supported Measurement Units**:
  - `grams`
  - `ounces`
  - `fluid_ounces`
  - `cup`
  - `tsp` (teaspoon)
  - `tbsp` (tablespoon)
  - `scoop`
  - `item` (for discrete items like eggs, apples, etc.)
- **Nutrient Breakdown**: Per unit values for calories, protein, carbohydrates, and fat
- **Inflammatory Response**: `is_non_inflammatory` flag identifies foods that may reduce inflammation

### Nutritional Calculation Logic

- **Serving Size Calculation**: Total nutrients = per-unit values × serving size
- **Daily Intake tracking**: Sum nutritional values across all consumed food items
- **Goal Alignment**: Compare daily intake against user-defined nutritional goals

## Backend Automation & Analytics

### User Engagement Analytics

- **Activity Monitoring**: Automated tracking of user login patterns and app usage
- **Engagement Metrics**: Calculate user retention, session frequency, and feature adoption
- **Activity Status Reports**: Generate reports on active vs. inactive user segments
- **Trend Analysis**: Identify patterns in user behavior and engagement over time

### Automated User Management

- **Password Reset Cleanup**: Automated removal of expired password reset tokens
- **Session Management**: Automatic cleanup of expired JWT tokens and sessions
- **Account Maintenance**: Identify dormant accounts for potential re-engagement campaigns

### Data Management Automation

- **Survey Data Analysis**: Automated processing of survey responses for trend identification
- **Goal Progress Tracking**: Calculate progress metrics based on survey history
- **Nutritional Data Maintenance**: Automated validation and cleanup of food catalog entries

### Email Automation

- **Password Reset Emails**: Automated delivery of password reset instructions via SendGrid
- **Engagement Emails**: Future implementation for user re-engagement campaigns
- **System Notifications**: Automated alerts for system maintenance and updates

### Analytics Views & Reports

- **Active Users View**: Pre-built database view for user activity status categorization
- **User Engagement Dashboard**: Metrics on daily, weekly, and monthly active users
- **Goal Achievement Analytics**: Track completion rates and popular goal combinations
- **System Health Monitoring**: Automated monitoring of database performance and API response times

## Business Logic & Guidelines

- **Nutritional Tracking**: Calculate total nutrient intake by multiplying per-unit values by serving size.
- **Goal Management**: Support setting and tracking custom fitness goals, including nutrient intake, weight management, and exercise targets.
- **User Sessions**: Maintain user authentication state through JWT tokens with appropriate expiration.
- **Activity Analytics**: Use debounced activity tracking for user engagement analysis while minimizing database writes.
