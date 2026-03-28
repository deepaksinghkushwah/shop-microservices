#!/bin/bash

# Docker Compose Stop Script
# Usage: ./docker-stop.sh

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

cd "$SCRIPT_DIR"

print_header "Stopping Microservices"

if docker compose down; then
    print_success "All services stopped"
    echo ""
    print_info "Containers have been stopped and removed"
    echo ""
else
    print_error "Failed to stop services"
    exit 1
fi
