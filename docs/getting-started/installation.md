# SREDIAG Installation Guide

This guide provides detailed instructions for installing SREDIAG in various environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation Methods](#installation-methods)
3. [Platform-Specific Instructions](#platform-specific-instructions)
4. [Verification](#verification)
5. [Troubleshooting](#troubleshooting)

## Prerequisites

Before installing SREDIAG, ensure your system meets these requirements:

### System Requirements

- CPU: 2+ cores recommended
- RAM: 2GB minimum, 4GB recommended
- Disk: 1GB free space
- OS: Linux (kernel 4.19+), macOS 10.15+, or Windows 10/Server 2019+

### Software Requirements

- Go 1.21 or later
- Docker 20.10.0+ (for containerized deployment)
- Kubernetes 1.24+ (for Kubernetes deployment)
- Git (for source installation)

## Installation Methods

### 1. Using Go

```bash
# Install the latest version
go install github.com/yourusername/srediag/cmd/srediag@latest

# Install a specific version
go install github.com/yourusername/srediag/cmd/srediag@v1.0.0
```

### 2. Using Docker

```bash
# Pull the latest image
docker pull srediag/srediag:latest

# Run the container
docker run -d \
  --name srediag \
  -p 8080:8080 \
  -p 9090:9090 \
  -v /path/to/config:/etc/srediag \
  srediag/srediag:latest
```

### 3. Using Helm (Kubernetes)

```bash
# Add SREDIAG Helm repository
helm repo add srediag https://charts.srediag.io
helm repo update

# Install SREDIAG
helm install srediag srediag/srediag \
  --namespace monitoring \
  --create-namespace \
  --values values.yaml
```

### 4. From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/srediag.git
cd srediag

# Build the binary
make build

# Install to system
make install
```

## Platform-Specific Instructions

### Linux

1. **System Packages**

   ```bash
   # Ubuntu/Debian
   apt-get update && apt-get install -y \
     build-essential \
     pkg-config \
     libsystemd-dev

   # RHEL/CentOS
   yum groupinstall -y "Development Tools"
   yum install -y systemd-devel
   ```

2. **Systemd Service**

   ```bash
   # Create service file
   cat > /etc/systemd/system/srediag.service << EOF
   [Unit]
   Description=SREDIAG Diagnostic Service
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/srediag
   Restart=always
   User=srediag

   [Install]
   WantedBy=multi-user.target
   EOF

   # Start service
   systemctl daemon-reload
   systemctl enable --now srediag
   ```

### macOS

```bash
# Using Homebrew
brew tap srediag/tools
brew install srediag
```

### Windows

```powershell
# Using Chocolatey
choco install srediag

# Using Scoop
scoop bucket add srediag https://github.com/yourusername/srediag-bucket
scoop install srediag
```

## Verification

After installation, verify SREDIAG is working correctly:

1. **Check Version**

   ```bash
   srediag version
   ```

2. **Health Check**

   ```bash
   srediag health check
   ```

3. **Test Configuration**

   ```bash
   srediag test config.yaml
   ```

## Troubleshooting

### Common Issues

1. **Permission Errors**

   ```bash
   # Fix directory permissions
   sudo chown -R srediag:srediag /etc/srediag
   sudo chmod 755 /etc/srediag
   ```

2. **Port Conflicts**

   ```bash
   # Check port usage
   sudo lsof -i :8080
   sudo lsof -i :9090
   ```

3. **Missing Dependencies**

   ```bash
   # Install required libraries
   sudo apt-get install -y \
     libssl-dev \
     zlib1g-dev
   ```

### Getting Help

- Check the [Troubleshooting Guide](../reference/troubleshooting.md)
- Join our [Community Discord](https://discord.gg/srediag)
- Open an issue on [GitHub](https://github.com/yourusername/srediag/issues)

## Next Steps

- [Quick Start Guide](quickstart.md)
- [Basic Configuration](../configuration/README.md)
- [Security Setup](../security/README.md)
- [Monitoring Setup](../configuration/telemetry.md)
