# Health Monitoring Dashboard

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

## Quick Start

### Using Docker Compose (Recommended)

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
   go run main.go
   ```
4. Open your browser and navigate to `http://localhost:8080`

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

## Development

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

### Dependencies

- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support

### Building

```bash
go build -o health-monitor main.go
```

### Running Tests

```bash
go test ./...
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on the project repository.
