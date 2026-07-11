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

## Panduan Lengkap: Dari Setup hingga Dashboard Grafana

### 1. Deploy Aplikasi ke Kubernetes

#### 1.1 Build Docker Image

```bash
docker build -t matric-api:latest .
```

#### 1.2 Load Image ke containerd k3s

```bash
docker save matric-api:latest | sudo k3s ctr images import -
```

> **Catatan:** Diperlukan akses `sudo` untuk mengimpor image ke containerd k3s.

#### 1.3 Deploy via Helm

```bash
helm upgrade --install matric-api ./helm/matric-api \
  -n matric-api --create-namespace
```

Verifikasi pod sudah running:

```bash
kubectl -n matric-api get pods -w
```

Tunggu hingga semua pod berstatus `Running` (3/3).

### 2. Setup Monitoring Stack

> **Jika sudah terinstall, lewati step ini.**

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace
```

Verifikasi:

```bash
kubectl -n monitoring get pods
```

### 3. Generate Traffic Test

Buat port-forward dan kirim request ke semua endpoint:

```bash
# Terminal 1: port-forward
kubectl -n matric-api port-forward svc/matric-api 8080:80
```

```bash
# Terminal 2: generate traffic (50 request acak)
for i in $(seq 1 50); do
  case $((RANDOM % 5)) in
    0) curl -s http://localhost:8080/ > /dev/null ;;
    1) curl -s http://localhost:8080/health > /dev/null ;;
    2) curl -s http://localhost:8080/api/matrix > /dev/null ;;
    3) curl -s -X POST http://localhost:8080/api/matrix > /dev/null ;;
    4) curl -s http://localhost:8080/metrics > /dev/null ;;
  esac
  sleep 0.3
done
echo "Selesai: 50 request terkirim"
```

Verifikasi metric tercatat:

```bash
curl -s http://localhost:8080/metrics | grep http_requests_total
```

### 4. Verifikasi di Prometheus

```bash
# Port-forward ke Prometheus UI
kubectl -n monitoring port-forward svc/monitoring-kube-prometheus-prometheus \
  9090:9090
```

Buka `http://localhost:9090` di browser:

1. **Status → Targets** — cari target `matric-api`, pastikan status **UP**
2. **Graph** — query: `http_requests_total` → Execute

### 5. Akses Grafana

```bash
# Port-forward ke Grafana
kubectl -n monitoring port-forward svc/monitoring-grafana 3000:80
```

Buka `http://localhost:3000` di browser.

**Login:**

| Field    | Value                                |
|----------|--------------------------------------|
| Username | `admin`                              |
| Password | `kubectl -n monitoring get secret monitoring-grafana -o jsonpath="{.data.admin-password}" \| base64 -d` |

### 6. Import Dashboard matric-api

1. ☰ (hamburger menu) → **Dashboards** → **New** → **Import**
2. Upload file `helm/matric-api/dashboards/matric-api-overview.json`
3. Pilih datasource **Prometheus** → **Import**

Dashboard akan muncul dengan panel-panel berikut:

| Baris | Panel | Source Data |
|-------|-------|-------------|
| 1 | Request Rate | `rate(http_requests_total[5m])` |
| 1 | Active Go Routines | `go_goroutines` |
| 1 | Memory Usage | `process_resident_memory_bytes` |
| 2 | Request Latency (p99/p95/p50) | `histogram_quantile(...)` |
| 2 | Requests by Endpoint | `sum by(path)(...)` |
| 2 | HTTP Status Breakdown | `sum by(status)(...)` |
| 3 | CPU Usage | `rate(process_cpu_seconds_total[5m])` |
| 3 | Memory Over Time | `go_memstats_heap_*` |

### 7. Verifikasi Dashboard

Jika panel masih kosong:

1. Generate traffic lagi (ulangi **step 3**)
2. Tunggu Prometheus scrape (≈30 detik)
3. Refresh dashboard (icon 🔄 di kanan atas)

Data akan mulai muncul dalam 1-2 menit setelah request pertama.

> **Auto-provisioning:** Jika Grafana menggunakan k8s-sidecar dengan label `grafana_dashboard: "1"`, dashboard akan otomatis muncul saat deploy Helm tanpa perlu import manual:
>
> ```bash
> helm upgrade matric-api ./helm/matric-api -n matric-api
> ```