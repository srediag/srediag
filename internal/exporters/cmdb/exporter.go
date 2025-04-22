package cmdb

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

// Config holds the configuration for the CMDB exporter
type Config struct {
	// Endpoint is the CMDB API endpoint
	Endpoint string `mapstructure:"endpoint"`

	// AuthToken is the authentication token
	AuthToken string `mapstructure:"auth_token"`

	// BatchSize is the number of items to send in one request
	BatchSize int `mapstructure:"batch_size"`

	// FlushInterval is how often to flush the buffer
	FlushInterval time.Duration `mapstructure:"flush_interval"`
}

// ConfigurationItem represents a CMDB configuration item
type ConfigurationItem struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Name          string                 `json:"name"`
	Environment   string                 `json:"environment"`
	State         string                 `json:"state"`
	Hash          string                 `json:"hash,omitempty"`
	LastModified  time.Time              `json:"last_modified"`
	Attributes    map[string]interface{} `json:"attributes"`
	Relationships []Relationship         `json:"relationships,omitempty"`
}

// Relationship represents a relationship between CIs
type Relationship struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

// Exporter implements the CMDB exporter
type Exporter struct {
	logger   *zap.Logger
	config   Config
	client   *http.Client
	tracer   trace.Tracer
	buffer   []*ConfigurationItem
	mu       sync.Mutex
	stopChan chan struct{}
}

// NewExporter creates a new CMDB exporter
func NewExporter(config Config, logger *zap.Logger, tracer trace.Tracer) (*Exporter, error) {
	if config.Endpoint == "" {
		return nil, fmt.Errorf("CMDB endpoint is required")
	}

	if config.BatchSize == 0 {
		config.BatchSize = 100
	}

	if config.FlushInterval == 0 {
		config.FlushInterval = 30 * time.Second
	}

	return &Exporter{
		logger:   logger,
		config:   config,
		client:   &http.Client{Timeout: 30 * time.Second},
		tracer:   tracer,
		buffer:   make([]*ConfigurationItem, 0, config.BatchSize),
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

// Export exports a configuration item to CMDB
func (e *Exporter) Export(ctx context.Context, ci *ConfigurationItem) error {
	ctx, span := e.tracer.Start(ctx, "cmdb.export")
	defer span.End()

	span.SetAttributes(
		attribute.String("ci.id", ci.ID),
		attribute.String("ci.type", ci.Type),
		attribute.String("ci.name", ci.Name),
	)

	e.mu.Lock()
	e.buffer = append(e.buffer, ci)
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
				e.logger.Error("failed to flush CMDB buffer", zap.Error(err))
			}
		}
	}
}

// flush sends buffered items to CMDB
func (e *Exporter) flush(ctx context.Context) error {
	e.mu.Lock()
	if len(e.buffer) == 0 {
		e.mu.Unlock()
		return nil
	}

	items := make([]*ConfigurationItem, len(e.buffer))
	copy(items, e.buffer)
	e.buffer = e.buffer[:0]
	e.mu.Unlock()

	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
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
		return fmt.Errorf("CMDB API error: status=%d", resp.StatusCode)
	}

	return nil
}
