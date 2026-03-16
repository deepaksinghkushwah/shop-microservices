# Build and Test Scripts

This document explains how to use the build and test scripts for the microservices project.

## 📋 Quick Start

### Run All Tests
```bash
./test.sh
```

### Build All Services
```bash
./build.sh
```

### Start All Services
```bash
cd dist
./start_all.sh
```

## 🧪 Test Script (`test.sh`)

Runs unit tests and integration tests for all services or a specific service.

### Usage

```bash
./test.sh [SERVICE_NAME] [TEST_ARGS]
```

### Examples

```bash
# Run all tests
./test.sh

# Run tests with verbose output
./test.sh -v

# Run tests with coverage
./test.sh -cover

# Run tests for specific service
./test.sh auth-service

# Run auth-service tests with verbose output
./test.sh auth-service -v

# Run with timeout
./test.sh -timeout 60s
```

### Available Services
- `auth-service` - Authentication microservice
- `catalog-service` - Product catalog microservice  
- `order-service` - Order management microservice

### Features
- Colored output for easy visibility
- Test summary at the end
- Error reporting with service names
- Compatible with all Go test flags

### Output
```
========================================
Running tests for auth-service
========================================
...test output...
✓ auth-service tests passed

========================================
Running tests for catalog-service
========================================
...test output...
✓ catalog-service tests passed

========================================
Running tests for order-service
========================================
...test output...
✓ order-service tests passed

========================================
Test Summary
========================================
Total services tested: 3
✓ All tests passed!
```

---

## 🔨 Build Script (`build.sh`)

Compiles all services into executable binaries and creates a distribution package with:
- Service binaries
- SQLite databases with migrations
- Startup/shutdown scripts
- API documentation
- Logging infrastructure

### Usage

```bash
./build.sh [clean] [SERVICE_NAME]
```

### Examples

```bash
# Build all services to dist folder
./build.sh

# Clean and rebuild all services
./build.sh clean

# Build only auth-service
./build.sh auth-service

# Clean and build only catalog-service
./build.sh clean catalog-service
```

### Build Output Structure

```
dist/
├── auth                      # Auth service binary
├── catalog                   # Catalog service binary
├── order                      # Order service binary
├── data/                      # Directory for SQLite databases
│   ├── auth-service/auth.db
│   ├── catalog-service/catalog.db
│   └── order-service/order.db
├── logs/                      # Service logs directory
├── auth_docs/                 # Auth service Swagger docs
├── catalog_docs/              # Catalog service Swagger docs
├── order_docs/                # Order service Swagger docs
├── .env                       # Environment configuration (copied from root)
├── .env.example               # Example environment file
├── README.md                  # Distribution README
├── start_all.sh              # Start all services script
├── stop_all.sh               # Stop all services script
└── init_db.sh                # Database initialization script
```

### Build Steps

1. **Clean (if requested)** - Removes old dist folder
2. **Create directories** - Sets up dist structure
3. **Copy configuration** - Copies .env file
4. **Build binaries** - Compiles each service
5. **Copy documentation** - Includes Swagger specs
6. **Create helper scripts** - start_all.sh, stop_all.sh, init_db.sh
7. **Create README** - Distribution README with instructions

### Environment Setup

After building, configure the environment:

```bash
cd dist
cp .env.example .env
# Edit .env with your configuration
```

### Features
- Multi-service support
- Clean builds with `clean` flag
- Service-specific builds
- Automatic database setup
- Pre-configured startup scripts
- Comprehensive error reporting
- Colored output for clarity

---

## 📦 Using dist/ Directory

After building with `./build.sh`, use the `dist/` directory to run your services.

### Starting Services

```bash
cd dist
./start_all.sh
```

This will:
1. Initialize databases if needed
2. Start each service as a background process
3. Log output to `logs/` directory
4. Display service URLs

### Stopping Services

```bash
cd dist
./stop_all.sh
```

Or manually:
```bash
# Kill specific service
pkill -f "^\./auth$"
pkill -f "^\./catalog$"
pkill -f "^\./order$"

# Or use foreground and Ctrl+C
```

### Accessing Services

After starting, services are available at:

- **Auth Service** - http://localhost:8081
  - Swagger UI: http://localhost:8081/swagger
  
- **Catalog Service** - http://localhost:8082
  - Swagger UI: http://localhost:8082/swagger
  
- **Order Service** - http://localhost:8083
  - Swagger UI: http://localhost:8083/swagger

### Viewing Logs

```bash
# Watch all logs
./dist/logs

# Or use tail
tail -f dist/logs/auth.log
tail -f dist/logs/catalog.log
tail -f dist/logs/order.log

# View specific log
cat dist/logs/auth.log
```

### Database Files

SQLite database files are stored in `dist/data/`:

```bash
# Check database
ls -lh dist/data/*/
sqlite3 dist/data/auth-service/auth.db ".tables"
sqlite3 dist/data/catalog-service/catalog.db ".tables"
sqlite3 dist/data/order-service/order.db ".tables"
```

---

## 🛠️ Makefile Targets

As an alternative to running scripts directly, use the Makefile:

### Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Test specific service
make test-auth-service
make test-catalog-service
make test-order-service
```

### Building

```bash
# Build all
make build
make dist        # Alias for build

