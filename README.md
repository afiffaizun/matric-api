# matric-api

REST API sederhana untuk manajemen matrix, dibangun dengan Go dan Gin framework.

## Tech Stack

- **Go** (1.25+)
- **Gin** (HTTP framework)
- **Docker** (multi-stage build)

## Struktur Project

```
.
├── cmd/server/main.go           # Entry point server
├── internal/handler/handlers.go # HTTP handlers
├── Dockerfile                   # Docker build config
├── go.mod
└── go.sum
```

## API Endpoints

| Method | Endpoint       | Description       |
|--------|----------------|-------------------|
| GET    | /              | Health check      |
| GET    | /api/matrix    | Get matrix (TODO) |
| POST   | /api/matrix    | Create matrix     |

## Cara Jalankan

### Local

```bash
go run cmd/server/main.go
```

### Docker

```bash
docker build -t matric-api:latest .
docker run -p 8080:8080 matric-api:latest
```

### Test

```bash
curl http://localhost:8080/
```
