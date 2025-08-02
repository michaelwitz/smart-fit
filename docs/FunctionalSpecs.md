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

## User Authentication & Management

### USERS Table
The `USERS` table stores user account information for authentication and profile management.

### Authentication Rules (MVP)
- **Login Identifier**: Email address is the unique login identifier
- **Email Immutability**: Users cannot change their email address in MVP
- **Password Management**: Users can update their password through upsert operations
- **Profile Updates**: All other profile fields (name, location, preferences) can be updated

### User Profile Fields
- `full_name`: User's display name
- `email`: Login identifier (immutable in MVP) 
- `password`: Hashed password for authentication
- `phone_number`: Optional contact information
- `identify_as`: Gender identity or preferred identification
- `city`, `state_province`, `country_code`: Location information
- `timezone`: IANA timezone format (e.g., "America/New_York")
- `utc_offset`: UTC offset in hours (e.g., -8 for PST, +5 for EST)

### Authentication Flow
1. **Registration/Update**: Use upsert endpoint with email as key
2. **Login**: Email + password verification through API service
3. **Session Management**: JWT tokens issued by API service
4. **Profile Management**: Updates allowed for all fields except email

## Business Logic & Guidelines
- **Nutritional Tracking**: Calculate total nutrient intake by multiplying per-unit values by serving size.
- **Goal Management**: Support setting and tracking custom fitness goals, including nutrient intake, weight management, and exercise targets.
- **User Sessions**: Maintain user authentication state through JWT tokens with appropriate expiration.
