#!/bin/bash

# Build script - builds all services with Docker support
# Usage: ./build.sh [docker|clean] [SERVICE_NAME]
# Examples:
#   ./build.sh                        # build all services locally
#   ./build.sh docker                 # build all services with Docker
#   ./build.sh clean                  # clean dist folder and rebuild locally
#   ./build.sh docker clean           # rebuild Docker images
#   ./build.sh auth-service           # build only auth-service locally
#   ./build.sh docker auth-service    # build only auth-service Docker image

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$SCRIPT_DIR/dist"
SERVICES=("auth-service" "catalog-service" "order-service")
SERVICE_PORTS=("8081" "8082" "8083")
CLEAN_BUILD=false
DOCKER_BUILD=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}$1${NC}"
    echo -e "${YELLOW}========================================${NC}"
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

# Parse arguments
for arg in "$@"; do
    if [ "$arg" = "docker" ]; then
        DOCKER_BUILD=true
    elif [ "$arg" = "clean" ]; then
        CLEAN_BUILD=true
    elif [ -d "services/$arg" ]; then
        SERVICES=("$arg")
    elif [[ "$arg" != "-"* ]]; then
        print_error "Unknown service: $arg. Available: auth-service, catalog-service, order-service"
        exit 1
    fi
done

cd "$SCRIPT_DIR"

# Function to build Docker images
build_docker_images() {
    print_header "Building Docker Images"
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    # Check if Docker daemon is running
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker daemon is not running"
        exit 1
    fi
    
    # Clean Docker images if requested
    if [ "$CLEAN_BUILD" = true ]; then
        print_info "Cleaning Docker images and containers..."
        docker compose down --rmi all 2>/dev/null || true
        print_success "Cleaned Docker resources"
    fi
    
    # Build Docker images
    if docker compose build; then
        print_success "All Docker images built successfully"
        echo ""
        echo -e "${GREEN}Docker images ready!${NC}"
        echo ""
        echo "Next steps:"
        echo "  1. Start services: docker compose up -d"
        echo "  2. View logs: docker compose logs -f"
        echo "  3. Stop services: docker compose down"
    else
        print_error "Failed to build Docker images"
        exit 1
    fi
}

# Function to build local binaries
build_local_binaries() {
    # Clean if requested
    if [ "$CLEAN_BUILD" = true ]; then
        print_header "Cleaning dist directory"
        if [ -d "$DIST_DIR" ]; then
            rm -rf "$DIST_DIR"
            print_success "Cleaned dist directory"
        fi
    fi

    # Create dist structure
    print_header "Setting up dist directory"
    mkdir -p "$DIST_DIR"
    print_success "Created dist directory at $DIST_DIR"

    failed_builds=()
    successful_builds=()

    # Build each service
    for service in "${SERVICES[@]}"; do
        if [ ! -d "services/$service" ]; then
            print_error "Service directory not found: services/$service"
            failed_builds+=("$service")
            continue
        fi

        print_header "Building $service"

        service_dir="services/$service"
        binary_name="${service%-service}"
        service_dist_dir="$DIST_DIR/$service"
        output_path="$service_dist_dir/$binary_name"

        # Prepare service dist layout
        mkdir -p "$service_dist_dir"
        mkdir -p "$service_dist_dir/data"
        mkdir -p "$service_dist_dir/logs"

        # Build the service
        if CGO_ENABLED=1 go build -o "$output_path" "$service_dir/cmd/server/main.go"; then
            print_success "Built $service binary: $output_path"
            successful_builds+=("$service")

            # Copy service-specific files if they exist
            if [ -d "$service_dir/docs" ]; then
                mkdir -p "$service_dist_dir/docs"
                cp -r "$service_dir/docs"/* "$service_dist_dir/docs/" 2>/dev/null || true
                print_success "Copied API docs for $service"
            fi

            # Create service-specific env files (PostgreSQL)
            port_var="${binary_name^^}_SERVICE_PORT"
            case "$service" in
                auth-service) port_value="8081"; db_name="auth" ;; 
                catalog-service) port_value="8082"; db_name="catalog" ;; 
                order-service) port_value="8083"; db_name="order" ;; 
                *) port_value="8080"; db_name="app" ;; 
            esac

            cat > "$service_dist_dir/.env" <<EOF
# Environment variables for $service
# PostgreSQL Configuration

$port_var=$port_value
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=$db_name
DB_SSLMODE=disable
EOF

            cat > "$service_dist_dir/.env.example" <<EOF
# Example env for $service

$port_var=$port_value
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=$db_name
DB_SSLMODE=disable
EOF

            print_success "Created env files for $service"
            print_success "Prepared dist layout for $service"

        else
            print_error "Failed to build $service"
            failed_builds+=("$service")
        fi

        echo ""
    done

    # Create database initialization script
    print_header "Creating database initialization utility"

    cat > "$DIST_DIR/init_db.sh" << 'DBINIT'
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
DBINIT

    chmod +x "$DIST_DIR/init_db.sh"
    print_success "Created database initialization script"
}

# Main execution logic
if [ "$DOCKER_BUILD" = true ]; then
    build_docker_images
else
    build_local_binaries
fi

print_header "Build Complete!"
print_success "All services built successfully"

