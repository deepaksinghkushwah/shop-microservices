#!/bin/bash
# Stop all microservices by port

# Function to kill process on a given port
kill_by_port() {
    local port=$1
    local service=$2
    
    # Find and kill process listening on port
    if command -v lsof &> /dev/null; then
        local pid=$(lsof -ti:$port 2>/dev/null)
        if [ -n "$pid" ]; then
            kill -9 $pid 2>/dev/null
            echo "Stopped $service (PID: $pid) listening on port $port"
        fi
    else
        # Fallback: use fuser if lsof is not available
        if command -v fuser &> /dev/null; then
            fuser -k $port/tcp 2>/dev/null && echo "Stopped service listening on port $port"
        else
            echo "Warning: lsof or fuser not found. Cannot stop service on port $port"
        fi
    fi
}

echo "Stopping all microservices..."

# Kill services by their ports
kill_by_port 8081 "auth-service"
kill_by_port 8082 "catalog-service"
kill_by_port 8083 "order-service"

sleep 1

echo "All services stopped"
