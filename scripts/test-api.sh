#!/bin/bash

# Smart Fit API Test Script
# This script tests the basic API endpoints to verify database connectivity

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_BASE="http://localhost:8080"
USER_SERVICE="http://localhost:8082"
MEAL_SERVICE="http://localhost:8083"

echo -e "${YELLOW}🧪 Smart Fit API Test Suite${NC}"
echo "=================================="

# Function to test an endpoint
test_endpoint() {
    local name="$1"
    local url="$2"
    local method="${3:-GET}"
    local data="${4:-}"
    
    echo -n "Testing $name... "
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    else
        response=$(curl -s -w "%{http_code}" -X "$method" "$url")
    fi
    
    # Extract status code (last 3 characters)
    status_code="${response: -3}"
    # Extract response body (everything except last 3 characters)
    response_body="${response%???}"
    
    if [ "$status_code" = "200" ] || [ "$status_code" = "201" ]; then
        echo -e "${GREEN}✅ PASS${NC} ($status_code)"
        if [ -n "$response_body" ]; then
            echo "   Response: $response_body"
        fi
    else
        echo -e "${RED}❌ FAIL${NC} ($status_code)"
        if [ -n "$response_body" ]; then
            echo "   Error: $response_body"
        fi
    fi
}

# Test health endpoints
echo -e "\n${YELLOW}🏥 Health Checks${NC}"
test_endpoint "API Service Health" "$API_BASE/health"
test_endpoint "User Service Health" "$USER_SERVICE/health"

# Test user creation
echo -e "\n${YELLOW}👤 User Management${NC}"
test_endpoint "Create User" "$USER_SERVICE/users" "POST" '{"fullName":"Test User","email":"test@example.com","password":"testpass123"}'

# Test user verification
echo -e "\n${YELLOW}🔐 Authentication${NC}"
test_endpoint "Verify User" "$USER_SERVICE/users/verify" "POST" '{"email":"test@example.com","password":"testpass123"}'

# Test goals endpoint
echo -e "\n${YELLOW}🎯 Goals${NC}"
test_endpoint "Get All Goals" "$USER_SERVICE/goals"

# Test error handling
echo -e "\n${YELLOW}⚠️  Error Handling${NC}"
test_endpoint "Invalid JSON" "$USER_SERVICE/users" "POST" "invalid json"
test_endpoint "Missing Fields" "$USER_SERVICE/users" "POST" '{"fullName":"Test User"}'

echo -e "\n${YELLOW}📊 Test Summary${NC}"
echo "=================================="
echo "All tests completed!"
echo ""
echo "To run individual tests:"
echo "  curl $API_BASE/health"
echo "  curl $USER_SERVICE/health"
echo "  curl -X POST $USER_SERVICE/users -H 'Content-Type: application/json' -d '{\"fullName\":\"Test User\",\"email\":\"test@example.com\",\"password\":\"testpass123\"}'"
echo ""
echo "Note: Make sure your services are running with:"
echo "  docker-compose up -d"
echo "  # or"
echo "  make up"
