# mini-redis-go

A minimal Redis-like persistent key-value store written in Go, with a simple REST API.

---

## Features

- **Persistent Storage:**  
  Data is stored in a log file and replayed on startup for durability.

- **REST API:**  
  - `GET /kv/{key}`: Retrieve a value by key  
  - `POST /kv/{key}`: Set a value for a key  
  - `DELETE /kv/{key}`: Delete a key

- **Concurrency Safe:**  
  Uses `sync.RWMutex` for safe concurrent access.

- **Structured Logging:**  
  Uses Go's `log/slog` for structured logs.

- **Unit & Integration Tests:**  
  Includes tests for storage and API layers.

---

## Getting Started

### Prerequisites

- Go 1.22 or newer recommended

### Build & Run

```sh
go build -o mini-redis-go ./cmd
./mini-redis-go
```

The server listens on port `9686` by default.

---

## API Usage

### Set a Key

```sh
curl -X POST http://localhost:9686/kv/mykey -d 'myvalue'
```

### Get a Key

```sh
curl http://localhost:9686/kv/mykey
```

### Delete a Key

```sh
curl -X DELETE http://localhost:9686/kv/mykey
```

---

## Testing

Run all tests:

```sh
go test ./...
```

Run with race detector:

```sh
go test -race ./...
```

---

## Project Structure

```
cmd/                # Main server and API handlers
internal/storage/   # Persistent storage implementation
```

---

## Possible Enhancements (TBD)

- Add TTL (expiration) support for keys
- Return 404 for missing keys instead of 500
- Add endpoint to list all keys
- Add authentication
- Docker support
