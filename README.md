# Go HTTP Server

A clean, well-structured HTTP server application written in Go following Layered Architecture and Domain-Driven Design (DDD).

## Features
- **Layered Architecture (DDD)**: Clear separation between domain, application, infrastructure, and presentation layers.
- **Dockerized Environment**: All commands and application runtime run fully inside Docker containers.
- **Graceful Shutdown**: Proper resource cleanup on shutdown.
- **Dependency Injection**: Interfaces for testability and flexibility.

## Project Structure
```text
.
├── cmd/                      # Application entry points
│   └── server/               # Main executable
├── internal/                 # Internal modules
│   ├── application/          # Application services and use cases
│   ├── config/               # Configuration management
│   ├── domain/               # Domain models, rules, and repository interfaces
│   ├── infrastructure/       # Implementations for external interactions (e.g., storage)
│   │   ├── persistence/      # File persistence implementation
│   │   └── repository/       # Memory and filesystem data repositories
│   └── presentation/         # Protocol-specific handlers (REST API)
│       └── http/             # HTTP Handlers
```

## Running the Application

All operations are designed to be run within Docker via Make commands.

### Build and Start
Start the server and its dependencies:
```bash
make app.build
```

### Stopping
Stop and remove containers:
```bash
make app.down
```

### Viewing Logs
```bash
make app.logs
```

### Testing
Run all unit and integration tests inside the Docker container:
```bash
make app.test
```

## Configuration
The application is configured through environment variables (or `docker-compose.yml`):
- `PORT`: Server port (default: `8000`)
- `FILENAME`: Log file name (default: `timestamps.log`)
- `THRESHOLD`: Timestamp expiration threshold in seconds (default: `60`)

## Usage
To record a timestamp, send a GET request:
```bash
curl http://localhost:8000/
```

**Response**: (Returns the total count of non-expired timestamps)
```json
{"count": 1}
```

## License
MIT
