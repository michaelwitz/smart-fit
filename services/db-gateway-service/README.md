# DB Gateway Service

A gRPC service that provides database access for the Smart Fit application.

## Overview

This service acts as a data access layer, exposing gRPC endpoints for database operations. It implements the repository pattern and provides a clean interface for other services to interact with the database.

## Structure

```
db-gateway-service/
├── main.go                      # Service entry point
├── internal/                    # Private implementation (Go enforced)
│   ├── database/               # Database connection management
│   │   └── connection.go       # Connection pool implementation
│   └── services/               # gRPC service implementations
│       ├── user_service.go     # UserService implementation
│       └── user_service_test.go # Unit tests
├── proto/                       # Generated protobuf files
│   ├── user.pb.go              # User message definitions
│   └── user_grpc.pb.go         # User service definitions
└── sql/                        # SQL repositories
    └── user-service/
        └── users.go            # User repository implementation
```

## Environment Variables

The service requires the following environment variables:

- `SERVICE_PORT` - gRPC server port (default: 8086)
- `DB_HOST` - PostgreSQL host (required)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database user (required)
- `DB_PASSWORD` - Database password (required)
- `DB_NAME` - Database name (required)

## Running the Service

### Local Development

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=smartfit

# Run the service
go run main.go
```

### Docker

```bash
docker build -t db-gateway-service .
docker run -p 8086:8086 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=yourpassword \
  -e DB_NAME=smartfit \
  db-gateway-service
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific test package
go test ./internal/services -v
```

## gRPC Endpoints

### UserService

- `CreateUser` - Create a new user
- `GetUserByID` - Retrieve a user by ID
- `GetAllUsers` - Retrieve all users
- `UpdateUser` - Update an existing user
- `DeleteUser` - Delete a user
- `VerifyUser` - Verify user credentials
- `UpsertUser` - Create or update a user

## Development

### Regenerating Protocol Buffers

If you need to regenerate the protobuf files:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       ../../proto/user.proto
```

### Adding New Services

1. Create the repository in `sql/<service-name>/`
2. Create the gRPC service implementation in `internal/services/`
3. Wire the service in `main.go`
4. Add tests in `internal/services/`

## Architecture Notes

### Why `/internal`?

The `/internal` directory is a Go convention that provides compile-time enforcement of package privacy. Packages inside `/internal` can only be imported by code in the parent directory tree. This ensures:

1. **Clear API boundaries** - Only the gRPC interface is public
2. **Implementation hiding** - Other services cannot depend on internal implementation details
3. **Easier refactoring** - Internal code can be changed without breaking external consumers

### Repository Pattern

The service uses the repository pattern to separate business logic from data access:

- **Repository** (`sql/user-service/users.go`) - Handles database operations
- **Service** (`internal/services/user_service.go`) - Implements business logic and gRPC interface
- **Proto** (`proto/`) - Defines the contract between services

This separation allows for:
- Easy testing with mocked repositories
- Database implementation changes without affecting business logic
- Clear separation of concerns
