#!/bin/bash

# Docker Compose Start Script
# Usage: ./docker-start.sh [service]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${YELLOW}════════════════════════════════════════════════════════════${NC}"
    echo -e "${YELLOW}$1${NC}"
    echo -e "${YELLOW}════════════════════════════════════════════════════════════${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker daemon is not running"
    exit 1
fi

cd "$SCRIPT_DIR"

print_header "Starting Microservices with Docker Compose"

# Build and start services
if docker compose up -d; then
    print_success "All services started successfully"
    echo ""
    sleep 3
    
    # Check service health
    print_info "Checking service health..."
    docker compose ps
    
    echo ""
    print_info "Services are starting up. Waiting for them to become healthy..."
    sleep 5
    
    echo ""
    echo -e "${GREEN}Services available at:${NC}"
    echo "  • Auth Service: http://localhost:8081"
    echo "    - Swagger UI: http://localhost:8081/swagger/index.html"
    echo ""
    echo "  • Catalog Service: http://localhost:8082"
    echo "    - Swagger UI: http://localhost:8082/swagger/index.html"
    echo ""
    echo "  • Order Service: http://localhost:8083"
    echo "    - Swagger UI: http://localhost:8083/swagger/index.html"
    echo ""
    echo "  • PostgreSQL: localhost:5432"
    echo ""
    
    echo -e "${GREEN}Useful commands:${NC}"
    echo "  docker compose logs -f              # View all logs"
    echo "  docker compose logs -f auth-service # View auth-service logs"
    echo "  docker compose ps                   # Show container status"
    echo "  docker compose down                 # Stop all services"
    echo ""
else
    print_error "Failed to start services"
    exit 1
fi
