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
