#!/bin/bash
# Initialize PostgreSQL databases for microservices

# Database credentials from .env
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"

export PGPASSWORD="$DB_PASSWORD"

echo "Initializing PostgreSQL databases..."

# Create databases
for db_name in auth catalog order; do
    echo "Creating database: $db_name"
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -tc "SELECT 1 FROM pg_database WHERE datname = '$db_name'" | grep -q 1 || \
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $db_name"
done

echo "Database initialization complete"
