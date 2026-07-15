#!/bin/bash
set -e

echo "Starting Laundry Management System in development mode..."

# Check if PostgreSQL is running
if ! docker compose ps postgres 2>/dev/null | grep -q "Up"; then
    echo "Starting PostgreSQL container..."
    docker compose up -d postgres
    echo "Waiting for PostgreSQL to be ready..."
    sleep 3
fi

# Run with Air (hot reload)
echo "Starting application with Air..."
air -c .air.toml
