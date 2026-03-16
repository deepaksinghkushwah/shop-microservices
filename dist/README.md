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
