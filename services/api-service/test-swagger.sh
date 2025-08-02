#!/bin/bash

# Smart Fit Girl API Service - Swagger Documentation Test Script

API_BASE="http://localhost:8080"

echo "🧪 Testing Smart Fit Girl API Service with Swagger Documentation"
echo "=================================================================="

# Test if service is running
echo "1. Testing service health..."
HEALTH_RESPONSE=$(curl -s "${API_BASE}/health")
if [[ $? -eq 0 ]]; then
    echo "✅ Service is running: $HEALTH_RESPONSE"
else
    echo "❌ Service is not running or not accessible"
    exit 1
fi

# Test Swagger JSON endpoint
echo ""
echo "2. Testing Swagger JSON documentation..."
SWAGGER_JSON=$(curl -s "${API_BASE}/swagger/doc.json")
if [[ $? -eq 0 && "$SWAGGER_JSON" != *"404"* ]]; then
    echo "✅ Swagger JSON is available"
    echo "   First few lines:"
    echo "$SWAGGER_JSON" | head -5
else
    echo "❌ Swagger JSON not available"
    echo "   Response: $SWAGGER_JSON"
fi

# Test Swagger UI endpoint
echo ""
echo "3. Testing Swagger UI..."
SWAGGER_UI=$(curl -s "${API_BASE}/swagger/index.html" | head -1)
if [[ $? -eq 0 && "$SWAGGER_UI" != *"404"* ]]; then
    echo "✅ Swagger UI is available"
    echo "   Response starts with: $SWAGGER_UI"
else
    echo "❌ Swagger UI not available"
fi

# Test authentication endpoint
echo ""
echo "4. Testing authentication endpoint..."
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"sophia.woytowitz@gmail.com","password":"password"}')

if [[ $? -eq 0 ]]; then
    echo "✅ Login endpoint is accessible"
    if [[ "$LOGIN_RESPONSE" == *"token"* ]]; then
        echo "   ✅ Login successful - JWT token received"
        # Extract token for protected endpoint test
        TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    else
        echo "   ❌ Login failed: $LOGIN_RESPONSE"
    fi
else
    echo "❌ Login endpoint not accessible"
fi

# Test protected endpoint
if [[ -n "$TOKEN" ]]; then
    echo ""
    echo "5. Testing protected endpoint with JWT..."
    PROTECTED_RESPONSE=$(curl -s -X GET "${API_BASE}/api/protected" \
        -H "Authorization: Bearer $TOKEN")
    
    if [[ $? -eq 0 && "$PROTECTED_RESPONSE" == *"protected endpoint"* ]]; then
        echo "✅ Protected endpoint accessible with JWT"
        echo "   Response: $PROTECTED_RESPONSE"
    else
        echo "❌ Protected endpoint failed: $PROTECTED_RESPONSE"
    fi
fi

echo ""
echo "🎉 Testing complete!"
echo ""
echo "📖 Access the Swagger UI at: ${API_BASE}/swagger/index.html"
echo "📄 Access the Swagger JSON at: ${API_BASE}/swagger/doc.json"
echo ""
echo "🔧 Available endpoints:"
echo "   GET  /health           - Health check"
echo "   POST /auth/login       - User authentication"  
echo "   GET  /api/protected    - Protected endpoint (requires JWT)"
echo ""
echo "💡 To rebuild services with Swagger documentation:"
echo "   docker-compose down"
echo "   docker-compose build api-service"
echo "   docker-compose up -d"
