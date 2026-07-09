# matric-api

REST API sederhana untuk manajemen matrix, dibangun dengan Go dan Gin framework.

## Tech Stack

- **Go** (1.25+)
- **Gin** (HTTP framework)
- **Docker** / **Docker Compose**
- **Kubernetes** (k3s)
- **Kustomize** (opsional)

## Struktur Project

```
.
├── cmd/server/main.go            # Entry point server
├── internal/handler/handlers.go  # HTTP handlers
├── k8s/
│   ├── namespace.yaml            # Namespace matric-api
│   ├── configmap.yaml            # ConfigMap environment
│   ├── deployment.yaml           # Deployment 3 replica
│   ├── service.yaml              # Service ClusterIP
│   ├── service-nodeport.yaml     # Service NodePort :30080
│   └── ingress.yaml              # Ingress matric-api.local
├── docker-compose.yml            # Docker Compose config
├── Dockerfile                    # Multi-stage build
├── go.mod
└── go.sum
```

## API Endpoints

| Method | Endpoint       | Description       |
|--------|----------------|-------------------|
| GET    | /              | Health check      |
| GET    | /health        | Health check      |
| GET    | /api/matrix    | Get matrix (TODO) |
| POST   | /api/matrix    | Create matrix     |

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

### Kubernetes (k3s)

```bash
# Load image ke containerd k3s
docker save matric-api:latest | sudo k3s ctr images import -

# Apply semua resource
kubectl apply -f k8s/namespace.yaml
sleep 3
kubectl apply -f k8s/
```

### Test

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
```

## Akses via Kubernetes

| Service         | Type       | Port              |
|-----------------|------------|-------------------|
| matric-api      | ClusterIP  | 80                |
| matric-api-nodeport | NodePort | 8080:30080     |
| ingress         | matric-api.local | 80          |

```bash
# Via NodePort
curl http://localhost:30080/

# Via Ingress
curl -H "Host: matric-api.local" http://192.168.100.16/
```