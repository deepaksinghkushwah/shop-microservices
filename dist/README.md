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
