# Go HTTP Server

A sliding-window request counter. Each `GET` records the current time and returns how many requests arrived within the last `THRESHOLD` seconds. Timestamps are persisted to a file, so the count survives restarts.

Standard library only, layered as domain / application / infrastructure / presentation.

## Project Structure
```text
cmd/server/                  # Entry point, wiring, graceful shutdown
internal/
├── config/                  # Environment configuration
├── domain/                  # Repository interface
├── application/             # TimestampService (use cases)
├── infrastructure/
│   ├── persistence/         # Atomic file read/write
│   └── repository/          # In-memory store backed by the file
└── presentation/http/       # HTTP handlers and routes
```

## Running

Everything runs in Docker. `make help` lists all targets.

```bash
make build    # build and start
make logs     # follow logs
make test     # go vet + go test -race in a container
make down     # stop and remove
```

## Configuration
Environment variables (set them in your shell or a `.env` file for `docker compose`):
- `PORT`: Host port (default: `8000`)
- `ROUTE`: Path of the counter endpoint (default: `/`)
- `THRESHOLD`: Window size in seconds (default: `60`)
- `FILENAME`: Storage file (default: `timestamps.log`; in Docker it lives on the `data` volume at `/app/data/timestamps.log`)

## Usage
```bash
curl http://localhost:8000/
# {"count":1}
```
`GET /health` returns `OK`. Other methods on the counter route return `405`.

## License
MIT
