FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-w -s' -o app ./cmd/server

FROM alpine:3.20
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=builder /build/app .
# data/ is the volume mount point; a fresh named volume inherits this ownership.
RUN mkdir data && chown -R appuser:appuser /app
USER appuser
EXPOSE 8000
# wget comes with busybox in the base image.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8000/health || exit 1
CMD ["./app"]
