# SPDX-License-Identifier: Apache-2.0

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Set build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install build tools
RUN go install github.com/magefile/mage@latest

# Copy source code
COPY . .

# Build the application with build info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    mage build

# Verify the binary was created
RUN ls -la bin/ && \
    file bin/srediag && \
    ./bin/srediag --help

# Final stage - distroless for security
FROM gcr.io/distroless/static-debian12:nonroot

# Add metadata labels
LABEL org.opencontainers.image.title="SREDIAG" \
      org.opencontainers.image.description="System Resource and Environment Diagnostics Agent" \
      org.opencontainers.image.source="https://github.com/srediag/srediag" \
      org.opencontainers.image.vendor="Integra SH" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.documentation="https://github.com/srediag/srediag/blob/main/README.md"

# Copy binary from builder stage
COPY --from=builder /app/bin/srediag /srediag

# Copy configuration files
COPY --from=builder /app/configs /configs

# Copy CA certificates for TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set working directory
WORKDIR /

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD ["/srediag", "service", "health"]

# Expose common ports
EXPOSE 4317 4318 8889 13133

# Use non-root user (from distroless)
USER nonroot:nonroot

# Set entrypoint
ENTRYPOINT ["/srediag"]
CMD ["service", "start", "--config", "/configs/srediag.yaml"]
