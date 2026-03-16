# CLI Reference Guide

This guide provides comprehensive command-line interface (CLI) documentation for building, testing, and running the microservices project.

## 📁 Files Overview

### Created Files

| File | Purpose | Executable |
|------|---------|-----------|
| `test.sh` | Run tests for all or specific services | ✓ Yes |
| `build.sh` | Build all services into dist directory | ✓ Yes |
| `Makefile` | Make-based command interface | - |
| `SCRIPTS.md` | Detailed documentation of scripts | - |
| `QUICK_START.sh` | Quick reference guide | ✓ Yes |
| `CLI.md` | This file - CLI reference | - |

### Generated Files (in dist/)

```
dist/
├── auth               # Compiled auth service
├── catalog            # Compiled catalog service
├── order              # Compiled order service
├── start_all.sh       # Auto-generated startup script
├── stop_all.sh        # Auto-generated stop script
├── init_db.sh         # Auto-generated DB init script
├── README.md          # Distribution README
├── .env               # Environment config (copied)
└── data/              # SQLite databases directory
    ├── auth-service/
    ├── catalog-service/
    └── order-service/
```

---

## 🧪 Test Commands

### Basic Testing

```bash
# Run all tests
./test.sh

# Run with verbose output
./test.sh -v

# Run with coverage
./test.sh -cover

# Run with timeout
./test.sh -timeout 60s
```

### Service-Specific Testing

```bash
# Test auth service only
./test.sh auth-service

# Test catalog service only
./test.sh catalog-service

# Test order service only
./test.sh order-service

# Test with verbose and timeout
./test.sh auth-service -v -timeout 30s
```

### Make-Based Testing

```bash
make test                 # All tests
make test-verbose         # Verbose output
make test-auth-service    # Auth service only
make test-catalog-service # Catalog service only
make test-order-service   # Order service only
```

---

## 🔨 Build Commands

### Full Build

```bash
# Build all services to dist folder
./build.sh

# Clean and rebuild
./build.sh clean

# Build with Makefile
make build
make build-clean
```

### Service-Specific Builds

```bash
# Build single service
./build.sh auth-service
./build.sh catalog-service
./build.sh order-service

# Clean and build single service
./build.sh clean auth-service

# Make equivalents
make build-auth-service
make build-catalog-service
make build-order-service
```

### Build Output

After building, check the dist folder:

```bash
# List built binaries
ls -lh dist/

# Check database directories
ls -lh dist/data/*/

# View startup scripts
ls -lh dist/*.sh
```

---

## 🚀 Running Services

### From dist Folder

```bash
# Navigate to dist
cd dist

# Initialize databases (if needed)
./init_db.sh

# Start all services in background
./start_all.sh

# Stop all services
./stop_all.sh
```

### Using Make

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

### Verification Commands

```bash
# Check if services are running
ps aux | grep -E "auth|catalog|order"

# Check port usage
lsof -i :8081 -i :8082 -i :8083

# Test service endpoints
curl http://localhost:8081/
curl http://localhost:8082/
curl http://localhost:8083/

# Access Swagger documentation
# Auth: http://localhost:8081/swagger
# Catalog: http://localhost:8082/swagger
# Order: http://localhost:8083/swagger
```

---

## 📊 Common Workflows

### Development Workflow

```bash
# 1. Run tests
./test.sh

# 2. Build fresh binaries
./build.sh clean

# 3. Start services
cd dist && ./start_all.sh

# 4. Monitor logs
tail -f logs/*.log

# 5. Test API endpoints
curl http://localhost:8082/api/products

# 6. Stop services when done
./stop_all.sh
```

### CI/CD Workflow

```bash
# Run security and build checks
./test.sh -v

# Build for deployment
./build.sh

# Verify binaries
file dist/auth dist/catalog dist/order

# Package for deployment
tar -czf dist.tar.gz dist/
```

### Database Reset

```bash
# Option 1: Reset individual databases
rm dist/data/auth-service/auth.db
rm dist/data/catalog-service/catalog.db
rm dist/data/order-service/order.db

# Option 2: Reset all databases
rm -rf dist/data/*
cd dist && ./init_db.sh

# Option 3: Complete clean rebuild
./build.sh clean
cd dist && ./start_all.sh
```

### Quick Restart

```bash
# Option 1: Using make
make dev/restart

# Option 2: Using scripts
cd dist && ./stop_all.sh && sleep 2 && ./start_all.sh

# Option 3: Kill and restart
pkill -f "^\./auth$|^\./catalog$|^\./order$"
cd dist && ./start_all.sh
```

---

## 📈 Advanced Usage

### Custom Test Timeouts

```bash
# Increase timeout for slow tests
./test.sh -timeout 120s

# Or for single service
./test.sh order-service -timeout 60s
```

### Parallel Test Execution

