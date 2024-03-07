# FROM golang:1.22 as builder

# WORKDIR /app
# COPY go.mod .

# RUN go mod download
# COPY . .
# ENV GOCACHE=/root/.cache/go-build
# RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app

# FROM ubuntu:22.04

# RUN mkdir /app
# WORKDIR /app
# COPY --from=builder /app/app .
# ENTRYPOINT ["./app"] 




# FROM golang:1.22

# WORKDIR /app

# # Install Go dependencies
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     gcc \
#     && rm -rf /var/lib/apt/lists/*

# # Copy the source code
# COPY . .

# # Run tests
# CMD ["go", "test", "-v", "./..."]


FROM golang:1.22

WORKDIR /app

# Install Go dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy the source code
COPY . .

# Build the application
RUN go build -o app

# Run the application
CMD ["./app"]
