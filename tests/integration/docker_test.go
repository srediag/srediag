//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDockerBuild tests that the Docker image builds successfully
func TestDockerBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Build the Docker image
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", "srediag:test", ".")
	cmd.Dir = "../../"

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker build output: %s", string(output))
		t.Fatalf("Failed to build Docker image: %v", err)
	}

	t.Logf("Docker image built successfully")

	// Clean up the image after test
	defer func() {
		cleanupCmd := exec.Command("docker", "rmi", "srediag:test")
		cleanupCmd.Run()
	}()

	// Test that the image was created
	cmd = exec.Command("docker", "images", "srediag:test", "--format", "{{.Repository}}:{{.Tag}}")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Failed to list Docker images")

	assert.Contains(t, string(output), "srediag:test", "Docker image should exist")
}

// TestDockerRun tests running the Docker container
func TestDockerRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// First build the image
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	buildCmd := exec.CommandContext(ctx, "docker", "build", "-t", "srediag:test", ".")
	buildCmd.Dir = "../../"
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build Docker image for run test")

	// Clean up the image after test
	defer func() {
		cleanupCmd := exec.Command("docker", "rmi", "srediag:test")
		cleanupCmd.Run()
	}()

	// Start the container
	containerName := "srediag-test-" + fmt.Sprintf("%d", time.Now().Unix())

	runCtx, runCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer runCancel()

	runCmd := exec.CommandContext(runCtx, "docker", "run", "--rm", "-d",
		"-p", "14317:4317", // Use different host ports to avoid conflicts
		"-p", "14318:4318",
		"-p", "18889:8889",
		"-p", "23133:13133",
		"--name", containerName,
		"srediag:test")

	err = runCmd.Run()
	require.NoError(t, err, "Failed to start Docker container")

	// Clean up container
	defer func() {
		stopCmd := exec.Command("docker", "stop", containerName)
		stopCmd.Run()
	}()

	// Wait for container to start
	time.Sleep(10 * time.Second)

	// Test that container is running
	statusCmd := exec.Command("docker", "ps", "--filter", "name="+containerName, "--format", "{{.Status}}")
	output, err := statusCmd.CombinedOutput()
	require.NoError(t, err, "Failed to check container status")

	assert.Contains(t, string(output), "Up", "Container should be running")

	// Test health endpoint (using host port 23133)
	t.Run("ContainerHealthCheck", func(t *testing.T) {
		// Wait a bit more for service to be ready
		time.Sleep(5 * time.Second)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get("http://localhost:23133/")

		if err != nil {
			// Check container logs for debugging
			logsCmd := exec.Command("docker", "logs", containerName)
			logs, _ := logsCmd.CombinedOutput()
			t.Logf("Container logs: %s", string(logs))

			t.Fatalf("Health check request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Container health check should return 200")
	})
}

// TestDockerCompose tests the docker-compose setup
func TestDockerCompose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Check if docker-compose is available
	if _, err := exec.LookPath("docker-compose"); err != nil {
		t.Skip("docker-compose not available, skipping test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start services
	upCmd := exec.CommandContext(ctx, "docker-compose", "-f", "../../docker-compose.yml", "up", "-d", "srediag", "jaeger")
	upCmd.Dir = "../../"

	output, err := upCmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker compose up output: %s", string(output))
		t.Fatalf("Failed to start docker-compose services: %v", err)
	}

	// Clean up services after test
	defer func() {
		downCmd := exec.Command("docker-compose", "-f", "../../docker-compose.yml", "down", "-v")
		downCmd.Dir = "../../"
		downCmd.Run()
	}()

	// Wait for services to start
	time.Sleep(30 * time.Second)

	// Check that services are running
	psCmd := exec.Command("docker-compose", "-f", "../../docker-compose.yml", "ps")
	psCmd.Dir = "../../"
	output, err = psCmd.CombinedOutput()
	require.NoError(t, err, "Failed to check docker-compose services")

	servicesOutput := string(output)
	t.Logf("Docker compose services: %s", servicesOutput)

	// Check specific services
	assert.Contains(t, servicesOutput, "srediag", "SREDIAG service should be running")
	assert.Contains(t, servicesOutput, "jaeger", "Jaeger service should be running")

	// Test SREDIAG health
	t.Run("DockerComposeHealthCheck", func(t *testing.T) {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get("http://localhost:13133/")

		if err != nil {
			// Get container logs for debugging
			logsCmd := exec.Command("docker-compose", "-f", "../../docker-compose.yml", "logs", "srediag")
			logsCmd.Dir = "../../"
			logs, _ := logsCmd.CombinedOutput()
			t.Logf("SREDIAG logs: %s", string(logs))

			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "SREDIAG health check should return 200")
	})

	// Test Jaeger UI
	t.Run("JaegerUI", func(t *testing.T) {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get("http://localhost:16686/")

		if err != nil {
			t.Logf("Jaeger UI request failed: %v", err)
			return // Jaeger might take longer to start
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Jaeger UI should be accessible")
	})
}

// Helper function to check if a container is running
func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "--filter", "name="+containerName, "--format", "{{.Names}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), containerName)
}
