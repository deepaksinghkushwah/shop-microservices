#!/bin/bash

# Build script - builds all services and sets up databases with migrations
# Usage: ./build.sh [clean] [SERVICE_NAME]
# Examples:
#   ./build.sh                 # build all services
#   ./build.sh clean           # clean dist folder and rebuild
#   ./build.sh auth-service    # build only auth-service
#   ./build.sh clean auth-service  # clean and rebuild only auth-service

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$SCRIPT_DIR/dist"
SERVICES=("auth-service" "catalog-service" "order-service")
SERVICE_PORTS=("8081" "8082" "8083")
CLEAN_BUILD=false

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
    if [ "$arg" = "clean" ]; then
        CLEAN_BUILD=true
    elif [ -d "services/$arg" ]; then
        SERVICES=("$arg")
    elif [[ "$arg" != "-"* ]]; then
        print_error "Unknown service: $arg. Available: auth-service, catalog-service, order-service"
        exit 1
    fi
done

cd "$SCRIPT_DIR"

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

# Create database directory
mkdir -p "$DIST_DIR/data"
print_success "Created data directory for databases"

# Create migrations directory
mkdir -p "$DIST_DIR/logs"
print_success "Created logs directory"

# Copy .env file if it exists
if [ -f "$SCRIPT_DIR/.env" ]; then
    cp "$SCRIPT_DIR/.env" "$DIST_DIR/.env"
    print_success "Copied .env file to dist"
else
    print_info "No .env file found in root directory"
fi

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
    output_path="$DIST_DIR/$binary_name"

    # Build the service
    if go build -o "$output_path" "$service_dir/cmd/server/main.go"; then
        print_success "Built $service binary: $output_path"
        successful_builds+=("$service")

        # Copy service-specific files if they exist
        if [ -d "$service_dir/docs" ]; then
            mkdir -p "$DIST_DIR/${service}_docs"
            cp -r "$service_dir/docs"/* "$DIST_DIR/${service}_docs/" 2>/dev/null || true
            print_success "Copied API docs for $service"
        fi

        # Create service data directory
        mkdir -p "$DIST_DIR/data/$service"
        print_success "Created data directory for $service"

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
DBINIT

chmod +x "$DIST_DIR/init_db.sh"
print_success "Created database initialization script"

# Create startup script
print_header "Creating service startup script"

cat > "$DIST_DIR/start_all.sh" << 'STARTUP'
#!/bin/bash
# Start all microservices

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Starting all microservices...${NC}"
echo ""

# Initialize databases if needed
if [ -f "./init_db.sh" ]; then
    ./init_db.sh
    echo ""
fi

# Start each service in background
services=("auth" "catalog" "order")
pids=()

for service in "${services[@]}"; do
    if [ -f "./$service" ]; then
        echo -e "${GREEN}Starting $service service...${NC}"
        ./$service > logs/${service}.log 2>&1 &
        pids+=($!)
        sleep 1
    fi
done

echo ""
echo -e "${GREEN}All services started with PIDs: ${pids[@]}${NC}"
echo "Logs available in logs/ directory"
echo ""
echo "Services:"
echo "  - Auth Service: http://localhost:8081/swagger"
echo "  - Catalog Service: http://localhost:8082/swagger"
echo "  - Order Service: http://localhost:8083/swagger"
echo ""
echo "To stop services, run: kill ${pids[@]}"
STARTUP

chmod +x "$DIST_DIR/start_all.sh"
print_success "Created service startup script"

# Create stop script
cat > "$DIST_DIR/stop_all.sh" << 'STOPSCRIPT'
#!/bin/bash
# Stop all microservices

pkill -f "^\./auth$" || true
pkill -f "^\./catalog$" || true
pkill -f "^\./order$" || true

echo "All services stopped"
STOPSCRIPT

chmod +x "$DIST_DIR/stop_all.sh"
print_success "Created service stop script"

# Create a README for the dist folder
cat > "$DIST_DIR/README.md" << 'README'
# Microservices Build Output

This directory contains built microservice binaries with initialized databases.

## Structure
```
dist/
├── auth              # Auth service binary
├── catalog           # Catalog service binary
├── order             # Order service binary
├── data/             # SQLite databases
│   ├── auth-service/
│   ├── catalog-service/
│   └── order-service/
├── logs/             # Service logs
├── *_docs/           # API documentation (Swagger)
├── .env              # Environment configuration
├── start_all.sh      # Start all services
├── stop_all.sh       # Stop all services
└── init_db.sh        # Initialize databases
```

## Usage

### Start all services
```bash
./start_all.sh
```

### Stop all services
```bash
./stop_all.sh
```

### View logs
```bash
tail -f logs/auth.log
tail -f logs/catalog.log
tail -f logs/order.log
```

## Service Endpoints

- **Auth Service** - http://localhost:8081
  - Swagger UI: http://localhost:8081/swagger
  
- **Catalog Service** - http://localhost:8082
  - Swagger UI: http://localhost:8082/swagger
  
- **Order Service** - http://localhost:8083
  - Swagger UI: http://localhost:8083/swagger

## Database Files

Databases are SQLite files stored in `data/`:
- `data/auth-service/auth.db`
- `data/catalog-service/catalog.db`
- `data/order-service/order.db`

Migrations are automatically run on first service startup.

## Environment Variables

Configure services using `.env` file in this directory:
```
AUTH_SERVICE_PORT=8081
CATALOG_SERVICE_PORT=8082
ORDER_SERVICE_PORT=8083
```

## Building from Source

To rebuild from source:
```bash
cd ..
./build.sh              # Build all services
./build.sh clean        # Clean and rebuild
./build.sh auth-service # Build single service
```
README

print_success "Created dist/README.md"

# Create .env.example if .env doesn't exist in dist
if [ ! -f "$DIST_DIR/.env" ]; then
    cat > "$DIST_DIR/.env.example" << 'ENVFILE'
# Service Ports
AUTH_SERVICE_PORT=8081
CATALOG_SERVICE_PORT=8082
ORDER_SERVICE_PORT=8083

# Database Configuration (if needed)
# DATABASE_HOST=localhost
# DATABASE_PORT=5432
# DATABASE_USER=postgres
# DATABASE_PASSWORD=password
ENVFILE

    print_info "Created .env.example - copy to .env and configure if needed"
fi

# Build summary
print_header "Build Summary"

if [ ${#successful_builds[@]} -gt 0 ]; then
    echo -e "${GREEN}Successfully built:${NC}"
    for service in "${successful_builds[@]}"; do
        echo "  ✓ $service"
    done
fi

if [ ${#failed_builds[@]} -gt 0 ]; then
    echo -e "${RED}Failed builds:${NC}"
    for service in "${failed_builds[@]}"; do
        echo "  ✗ $service"
    done
    echo ""
    print_error "Build completed with errors"
    exit 1
else
    echo ""
    print_success "All services built successfully!"
    echo ""
    echo "Next steps:"
    echo "  1. Configure environment: cd $DIST_DIR && cp .env.example .env"
    echo "  2. Start services: cd $DIST_DIR && ./start_all.sh"
    echo "  3. View logs: tail -f $DIST_DIR/logs/*.log"
    echo ""
    exit 0
fi
