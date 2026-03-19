#!/bin/bash
# Initialize databases with migrations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICES=("auth-service:auth.db" "catalog-service:catalog.db" "order-service:order.db")

for service_db in "${SERVICES[@]}"; do
    IFS=':' read -r service db <<< "$service_db"

    # Prefer a service-specific DB_PATH from that service's .env file
    env_file="$SCRIPT_DIR/$service/.env"
    db_path=""

    if [ -f "$env_file" ]; then
        db_path=$(grep -m1 '^DB_PATH=' "$env_file" | cut -d'=' -f2-)
    fi

    # Default to the service data directory if DB_PATH is not set.
    if [ -z "$db_path" ]; then
        db_path="$SCRIPT_DIR/$service/data/$db"
    else
        # If DB_PATH is a relative path, resolve it relative to the service folder.
        case "$db_path" in
            /*) ;; # absolute path already
            *) db_path="$SCRIPT_DIR/$service/$db_path" ;;
        esac
    fi

    if [ -f "$db_path" ]; then
        echo "Database already exists: $db_path"
    else
        echo "Creating database: $db_path"
        mkdir -p "$(dirname "$db_path")"
        touch "$db_path"
    fi
done

echo "Database initialization complete"
