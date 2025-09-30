# GitHub Actions CI/CD Setup

This repository includes automated CI/CD workflows for building and pushing Docker images to Docker Hub.

## Workflows

### Docker Build and Push (`docker-build-push.yml`)

This workflow automatically:
- Builds multi-architecture Docker images (linux/amd64, linux/arm64)
- Pushes to Docker Hub using the exact tag name you push
- Runs only on tag pushes
- Uses GitHub Actions cache for faster builds

## Required Secrets

To use this workflow, you need to set up the following secrets in your GitHub repository:

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Add the following repository secrets:

### `DOCKER_USERNAME`
Your Docker Hub username.

### `DOCKER_PASSWORD`
Your Docker Hub password or access token.

> **Note**: For better security, use a Docker Hub access token instead of your password. You can create one at [Docker Hub Account Settings](https://hub.docker.com/settings/security).

## Simple Tag-Based Releases

The workflow uses a simple approach - whatever tag you push becomes the Docker image tag.

### Creating a Release

Simply create and push any tag:

```bash
# Create and push a version tag
git tag v1.2.3
git push origin v1.2.3

# Or create a date-based tag
git tag release-2024-01-15
git push origin release-2024-01-15

# Or any tag format you prefer
git tag latest
git push origin latest
```

## Workflow Triggers

The workflow runs on:
- **Tag push**: Builds and pushes with the exact tag name

## Build Arguments

The Docker build includes the following build arguments:
- `VERSION`: The tag name you pushed
- `COMMIT_SHA`: Git commit hash
- `BUILD_DATE`: ISO 8601 build timestamp

## Image Tags

Images are tagged with your tag name plus `latest`:
- `your-username/health-caretaker:v1.2.3` (if you push tag `v1.2.3`)
- `your-username/health-caretaker:latest` (always updated to the most recent build)
- `your-username/health-caretaker:release-2024-01-15` (if you push tag `release-2024-01-15`)
- `your-username/health-caretaker:latest` (updated to release-2024-01-15)

## Usage

After a successful build, you can pull and run the image:

```bash
# Pull the latest version
docker pull your-username/health-caretaker:latest

# Pull a specific version
docker pull your-username/health-caretaker:v1.2.3

# Run the container with default settings
docker run -d \
  --name health-caretaker \
  -p 8080:8080 \
  -p 9091:9091 \
  your-username/health-caretaker:latest

# Run with custom environment variables
docker run -d \
  --name health-caretaker \
  -p 3000:3000 \
  -p 9092:9092 \
  -e WEB_PORT=3000 \
  -e METRICS_PORT=9092 \
  -e METRICS_ENABLED=true \
  your-username/health-caretaker:latest
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WEB_PORT` | `8080` | Port for the web UI and API |
| `METRICS_PORT` | `9091` | Port for Prometheus metrics |
| `METRICS_ENABLED` | `true` | Enable/disable metrics endpoint |
| `METRICS_PATH` | `/metrics` | Path for metrics endpoint |

### Health Checks

The Docker image includes built-in health checks that verify the application is running correctly:

```bash
# Check container health status
docker ps
docker inspect health-caretaker --format='{{.State.Health.Status}}'
```

## Troubleshooting

### Common Issues

1. **Authentication failed**: Check that `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets are correctly set.

2. **Build failed**: Check the Actions logs for specific error messages.

3. **Tag not found**: Ensure the tag follows semantic versioning format (v1.2.3).

### Viewing Workflow Logs

1. Go to your repository on GitHub
2. Click on the **Actions** tab
3. Select the workflow run you want to inspect
4. Click on the job to see detailed logs

## Security

- The workflow uses GitHub Actions cache for faster builds
- Docker images are built with multi-stage builds for smaller size
- Non-root user is used in the final image
- Build arguments are properly escaped
