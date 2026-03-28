# Docker Deployment Guide

This guide covers deploying the microservices using Docker and Docker Compose for a true microservice architecture.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- At least 2GB RAM available
- Ports 8081, 8082, 8083, and 5432 available

## Quick Start

### 1. Build Docker Images

```bash
# Build all Docker images
./build.sh docker

# Build specific service image
./build.sh docker auth-service

# Clean and rebuild (remove old images)
./build.sh docker clean
```

### 2. Start Services

```bash
# Start all services with Docker Compose
./docker-start.sh

# Or use docker-compose directly
docker-compose up -d
```

### 3. Stop Services

```bash
# Stop all services
./docker-stop.sh

# Or use docker-compose directly
docker-compose down
```

### 4. View Logs

```bash
# View logs from all services
docker-compose logs -f

# View logs from specific service
docker-compose logs -f auth-service

# Or use the utility script
./docker-util.sh logs
./docker-util.sh logs catalog-service
```

## Docker Architecture

```
┌─────────────────────────────────────────────────┐
│              Docker Network                      │
│              (microservices)                      │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌────────────────────────────────────┐         │
│  │   PostgreSQL Container             │         │
│  │   @postgres:5432                   │         │
│  │   (Volume: postgres_data)          │         │
│  │                                    │         │
│  │  - auth (database)                 │         │
│  │  - catalog (database)              │         │
│  │  - order (database)                │         │
│  └────────────────────────────────────┘         │
│                    ▲                             │
│         ┌──────────┼──────────┐                 │
│         │          │          │                 │
│  ┌──────▼────┐ ┌───▼──────┐ ┌─▼──────────┐    │
│  │Auth       │ │Catalog   │ │Order       │    │
│  │Service    │ │Service   │ │Service     │    │
│  │:8081      │ │:8082     │ │:8083       │    │
│  └───────────┘ └──────────┘ └────────────┘    │
│                                                  │
└─────────────────────────────────────────────────┘
         │              │              │
         └──────────────┼──────────────┘
                        │
                   Host Machine
              (Ports 8081, 8082, 8083)
```

## Service Details

### Auth Service
- **Container Name**: shop-auth-service
- **Port**: 8081
- **Database**: auth
- **Health Check**: ✓ Every 30s
- **Dependencies**: PostgreSQL
- **URL**: http://localhost:8081
- **Swagger**: http://localhost:8081/swagger/index.html

### Catalog Service
- **Container Name**: shop-catalog-service
- **Port**: 8082
- **Database**: catalog
- **Health Check**: ✓ Every 30s
- **Dependencies**: PostgreSQL
- **URL**: http://localhost:8082
- **Swagger**: http://localhost:8082/swagger/index.html

### Order Service
- **Container Name**: shop-order-service
- **Port**: 8083
- **Database**: order
- **Health Check**: ✓ Every 30s
- **Dependencies**: PostgreSQL
- **URL**: http://localhost:8083
- **Swagger**: http://localhost:8083/swagger/index.html

### PostgreSQL Database
- **Container Name**: shop-postgres
- **Port**: 5432
- **User**: postgres
- **Password**: password (change in .env for production)
- **Volume**: postgres_data
- **Health Check**: ✓ Every 10s

## Utility Commands

```bash
# Show container status
./docker-util.sh status

# View service logs
./docker-util.sh logs
./docker-util.sh logs auth-service

# Rebuild images
./docker-util.sh rebuild

# Open shell in container
./docker-util.sh shell
./docker-util.sh shell catalog-service

# Execute command in container
./docker-util.sh exec auth-service psql

# Clean up (remove containers and volumes)
./docker-util.sh clean
```

## Database Initialization

The PostgreSQL database is automatically initialized with the following databases:
- `auth` - for Auth Service
- `catalog` - for Catalog Service
- `order` - for Order Service

Initialization is done through the `init-db.sql` script which runs when the PostgreSQL container starts.

## Configuration

### Environment Variables

Edit `.env` to configure:

```env
# Service ports
AUTH_SERVICE_PORT=8081
CATALOG_SERVICE_PORT=8082
ORDER_SERVICE_PORT=8083

# PostgreSQL
DB_HOST=postgres        # Docker network name
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password    # CHANGE FOR PRODUCTION
DB_SSLMODE=disable
```

### Service-Specific Configuration

Each service can be configured via Docker Compose environment variables in `docker-compose.yml`.

## Troubleshooting

### Containers won't start

```bash
# Check logs
docker-compose logs

# Ensure ports are not in use
lsof -i :8081
lsof -i :8082
lsof -i :8083
lsof -i :5432
```

### Database connection issues

```bash
# Check PostgreSQL container
docker-compose ps postgres

# Connect to PostgreSQL directly
docker-compose exec postgres psql -U postgres -d postgres

# View databases
\l

# Exit psql
\q
```

### Health check failing

```bash
# Check service logs
docker-compose logs auth-service

# Monitor container health
watch -n 1 'docker-compose ps'
```

## Production Deployment

For production:

1. **Change database password** in `.env`
2. **Use secrets** instead of plain text passwords
3. **Configure SSL/TLS** for PostgreSQL
4. **Use external database** instead of containerized
5. **Enable resource limits** in docker-compose.yml
6. **Use health checks** (already configured)
7. **Set up monitoring** and logging
8. **Use Docker registries** for image management
9. **Configure auto-restart** policies (already set)
10. **Implement backup strategy** for PostgreSQL volume

Example production environment:

```env
JWT_SECRET="your-secure-random-secret-key"
DB_PASSWORD="your-secure-database-password"
DB_HOST=prod-postgres.example.com
DB_SSLMODE=require
```

## Performance Tuning

### Enable Compose Caching

```bash
DOCKER_BUILDKIT=1 docker-compose build
```

### Monitor Resource Usage

```bash
docker stats
```

### Increase Container Resources

Edit `docker-compose.yml`:

```yaml
services:
  auth-service:
    # ... other config
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## Cleanup

```bash
# Stop and remove containers
docker-compose down

# Remove with volumes (data is deleted)
docker-compose down -v

# Remove unused images
docker image prune

# Complete cleanup
docker system prune -a
```

## Next Steps

- Configure monitoring and logging
- Set up CI/CD pipelines
- Deploy to production infrastructure
- Configure load balancing
- Set up API gateway
