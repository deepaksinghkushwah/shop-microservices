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
