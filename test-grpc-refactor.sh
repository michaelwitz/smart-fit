#!/bin/bash

echo "Testing gRPC refactoring..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test API service health endpoint (still HTTP)
echo -e "${YELLOW}Testing API service health endpoint...${NC}"
response=$(curl -s http://localhost:8080/health)
if [[ $response == *"healthy"* ]]; then
    echo -e "${GREEN}✓ API service health check passed${NC}"
else
    echo -e "${RED}✗ API service health check failed${NC}"
    echo "Response: $response"
fi

# Test login endpoint (API service should call user-service via gRPC)
echo -e "${YELLOW}Testing login endpoint (gRPC communication)...${NC}"
response=$(curl -s -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"sophia.woytowitz@gmail.com","password":"password"}')

if [[ $response == *"token"* ]]; then
    echo -e "${GREEN}✓ Login endpoint working (gRPC communication successful)${NC}"
    echo "Response contains JWT token"
else
    echo -e "${RED}✗ Login endpoint failed${NC}"
    echo "Response: $response"
fi

# Test with invalid credentials
echo -e "${YELLOW}Testing login with invalid credentials...${NC}"
response=$(curl -s -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"sophia.woytowitz@gmail.com","password":"wrongpassword"}')

if [[ $response == *"Invalid email or password"* ]]; then
    echo -e "${GREEN}✓ Invalid credentials handled correctly${NC}"
else
    echo -e "${RED}✗ Invalid credentials test failed${NC}"
    echo "Response: $response"
fi

echo -e "${YELLOW}Testing gRPC health check on user-service...${NC}"
# Use grpcurl if available to test gRPC health
if command -v grpcurl &> /dev/null; then
    grpcurl -plaintext localhost:8082 grpc.health.v1.Health/Check
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ User service gRPC health check passed${NC}"
    else
        echo -e "${RED}✗ User service gRPC health check failed${NC}"
    fi
else
    echo -e "${YELLOW}grpcurl not installed, skipping gRPC health check${NC}"
fi

echo -e "${GREEN}Testing complete!${NC}"
