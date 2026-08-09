# Integration Tests

This directory contains integration tests for SREDIAG that test the complete system functionality including service startup, Docker containers, and external integrations.

## Running Integration Tests

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Available ports: 4317, 4318, 8889, 13133 (or test variations)

### Running Tests

```bash
# Run all integration tests
go test -tags=integration ./tests/integration/...

# Run with verbose output
go test -tags=integration -v ./tests/integration/...

# Run specific test
go test -tags=integration -v ./tests/integration/ -run TestServiceIntegration

# Run tests with timeout
go test -tags=integration -timeout 10m ./tests/integration/...
```

### Test Categories

#### Service Integration Tests (`service_test.go`)

Tests the complete SREDIAG service lifecycle:

- **TestServiceIntegration**: Tests service startup, health checks, metrics endpoint, and OTLP endpoints
- **TestDiagnosticsIntegration**: Tests system, performance, and security diagnostics commands
- **TestPluginIntegration**: Tests plugin listing and management
- **TestBuildIntegration**: Tests build command functionality

#### Docker Integration Tests (`docker_test.go`)

Tests Docker-related functionality:

- **TestDockerBuild**: Tests that the Docker image builds successfully
- **TestDockerRun**: Tests running the Docker container and health checks
- **TestDockerCompose**: Tests the complete docker-compose setup with SREDIAG and Jaeger

### Test Environment

Integration tests use temporary ports and isolated environments to avoid conflicts:

- Service tests use default ports (4317, 4318, 8889, 13133)
- Docker run tests use offset ports (14317, 14318, 18889, 23133)
- All tests include proper cleanup procedures

### Debugging Tests

If tests fail, you can debug by:

1. **Check test output**: Tests log detailed output including command results
2. **Manual testing**: Run the same commands manually to reproduce issues
3. **Docker logs**: Check container logs when Docker tests fail
4. **Port conflicts**: Ensure no other services are using the test ports

### CI/CD Integration

Integration tests are designed to run in CI/CD environments:

- Tests skip automatically in short mode (`go test -short`)
- Tests include appropriate timeouts for CI environments
- Docker tests handle image cleanup automatically
- Tests provide detailed logging for debugging CI failures

### Test Data and Fixtures

Tests use:

- Default configuration files from `configs/` directory
- Temporary build outputs in `/tmp/srediag-test-build`
- Test-specific container names with timestamps to avoid conflicts

### Extending Tests

When adding new integration tests:

1. Use the `// +build integration` build tag
2. Include proper cleanup in `defer` statements
3. Use appropriate timeouts for external services
4. Add logging for debugging purposes
5. Test both success and failure scenarios where appropriate

### Common Issues

#### Port Conflicts
If tests fail with port binding errors, check that no other services are using the test ports:

```bash
# Check port usage
netstat -tulpn | grep ":4317\|:4318\|:8889\|:13133"

# Kill processes using the ports if needed
sudo fuser -k 4317/tcp 4318/tcp 8889/tcp 13133/tcp
```

#### Docker Issues
If Docker tests fail:

```bash
# Check Docker is running
docker ps

# Clean up test containers and images
docker container prune -f
docker image prune -f

# Check Docker Compose
docker-compose --version
```

#### Build Issues
If build tests fail, ensure Go modules are up to date:

```bash
go mod download
go mod verify
```

## Test Coverage

Integration tests provide coverage for:

- ✅ Service startup and shutdown
- ✅ Health check endpoints
- ✅ Metrics collection endpoints
- ✅ OTLP receiver endpoints
- ✅ Diagnostics commands
- ✅ Plugin management
- ✅ Build commands
- ✅ Docker image building
- ✅ Docker container execution
- ✅ Docker Compose orchestration
- ✅ Service integration with Jaeger

These tests complement unit tests to ensure the complete system works as expected in realistic deployment scenarios.