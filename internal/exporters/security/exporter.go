package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config holds the configuration for the security event exporter
type Config struct {
	// Endpoint is the SIEM API endpoint
	Endpoint string `mapstructure:"endpoint"`

	// AuthToken is the authentication token
	AuthToken string `mapstructure:"auth_token"`

	// BatchSize is the number of events to send in one request
	BatchSize int `mapstructure:"batch_size"`

	// FlushInterval is how often to flush the buffer
	FlushInterval time.Duration `mapstructure:"flush_interval"`

	// SeverityMapping maps event types to severity levels
	SeverityMapping map[string]string `mapstructure:"severity_mapping"`
}

// SecurityEvent represents a security-relevant event
type SecurityEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Type       string                 `json:"type"`
	Severity   string                 `json:"severity"`
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Action     string                 `json:"action"`
	Outcome    string                 `json:"outcome"`
	Hash       string                 `json:"hash,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
	RawData    interface{}            `json:"raw_data,omitempty"`
}

// Exporter implements the security event exporter
type Exporter struct {
	logger   *zap.Logger
	config   Config
	client   *http.Client
	tracer   trace.Tracer
	buffer   []*SecurityEvent
	mu       sync.Mutex
	stopChan chan struct{}
}

// NewExporter creates a new security event exporter
func NewExporter(config Config, logger *zap.Logger, tracer trace.Tracer) (*Exporter, error) {
	if config.Endpoint == "" {
		return nil, fmt.Errorf("SIEM endpoint is required")
	}

	if config.BatchSize == 0 {
		config.BatchSize = 100
	}

	if config.FlushInterval == 0 {
		config.FlushInterval = 10 * time.Second
	}

	if len(config.SeverityMapping) == 0 {
		config.SeverityMapping = map[string]string{
			"config_change": "INFO",
			"access_denied": "WARNING",
			"auth_failure":  "WARNING",
			"data_leak":     "CRITICAL",
			"malware":       "CRITICAL",
		}
	}

	return &Exporter{
		logger:   logger,
		config:   config,
		client:   &http.Client{Timeout: 30 * time.Second},
		tracer:   tracer,
		buffer:   make([]*SecurityEvent, 0, config.BatchSize),
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins the export process
func (e *Exporter) Start(ctx context.Context) error {
	go e.flushLoop(ctx)
	return nil
}

// Stop stops the export process
func (e *Exporter) Stop(ctx context.Context) error {
	close(e.stopChan)
	return e.flush(ctx)
}

// Export exports a security event
func (e *Exporter) Export(ctx context.Context, event *SecurityEvent) error {
	ctx, span := e.tracer.Start(ctx, "security.export")
	defer span.End()

	// Set severity based on mapping
	if severity, ok := e.config.SeverityMapping[event.Type]; ok {
		event.Severity = severity
	}

	span.SetAttributes(
		attribute.String("event.id", event.ID),
		attribute.String("event.type", event.Type),
		attribute.String("event.severity", event.Severity),
		attribute.String("event.source", event.Source),
	)

	e.mu.Lock()
	e.buffer = append(e.buffer, event)
	shouldFlush := len(e.buffer) >= e.config.BatchSize
	e.mu.Unlock()

	if shouldFlush {
		return e.flush(ctx)
	}

	return nil
}

// flushLoop periodically flushes the buffer
func (e *Exporter) flushLoop(ctx context.Context) {
	ticker := time.NewTicker(e.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopChan:
			return
		case <-ticker.C:
			if err := e.flush(ctx); err != nil {
				e.logger.Error("failed to flush security events", zap.Error(err))
			}
		}
	}
}

// flush sends buffered events to SIEM
func (e *Exporter) flush(ctx context.Context) error {
	e.mu.Lock()
	if len(e.buffer) == 0 {
		e.mu.Unlock()
		return nil
	}

	events := make([]*SecurityEvent, len(e.buffer))
	copy(events, e.buffer)
	e.buffer = e.buffer[:0]
	e.mu.Unlock()

	data, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.config.Endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if e.config.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.config.AuthToken))
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SIEM API error: status=%d", resp.StatusCode)
	}

	return nil
}
