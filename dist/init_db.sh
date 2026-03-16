#!/bin/bash
# Initialize databases with migrations

DATA_DIR="$(dirname "$0")/data"
SERVICES=("auth-service:auth.db" "catalog-service:catalog.db" "order-service:order.db")

for service_db in "${SERVICES[@]}"; do
    IFS=':' read -r service db <<< "$service_db"
    db_path="$DATA_DIR/$service/$db"
    
    if [ -f "$db_path" ]; then
        echo "Database already exists: $db_path"
    else
        echo "Creating database: $db_path"
        mkdir -p "$(dirname "$db_path")"
        touch "$db_path"
    fi
done

echo "Database initialization complete"
