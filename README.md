# Go HTTP Server

A clean, well-structured HTTP server application written in Go following best practices, clean code principles, and design patterns.

## Features

- **Clean Architecture**: Separation of concerns with handlers, services, and repositories
- **Dependency Injection**: Interfaces for testability and flexibility
- **Graceful Shutdown**: Proper resource cleanup on shutdown
- **Context Support**: All operations support context for cancellation and timeouts
- **Error Handling**: Proper error handling throughout the application
- **Configuration**: Environment variable-based configuration with sensible defaults
- **Comprehensive Tests**: Table-driven tests with proper isolation

## Project Structure

```
.
├── config/          # Configuration management
├── handler/         # HTTP handlers
├── service/         # Business logic layer
├── store/           # Data storage layer (interfaces and implementations)
├── persistence/    # File persistence layer
└── main.go         # Application entry point
```

## Installation

### Prerequisites

- Go 1.22 or later (required for latest features)
- Docker and Docker Compose (optional)

### Building

```bash
go build -o app
```

### Building Docker Container

```bash
make build
```

The Dockerfile uses a multi-stage build optimized for minimal memory and fastest boot:
- **Build stage**: Uses Go 1.22 with build cache mounts for faster builds
- **Final stage**: Minimal Alpine Linux image (~15MB) with only essential tools
- **Optimizations**: 
  - Build cache mounts for faster dependency downloads
  - Static binary compilation for portability
  - Stripped binaries (`-w -s`) for smaller size
  - Trimmed paths for reproducible builds
- **Security**: Runs as non-root user, minimal attack surface
- **Health checks**: Built-in `/health` endpoint for monitoring

## Configuration

The application can be configured using environment variables:

- `FILENAME`: Log file name (default: `timestamps.log`)
- `ADDRESS`: Server address (default: `localhost`)
- `ROUTE`: HTTP route path (default: `/`)
- `PORT`: Server port (default: `8000`)
- `THRESHOLD`: Timestamp expiration threshold in seconds (default: `60`)

### Example

```bash
export PORT=8080
export THRESHOLD=120
./app
```

## Running

### Local Development

```bash
go run main.go
```

### Docker

```bash
make build
```

The server will start at `http://localhost:8000` (or the configured port).

#### Docker Compose Configuration

The `docker-compose.yml` file supports environment variable configuration:

```bash
# Set custom port
export PORT=8080

# Set custom threshold
export THRESHOLD=120

# Build and run
make build
```

#### Docker Features

- **Multi-stage build**: Optimized image size (~15MB final image)
- **Build cache mounts**: Faster builds with cached dependencies
- **Resource limits**: Memory limited to 64MB, CPU to 0.5 cores
- **Non-root user**: Runs as `appuser` for security
- **Health checks**: Built-in `/health` endpoint for monitoring
- **Volume mounting**: Persists timestamp log file
- **Environment variables**: Configurable via docker-compose
- **Auto-restart**: Container restarts unless stopped
- **Read-only filesystem**: Enhanced security with tmpfs for writable areas
- **Fast boot**: Minimal image size and optimized startup

#### Go 1.22 Features Used

- **slices package**: Using `slices.Clip()` for efficient memory management
- **Improved for loops**: Better variable scoping in loops
- **Build optimizations**: `-trimpath` for reproducible builds
- **Performance**: Latest compiler optimizations

## Testing

### Run All Tests

```bash
make test
```

Or directly:

```bash
go test ./...
```

### Run Tests with Coverage

```bash
make test-coverage
```

This will run all tests and display coverage statistics.

### Generate HTML Coverage Report

```bash
make test-coverage-html
```

This generates a `coverage.html` file that you can open in your browser.

### View Coverage Report in Browser

```bash
make test-coverage-view
```

This generates the HTML report and automatically opens it in your default browser.

### Current Test Coverage

- **Overall**: 86.4% of statements
- **Service Layer**: 95.0% coverage
- **Store Layer**: 92.1% coverage
- **Persistence Layer**: 78.8% coverage

## Usage

### Record a Timestamp

```bash
curl http://localhost:8000/
```

Response:
```json
{"count": 1}
```

The server records the current timestamp and returns the count of valid (non-expired) timestamps.

## Architecture

### Design Patterns Used

1. **Dependency Injection**: All dependencies are injected through constructors
2. **Interface Segregation**: Small, focused interfaces for each responsibility
3. **Repository Pattern**: Store interface abstracts data access
4. **Service Layer**: Business logic separated from HTTP handlers
5. **Factory Pattern**: Constructor functions for creating instances

### Best Practices

- **Error Handling**: All errors are properly handled and returned
- **Context Usage**: All operations support context for cancellation
- **Resource Management**: Proper cleanup of resources (files, connections)
- **Testability**: Interfaces allow easy mocking for testing
- **Separation of Concerns**: Clear boundaries between layers
- **No Global State**: All state is managed through dependency injection

## Development

### Code Style

The project follows Go best practices:
- Use `gofmt` for formatting
- Follow Go naming conventions
- Write comprehensive tests
- Document exported functions and types

### Adding New Features

1. Define interfaces for new functionality
2. Implement the interfaces
3. Add tests with table-driven approach
4. Update documentation

## License

MIT
