# Docker Compose Setup for NGINX Prometheus Exporter

This Docker Compose configuration provides a complete monitoring stack for NGINX using Prometheus and Grafana.

## Table of Contents

- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Services](#services)
  - [nginx-exporter](#nginx-exporter)
  - [nginx (Main Server)](#nginx-main-server)
  - [app (Sample Application)](#app-sample-application)
  - [prometheus](#prometheus)
  - [grafana](#grafana)
- [Configuration Files](#configuration-files)
- [Testing](#testing)
- [Stopping](#stopping)
- [Troubleshooting](#troubleshooting)
- [Network](#network)
- [Volumes](#volumes)

## Architecture

The setup includes:
- **nginx**: Main NGINX server with stub_status enabled (for demo purposes)
- **app**: Sample web application served by NGINX (for demo purposes)
- **nginx-exporter**: NGINX Prometheus Exporter that scrapes metrics from NGINX
- **prometheus**: Prometheus time-series database for metrics collection
- **grafana**: Grafana dashboard for visualization

**Note**: If you have an existing NGINX setup, you only need to run the `nginx-exporter` service. Make sure your NGINX has stub_status enabled and is accessible at the configured scrape URI.

## Quick Start

1. **Navigate to the project root directory**:
   ```bash
   cd /path/to/nginx-prometheus-exporter
   ```

2. **For existing NGINX setup - Start only the nginx-exporter**:
   ```bash
   docker-compose up -d nginx-exporter
   ```
   This assumes you already have NGINX running with stub_status enabled at `http://nginx:8081/stub_status`.

3. **For complete demo setup - Start the full monitoring stack**:
   ```bash
   docker-compose up -d
   ```
   This starts all services including a sample NGINX server, Prometheus, and Grafana.

4. **Access the services**:
   - **Sample App**: http://localhost:8080
   - **NGINX Server**: http://localhost:80
   - **NGINX Metrics**: http://localhost:9113/metrics
   - **NGINX Stub Status**: http://localhost:8081/stub_status
   - **Prometheus**: http://localhost:9090
   - **Grafana**: http://localhost:3000 (admin/admin123)

## Services

### nginx-exporter
- **Image**: `nginx/nginx-prometheus-exporter:latest`
- **Port**: 9113
- **Function**: Scrapes NGINX stub_status and converts to Prometheus metrics

### nginx (Main Server)
- **Image**: `nginx:alpine`
- **Ports**: 80, 8081 (stub_status)
- **Function**: Proxies to sample app and provides metrics endpoint

### app (Sample Application)
- **Image**: `nginx:alpine`
- **Port**: 8080
- **Function**: Serves sample web content

### prometheus
- **Image**: `prom/prometheus:latest`
- **Port**: 9090
- **Function**: Collects and stores metrics from nginx-exporter

### grafana
- **Image**: `grafana/grafana:latest`
- **Port**: 3000
- **Function**: Provides visualization dashboards

## Configuration Files

- `nginx/nginx.conf`: Main NGINX configuration with stub_status
- `app/nginx.conf`: Sample application NGINX configuration
- `prometheus/prometheus.yml`: Prometheus scraping configuration
- `grafana/provisioning/`: Grafana datasource and dashboard configuration

## Testing

1. **Check nginx-exporter metrics**:
   ```bash
   curl http://localhost:9113/metrics | grep nginx
   ```

2. **Check stub_status directly**:
   ```bash
   curl http://localhost:8081/stub_status
   ```

3. **Generate traffic**:
   ```bash
   for i in {1..10}; do curl -s http://localhost:80 > /dev/null; done
   ```

## Stopping

Stop all services:
```bash
docker-compose down
```

Stop and remove volumes:
```bash
docker-compose down -v
```

## Troubleshooting

1. **Check container logs**:
   ```bash
   docker-compose logs nginx-exporter
   docker-compose logs nginx
   ```

2. **Verify container status**:
   ```bash
   docker-compose ps
   ```

3. **Test internal connectivity**:
   ```bash
   docker exec nginx-server curl -s http://localhost:8081/stub_status
   ```

## Network

All services run on the `nginx-monitoring` bridge network for internal communication.

## Volumes

- `prometheus_data`: Persistent storage for Prometheus metrics
- `grafana_data`: Persistent storage for Grafana configurations and dashboards