# Build with clean
make build-clean

# Build specific service
make build-auth-service
make build-catalog-service
make build-order-service
```

### Distribution Management

```bash
# Start services
make dist/start

# Stop services
make dist/stop

# View logs
make dist/logs

# Check status
make dist/status
```

### Development Workflow

```bash
# Full dev setup: build, configure, start
make dev

# Just build for dev
make dev/build

# Run tests
make dev/test

# Start services
make dev/start

# Stop services
make dev/stop

# Restart services
make dev/restart
```

### Cleanup

```bash
# Remove dist folder
make clean
```

### Show Help

```bash
make help
```

---

## 🚀 Complete Workflow Example

### 1. Run Tests

```bash
./test.sh -v
# or
make test
```

### 2. Build Services

```bash
./build.sh clean
# or
make build-clean
```

### 3. Start Services

```bash
cd dist
./start_all.sh
# or from root: make dist/start
```

### 4. Test Endpoints

```bash
# Auth Service
curl http://localhost:8081/swagger

# Catalog Service
curl http://localhost:8082/swagger

# Order Service
curl http://localhost:8083/swagger
```

### 5. View Logs

```bash
tail -f dist/logs/*.log
# or
make dist/logs
```

### 6. Stop Services

```bash
cd dist
./stop_all.sh
# or from root: make dist/stop
```

---

## 🔧 Configuration

### Environment Variables

Create or edit `.env` in the root directory:

```bash
# Service ports
AUTH_SERVICE_PORT=8081
CATALOG_SERVICE_PORT=8082
ORDER_SERVICE_PORT=8083

# Optional: Database configuration
# DATABASE_HOST=localhost
# DATABASE_PORT=5432
# DATABASE_USER=postgres
# DATABASE_PASSWORD=password
```

This will be automatically copied to `dist/.env` during build.

### Modifying Build Output

To change the dist directory location, edit the scripts:

**test.sh:**
```bash
DIST_DIR="$SCRIPT_DIR/dist"  # Change this path
```

**build.sh:**
```bash
DIST_DIR="$SCRIPT_DIR/dist"  # Change this path
```

---

## 📊 Script Output Examples

### Test Output

```
========================================
Running tests for auth-service
========================================
=== RUN   TestService
--- PASS: TestService (0.001s)
ok      github.com/deepaksinghkushwah/shop-microservices/services/auth-service/...   0.005s
✓ auth-service tests passed

========================================
Test Summary
========================================
Total services tested: 3
✓ All tests passed!
```

### Build Output

```
========================================
Setting up dist directory
========================================
✓ Created dist directory at /path/to/dist
✓ Created data directory for databases

========================================
Building auth-service
========================================
✓ Built auth-service binary: /path/to/dist/auth
✓ Copied API docs for auth-service
✓ Created data directory for auth-service

... (catalog and order services) ...

========================================
Build Summary
========================================
Successfully built:
  ✓ auth-service
  ✓ catalog-service
  ✓ order-service

✓ All services built successfully!

Next steps:
  1. Configure environment: cd /path/to/dist && cp .env.example .env
  2. Start services: cd /path/to/dist && ./start_all.sh
  3. View logs: tail -f /path/to/dist/logs/*.log
```

---

## ⚠️ Troubleshooting

### Tests Fail

1. Ensure all dependencies are installed:
   ```bash
   go mod tidy
   ```

2. Check Go version (requires Go 1.20+):
   ```bash
   go version
   ```

3. Run with verbose output:
   ```bash
   ./test.sh -v
   ```

### Build Fails

1. Check Go is installed:
   ```bash
   go version
   ```

2. Try clean rebuild:
   ```bash
   ./build.sh clean
   ```

3. Check for compilation errors:
   ```bash
   go build ./...
   ```

### Services Won't Start

1. Check if ports are in use:
   ```bash
   lsof -i :8081 -i :8082 -i :8083
   ```

2. Change ports in `.env`:
   ```bash
   AUTH_SERVICE_PORT=8091
   CATALOG_SERVICE_PORT=8092
   ORDER_SERVICE_PORT=8093
   ```

3. Check logs:
   ```bash
   cat dist/logs/*.log
   tail -f dist/logs/*.log
   ```

### Database Issues

1. Reset databases:
   ```bash
   rm -rf dist/data/*
   cd dist && ./init_db.sh
   ```

2. Verify database access:
   ```bash
   sqlite3 dist/data/auth-service/auth.db ".tables"
   ```

---

## 📝 Notes

- Tests have a default 30-second timeout. Modify with: `./test.sh -timeout 60s`
- Services are built as static binaries in `dist/` - they can be moved or deployed
- Logs are stored in `dist/logs/` - check them for service errors
- Databases are SQLite files - easily backup by copying `dist/data/` directory
- The `.env` file in `dist/` overrides root `.env` - edit it after building

---

## 🔗 Related Files

- `test.sh` - Test execution script
- `build.sh` - Build and distribution script
- `Makefile` - Alternative command interface
- `dist/README.md` - Distribution-specific documentation
- `dist/start_all.sh` - Service startup script (auto-generated)
- `dist/stop_all.sh` - Service shutdown script (auto-generated)