```bash
# Run tests in parallel with fail-fast
go test -race -parallel 4 ./...

# Or use the script with Go flags
./test.sh -race -parallel 2
```

### Conditional Builds

```bash
# Get build info
go version
go env GOROOT

# Build for specific OS/arch (if needed)
# GOOS=linux GOARCH=amd64 go build -o dist/auth-linux services/auth-service/cmd/server/main.go

# Check binary info
file dist/auth dist/catalog dist/order
```

### Environment Configuration

```bash
# Edit .env before building
cat > .env << EOF
AUTH_SERVICE_PORT=8081
CATALOG_SERVICE_PORT=8082
ORDER_SERVICE_PORT=8083
EOF

# Build will copy .env to dist/
./build.sh

# Override in dist/
cd dist
cp .env.example .env
# Edit .env
./start_all.sh
```

### Log Analysis

```bash
# View auth service logs
cat dist/logs/auth.log

# Watch catalog service in real-time
tail -f dist/logs/catalog.log

# Search for errors
grep "ERROR" dist/logs/*.log

# Extract specific errors
grep -A 5 "panic" dist/logs/*.log
```

---

## 🛠️ Makefile Quick Reference

```bash
# Show all available targets
make help

# Testing
make test                 # Run all tests
make test-verbose         # Verbose test output
make test-auth-service    # Test single service

# Building
make build                # Build all services
make build-clean          # Clean build
make dist                 # Alias for build

# Distribution
make dist/start           # Start services
make dist/stop            # Stop services
make dist/logs            # Watch logs
make dist/status          # Check status

# Development
make dev                  # Full dev setup
make dev/build            # Dev build
make dev/test             # Run tests
make dev/start            # Start for dev
make dev/stop             # Stop services
make dev/restart          # Restart services

# Cleanup
make clean                # Remove dist folder
```

---

## ⚡ Shortcut Aliases

Create shell aliases for common commands:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias mtest='./test.sh'
alias mtestv='./test.sh -v'
alias mbuild='./build.sh'
alias mbuild-clean='./build.sh clean'
alias mstart='cd dist && ./start_all.sh'
alias mstop='cd dist && ./stop_all.sh'
alias mlogs='tail -f dist/logs/*.log'
alias mclean='rm -rf dist'

# Then use:
mtest                   # Run all tests
mtestv                  # Verbose tests
mbuild                  # Build services
mstart                  # Start services
mstop                   # Stop services
mlogs                   # Watch logs
mclean                  # Clean dist
```

---

## 🐛 Troubleshooting

### Build Failures

```bash
# Clean Go cache
go clean -cache

# Update dependencies
go mod tidy

# Rebuild from scratch
./build.sh clean
```

### Test Failures

```bash
# Run with verbose output
./test.sh -v

# Check for race conditions
./test.sh -race

# Increase timeout
./test.sh -timeout 120s
```

### Service Won't Start

```bash
# Check if ports are available
lsof -i :8081

# View service logs
cat dist/logs/auth.log

# Check .env configuration
cat dist/.env

# Try manual start
cd dist && ./auth
```

### Database Issues

```bash
# Check database exists
ls -lh dist/data/*/

# Query database directly
sqlite3 dist/data/auth-service/auth.db ".tables"

# Reset and reinitialize
rm -rf dist/data/*
cd dist && ./init_db.sh && ./start_all.sh
```

---

## 📚 Additional Resources

For detailed information, see:

- **SCRIPTS.md** - Full documentation of test.sh and build.sh
- **dist/README.md** - Distribution-specific documentation
- **Makefile** - All available make targets
- **QUICK_START.sh** - Quick reference guide (run to display)

Run any of these:

```bash
cat SCRIPTS.md           # Detailed script documentation
cat dist/README.md       # Distribution guide
make help                # Makefile help
./QUICK_START.sh         # Quick reference
```

---

## 💡 Best Practices

1. **Always test before building:** `./test.sh && ./build.sh`
2. **Use clean builds for deployment:** `./build.sh clean`
3. **Monitor logs during development:** `tail -f dist/logs/*.log`
4. **Check service status before making requests:** `make dist/status`
5. **Backup databases before major changes:** `cp -r dist/data dist/data.backup`
6. **Document any custom environment variables:** Edit `.env` and commit to git
7. **Use make targets** instead of scripts directly for consistency
8. **Review SCRIPTS.md** for comprehensive documentation

---

## 📋 Summary

| Task | Command |
|------|---------|
| Test everything | `./test.sh` |
| Build everything | `./build.sh` |
| Start services | `make dist/start` |
| Stop services | `make dist/stop` |
| View logs | `tail -f dist/logs/*.log` |
| Check status | `make dist/status` |
| Full setup | `make dev` |
| Reset & rebuild | `make build-clean` |
| Show help | `make help` |

Happy building! 🚀
