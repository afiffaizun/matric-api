# matric-api

REST API sederhana untuk manajemen matrix, dibangun dengan Go dan Gin framework.

## Tech Stack

- **Go** (1.25+)
- **Gin** (HTTP framework)
- **Prometheus** (metrics)
- **Docker** / **Docker Compose**
- **Kubernetes** (k3s)
- **Helm** (chart)

## Struktur Project

```
.
├── cmd/server/main.go                 # Entry point server
├── internal/
│   ├── handler/handlers.go            # HTTP handlers
│   └── middleware/metrics.go          # Prometheus metrics middleware
├── helm/matric-api/                   # Helm chart
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── _helpers.tpl
│       ├── configmap.yaml
│       ├── deployment.yaml
│       ├── httproute.yaml
│       ├── ingress.yaml
│       ├── service.yaml
│       ├── servicemonitor.yaml
│       └── NOTES.txt
├── docker-compose.yml                 # Docker Compose config
├── Dockerfile                         # Multi-stage build
├── go.mod
├── go.sum
└── .gitignore
```

## API Endpoints

| Method | Endpoint       | Description              |
|--------|----------------|--------------------------|
| GET    | /              | Health check             |
| GET    | /health        | Health check             |
| GET    | /api/matrix    | Get matrix (TODO)        |
| POST   | /api/matrix    | Create matrix            |
| GET    | /metrics       | Prometheus metrics       |

## Cara Jalankan

### Local

```bash
go run cmd/server/main.go
```

### Docker Compose

```bash
docker compose up -d
```

### Docker

```bash
docker build -t matric-api:latest .
docker run -p 8080:8080 matric-api:latest
```

### Helm (Kubernetes)

```bash
# Install
helm install matric-api ./helm/matric-api -n matric-api --create-namespace

# Upgrade
helm upgrade matric-api ./helm/matric-api -n matric-api

# Uninstall
helm uninstall matric-api -n matric-api
```

## Akses via Kubernetes

| Service         | Type       | Port              |
|-----------------|------------|-------------------|
| matric-api      | ClusterIP  | 80                |
| matric-api-nodeport | NodePort | 8080:30080     |
| ingress         | matric-api.local | 80          |

```bash
# Port-forward
kubectl -n matric-api port-forward svc/matric-api 8080:80

# Test
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## Prometheus Monitoring

Prometheus metrics tersedia di endpoint `/metrics`:

```bash
curl http://localhost:8080/metrics | head -20
```

Untuk auto-discovery di Kubernetes, Helm chart menyertakan **ServiceMonitor** (perlu Prometheus Operator):

```yaml
# Diaktifkan via values.yaml
serviceMonitor:
  enabled: true
  interval: 15s
  path: /metrics
```