# API Tests

This directory contains Postman collections and environments for testing the Smart Fit Girl APIs independently of the React frontend.

## API-First Development Approach

We follow an API-first approach where:
1. Design and implement each API endpoint
2. Test thoroughly with Postman
3. Then build the corresponding React UI components

## Directory Structure

```
api-tests/
├── collections/     # Postman collection files (.json)
├── environments/    # Postman environment files (.json)
├── data/           # Test data files (JSON, CSV, etc.)
└── scripts/        # Custom test scripts and utilities
```

## API Endpoints Structure

### User Service (Port 8082)
- `POST /users/register` - User registration
- `POST /users/login` - User authentication
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `GET /users/settings` - Get user settings
- `PUT /users/settings` - Update user settings

### Meal Service (Port 8083)
- `GET /meals/plans` - Get meal plans
- `POST /meals/plans` - Create meal plan
- `PUT /meals/plans/:id` - Update meal plan
- `DELETE /meals/plans/:id` - Delete meal plan

### Tracking Service (Port 8084)
- `POST /tracking/checkin` - Daily check-in (weight, mood, sleep)
- `GET /tracking/checkins` - Get check-in history
- `POST /tracking/evening` - Evening check-in (steps, activities)
- `GET /tracking/stats` - Get tracking statistics

### API Service (Port 8080)
All endpoints are accessible through the main API service which routes to appropriate microservices.

## Testing Workflow

1. Start backend services: `docker-compose up -d`
2. Import Postman collections from `collections/` directory
3. Set up environment variables in Postman
4. Test each endpoint thoroughly
5. Document any issues or changes needed
6. Once API is stable, implement the corresponding UI component

## Environment Variables for Postman

```json
{
  "api_base_url": "http://localhost:8080",
  "user_service_url": "http://localhost:8082",
  "meal_service_url": "http://localhost:8083",
  "tracking_service_url": "http://localhost:8084"
}
```
