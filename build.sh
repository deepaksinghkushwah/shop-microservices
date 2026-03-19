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

# Service output directories will be created per service (binary, docs, data, logs)
# Each service has its own .env file for configuration

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

        # Create service-specific env files (ports + db path)
        port_var="${binary_name^^}_SERVICE_PORT"
        case "$service" in
            auth-service) port_value="8081" ;; 
            catalog-service) port_value="8082" ;; 
            order-service) port_value="8083" ;; 
            *) port_value="8080" ;; 
        esac

        cat > "$service_dist_dir/.env" <<EOF
# Environment variables for $service
# Adjust as needed for your deployment.

$port_var=$port_value
DB_PATH=data/$binary_name.db
EOF

        cat > "$service_dist_dir/.env.example" <<EOF
# Example env for $service

$port_var=$port_value
DB_PATH=data/$binary_name.db
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
services=("auth-service" "catalog-service" "order-service")
pids=()

for service in "${services[@]}"; do
    bin_name="${service%-service}"
    service_dir="./${service}"

    if [ -f "$service_dir/$bin_name" ]; then
        mkdir -p "$service_dir/logs"
        echo -e "${GREEN}Starting ${service}...${NC}"
        pushd "$service_dir" > /dev/null
        ./$bin_name > "logs/${bin_name}.log" 2>&1 &
        pids+=($!)
        popd > /dev/null
        sleep 1
    fi
done

echo ""
echo -e "${GREEN}All services started with PIDs: ${pids[@]}${NC}"
echo "Logs available in each service's logs/ directory"
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
STOPSCRIPT

chmod +x "$DIST_DIR/stop_all.sh"
print_success "Created service stop script"

# Create a README for the dist folder
cat > "$DIST_DIR/README.md" << 'README'
# Microservices Build Output

This directory contains built microservice binaries, API docs, and databases organized per service.

## Structure
```
dist/
├── auth-service/
│   ├── auth
│   ├── docs/
│   ├── data/
│   │   └── auth.db
│   └── logs/
│       └── auth.log
├── catalog-service/
│   ├── catalog
│   ├── docs/
│   ├── data/
│   │   └── catalog.db
│   └── logs/
│       └── catalog.log
├── order-service/
│   ├── order
│   ├── docs/
│   ├── data/
│   │   └── order.db
│   └── logs/
│       └── order.log
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
tail -f auth-service/logs/auth.log
tail -f catalog-service/logs/catalog.log
tail -f order-service/logs/order.log
```

## Service Endpoints

- **Auth Service** - http://localhost:8081
  - Swagger UI: http://localhost:8081/swagger
  
- **Catalog Service** - http://localhost:8082
  - Swagger UI: http://localhost:8082/swagger
  
- **Order Service** - http://localhost:8083
  - Swagger UI: http://localhost:8083/swagger

## Database Files

Each service uses its own SQLite database under its `data/` directory:
- `auth-service/data/auth.db`
- `catalog-service/data/catalog.db`
- `order-service/data/order.db`

Migrations are automatically run on first service startup.

## Environment Variables

Each service has its own configuration file at `<service>/.env`.

Example (`auth-service/.env`):
```
AUTH_SERVICE_PORT=8081
DB_PATH=data/auth.db
```

Example (`catalog-service/.env`):
```
CATALOG_SERVICE_PORT=8082
DB_PATH=data/catalog.db
```

Example (`order-service/.env`):
```
ORDER_SERVICE_PORT=8083
DB_PATH=data/order.db
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
    echo "  1. Configure each service: edit dist/<service>/.env for each service"
    echo "  2. Start services: cd $DIST_DIR && ./start_all.sh"
    echo "  3. View logs: tail -f $DIST_DIR/<service>/logs/*.log (e.g. $DIST_DIR/auth-service/logs/auth.log)"
    echo ""
    exit 0
fi
