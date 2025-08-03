# Smart Fit Girl - Fitness Tracking App

A comprehensive fitness tracking and advice application built with a microservices architecture.

## Architecture Overview

### Backend Services
- **api-service**: Main API gateway using Go + Gin (authentication, JWT, middleware, client routing)
- **user-service**: User registration, profile management, and user data operations (Go + stdlib http)
- **meal-service**: Meal planning, nutrition tracking, and meal data operations (Go + stdlib http)
- **tracking-service**: Daily check-ins, weight, mood, sleep, activity tracking (Go + stdlib http)

### Frontend
- **web-client**: React.js + Mantine UI library for mobile-first responsive design
  - **Mobile Detection**: `react-device-detect` for reliable device type detection
  - **Responsive Design**: Mobile-first CSS with `100dvh` viewport units and CSS safe-area insets
  - **Progressive Web App**: iOS Safari meta tags, manifest.json, and fullscreen support
  - **Touch Optimization**: Enhanced mobile scrolling and address bar hiding for iOS Safari

### Database
- PostgreSQL running in Docker container

### Analytics
- PostHog integration for usage analytics and user behavior tracking

## Tech Stack

- **API Gateway**: Go 1.23 + Gin framework (authentication, middleware)
- **Microservices**: Go 1.23 + standard http package (business logic)
- **Database**: PostgreSQL + SQLBoiler ORM
- **Frontend**: React.js + Mantine UI + react-device-detect for mobile-first responsive design
- **Containerization**: Docker + Docker Compose
  - **Build Images**: `golang:1.23-bookworm` (Debian-based for compatibility)
  - **Runtime Images**: `gcr.io/distroless/static-debian12:nonroot` (security-focused, minimal attack surface)
- **Fault Tolerance**: failsafe-go library
- **Communication**: HTTP/REST APIs
- **Documentation**: Swagger/OpenAPI 3.0 with interactive UI
- **Analytics**: PostHog (client-side only)

## 🛡️ Security Implementation

This application implements multiple layers of security following industry best practices:

### **Container Security**
- ✅ **Distroless Runtime Images**: Using `gcr.io/distroless/static-debian12:nonroot`
  - No shell, package managers, or unnecessary binaries
  - Minimal attack surface with only essential runtime components
  - Non-root user execution by default
- ✅ **Multi-Stage Docker Builds**: Build with full toolchain, deploy with minimal runtime
- ✅ **Latest Go Version**: Go 1.23 with latest security patches and improvements

### **Authentication & Authorization** 
- ✅ **JWT Tokens**: Secure, stateless authentication with configurable expiration
- ✅ **bcrypt Password Hashing**: Industry-standard password encryption (cost factor 10+)
- ✅ **Bearer Token Authentication**: Standard HTTP Authorization header implementation
- ✅ **Protected Routes**: Middleware-based route protection with token validation

### **Database Security**
- ✅ **Parameterized Queries**: All SQL queries use parameters to prevent injection attacks
- ✅ **Connection Pooling**: Secure database connection management
- ✅ **Environment-Based Credentials**: Database credentials via environment variables
- ✅ **PostgreSQL**: Production-grade database with built-in security features

### **Application Security**
- ✅ **Input Validation**: Request payload validation using Go struct tags
- ✅ **CORS Configuration**: Controlled cross-origin resource sharing
- ✅ **Error Handling**: Sanitized error responses (no sensitive data leakage)
- ✅ **Structured Logging**: Secure logging without credential exposure

### **Infrastructure Security**
- ✅ **Environment Variable Management**: All secrets via environment variables
- ✅ **Docker Network Isolation**: Services communicate via isolated Docker network
- ✅ **Non-Root Container Execution**: All containers run as unprivileged users
- ✅ **Minimal Dependencies**: Reduced attack surface with essential dependencies only

### **Development Security**
- ✅ **Gitignored Secrets**: All `.env` files excluded from version control
- ✅ **Example Templates**: `.env.example` files with placeholder values
- ✅ **Secret Rotation Ready**: Environment-based configuration supports easy rotation
- ✅ **Development vs Production**: Clear separation of development and production configs

### **API Security**
- ✅ **Rate Limiting Ready**: Architecture supports rate limiting implementation
- ✅ **HTTPS Ready**: TLS termination at load balancer level
- ✅ **API Documentation Security**: Swagger UI with authentication integration
- ✅ **Microservices Isolation**: Each service handles specific business logic only

