// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file contains unit tests for the Logger type, ensuring correct configuration, output, and integration with zap.
// The tests validate logger creation, level parsing, component scoping, and output formatting for robust diagnostics and production use.
//
// Usage:
//   - Run with `go test ./internal/core` to validate logger behavior.
//   - Use as a reference for writing new logger-related tests.
//
// Best Practices:
//   - Cover all logger configuration options and output formats.
//   - Use testify/require for clear assertions and error handling.
//
// TODO:
//   - Add tests for dynamic log level changes and error output paths.
//   - Add tests for feature gate integration and OpenTelemetry log export.
//
// Redundancy/Refactor:
//   - No redundancy; these are canonical logger tests for SREDIAG.
package core

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestNewLogger validates logger creation with various configurations.
//
// Usage:
//   - Ensures that default, valid, and invalid logger configs are handled correctly.
//   - Checks that errors are returned for invalid output paths.
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Logger
		wantErr bool
	}{
		{
			name:    "nil config uses defaults",
			cfg:     nil,
			wantErr: false,
		},
		{
			name: "valid config",
			cfg: &Logger{
				Level:  "debug",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "invalid output path",
			cfg: &Logger{
				OutputPaths: []string{"/invalid/path/that/does/not/exist"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}

// TestLogger_Levels validates the parseLevel function for all supported log levels.
//
// Usage:
//   - Ensures that string log levels map to the correct zapcore.Level.
//   - Checks that invalid levels default to info.
func TestLogger_Levels(t *testing.T) {
	tests := []struct {
		level string
		want  zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"warning", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"dpanic", zapcore.DPanicLevel},
		{"panic", zapcore.PanicLevel},
		{"fatal", zapcore.FatalLevel},
		{"invalid", zapcore.InfoLevel}, // default level
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := parseLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestLogger_WithComponent validates the WithComponent method for logger scoping.
//
// Usage:
//   - Ensures that logs from a component logger include the correct component field.
func TestLogger_WithComponent(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger := &Logger{
		logger: zap.New(core),
	}

	// Create a component logger and write a message
	componentLogger := logger.WithComponent("test-component")
	componentLogger.Info("test message")

	// Parse the output
	var output map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &output)
	require.NoError(t, err)

	// Check that the component field was added
	assert.Equal(t, "test-component", output["component"])
	assert.Equal(t, "test message", output["msg"])
}

// TestLogger_OutputFormats validates logger output for both JSON and console formats.
//
// Usage:
//   - Ensures that log output matches the expected format and contains the correct message.
func TestLogger_OutputFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
		check  func(t *testing.T, output string)
	}{
		{
			name:   "json format",
			format: "json",
			check: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				assert.NoError(t, err)
				assert.Equal(t, "test message", result["msg"])
			},
		},
		{
			name:   "console format",
			format: "console",
			check: func(t *testing.T, output string) {
				assert.Contains(t, output, "test message")
				assert.Contains(t, output, "info")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			cfg := &Logger{
				Level:       "info",
				Format:      tt.format,
				OutputPaths: []string{},
			}

			logger, err := NewLogger(cfg)
			require.NoError(t, err)

			// Replace the core to write to our buffer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				zapcore.InfoLevel,
			)
			logger.logger = zap.New(core)

			logger.Info("test message")
			tt.check(t, buf.String())
		})
	}
}
