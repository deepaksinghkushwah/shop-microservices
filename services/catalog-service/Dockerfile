# Multi-stage build for catalog-service
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Set GOTOOLCHAIN to auto to allow downloading newer modules
ENV GOTOOLCHAIN=auto

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod files from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code - copy entire project
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o catalog-service ./services/catalog-service/cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /build/catalog-service .

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/health || exit 1

EXPOSE 8082

CMD ["./catalog-service"]
