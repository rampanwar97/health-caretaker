# 🛡️ Distroless Docker Images

This project uses Google's [distroless](https://github.com/GoogleContainerTools/distroless) base images for enhanced security and minimal attack surface.

## 🎯 Benefits of Distroless Images

### Security
- **No shell**: No `/bin/sh` or `/bin/bash` - prevents shell injection attacks
- **No package manager**: No `apt`, `yum`, or `apk` - prevents package manager exploits
- **Minimal attack surface**: Only essential runtime libraries included
- **Non-root user**: Runs as `nonroot` user by default
- **Regular updates**: Google maintains and updates the base images

### Size
- **Smaller images**: Significantly smaller than full OS images
- **No unnecessary tools**: Only runtime dependencies included
- **Optimized layers**: Efficient layer caching and distribution

### Compliance
- **CVE scanning**: Regular vulnerability scanning by Google
- **SBOM support**: Software Bill of Materials included
- **Supply chain security**: Minimal dependencies reduce supply chain risks

## 📦 Available Distroless Variants

### 1. Static (Default) - `Dockerfile`
```dockerfile
FROM gcr.io/distroless/static-debian12:nonroot
```
- **Use case**: Fully static Go binaries (CGO_ENABLED=0)
- **Size**: Smallest (~2MB base)
- **Dependencies**: None (fully static)
- **Best for**: Production deployments, maximum security

### 2. Base (Alternative) - `Dockerfile.distroless-base`
```dockerfile
FROM gcr.io/distroless/base-debian12:nonroot
```
- **Use case**: Go binaries that need glibc
- **Size**: Larger (~20MB base)
- **Dependencies**: glibc, basic runtime libraries
- **Best for**: Compatibility with CGO or external libraries

## 🚀 Building Distroless Images

### Single Architecture
```bash
# Build static distroless image (default)
make docker-build

# Build base distroless image (with glibc)
make docker-build-base
```

### Multi-Architecture
```bash
# Build multi-arch static distroless
make docker-build-multi

# Build multi-arch base distroless
make docker-build-multi-base
```

### Manual Docker Commands
```bash
# Static distroless
docker build -t health-caretaker:latest .

# Base distroless
docker build -f Dockerfile.distroless-base -t health-caretaker-base:latest .

# Multi-architecture static
docker buildx build --platform linux/amd64,linux/arm64 -t health-caretaker:latest .

# Multi-architecture base
docker buildx build --platform linux/amd64,linux/arm64 -f Dockerfile.distroless-base -t health-caretaker-base:latest .
```

## 🏃‍♂️ Running Distroless Images

### Static Distroless
```bash
# Using Make
make docker-run

# Using Docker
docker run -p 8080:8080 -p 9091:9091 health-caretaker:latest
```

### Base Distroless
```bash
# Using Make
make docker-run-base

# Using Docker
docker run -p 8080:8080 -p 9091:9091 health-caretaker-base:latest
```

## 🔍 Image Comparison

| Image Type | Base Size | Total Size | Security | Compatibility |
|------------|-----------|------------|----------|---------------|
| **Static Distroless** | ~2MB | ~15MB | ⭐⭐⭐⭐⭐ | Go static binaries only |
| **Base Distroless** | ~20MB | ~35MB | ⭐⭐⭐⭐ | Go + glibc dependencies |
| **Alpine** | ~5MB | ~25MB | ⭐⭐⭐ | Full OS with package manager |
| **Ubuntu** | ~70MB | ~90MB | ⭐⭐ | Full OS with many tools |

## 🛠️ Troubleshooting

### Static Binary Issues
If you get errors like "no such file or directory" or "exec format error":
```bash
# Check if binary is static
file health-caretaker
# Should show: "statically linked"

# Check dependencies
ldd health-caretaker
# Should show: "not a dynamic executable"

# Use base distroless instead
make docker-build-base
```

### Debugging Distroless Containers
Since distroless images have no shell, debugging is limited:

```bash
# Check container logs
docker logs <container-id>

# Check if container is running
docker ps

# Inspect container
docker inspect <container-id>

# For debugging, temporarily use a debug image
docker run --rm -it --entrypoint sh gcr.io/distroless/base-debian12:debug
```

### Health Check Issues
If health checks fail:
```bash
# Test health check manually
docker run --rm health-caretaker:latest /health-caretaker -version

# Check if ports are exposed correctly
docker run -p 8080:8080 -p 9091:9091 health-caretaker:latest
```

## 🔒 Security Best Practices

### 1. Use Static Distroless (Default)
- Prefer `gcr.io/distroless/static-debian12:nonroot`
- Ensures fully static binary with no dependencies

### 2. Regular Updates
```bash
# Pull latest distroless images
docker pull gcr.io/distroless/static-debian12:nonroot
docker pull gcr.io/distroless/base-debian12:nonroot
```

### 3. Scan for Vulnerabilities
```bash
# Using Trivy
trivy image health-caretaker:latest

# Using Docker Scout
docker scout cves health-caretaker:latest
```

### 4. Non-Root User
Both distroless variants run as `nonroot` user by default:
- UID: 65532
- GID: 65532
- No shell access
- Limited file system permissions

## 📊 Performance Impact

### Startup Time
- **Static distroless**: Fastest startup (no library loading)
- **Base distroless**: Slightly slower (glibc loading)
- **Alpine/Ubuntu**: Slowest (full OS initialization)

### Memory Usage
- **Static distroless**: Lowest memory footprint
- **Base distroless**: Slightly higher (glibc overhead)
- **Alpine/Ubuntu**: Highest (full OS overhead)

### Network Performance
- No impact on network performance
- Same Go runtime and application code

## 🚀 Production Recommendations

### For Maximum Security
```bash
# Use static distroless
make docker-build-multi
```

### For Compatibility
```bash
# Use base distroless if you need glibc
make docker-build-multi-base
```

### For CI/CD
```bash
# Build both variants for testing
make docker-build
make docker-build-base

# Use multi-arch for production
make docker-build-multi
```

## 📚 References

- [Google Distroless](https://github.com/GoogleContainerTools/distroless)
- [Distroless Images](https://console.cloud.google.com/gcr/images/distroless)
- [Container Security Best Practices](https://cloud.google.com/architecture/best-practices-for-operating-containers)
- [Go Static Binaries](https://golang.org/cmd/cgo/)
