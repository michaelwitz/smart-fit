# Ground Rules for Smart Fit Project

## Technologies Used

- **API Service**: Go + Gin framework (authentication, routing, middleware)
- **Microservices**: Go + gRPC (business logic only)
- **Database**: PostgreSQL + sqlx
- **Frontend**: React.js + Mantine UI
- **Containerization**: Docker + Docker Compose
- **Communication**: HTTP/REST APIs
- **Fault Tolerance**: failsafe-go library
- **Analytics**: PostHog (client-side only)

## Architecture Principles

- **API Service (Gin)**: Centralized authentication, JWT handling, middleware, client routing
- **Microservices (stdlib http)**: Pure business logic, no authentication, direct database access
- **Gateway Pattern**: API service acts as gateway, microservices remain simple and testable
- **API-First Development**: Design and test APIs before building UI components

## Database Conventions

- **Table Names**: ALL_CAPS (e.g., FOOD_CATALOG, USER_PROFILES)
- **Column Names**: snake_case (e.g., food_name, created_at, nutritional_info)
- **Boolean Columns**: Prefixed with `is_` (e.g., is_non_inflammatory, is_active)
- **Go Struct Mapping**: map snake_case columns to camelCase struct fields in the SQL statements
  - `food_name` → `FoodName`
  - `created_at` → `CreatedAt`
  - `nutritional_info` → `NutritionalInfo`

## JSON API Conventions

- **JSON Field Names**: camelCase (e.g., fullName, createdAt, phoneNumber)
- **Go Struct Fields**: PascalCase (e.g., FullName, CreatedAt, PhoneNumber)
- **JSON Tags**: Map Go PascalCase to JSON camelCase
  - `FullName string \`json:"fullName"\``
  - `CreatedAt time.Time \`json:"createdAt"\``
  - `PhoneNumber *string \`json:"phoneNumber"\``

## Database Modification Procedures

When adding columns or changing table structure:

1. **Drop constraints**: Remove foreign keys and other constraints
2. **Truncate table**: Clear existing data
3. **Apply schema changes**: Add/modify columns
4. **Reload seed data**: Insert updated sample data
5. **Restore constraints**: Add back foreign keys and constraints

## Database Schema Management Philosophy

**Development Phase (Current)**:

- **No migration files** - We recreate and reseed the database on schema changes
- **Single source of truth**: `database/schemas/complete_database_schema.sql` contains the full schema
- **Seed data**: `database/seeds/` contains all initial data
- **Benefits**: Simpler development, no schema drift, consistent state across environments

**Production Phase (Future)**:

- Migration files will be introduced when we need to preserve production data
- Schema evolution tracking will become necessary
- Rollback capabilities will be implemented

**Current Setup Process**:

1. Drop and recreate database container
2. Run complete schema file
3. Load all seed data
4. Verify schema matches application expectations

## Testing Approach

- **Unit Tests**:
  - Use `stretchr/testify` for assertions
  - Format with `gotestfmt`
  - Follow table-driven test patterns
  - Test business logic in isolation
- **Integration Tests**:
  - Utilize Postman for API testing
  - Test microservices independently (no auth required)
  - Test complete flows through API service

## Development Practices

- **Environment Management**: Use `.env` files for configuration and secrets
- **Version Control**: Commit early, commit often
- **Documentation**: Keep README and inline comments up to date

## Development Workflow

1. Design and implement APIs first
2. Develop frontend using stable APIs
3. Focus on modular and reusable code
4. Ensure tests are written and passing before merging into main branch

## Communication

- Regular updates in team meetings every Monday
- Use Slack for daily communication and GitHub issues for tracking work
