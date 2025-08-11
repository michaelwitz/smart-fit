# Smart Fit Testing Guide

This document explains how to test the Smart Fit application using both Go tests and curl commands.

## 🧪 Testing Approaches

### 1. **Go Tests (Recommended for Development)**

- **Location**: `services/api-service/api_test.go` (API Gateway tests)
- **Type**: Integration tests using Go's `testing` package
- **Benefits**: Fast, automated, type-safe, CI/CD ready

### 2. **Postman Tests (Manual Testing)**

- **Location**: `api-tests/` directory
- **Type**: Manual API testing with Postman
- **Benefits**: Visual, easy to modify, good for debugging

### 3. **Curl Tests (Quick Verification)**

- **Location**: `scripts/test-api.sh`
- **Type**: Bash script with curl commands
- **Benefits**: Quick, no GUI needed, good for CI/CD

## 🚀 Running Tests

### **Go Tests (Fastest)**

```bash
# Test API gateway (end-to-end testing)
cd services/api-service && go test -v

# Test specific function
go test -v -run TestHealthEndpoint

# Test with coverage
go test -v -cover
```

### **Curl Tests (Quick)**

```bash
# Run the test script
./scripts/test-api.sh

# Or run individual tests
curl http://localhost:8080/health
curl http://localhost:8082/health
```

### **Postman Tests (Manual)**

1. Import `api-tests/collections/Smart_Fit_Girl_API.postman_collection.json`
2. Import `api-tests/environments/Local_Development.postman_environment.json`
3. Set "Local Development" as active environment
4. Run the collection

## 🏗️ Test Structure

### **Go Test Files**

The API service has `api_test.go` that:

- Creates a test server using `httptest` and Gin
- Tests the API gateway endpoints
- Verifies routing, middleware, and authentication
- Tests end-to-end API functionality with real database users

### **Test Coverage**

- ✅ **Health endpoints** - API gateway availability
- ✅ **Authentication** - JWT login with real test users
- ✅ **Protected endpoints** - JWT token validation
- ✅ **Error handling** - Invalid JSON, wrong credentials
- ✅ **CORS middleware** - Cross-origin request handling

## 🔧 Test Environment

### **Database Requirements**

- PostgreSQL running on localhost:5432
- Database: `smartfitgirl`
- User: `smartfit`
- Password: `smartfit123`

### **Service Requirements**

- All services running via Docker Compose
- Environment variables properly set in `.env`

## 📊 Test Results

### **Go Tests Output**

```
=== RUN   TestHealthEndpoint
--- PASS: TestHealthEndpoint (0.02s)
=== RUN   TestLoginEndpoint
--- PASS: TestLoginEndpoint (0.15s)
=== RUN   TestProtectedEndpoint
--- PASS: TestProtectedEndpoint (0.05s)
=== RUN   TestErrorHandling
--- PASS: TestErrorHandling (0.03s)
=== RUN   TestCORSHeaders
--- PASS: TestCORSHeaders (0.02s)
PASS
ok      api-service    0.270s
```

### **Curl Tests Output**

```
🧪 Smart Fit API Test Suite
==================================
🏥 Health Checks
Testing API Service Health... ✅ PASS (200)
Testing User Service Health... ✅ PASS (200)
👤 User Management
Testing Create User... ✅ PASS (201)
```

## 🚨 Troubleshooting

### **Common Issues**

1. **Database Connection Failed**

   ```bash
   # Check if PostgreSQL is running
   docker-compose ps postgres

   # Check database connection
   psql -h localhost -U smartfit -d smartfitgirl
   ```

2. **Service Not Responding**

   ```bash
   # Check service status
   docker-compose ps

   # Check service logs
   docker-compose logs user-service
   ```

3. **Environment Variables Missing**

   ```bash
   # Verify .env file exists
   ls -la .env

   # Check environment variables
   docker-compose exec user-service env | grep DB
   ```

### **Test Database Setup**

```bash
# Reset database
make reset-db

# Setup fresh database
make setup-db

# Verify seed data
psql -h localhost -U smartfit -d smartfitgirl -c "SELECT COUNT(*) FROM USERS;"
```

## 🔄 Continuous Testing

### **Pre-commit Hook**

```bash
# Add to .git/hooks/pre-commit
#!/bin/bash
cd services/user-service && go test -v
cd ../meal-service && go test -v
# ... other services
```

### **CI/CD Pipeline**

```yaml
# Example GitHub Actions
- name: Run Tests
  run: |
    cd services/user-service && go test -v
    cd ../meal-service && go test -v
    cd ../check-in-service && go test -v
    cd ../survey-service && go test -v
```

## 📝 Adding New Tests

### **Go Tests**

1. Add test function to `services/api-service/api_test.go`
2. Follow naming convention: `TestFunctionName`
3. Use `httptest.NewServer` with Gin router for HTTP testing
4. Test both success and error cases
5. Use real test users from seed data for authentication tests

### **Curl Tests**

1. Add new test to `scripts/test-api.sh`
2. Use the `test_endpoint` function
3. Include proper error handling
4. Add to appropriate test category

## 🎯 Best Practices

1. **Test Real Database**: Use actual database, not mocks
2. **Clean State**: Reset database between test runs
3. **Comprehensive Coverage**: Test success, error, and edge cases
4. **Fast Execution**: Keep tests under 5 seconds total
5. **Clear Assertions**: Use descriptive error messages

## 🔗 Related Files

- `services/api-service/api_test.go` - Go integration tests for API gateway
- `scripts/test-api.sh` - Curl test script
- `api-tests/` - Postman collections
- `docker-compose.yml` - Test environment setup
- `.env` - Test configuration
