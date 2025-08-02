# Smart Fit Girl - Fitness Tracking App

A comprehensive fitness tracking and advice application built with a microservices architecture.

## Architecture Overview

### Backend Services (Go + Gin)
- **api-service**: Main API service handling all client requests (web + future mobile)
- **user-service**: User authentication, registration, profile management, and user data operations
- **meal-service**: Meal planning, nutrition tracking, and meal data operations
- **tracking-service**: Daily check-ins, weight, mood, sleep, activity tracking, and tracking data operations

### Frontend
- **web-client**: React.js + Mantine UI library for mobile-first responsive design

### Database
- PostgreSQL running in Docker container

### Analytics
- PostHog integration for usage analytics and user behavior tracking

## Tech Stack

- **Backend**: Go + Gin framework
- **Database**: PostgreSQL + SQLBoiler ORM
- **Frontend**: React.js + Mantine UI
- **Containerization**: Docker + Docker Compose
- **Fault Tolerance**: failsafe-go library
- **Communication**: HTTP/REST APIs
- **Analytics**: PostHog

## Quick Start

### Development Setup
1. **Run the setup script**
   ```bash
   ./scripts/setup-dev.sh
   ```
2. **Update your `.env` file** with your actual credentials and PostHog API key.

### Running the Application
- **Start all services**
  ```bash
  make up
  ```
- **Stop all services**
  ```bash
  make down
  ```
- **Check logs**
  ```bash
  make logs
  ```

### Test the APIs with Postman
- Use the base URL: `http://localhost:8080`
- Ensure all services are running via `docker-compose ps`

### Service Endpoints
- **API Service**:      `http://localhost:8080`
- **User Service**:     `http://localhost:8082`
- **Meal Service**:     `http://localhost:8083`
- **Tracking Service**: `http://localhost:8084`
- **Web Client**:       `http://localhost:3000`

### Need Help?
For any issues, check the service logs or consult the `/docs` directory for more detailed documentation.

## Environment Variables

### Required Setup
Before running the application, you must configure your environment variables:

1. **Run the setup script** (creates `.env` files from templates):
   ```bash
   ./scripts/setup-dev.sh
   ```

2. **Update the main `.env` file** with your values:
   ```bash
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=smartfitgirl
   DB_USER=smartfit
   DB_PASSWORD=your-secure-password
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key
   
   # PostHog Analytics (for React client)
   REACT_APP_POSTHOG_KEY=your-posthog-api-key
   REACT_APP_POSTHOG_HOST=https://app.posthog.com
   ```

### Environment Files Structure
- **`.env`** - Main environment file (gitignored)
- **`.env.example`** - Template with placeholder values
- **`services/*//.env.example`** - Service-specific templates

### Security Notes
- ✅ Environment files are gitignored for security
- ✅ Never commit actual credentials to version control
- ✅ Use cloud secrets management in production
- ✅ PostHog is only configured client-side for simplicity

## API Testing (API-First Development)

We follow an API-first approach where APIs are designed and tested before building UI components.

### Postman Collections
- Collections are stored in `/api-tests/collections/`
- Environment files in `/api-tests/environments/`
- Test data in `/api-tests/data/`

### Testing Workflow
1. **Start backend services**: `make up`
2. **Import Postman collections** from `api-tests/collections/`
3. **Configure Postman environment** with:
   ```json
   {
     "api_base_url": "http://localhost:8080",
     "user_service_url": "http://localhost:8082",
     "meal_service_url": "http://localhost:8083",
     "tracking_service_url": "http://localhost:8084"
   }
   ```
4. **Test each endpoint** thoroughly
5. **Document any issues** or required changes
6. **Once API is stable**, implement corresponding React UI components

### Direct Service Testing
You can also test individual services directly:
- **User Service**: `http://localhost:8082/users/health`
- **Meal Service**: `http://localhost:8083/meals/health`
- **Tracking Service**: `http://localhost:8084/tracking/health`

## Development

See individual service READMEs for specific development instructions:
- [API Service](./services/api-service/README.md)
- [User Service](./services/user-service/README.md)
- [Meal Service](./services/meal-service/README.md)
- [Tracking Service](./services/tracking-service/README.md)
- [Web Client](./web-client/README.md)

## Implementation Plan

### Phase 1: Foundation ⏳
- [x] Set up monorepo structure
- [ ] Configure Docker and Docker Compose
- [ ] Implement basic Web API with health checks
- [ ] Set up PostgreSQL with initial schema
- [ ] Create basic Data Service with database connection
- [ ] Configure failsafe-go for service communication

### Phase 2: User Management
- [ ] User registration and authentication
- [ ] User profile management
- [ ] Settings and preferences system (system + personal fitness goals)
- [ ] Basic security implementation
- [ ] PostHog integration for user analytics

### Phase 3: Core Tracking
- [ ] Daily check-in functionality (weight, mood, sleep)
- [ ] Data storage and retrieval
- [ ] Basic web client with check-in forms
- [ ] Track user engagement with PostHog

### Phase 4: Meal Planning
- [ ] Meal planning service
- [ ] Weekly meal plan management
- [ ] Meal planning UI components
- [ ] Analytics on meal planning usage

### Phase 5: Analytics & Evening Check-ins
- [ ] Evening check-ins (steps, activities TBD)
- [ ] Data visualization (charts and graphs for historical data)
- [ ] PostHog dashboards for app usage analytics
- [ ] UI/UX improvements based on analytics insights

### Phase 6: Testing & Polish
- [ ] Comprehensive testing across all services
- [ ] Performance optimization
- [ ] Bug fixes and final polish
- [ ] Deployment preparation