> **Production Recommendation**: Use cloud-native secret management (AWS Secrets Manager, Azure Key Vault, etc.) and implement additional security headers, rate limiting, and monitoring.

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
- **Web Client**:       `http://localhost:5050`
- **API Documentation (Swagger)**: `http://localhost:8080/swagger/index.html` 📖

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
   
   # SendGrid Email Service (for password reset functionality)
   SENDGRID_API_KEY=your-sendgrid-api-key
   
   # PostHog Analytics (for React client)
   REACT_APP_POSTHOG_KEY=your-posthog-api-key
   REACT_APP_POSTHOG_HOST=https://app.posthog.com
   ```

### Environment Files Structure
- **`.env`** - Main environment file (gitignored)
- **`.env.example`** - Template with placeholder values
- **`services/*//.env.example`** - Service-specific templates

### SendGrid Email Service Setup

The application uses SendGrid for sending password reset emails. To configure:

1. **Create a SendGrid account** at [sendgrid.com](https://sendgrid.com)

2. **Generate an API Key**:
   - Go to Settings → API Keys in your SendGrid dashboard
   - Click "Create API Key"
   - Choose "Restricted Access" for security
   - Grant the following permissions:
     - Mail Send: Full Access
   - Copy the generated API key

3. **Update your `.env` file**:
   ```bash
   SENDGRID_API_KEY=SG.your-actual-api-key-here
   ```

4. **Verify Sender Identity** (Required for production):
   - Go to Settings → Sender Authentication
   - Either verify a single sender email or authenticate your domain
   - For development, you can use the single sender verification

5. **Update Email Templates** (Optional):
   - The current implementation uses basic HTML templates
   - You can customize the email content in `services/user-service/handlers.go`
   - For production, consider using SendGrid's Dynamic Templates

**Testing Password Reset**:
```bash
# Request password reset
curl -X POST http://localhost:8082/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# Reset password with token (check your email for the token)
curl -X POST http://localhost:8082/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token":"your-reset-token-from-email",
    "newPassword":"newpassword123"
  }'
```

### Security Notes
- ✅ Environment files are gitignored for security
- ✅ Never commit actual credentials to version control
- ✅ Use cloud secrets management in production
- ✅ PostHog is only configured client-side for simplicity
- ✅ SendGrid API keys should be rotated regularly
- ✅ Password reset tokens expire after 1 hour for security

## API Testing (API-First Development)

We follow an API-first approach where APIs are designed and tested before building UI components.

### 📖 Interactive API Documentation (Swagger)

The API service includes comprehensive Swagger documentation with interactive testing:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/doc.json`
- **Features**: 
  - Interactive endpoint testing
  - JWT Bearer token authentication
  - Request/response examples
  - OpenAPI 3.0 specification

**Quick Test with Swagger:**
1. Start services: `make up`
2. Open Swagger UI: `http://localhost:8080/swagger/index.html`
3. Test `/auth/login` to get JWT token
4. Click "Authorize" and enter `Bearer <your-token>`
5. Test protected endpoints

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

## Developer Setup and Notes

### Local Development Ports (Docker Desktop)
When running locally with `docker-compose`, services are accessible on these ports:

| Service | Container Port | Host Port | Local URL |
|---------|----------------|-----------|----------|
| **API Service** (Gateway) | 8080 | 8080 | http://localhost:8080 |
| **User Service** | 8080 | 8082 | http://localhost:8082 |
| **Meal Service** | 8080 | 8083 | http://localhost:8083 |
| **Tracking Service** | 8080 | 8084 | http://localhost:8084 |
| **Web Client** (React) | 3000 | 5050 | http://localhost:5050 |
| **PostgreSQL** | 5432 | 5432 | localhost:5432 |

> **Note**: These ports are for local development only. Production/cloud deployments will use different port configurations.

### Development Workflow
1. **Start services**: `make up` or `docker-compose up -d`
2. **View logs**: `make logs` or `docker-compose logs -f [service-name]`
3. **Stop services**: `make down` or `docker-compose down`
4. **Rebuild after code changes**: `docker-compose build [service-name]`

### Testing Individual Services
During development, you can test services directly:

#### Health Check Endpoints
```bash
curl http://localhost:8080/health  # API Service
curl http://localhost:8082/health  # User Service
curl http://localhost:8083/health  # Meal Service
curl http://localhost:8084/health  # Tracking Service
```

#### User Service Examples
```bash
# Get all users
curl http://localhost:8082/users

# Create a user with full profile
curl -X POST http://localhost:8082/users \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Jane Smith",
    "email": "jane@example.com",
    "password": "password123",
    "phoneNumber": "+1234567890",
    "identifyAs": "She/Her",
    "city": "New York",
    "stateProvince": "NY",
    "postalCode": "10001",
    "countryCode": "US",
    "locale": "en-US",
    "timezone": "America/New_York",
    "utcOffset": -5
  }'

# Get user by ID
curl http://localhost:8082/users/1

# Update user (partial update)
curl -X PUT http://localhost:8082/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "city": "San Francisco",
    "stateProvince": "CA",
    "postalCode": "94105",
    "timezone": "America/Los_Angeles",
    "utcOffset": -8
  }'

# Upsert user (create or update by email)
curl -X POST http://localhost:8082/users/upsert \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "John Doe Updated",
    "email": "john@example.com",
    "password": "newpassword123",
    "city": "Chicago",
    "locale": "es-US"
  }'

# Verify user credentials
curl -X POST http://localhost:8082/users/verify \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"password123"}'

# Get all available goals (grouped by category)
curl http://localhost:8082/goals

# Create a fitness survey for user (includes goals selection)
curl -X POST http://localhost:8082/users/1/survey \
  -H "Content-Type: application/json" \
  -d '{
    "currentWeight": 165.5,
    "targetWeight": 150.0,
    "activityLevel": 4,
    "goalIds": [1, 5, 7, 9]
  }'

# Get user's latest survey with goals
curl http://localhost:8082/users/1/survey/latest

# Test user Sophia Woytowitz (from seed data)
curl -X POST http://localhost:8082/users/1/survey \
  -H "Content-Type: application/json" \
  -d '{
    "currentWeight": 165.5,
    "targetWeight": 150.0,
    "activityLevel": 4,
    "goalIds": [1, 5, 7, 9]
  }'

# Delete user
curl -X DELETE http://localhost:8082/users/1
```

#### API Service Authentication
```bash
# Login (get JWT token)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"password123"}'

# Example response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": 1,
#     "fullName": "Jane Smith",
#     "email": "jane@example.com",
#     "locale": "en-US",
#     "timezone": "America/New_York"
#   }
# }

# Use JWT token for protected endpoints
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
curl -H "Authorization: Bearer $JWT_TOKEN" \
  http://localhost:8080/api/users/profile
```

**💡 Tip**: Use the interactive [Swagger UI](http://localhost:8080/swagger/index.html) for easier API testing with built-in authentication!

#### Meal Service Examples (when implemented)
```bash
# Get all meals
curl http://localhost:8083/meals

# Create a meal plan
curl -X POST http://localhost:8083/meals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "name": "Healthy Breakfast",
    "calories": 350,
    "protein": 20,
    "carbs": 45,
    "fat": 12
  }'
```

#### Tracking Service Examples (when implemented)
```bash
# Daily check-in
curl -X POST http://localhost:8084/tracking/checkin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "weight": 150.5,
    "mood": 8,
    "sleep_hours": 7.5,
    "energy_level": 7
  }'

# Get tracking history
curl -H "Authorization: Bearer $JWT_TOKEN" \
  http://localhost:8084/tracking/history?days=30
```

### Database Access
```bash
# Connect to PostgreSQL (password: smartfit123)
psql -h localhost -p 5432 -U smartfit -d smartfitgirl

# View users table
SELECT * FROM USERS;

# Create password reset table (if not already created)
\i services/user-service/password_resets_table.sql

# View password reset tokens
SELECT * FROM PASSWORD_RESETS;
```

### Service Architecture Notes
- **API Service**: Acts as gateway, handles authentication, routes requests
- **Microservices**: Handle specific business logic (users, meals, tracking)
- **Internal Communication**: Services communicate via Docker network using service names
- **External Access**: Use localhost ports for development/testing
- **Database**: Shared PostgreSQL instance across all services
- **Container Security**: 
  - Multi-stage builds with Debian for compatibility during build
  - Distroless runtime images for minimal attack surface
  - Non-root user execution for enhanced security
  - No shell or package managers in production containers

### Key Features Implementation
- **Authentication**: JWT tokens with bcrypt password hashing
- **International Support**: User profiles include `locale`, `postal_code`, `timezone`, `utc_offset`
- **Security**: 
  - Parameterized SQL queries prevent injection attacks
  - Distroless container images for minimal attack surface
  - Non-root container execution
  - JWT-based authentication with secure token handling
- **CRUD Operations**: Complete user management with upsert functionality
- **Service Discovery**: Docker Compose networking for inter-service communication
- **API Documentation**: Interactive Swagger UI with OpenAPI 3.0 specification

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
