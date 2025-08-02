#!/bin/bash

# Smart Fit Girl - Development Environment Setup

echo "🏋️‍♀️ Setting up Smart Fit Girl development environment..."

# Check if .env exists, if not copy from .env.example
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created! Please update with your actual values."
else
    echo "ℹ️  .env file already exists"
fi

# Create individual service .env files if they don't exist
for service in user-service meal-service tracking-service api-service; do
    if [ ! -f "services/$service/.env" ]; then
        echo "📝 Creating .env for $service..."
        cp "services/$service/.env.example" "services/$service/.env"
    fi
done

echo ""
echo "🚀 Environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your actual database credentials"
echo "2. Add your PostHog API key to .env"
echo "3. Run: make up    (or docker-compose up -d)"
echo "4. Test APIs with Postman at http://localhost:8080"
echo ""
echo "Service URLs:"
echo "  - API Service:      http://localhost:8080"
echo "  - User Service:     http://localhost:8082"
echo "  - Meal Service:     http://localhost:8083"
echo "  - Tracking Service: http://localhost:8084"
echo "  - Web Client:       http://localhost:3000"
echo "  - PostgreSQL:       localhost:5432"
