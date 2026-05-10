# example-app
A production-ready Go service template featuring a complete observability stack (Prometheus/Grafana), automated CI/CD via GitHub Actions, and Kubernetes deployment.

# Observability
- Latency (Request latency): a histogram of request durations by endpoint
- Traffic (Request volume): a total request counter by endpoint and response code
- Errors (Status Codes): a total request counter by endpoint and response code
- Saturation: default metrics from Prometheus & Go such as CPU utilization, RAM utilization, Go Goroutines, Go Heap size, file description utilization, etc.

## Local observability with Docker Compose

Docker-Compose runs:

- the Go app on `http://localhost:8080`
- Prometheus on `http://localhost:9090`
- Grafana on `http://localhost:3000`

### Run the stack

```bash
docker compose up --build
```

Then open:

- App health endpoint: `http://localhost:8080/health`
- App metrics endpoint: `http://localhost:8080/metrics`
- Prometheus UI: `http://localhost:9090`
- Grafana UI: `http://localhost:3000` using `admin` / `admin`

### Details about the Configuration

- Metrics storage uses both a time-based retention limit (`7d`) and a size cap (`512MB`) to keep local disk usage predictable
- Grafana is provisioned automatically so the Prometheus datasource and dashboard are available on first startup
- Docker volumes persist Prometheus and Grafana data across restarts
- The app includes a healthcheck so Prometheus waits until the service is actually ready to scrape

### How to generate metrics locally

In another terminal, send a requests in a loop:

```bash
for i in $(seq 1 40); do curl -s http://localhost:8080/health > /dev/null; done
```

After that, you should be able to query these in Prometheus or see them in Grafana:

- `http_requests_total`
- `http_request_duration_seconds`
- `process_cpu_seconds_total`
- `process_resident_memory_bytes`
- `go_goroutines`
- `go_gc_duration_seconds`
- `go_memstats_heap_sys_bytes`

### Stop the stack

```bash
docker compose down
```

To also remove persisted Prometheus and Grafana data:

```bash
docker compose down -v
```
