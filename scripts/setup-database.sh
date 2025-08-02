#!/bin/bash

# Smart Fit Girl - Database Setup Script

echo "🗄️  Setting up Smart Fit Girl database..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi

# Start PostgreSQL container
echo "🐘 Starting PostgreSQL container..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 10

# Run migrations
echo "🔧 Running database migrations..."
docker-compose exec postgres psql -U smartfit -d smartfitgirl -f /docker-entrypoint-initdb.d/001_create_food_catalog.up.sql

# Load seed data
echo "🌱 Loading seed data..."
docker-compose exec postgres psql -U smartfit -d smartfitgirl -c "$(cat database/seeds/001_food_catalog_seeds.sql)"

echo ""
echo "✅ Database setup complete!"
echo ""
echo "You can now:"
echo "1. Connect to database: docker-compose exec postgres psql -U smartfit -d smartfitgirl"
echo "2. Query food catalog: SELECT * FROM FOOD_CATALOG LIMIT 10;"
echo "3. Start other services: make up"
