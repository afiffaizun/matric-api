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
│   ├── dashboards/
│   │   └── matric-api-overview.json   # Grafana dashboard
│   └── templates/
│       ├── _helpers.tpl
│       ├── configmap.yaml
│       ├── dashboard.yaml             # Dashboard provisioning
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

## Grafana

### Install (jika belum ada)

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

### Akses Grafana

```bash
# Port-forward ke Grafana
kubectl -n monitoring port-forward svc/monitoring-grafana 3000:80
```

Buka `http://localhost:3000` di browser.

**Login:**
- Username: `admin`
- Password:

```bash
kubectl -n monitoring get secret monitoring-grafana -o jsonpath="{.data.admin-password}" | base64 -d
```

## Grafana Dashboard

Dashboard **matric-api Overview** tersedia di `helm/matric-api/dashboards/matric-api-overview.json`.

### Import manual
1. Buka Grafana → ☰ → **Dashboards** → **New** → **Import**
2. Upload file `helm/matric-api/dashboards/matric-api-overview.json`
3. Pilih datasource **Prometheus** → **Import**

### Auto-provisioning (via Helm)
Jika Grafana menggunakan k8s-sidecar dengan label `grafana_dashboard: "1"`, dashboard akan otomatis muncul setelah deploy Helm chart:

```bash
helm upgrade matric-api ./helm/matric-api -n matric-api
```