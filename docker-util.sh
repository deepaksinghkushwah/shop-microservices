#!/bin/bash

# Docker utility commands for microservices
# Usage: ./docker-util.sh [command]
# Commands:
#   logs [service]           - Show logs for all or specific service
#   status                   - Show container status
#   rebuild [service]        - Rebuild Docker images
#   clean                    - Remove all images and volumes
#   ps                       - Show running processes

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

COMMAND="${1:-help}"
SERVICE="${2:-}"

case "$COMMAND" in
    logs)
        print_header "Service Logs"
        if [ -n "$SERVICE" ]; then
            docker compose logs -f "$SERVICE"
        else
            docker compose logs -f
        fi
        ;;
    status|ps)
        print_header "Container Status"
        docker compose ps
        ;;
    rebuild)
        print_header "Rebuilding Docker Images"
        if [ -n "$SERVICE" ]; then
            docker compose build --no-cache "$SERVICE"
        else
            docker compose build --no-cache
        fi
        print_success "Docker images rebuilt"
        ;;
    clean)
        print_header "Cleaning Docker Resources"
        print_info "Stopping services..."
        docker compose down -v 2>/dev/null || true
        print_info "Removing images..."
        docker compose rm -f 2>/dev/null || true
        print_success "Docker resources cleaned"
        ;;
    exec)
        print_info "Executing command in container"
        # Usage: ./docker-util.sh exec auth-service sh
        SERVICE="${2}"
        COMMAND="${@:3}"
        docker compose exec "$SERVICE" $COMMAND
        ;;
    shell)
        print_info "Opening shell in $SERVICE container"
        SERVICE="${2:-auth-service}"
        docker compose exec "$SERVICE" sh
        ;;
    *)
        echo "Docker Utility Commands"
        echo ""
        echo "Usage: ./docker-util.sh [command] [options]"
        echo ""
        echo "Commands:"
        echo "  logs [service]      - Show logs (press Ctrl+C to exit)"
        echo "  status              - Show container status"
        echo "  ps                  - Show running containers"
        echo "  rebuild [service]   - Rebuild Docker images"
        echo "  clean               - Remove all containers and volumes"
        echo "  exec <service> cmd  - Execute command in container"
        echo "  shell [service]     - Open shell in container (default: auth-service)"
        echo "  help, -h, --help    - Show this help message"
        echo ""
        echo "Examples:"
        echo "  ./docker-util.sh logs                   # Show all service logs"
        echo "  ./docker-util.sh logs auth-service      # Show auth-service logs"
        echo "  ./docker-util.sh status                 # Show container status"
        echo "  ./docker-util.sh rebuild                # Rebuild all images"
        echo "  ./docker-util.sh shell auth-service     # Open shell in auth-service"
        echo ""
        ;;
esac
