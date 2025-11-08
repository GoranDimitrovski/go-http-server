# Build stage - optimized for cache and speed
FROM golang:1.22-alpine AS builder

# Install minimal build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum* ./

# Download dependencies with build cache mount for faster builds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build static binary with optimizations
# -ldflags='-w -s' strips debug info and symbol table
# -trimpath removes file system paths from the binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o app .

# Final stage - minimal alpine with only essential tools
# This provides the smallest size with health check capability
FROM alpine:latest

# Install only wget for health checks (minimal size)
RUN apk add --no-cache ca-certificates wget && \
    rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy the binary
COPY --from=builder /build/app .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

# Health check using health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# Run the application
CMD ["./app"]
