# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for building
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for building
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' -X 'main.Version=${VERSION}' -X 'main.CommitSHA=${COMMIT_SHA}' -X 'main.BuildDate=${BUILD_DATE}'" \
    -a -installsuffix cgo \
    -o health-caretaker \
    ./cmd/server

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /app/health-caretaker /health-caretaker

# Copy static files
COPY --from=builder /app/static /static

# Copy config file
COPY --from=builder /app/config.json /config.json

# Use non-root user
USER appuser

# Expose ports
EXPOSE 8080 9091

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/health-caretaker", "-version"] || exit 1

# Run the application
ENTRYPOINT ["/health-caretaker"]
