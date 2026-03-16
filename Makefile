# Microservices Build System
#
# Available targets:
#   make test              - Run all tests
#   make test-verbose      - Run all tests with verbose output  
#   make test-<service>    - Run tests for specific service
#   make build             - Build all services
#   make build-clean       - Clean and rebuild all services
#   make build-<service>   - Build specific service only
#   make dist              - Alias for build
#   make dist/start        - Start all services from dist
#   make dist/stop         - Stop all services
#   make dist/logs         - Watch all service logs
#   make clean             - Remove dist folder
#   make help              - Show this help message

.PHONY: help test test-verbose build build-clean clean dist help
.PHONY: dist/start dist/stop dist/logs dist/status
.PHONY: test-auth-service test-catalog-service test-order-service
.PHONY: build-auth-service build-catalog-service build-order-service

# Default target
.DEFAULT_GOAL := help

SHELL := /bin/bash
DIST_DIR := ./dist
SERVICES := auth-service catalog-service order-service

# Help target
help:
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║       Microservices Build System - Available Commands      ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📋 Testing Commands:"
	@echo "  make test              Run all tests"
	@echo "  make test-verbose      Run all tests with verbose output"
	@echo "  make test-auth-service         Test auth service only"
	@echo "  make test-catalog-service      Test catalog service only"
	@echo "  make test-order-service        Test order service only"
	@echo ""
	@echo "🔨 Build Commands:"
	@echo "  make build             Build all services"
	@echo "  make build-clean       Clean and rebuild all services"
	@echo "  make build-auth-service        Build auth service only"
	@echo "  make build-catalog-service     Build catalog service only"
	@echo "  make build-order-service       Build order service only"
	@echo ""
	@echo "🚀 Distribution Commands:"
	@echo "  make dist              Build all services (alias: make build)"
	@echo "  make dist/start        Start all services from dist folder"
	@echo "  make dist/stop         Stop all running services"
	@echo "  make dist/logs         Watch logs from all services"
	@echo "  make dist/status       Show running services status"
	@echo ""
	@echo "🧹 Cleanup Commands:"
	@echo "  make clean             Remove dist folder and rebuild"
	@echo ""

# Testing targets
test:
	@echo "🧪 Running tests for all services..."
	@./test.sh

test-verbose:
	@echo "🧪 Running tests with verbose output..."
	@./test.sh -v

test-auth-service:
	@echo "🧪 Testing auth-service..."
	@./test.sh auth-service

test-catalog-service:
	@echo "🧪 Testing catalog-service..."
	@./test.sh catalog-service

test-order-service:
	@echo "🧪 Testing order-service..."
	@./test.sh order-service

# Build targets
build:
	@echo "🔨 Building all services..."
	@./build.sh

build-clean:
	@echo "🔨 Cleaning and building all services..."
	@./build.sh clean

build-auth-service:
	@echo "🔨 Building auth-service..."
	@./build.sh auth-service

build-catalog-service:
	@echo "🔨 Building catalog-service..."
	@./build.sh catalog-service

build-order-service:
	@echo "🔨 Building order-service..."
	@./build.sh order-service

# Dist alias
dist: build

# Distribution targets
dist/start:
	@if [ ! -f "$(DIST_DIR)/start_all.sh" ]; then \
		echo "❌ Build not found. Run 'make build' first"; exit 1; \
	fi
	@echo "🚀 Starting all services..."
	@cd $(DIST_DIR) && ./start_all.sh

dist/stop:
	@if [ ! -f "$(DIST_DIR)/stop_all.sh" ]; then \
		echo "❌ Build not found. Run 'make build' first"; exit 1; \
	fi
	@echo "⏹️  Stopping all services..."
	@cd $(DIST_DIR) && ./stop_all.sh

dist/logs:
	@echo "📋 Tailing logs from all services..."
	@if [ -d "$(DIST_DIR)/logs" ]; then \
		tail -f $(DIST_DIR)/logs/*.log 2>/dev/null || echo "No logs found yet"; \
	else \
		echo "Logs directory not found"; \
	fi

dist/status:
	@echo "📊 Service Status:"
	@pgrep -f "^\./auth$$" > /dev/null && echo "✓ Auth service running (PID: $$(pgrep -f '^\./auth$$'))" || echo "✗ Auth service not running"
	@pgrep -f "^\./catalog$$" > /dev/null && echo "✓ Catalog service running (PID: $$(pgrep -f '^\./catalog$$'))" || echo "✗ Catalog service not running"
	@pgrep -f "^\./order$$" > /dev/null && echo "✓ Order service running (PID: $$(pgrep -f '^\./order$$'))" || echo "✗ Order service not running"

# Clean target
clean:
	@echo "🧹 Removing dist folder..."
	@rm -rf $(DIST_DIR)
	@echo "✓ Cleaned"

# Development workflow targets
.PHONY: dev dev/build dev/test dev/start dev/restart

dev: dev/build dev/start
	@echo "✓ Development environment ready"

dev/build:
	@echo "🔨 Building for development..."
	@./build.sh clean

dev/test:
	@echo "🧪 Running tests..."
	@./test.sh

dev/start: | dev/build
	@echo "🚀 Starting services..."
	@cd $(DIST_DIR) && ./start_all.sh

dev/restart: | dev/stop dev/start
	@echo "♻️  Services restarted"

dev/stop:
	@cd $(DIST_DIR) && ./stop_all.sh 2>/dev/null || true

# Quick targets
.PHONY: quick-test quick-build quick-start quick-stop

quick-test:
	@go test ./... -timeout 10s -short

quick-build:
	@./build.sh

quick-start: | dist/start

quick-stop: | dist/stop
