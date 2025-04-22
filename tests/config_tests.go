// tests/config_test.go
package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/srediag/srediag/internal/config"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := []byte(`
version: "v0.1.0"
debug: false
log_level: "info"
service:
  name: "srediag-test"
  environment: "test"
telemetry:
  enabled: true
  service_name: "srediag-test"
  endpoint: "localhost:4317"
`)

	err := os.WriteFile(configPath, configContent, 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Load the configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test configuration values
	if cfg.Version != "v0.1.0" {
		t.Errorf("expected version v0.1.0, got %s", cfg.Version)
	}

	if cfg.Service.Name != "srediag-test" {
		t.Errorf("expected service name srediag-test, got %s", cfg.Service.Name)
	}

	if cfg.Telemetry.Endpoint != "localhost:4317" {
		t.Errorf("expected telemetry endpoint localhost:4317, got %s", cfg.Telemetry.Endpoint)
	}
}
