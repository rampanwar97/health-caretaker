# Health Caretaker Dashboard

A modern, real-time health monitoring application built with Go that allows you to monitor HTTP/HTTPS endpoints and view their status through a beautiful web interface.

## Features

- 🚀 **Real-time Monitoring**: Monitor multiple HTTP/HTTPS endpoints simultaneously
- 📊 **Beautiful Dashboard**: Modern, responsive web UI with real-time updates
- 🔄 **WebSocket Support**: Live status updates without page refresh
- ⚡ **Fast & Lightweight**: Built with Go for high performance
- 🐳 **Docker Ready**: Easy deployment with Docker and Docker Compose
- 🔧 **Configurable**: Customizable check intervals, timeouts, and HTTP methods
- 🏷️ **Custom Labels**: Add custom labels to metrics for better organization and filtering
- 📈 **Prometheus Metrics**: Export metrics in Prometheus format with detailed labels
- 📱 **Mobile Friendly**: Responsive design that works on all devices
- 🚀 **CI/CD Ready**: GitHub Actions workflow for automated Docker builds and releases

## Quick Start

### Using Docker Hub (Recommended)

Pull and run the latest image from Docker Hub:

```bash
# Pull the latest image
docker pull your-username/health-caretaker:latest

# Run with default settings
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  your-username/health-caretaker:latest

# Or run with custom environment variables
docker run -d \
  --name health-caretaker \
  -p 3000:3000 \
  -p 9092:9092 \
  -e WEB_PORT=3000 \
  -e METRICS_PORT=9092 \
  -e METRICS_ENABLED=true \
  your-username/health-caretaker:latest
```

### Using Docker Compose

1. Clone or download this repository
2. Run the application:
   ```bash
   docker compose up -d
   ```
3. Open your browser and navigate to `http://localhost:8080`

### Manual Build and Run

1. Make sure you have Go 1.21+ installed
2. Download dependencies:
   ```bash
   go mod download
   ```
3. Run the application:
   ```bash
   go run ./cmd/server
   ```
4. Open your browser and navigate to `http://localhost:8080`

### Accessing the Application

After starting the application, you can access:

- **Web Dashboard**: `http://localhost:8080` - Main monitoring interface
- **Prometheus Metrics**: `http://localhost:9091/metrics` - Metrics endpoint
- **Health Check**: `http://localhost:8080/healthz` - Application health status
- **Readiness Check**: `http://localhost:8080/readyz` - Application readiness status

## Usage

### Adding Endpoints

1. Fill out the "Add New Endpoint" form:
   - **Name**: A friendly name for your endpoint
   - **URL**: The HTTP/HTTPS URL to monitor
   - **Method**: HTTP method (GET, POST, PUT, DELETE, HEAD)
   - **Check Interval**: How often to check (5-3600 seconds)
   - **Timeout**: Request timeout (1-60 seconds)

2. Click "Add Endpoint" to start monitoring

### Monitoring Features

- **Real-time Status**: See live status updates (UP/DOWN/CHECKING)
- **Response Time**: Monitor endpoint response times
- **Status Codes**: View HTTP status codes
- **Error Messages**: See detailed error information when endpoints are down
- **Manual Checks**: Trigger immediate checks with "Check Now" button
- **Remove Endpoints**: Delete endpoints you no longer need to monitor

### Status Indicators

- 🟢 **UP**: Endpoint is responding successfully (2xx-3xx status codes)
- 🔴 **DOWN**: Endpoint is not responding or returning error status codes
- 🟡 **CHECKING**: Currently performing a health check

## API Endpoints

The application provides a REST API for programmatic access:

- `GET /api/endpoints` - Get all monitored endpoints
- `POST /api/endpoints` - Add a new endpoint
- `DELETE /api/endpoints/{id}` - Remove an endpoint
- `POST /api/endpoints/{id}/check` - Manually trigger a check

### Example API Usage

Add a new endpoint:
```bash
curl -X POST http://localhost:8080/api/endpoints \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API",
    "url": "https://api.example.com/health",
    "method": "GET",
    "interval": 30,
    "timeout": 10
  }'
```

Get all endpoints:
```bash
curl http://localhost:8080/api/endpoints
```

## Configuration

The application uses a `config.json` file for configuration. You can modify this file to add your own endpoints or change default settings.

### Example Configuration

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

### Custom Labels

You can add custom labels to each endpoint for better organization and filtering in your monitoring system. Labels are included in all Prometheus metrics and can be used for:

- **Service Classification**: Group endpoints by service type
- **Environment Separation**: Distinguish between production, staging, development
- **Team Ownership**: Assign endpoints to specific teams
- **Criticality Levels**: Mark endpoints by importance
- **Custom Metadata**: Add any additional information

#### Label Examples

```json
{
  "name": "API Health Check",
  "url": "https://api.example.com/health",
  "method": "GET",
  "interval": 30,
  "timeout": 10,
  "labels": {
    "service": "api",
    "environment": "production",
    "team": "backend",
    "criticality": "high",
    "datacenter": "us-east-1",
    "version": "v2.1.0",
    "component": "authentication"
  }
}
```

These labels will appear in your Prometheus metrics like:
```
probe_success{name="API Health Check", url="https://api.example.com/health", service="api", environment="production", team="backend", criticality="high", datacenter="us-east-1", version="v2.1.0", component="authentication"} 1
```

