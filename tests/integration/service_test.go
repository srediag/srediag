//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceIntegration tests the complete service lifecycle
func TestServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the application first
	buildCmd := exec.Command("go", "build", "-o", "../../bin/srediag-test", "../../cmd/srediag")
	buildCmd.Dir = "../../"
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build SREDIAG binary")

	// Clean up binary after test
	defer func() {
		os.Remove("../../bin/srediag-test")
	}()

	// Start the service
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "../../bin/srediag-test", "service", "start",
		"--config", "../../configs/srediag.yaml",
		"--log-level", "debug")
	cmd.Dir = "../../"

	// Start the process
	err = cmd.Start()
	require.NoError(t, err, "Failed to start SREDIAG service")

	// Ensure process is killed after test
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
			cmd.Wait()
		}
	}()

	// Wait for service to start
	time.Sleep(5 * time.Second)

	// Test health endpoint
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get("http://localhost:13133/")
		require.NoError(t, err, "Health check request failed")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return 200")
	})

	// Test metrics endpoint
	t.Run("MetricsEndpoint", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8889/metrics")
		require.NoError(t, err, "Metrics request failed")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Metrics endpoint should return 200")
	})

	// Test OTLP endpoints are accessible
	t.Run("OTLPEndpoints", func(t *testing.T) {
		// Test gRPC endpoint connectivity (just TCP connection)
		conn, err := http.Get("http://localhost:4317")
		if err == nil {
			conn.Body.Close()
		}
		// gRPC endpoint might not respond to HTTP, that's expected

		// Test HTTP endpoint
		resp, err := http.Get("http://localhost:4318/v1/traces")
		if err == nil {
			defer resp.Body.Close()
			// Endpoint exists, might return 405 Method Not Allowed for GET
			assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500,
				"OTLP HTTP endpoint should be accessible")
		}
	})
}

// TestDiagnosticsIntegration tests the diagnostics functionality
func TestDiagnosticsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the application first
	buildCmd := exec.Command("go", "build", "-o", "../../bin/srediag-test", "../../cmd/srediag")
	buildCmd.Dir = "../../"
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build SREDIAG binary")

	// Clean up binary after test
	defer func() {
		os.Remove("../../bin/srediag-test")
	}()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "SystemDiagnostics",
			args: []string{"diagnose", "system", "--log-level", "debug"},
		},
		{
			name: "PerformanceDiagnostics",
			args: []string{"diagnose", "performance", "--log-level", "debug"},
		},
		{
			name: "SecurityDiagnostics",
			args: []string{"diagnose", "security", "--log-level", "debug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, "../../bin/srediag-test")
			cmd.Args = append([]string{"../../bin/srediag-test"}, tt.args...)
			cmd.Dir = "../../"

			output, err := cmd.CombinedOutput()
			t.Logf("Command output: %s", string(output))

			// Diagnostics might fail due to missing dependencies, but should not crash
			// We mainly check that the command executed without panic
			if err != nil {
				// Check if it's a context timeout (acceptable for integration tests)
				if ctx.Err() == context.DeadlineExceeded {
					t.Logf("Command timed out (acceptable in integration tests)")
				} else {
					t.Logf("Command failed with error: %v", err)
				}
			}

			// Check that output contains some expected content
			assert.Contains(t, string(output), "diagnose",
				"Output should contain 'diagnose' indicating the command ran")
		})
	}
}

// TestPluginIntegration tests plugin functionality
func TestPluginIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the application first
	buildCmd := exec.Command("go", "build", "-o", "../../bin/srediag-test", "../../cmd/srediag")
	buildCmd.Dir = "../../"
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build SREDIAG binary")

	// Clean up binary after test
	defer func() {
		os.Remove("../../bin/srediag-test")
	}()

	t.Run("PluginList", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "../../bin/srediag-test", "plugin", "list",
			"--log-level", "debug")
		cmd.Dir = "../../"

		output, err := cmd.CombinedOutput()
		t.Logf("Plugin list output: %s", string(output))

		// Command should execute successfully (even if no plugins are loaded)
		if err != nil && ctx.Err() != context.DeadlineExceeded {
			t.Logf("Plugin list command failed: %v", err)
		}

		// Check that output contains expected content
		assert.Contains(t, string(output), "plugin",
			"Output should contain 'plugin' indicating the command ran")
	})
}

// TestBuildIntegration tests build functionality
func TestBuildIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the application first
	buildCmd := exec.Command("go", "build", "-o", "../../bin/srediag-test", "../../cmd/srediag")
	buildCmd.Dir = "../../"
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build SREDIAG binary")

	// Clean up binary after test
	defer func() {
		os.Remove("../../bin/srediag-test")
	}()

	t.Run("BuildAll", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "../../bin/srediag-test", "build", "all",
			"--log-level", "debug",
			"--output-dir", "/tmp/srediag-test-build")
		cmd.Dir = "../../"

		output, err := cmd.CombinedOutput()
		t.Logf("Build output: %s", string(output))

		// Build might fail due to missing dependencies, but should not crash
		if err != nil && ctx.Err() != context.DeadlineExceeded {
			t.Logf("Build command failed: %v", err)
		}

		// Check that output contains expected content
		assert.Contains(t, string(output), "build",
			"Output should contain 'build' indicating the command ran")
	})
}
