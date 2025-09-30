# Health Caretaker

[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/r/ramp110397/health-caretaker)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![Helm](https://img.shields.io/badge/Helm-0F1689?style=for-the-badge&logo=Helm&labelColor=0F1689)](https://helm.sh/)

A modern, real-time health monitoring application built with Go that allows you to monitor HTTP/HTTPS endpoints and view their status through a beautiful web interface with Prometheus metrics support.

## 📋 Table of Contents

- [What is Health Caretaker?](#what-is-health-caretaker)
- [Features](#features)
- [Screenshots](#screenshots)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Guide](#usage-guide)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [API Reference](#api-reference)
- [Monitoring & Metrics](#monitoring--metrics)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## 🤔 What is Health Caretaker?

Health Caretaker is a comprehensive endpoint monitoring solution designed for modern applications. It provides:

- **Real-time monitoring** of HTTP/HTTPS endpoints
- **Beautiful web dashboard** with live status updates
- **Prometheus metrics** for integration with monitoring stacks
- **Custom labels** for advanced filtering and organization
- **Multiple deployment options** (Docker, Kubernetes, Helm)
- **Production-ready** with health checks and graceful shutdown

Perfect for monitoring APIs, microservices, databases, and any HTTP-based service in your infrastructure.

## ✨ Features

### Core Monitoring
- 🚀 **Real-time Monitoring**: Monitor multiple HTTP/HTTPS endpoints simultaneously
- 📊 **Beautiful Dashboard**: Modern, responsive web UI with real-time updates
- 🔄 **WebSocket Support**: Live status updates without page refresh
- ⚡ **Fast & Lightweight**: Built with Go for high performance
- 🔧 **Configurable**: Customizable check intervals, timeouts, and HTTP methods

### Advanced Features
- 🏷️ **Custom Labels**: Add custom labels to metrics for better organization and filtering
- 📈 **Prometheus Metrics**: Export metrics in Prometheus format with detailed labels
- 📱 **Mobile Friendly**: Responsive design that works on all devices
- 🔍 **Health Checks**: Built-in `/healthz` and `/readyz` endpoints
- 🛡️ **Security**: Non-root containers, read-only filesystem, security contexts

### Deployment & Operations
- 🐳 **Docker Ready**: Easy deployment with Docker and Docker Compose
- ☸️ **Kubernetes Native**: Complete Helm chart with production-ready configurations
- 🚀 **CI/CD Ready**: GitHub Actions workflow for automated Docker builds and releases
- 📊 **Monitoring Integration**: ServiceMonitor for Prometheus Operator
- 🔄 **Auto-scaling**: Horizontal Pod Autoscaler support

## 📸 Screenshots

### Web Dashboard
The modern, responsive web interface provides real-time monitoring with:
- Live endpoint status updates
- Response time tracking
- Custom label filtering
- Mobile-friendly design

### Prometheus Metrics
Rich metrics export including:
- `probe_success` - Endpoint availability (0/1)
- `probe_duration_seconds` - Response time
- `probe_interval_seconds` - Check interval
- Custom labels for filtering and grouping

## 🚀 Quick Start

### Option 1: Docker (Recommended for Testing)

```bash
# Pull and run the latest image
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  ramp110397/health-caretaker:latest

# Access the web UI
open http://localhost:8080

# View Prometheus metrics
curl http://localhost:9091/metrics
```

### Option 2: Docker Compose

```bash
# Clone the repository
git clone https://github.com/ramp110397/health-caretaker.git
cd health-caretaker

# Start the application
docker compose up -d

# Access the web UI
open http://localhost:8080
```

### Option 3: Kubernetes with Helm

```bash
# Add the Helm repository (if published)
helm repo add health-caretaker https://ramp110397.github.io/health-caretaker
helm repo update

# Install with default configuration
helm install my-health-caretaker health-caretaker/health-caretaker

# Or install from local chart
helm install my-health-caretaker ./helm/health-caretaker
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WEB_PORT` | `8080` | Port for the web UI and API |
| `METRICS_PORT` | `9091` | Port for Prometheus metrics |
| `METRICS_ENABLED` | `true` | Enable/disable metrics endpoint |
| `METRICS_PATH` | `/metrics` | Path for metrics endpoint |

### Configuration File (config.json)

The application uses a JSON configuration file to define endpoints to monitor:

```json
{
  "endpoints": [
    {
      "name": "Google",
      "url": "https://www.google.com",
      "method": "GET",
      "interval": 30,
      "timeout": 10,
      "labels": {
        "service": "search",
        "environment": "production",
        "team": "platform",
        "criticality": "high"
      }
    }
  ],
  "server": {
    "port": "8080"
  },
  "metrics": {
    "enabled": true,
    "path": "/metrics",
    "port": "9091"
  }
}
```

### Endpoint Configuration

Each endpoint can be configured with:

- **name**: Friendly name for the endpoint
- **url**: HTTP/HTTPS URL to monitor
- **method**: HTTP method (GET, POST, PUT, DELETE, HEAD)
- **interval**: Check interval in seconds (5-3600)
- **timeout**: Request timeout in seconds (1-60)
- **labels**: Custom key-value pairs for metrics filtering
- **probe_type**: Type of probe (optional, e.g., "livez", "readyz")

## 📖 Usage Guide

### Adding Endpoints via Web UI

1. Open the web dashboard at `http://localhost:8080`
2. Fill out the "Add New Endpoint" form:
   - **Name**: A friendly name for your endpoint
   - **URL**: The HTTP/HTTPS URL to monitor
   - **Method**: HTTP method (GET, POST, PUT, DELETE, HEAD)
   - **Check Interval**: How often to check (5-3600 seconds)
   - **Timeout**: Request timeout (1-60 seconds)
3. Click "Add Endpoint" to start monitoring

### Adding Endpoints via Configuration

Edit the `config.json` file and restart the application:

```json
{
  "endpoints": [
    {
      "name": "My API",
      "url": "https://api.example.com/health",
      "method": "GET",
      "interval": 30,
      "timeout": 10,
      "labels": {
        "service": "api",
        "environment": "production",
        "team": "backend"
      }
    }
  ]
}
```

### Monitoring Endpoints

- **Green**: Endpoint is healthy and responding
- **Red**: Endpoint is down or not responding
- **Yellow**: Endpoint is being checked
- **Gray**: Endpoint check failed or timed out

### Custom Labels

Add custom labels to organize and filter your endpoints:

```json
{
  "labels": {
    "service": "api",
    "environment": "production",
    "team": "backend",
    "criticality": "high",
    "region": "us-east-1"
  }
}
```

Use these labels in Prometheus queries:
```promql
# Filter by service
probe_success{service="api"}

# Filter by environment and criticality
probe_success{environment="production", criticality="high"}

# Calculate uptime percentage
avg_over_time(probe_success[5m])
```

## 🐳 Docker Deployment

### Basic Docker Run

```bash
# Run with default configuration
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  ramp110397/health-caretaker:latest
```

### Custom Configuration

```bash
# Run with custom ports
docker run -d \
  --name health-caretaker \
  -p 3000:3000 \
  -p 9092:9092 \
  -e WEB_PORT=3000 \
  -e METRICS_PORT=9092 \
  ramp110397/health-caretaker:latest

# Run with custom config file
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  -v /path/to/your/config.json:/config.json:ro \
  ramp110397/health-caretaker:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  health-caretaker:
    image: ramp110397/health-caretaker:latest
    ports:
      - "8080:8080"
      - "9091:9091"
    environment:
      - WEB_PORT=8080
      - METRICS_PORT=9091
      - METRICS_ENABLED=true
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "/health-caretaker", "-version"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## ☸️ Kubernetes Deployment

### Using Helm (Recommended)

```bash
# Install with default values
helm install my-health-caretaker ./helm/health-caretaker

# Install with custom values
helm install my-health-caretaker ./helm/health-caretaker \
  --set image.tag=v1.0.0 \
  --set replicaCount=3 \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=health-caretaker.example.com

# Install with production values
helm install my-health-caretaker ./helm/health-caretaker \
  -f ./helm/health-caretaker/values-production.yaml
```

### Using kubectl

```bash
# Apply the Kubernetes manifests
kubectl apply -f k8s/

# Check the deployment
kubectl get pods -l app=health-caretaker

# Port forward for testing
kubectl port-forward svc/health-caretaker 8080:8080
kubectl port-forward svc/health-caretaker 9091:9091
```

### Production Configuration

The Helm chart includes production-ready features:

- **Security**: Non-root containers, read-only filesystem
- **High Availability**: Pod disruption budget, anti-affinity
- **Auto-scaling**: Horizontal Pod Autoscaler
- **Monitoring**: ServiceMonitor for Prometheus
- **Ingress**: TLS termination, load balancing
- **Resource Management**: CPU and memory limits

## 🔌 API Reference

### REST API Endpoints

#### Get All Endpoints
```bash
GET /api/endpoints
```

#### Add New Endpoint
```bash
POST /api/endpoints
Content-Type: application/json

{
  "name": "My API",
  "url": "https://api.example.com/health",
  "method": "GET",
  "interval": 30,
  "timeout": 10,
  "labels": {
    "service": "api",
    "environment": "production"
  }
}
```

#### Delete Endpoint
```bash
DELETE /api/endpoints/{id}
```

#### Check Endpoint
```bash
POST /api/endpoints/{id}/check
```

### WebSocket API

Connect to `/ws` for real-time updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Endpoint update:', data);
};
```

### Health Check Endpoints

- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe

## 📊 Monitoring & Metrics

### Prometheus Metrics

The application exports the following metrics:

#### `probe_success`
- **Type**: Gauge
- **Description**: Displays whether the probe was successful (1 = up, 0 = down)
- **Labels**: `name`, `url`, plus any custom labels

#### `probe_duration_seconds`
- **Type**: Gauge
- **Description**: Returns how long the probe took to complete in seconds
- **Labels**: `name`, `url`, plus any custom labels

#### `probe_interval_seconds`
- **Type**: Gauge
- **Description**: The interval between probes in seconds
- **Labels**: `name`, `url`, plus any custom labels

### Example Prometheus Queries

```promql
# Overall uptime percentage
avg(probe_success) * 100

# Uptime by service
avg by (service) (probe_success) * 100

# Response time by endpoint
avg by (name) (probe_duration_seconds)

# Endpoints down for more than 5 minutes
probe_success == 0 and time() - probe_success < 300

# High response time alerts
probe_duration_seconds > 5
```

### Grafana Dashboard

Create a Grafana dashboard using these queries:

1. **Uptime Overview**: `avg(probe_success) * 100`
2. **Response Time**: `avg(probe_duration_seconds)`
3. **Endpoint Status**: `probe_success`
4. **Service Breakdown**: `avg by (service) (probe_success)`

### ServiceMonitor (Prometheus Operator)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: health-caretaker
spec:
  selector:
    matchLabels:
      app: health-caretaker
  endpoints:
  - port: metrics
    path: /metrics
    interval: 30s
```

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Docker (optional)
- Make (optional)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ramp110397/health-caretaker.git
cd health-caretaker

# Download dependencies
go mod download

# Build the application
make build

# Run the application
make run

# Run tests
make test

# Build Docker image
make docker-build
```

### Project Structure

```
health-caretaker/
├── cmd/server/           # Application entry point
├── internal/             # Internal packages
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── metrics/         # Prometheus metrics
│   ├── models/          # Data models
│   ├── monitor/         # Endpoint monitoring
│   └── server/          # HTTP server
├── pkg/                 # Reusable packages
│   ├── logger/          # Logging utilities
│   └── middleware/      # HTTP middleware
├── static/              # Web UI assets
├── helm/                # Helm chart
└── .github/             # GitHub Actions
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run tests
make docker-build  # Build Docker image
make compose-up    # Start with Docker Compose
make version       # Show version information
make lint          # Lint code
make fmt           # Format code
```

### Environment Variables for Development

```bash
# Copy example environment file
cp env.example .env

# Edit environment variables
WEB_PORT=8080
METRICS_PORT=9091
METRICS_ENABLED=true
METRICS_PATH=/metrics
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style

- Follow Go conventions
- Add tests for new features
- Update documentation
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Web UI uses modern HTML5, CSS3, and JavaScript
- Metrics compatible with [Prometheus](https://prometheus.io/)
- Containerized with [Docker](https://www.docker.com/)
- Orchestrated with [Kubernetes](https://kubernetes.io/)
- Packaged with [Helm](https://helm.sh/)

## 📞 Support

- 📖 [Documentation](https://github.com/ramp110397/health-caretaker#readme)
- 🐛 [Issue Tracker](https://github.com/ramp110397/health-caretaker/issues)
- 💬 [Discussions](https://github.com/ramp110397/health-caretaker/discussions)
- 📧 [Email Support](mailto:support@example.com)

---

**Made with ❤️ by [ramp110397](https://github.com/ramp110397)**