### Environment Variables

- `WEB_PORT`: Web server port (default: 8080)
- `METRICS_PORT`: Metrics server port (default: 9091)
- `METRICS_ENABLED`: Enable metrics server (default: true)
- `METRICS_PATH`: Metrics endpoint path (default: /metrics)
- `DEBUG`: Enable debug logging (default: false)

## Docker Deployment

### Build and Run with Docker

```bash
# Build the image
docker build -t health-monitor .

# Run the container
docker run -p 8080:8080 health-monitor
```

### Using Docker Compose

```bash
# Start the service
docker compose up -d

# View logs
docker compose logs -f

# Stop the service
docker compose down
```

## Security Considerations

- The application accepts self-signed certificates for HTTPS endpoints
- No authentication is implemented - consider adding authentication for production use
- The WebSocket connection allows all origins - restrict this for production

## Docker Hub Repository

### Available Images

The application is available on Docker Hub with the following tags:

- `your-username/health-caretaker:latest` - Always points to the most recent build
- `your-username/health-caretaker:v1.0.0` - Specific version tags
- `your-username/health-caretaker:release-2024-01-15` - Date-based releases

> **Note**: Every time you push a tag, both the specific tag and `latest` are updated automatically.

### Running from Docker Hub

#### Basic Usage

```bash
# Pull and run the latest image
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  your-username/health-caretaker:latest
```

#### With Custom Configuration

```bash
# Run with custom ports and settings
docker run -d \
  --name health-caretaker \
  -p 3000:3000 \
  -p 9092:9092 \
  -e WEB_PORT=3000 \
  -e METRICS_PORT=9092 \
  -e METRICS_ENABLED=true \
  -e METRICS_PATH=/custom-metrics \
  your-username/health-caretaker:latest
```

#### With Custom Config File

```bash
# Run with your own configuration file
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  -v /path/to/your/config.json:/config.json:ro \
  your-username/health-caretaker:latest
```

### Environment Variables

The application supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WEB_PORT` | `8080` | Port for the web UI and API |
| `METRICS_PORT` | `9091` | Port for Prometheus metrics |
| `METRICS_ENABLED` | `true` | Enable/disable metrics endpoint |
| `METRICS_PATH` | `/metrics` | Path for metrics endpoint |

#### Environment Variable Examples

```bash
# Production setup with custom ports
docker run -d \
  --name health-caretaker \
  -p 80:80 \
  -p 9090:9090 \
  -e WEB_PORT=80 \
  -e METRICS_PORT=9090 \
  -e METRICS_PATH=/prometheus/metrics \
  your-username/health-caretaker:latest

# Development setup
docker run -d \
  --name health-caretaker-dev \
  -p 3000:3000 \
  -p 9091:9091 \
  -e WEB_PORT=3000 \
  your-username/health-caretaker:latest
```

### Health Checks

The Docker image includes built-in health checks:

```bash
# Check container health
docker ps
docker inspect health-caretaker --format='{{.State.Health.Status}}'

# View health check logs
docker inspect health-caretaker --format='{{range .State.Health.Log}}{{.Output}}{{end}}'
```

### Docker Compose with Environment Variables

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  health-caretaker:
    image: your-username/health-caretaker:latest
    ports:
      - "${WEB_PORT:-8080}:${WEB_PORT:-8080}"
      - "${METRICS_PORT:-9091}:${METRICS_PORT:-9091}"
    environment:
      - WEB_PORT=${WEB_PORT:-8080}
      - METRICS_PORT=${METRICS_PORT:-9091}
      - METRICS_ENABLED=${METRICS_ENABLED:-true}
      - METRICS_PATH=${METRICS_PATH:-/metrics}
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "/health-caretaker", "-version"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

Then run with:

```bash
# With default settings
docker compose up -d

# With custom environment variables
WEB_PORT=3000 METRICS_PORT=9092 docker compose up -d
```

## CI/CD and Docker Builds

This project includes automated CI/CD workflows using GitHub Actions for building and pushing Docker images to Docker Hub.

### GitHub Actions Setup

1. **Set up Docker Hub secrets** in your GitHub repository:
   - Go to **Settings** → **Secrets and variables** → **Actions**
   - Add `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets

2. **Automatic builds** on:
   - Tag push (any tag format)

### Creating Releases

Simply create and push a tag to trigger the build:

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0

# Or create any tag format
git tag release-2024-01-15
git push origin release-2024-01-15
```

### Docker Images

After a successful build, the images will be available as:
- `your-username/health-caretaker:your-tag-name` (specific version)
- `your-username/health-caretaker:latest` (always points to the most recent build)

For example, when you push tag `v1.0.0`:
- `your-username/health-caretaker:v1.0.0`
- `your-username/health-caretaker:latest` (updated to v1.0.0)

For detailed CI/CD setup instructions, see [.github/README.md](.github/README.md).

## Development

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- Make (optional, for using Makefile commands)

### Dependencies

- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support

### Building

```bash
# Using Make (recommended)
make build

# Or manually
go build -o health-monitor ./cmd/server
```

### Running Tests

```bash
# Using Make
make test

# Or manually
go test ./...
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run tests
make docker-build  # Build Docker image
make compose-up    # Start with Docker Compose
make version       # Show version information
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on the project repository.
