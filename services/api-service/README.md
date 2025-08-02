# Smart Fit Girl API Service

This is the main API gateway service for the Smart Fit Girl application. It handles authentication, routing, and proxying requests to various microservices.

## 🚀 Features

- **JWT Authentication**: Secure token-based authentication
- **API Gateway**: Routes requests to appropriate microservices
- **Swagger Documentation**: Auto-generated OpenAPI 3.0 documentation
- **Health Checks**: Service health monitoring
- **CORS Support**: Cross-origin resource sharing for web clients
- **Middleware**: Authentication, logging, and error handling

## 📚 API Documentation

The API documentation is automatically generated using Swaggo and is available at:

**Swagger UI**: `http://localhost:8080/swagger/index.html`

### Available Endpoints

#### Health Check
- **GET** `/health` - Check if the API service is healthy

#### Authentication
- **POST** `/auth/login` - User login with email/password

#### Protected Routes
- **GET** `/api/protected` - Example protected endpoint (requires JWT)

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Swaggo CLI tool
- User service running (for authentication)

### Setup

1. **Install dependencies**:
   ```bash
   make deps
   ```

2. **Generate Swagger docs**:
   ```bash
   make docs
   ```

3. **Run the service**:
   ```bash
   make run
   # OR
   go run .
   ```

4. **Access Swagger UI**:
   Open `http://localhost:8080/swagger/index.html` in your browser

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVICE_PORT` | `8080` | Port for the API service |
| `USER_SERVICE_URL` | `http://user-service:8080` | URL of the user service |
| `JWT_SECRET` | `your-super-secret-jwt-key-change-in-production` | JWT signing secret |

### Make Commands

```bash
make help          # Show available commands
make build         # Build the application  
make run           # Run the application
make docs          # Generate Swagger documentation
make test          # Run tests
make fmt           # Format code
make clean         # Clean build artifacts
```

## 📖 Swagger Documentation Generation

This service uses [Swaggo](https://github.com/swaggo/swag) to automatically generate OpenAPI/Swagger documentation from Go annotations.

### Adding Documentation to New Endpoints

1. **Add Swagger comments** above your handler function:
   ```go
   // getUserProfile godoc
   // @Summary      Get User Profile
   // @Description  Get the authenticated user's profile information
   // @Tags         users
   // @Accept       json
   // @Produce      json
   // @Security     Bearer
   // @Success      200  {object}  UserProfile
   // @Failure      401  {object}  ErrorResponse
   // @Failure      404  {object}  ErrorResponse
   // @Router       /api/users/profile [get]
   func getUserProfile(c *gin.Context) {
       // Handler implementation
   }
   ```

2. **Define response structures** with JSON tags and examples:
   ```go
   type UserProfile struct {
       ID       int    `json:\"id\" example:\"1\"`
       Email    string `json:\"email\" example:\"user@example.com\"`
       Name     string `json:\"name\" example:\"John Doe\"`
   }
   ```

3. **Regenerate documentation**:
   ```bash
   make docs
   ```

### Swagger Annotations Reference

- `@Summary` - Brief description of the endpoint
- `@Description` - Detailed description
- `@Tags` - Groups endpoints in the UI
- `@Accept` - Content types the endpoint accepts
- `@Produce` - Content types the endpoint produces
- `@Param` - Request parameters (query, path, body, header)
- `@Success` - Success response definition
- `@Failure` - Error response definition
- `@Security` - Security requirements (Bearer token)
- `@Router` - URL path and HTTP method

### Security Configuration

JWT Bearer token authentication is configured in the main package comment:

```go
// SecurityDefinitions:
// Bearer:
//   type: apiKey
//   name: Authorization
//   in: header
//   description: \"JWT Authorization header using the Bearer scheme. Example: 'Bearer {token}'\"
```

## 🔐 Authentication Flow

1. **Login**: POST `/auth/login` with email/password
2. **Get JWT**: Receive JWT token in response
3. **Use Token**: Include `Authorization: Bearer <token>` header in subsequent requests
4. **Access Protected**: Access protected endpoints with valid JWT

### Example Login Request

```bash
curl -X POST http://localhost:8080/auth/login \\
  -H \"Content-Type: application/json\" \\
  -d '{
    \"email\": \"user@example.com\",
    \"password\": \"password123\"
  }'
```

### Example Protected Request

```bash
curl -X GET http://localhost:8080/api/protected \\
  -H \"Authorization: Bearer <your-jwt-token>\"
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Mobile Client  │    │  Other Services │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                  ┌──────────────▼──────────────┐
                  │     API Gateway Service     │
                  │   (Authentication & Routing) │
                  └──────────────┬──────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
┌─────────▼───────┐    ┌─────────▼───────┐    ┌─────────▼───────┐
│  User Service   │    │  Meal Service   │    │ Tracking Service│
│   (Port 8081)   │    │   (Port 8083)   │    │   (Port 8084)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🧪 Testing

### Manual Testing with Swagger UI

1. Start the API service: `make run`
2. Open Swagger UI: `http://localhost:8080/swagger/index.html`
3. Test the `/health` endpoint
4. Test login with valid credentials
5. Copy the JWT token from the login response
6. Use \"Authorize\" button in Swagger UI to set Bearer token
7. Test protected endpoints

### Automated Testing with Postman

See the `/api-tests` directory for comprehensive Postman collections that test all endpoints.

## 📦 Deployment

### Docker

```bash
# Build Docker image
make docker-build

# Run with Docker
docker run -p 8080:8080 \\
  -e USER_SERVICE_URL=http://user-service:8080 \\
  smart-fit-girl/api-service
```

### Production Build

```bash
# Build for Linux production
make build-prod
```

## 🔧 Configuration

### Service Discovery

The API service communicates with other microservices:

- **User Service**: User authentication and management
- **Meal Service**: Meal planning and nutrition (future)
- **Tracking Service**: Progress tracking (future)

### Load Balancing

In production, consider:
- Multiple API service instances behind a load balancer
- Service mesh for inter-service communication
- Circuit breakers for resilience

## 📝 Contributing

1. Add proper Swagger annotations to new endpoints
2. Update struct definitions with JSON tags and examples  
3. Regenerate documentation: `make docs`
4. Test endpoints with Swagger UI
5. Update this README if adding new features

## 🐛 Troubleshooting

### Common Issues

1. **Swagger UI not loading**: Ensure docs are generated with `make docs`
2. **JWT validation fails**: Check JWT_SECRET environment variable
3. **User service unavailable**: Ensure user-service is running
4. **CORS errors**: Check CORS middleware configuration

### Logs

The service logs important events:
- Service startup information
- Authentication attempts
- Service communication errors
- JWT token validation issues

---

## 📊 API Metrics

Consider integrating:
- Prometheus metrics
- Request/response logging
- Performance monitoring
- Error rate tracking
