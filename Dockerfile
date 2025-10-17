# Build stage
FROM golang:1.24.2-alpine3.20 AS builder

RUN apk update && apk add --no-cache \
    git \
    openssh \
    tzdata \
    build-base \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o task-queue \
    main.go

# Runtime stage with Redis
FROM alpine:latest

# Install runtime dependencies including Redis
RUN apk update && apk upgrade && \
    apk --update --no-cache add \
    tzdata \
    ca-certificates \
    curl \
    redis && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /app/task-queue /app/task-queue

# Create necessary directories
RUN mkdir -p /data/images /data/predicted /data/redis && \
    chown -R appuser:appuser /data

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'redis-server --daemonize yes --appendonly yes --dir /data/redis' >> /app/start.sh && \
    echo 'exec /app/task-queue' >> /app/start.sh && \
    chmod +x /app/start.sh && \
    chown appuser:appuser /app/start.sh

USER appuser

EXPOSE 8001 6379

ENTRYPOINT ["/app/start.sh"]