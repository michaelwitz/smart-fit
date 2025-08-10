# Docker Organization

## Structure

```
docker/
├── base/
│   └── Dockerfile.go-base          # Base Go image with common tools
├── services/                        # Service-specific Dockerfiles (future)
└── README.md                       # This file
```

## Philosophy

- **Consistency**: All Go services use the same base image and build process
- **Alpine-based**: Lightweight, secure, and consistent across all services
- **Latest versions**: Go 1.24.6, Node.js 20 LTS
- **Multi-stage builds**: Optimized for production with minimal runtime images

## Current Services

### Go Services (All use Go 1.24.6 + Alpine)
- **api-service**: Main API gateway with Gin framework
- **user-service**: User management and authentication
- **meal-service**: Meal planning and nutrition tracking
- **check-in-service**: Daily check-ins, weight, mood, sleep, activity tracking
- **survey-service**: Survey management, goals, user preferences

### Web Client
- **web-client**: React.js frontend with Node.js 20 LTS

### Database
- **postgres**: PostgreSQL 15 Alpine

## Service Ports

- **api-service**: 8080 (main gateway)
- **user-service**: 8082
- **meal-service**: 8083
- **check-in-service**: 8084
- **survey-service**: 8085
- **web-client**: 5050
- **postgres**: 5432

## Best Practices

1. **Always use Alpine** for consistency and size
2. **Latest stable versions** for Go and Node.js
3. **Multi-stage builds** to minimize final image size
4. **Non-root execution** when possible (future enhancement)
5. **Consistent port exposure** (8080 for Go services, 5050 for web)

## Future Enhancements

- Service-specific Dockerfiles in `docker/services/`
- Production-optimized images with distroless
- Security scanning and vulnerability checks
- Multi-architecture builds (ARM64 support)
