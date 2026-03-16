#!/bin/bash

# Test script - runs all tests for all services
# Usage: ./test.sh [SERVICE_NAME] [args...]
# Examples:
#   ./test.sh                 # run all tests
#   ./test.sh auth-service    # run auth-service tests only
#   ./test.sh -v              # run all tests with verbose output

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICES=("auth-service" "catalog-service" "order-service")
TEST_ARGS="${@}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Check if specific service is requested
if [ $# -gt 0 ] && [ -d "services/$1" ]; then
    SERVICES=("$1")
    TEST_ARGS="${@:2}"
elif [ $# -gt 0 ] && [[ "$1" != "-"* ]]; then
    print_error "Service '$1' not found. Available services: ${SERVICES[@]}"
    exit 1
fi

failed_tests=()
total_tests=0

cd "$SCRIPT_DIR"

# Run tests for each service
for service in "${SERVICES[@]}"; do
    if [ ! -d "services/$service" ]; then
        print_error "Service directory not found: services/$service"
        continue
    fi

    print_header "Running tests for $service"
    
    total_tests=$((total_tests + 1))
    
    if go test -v $TEST_ARGS ./services/$service/... -timeout 30s; then
        print_success "$service tests passed"
    else
        print_error "$service tests failed"
        failed_tests+=("$service")
    fi
    
    echo ""
done

# Summary
print_header "Test Summary"
echo "Total services tested: $total_tests"

if [ ${#failed_tests[@]} -eq 0 ]; then
    print_success "All tests passed!"
    exit 0
else
    print_error "Failed services: ${failed_tests[@]}"
    exit 1
fi
