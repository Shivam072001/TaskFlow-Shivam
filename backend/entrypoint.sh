#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

echo "Applying seed data..."
/app/seeder || echo "Seeding may have partially failed, continuing..."

echo "Starting server..."
exec /app/server
