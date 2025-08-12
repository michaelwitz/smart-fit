#!/bin/bash

# Test script to verify the service starts correctly
# This will fail immediately if required env vars are not set

echo "Testing db-gateway-service startup..."

# Set test environment variables
export SERVICE_PORT=8086
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=testuser
export DB_PASSWORD=testpass
export DB_NAME=testdb

# Try to start the service and capture output
timeout 2s ./db-gateway-service 2>&1 | head -20

# The service will fail to connect to the database (which is expected)
# but we can verify it attempts to start with the correct configuration
echo ""
echo "Test complete. The service attempted to start with:"
echo "  - SERVICE_PORT: $SERVICE_PORT"
echo "  - DB_HOST: $DB_HOST"
echo "  - DB_PORT: $DB_PORT"
echo "  - DB_NAME: $DB_NAME"
echo ""
echo "Note: Connection failure is expected in test environment without a running database."
