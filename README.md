http://localhost:8081/swagger/index.html : Auth Service
http://localhost:8082/swagger/index.html : Catalog Service
http://localhost:8083/swagger/index.html : Catalog Service

cd services/catalog-service
swag init -g cmd/server/main.go

cd services/auth-service
swag init -g cmd/server/main.go

cd services/order-service
swag init -g cmd/server/main.go


docker compose up --build -d