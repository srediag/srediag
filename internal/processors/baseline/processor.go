package baseline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config holds the configuration for the baseline processor
type Config struct {
	// MonitoredPaths are the file paths to monitor (supports glob patterns)
	MonitoredPaths []string `mapstructure:"monitored_paths"`

	// BaselineStore is the path to store baseline hashes
	BaselineStore string `mapstructure:"baseline_store"`

	// ScanInterval is how often to scan for changes
	ScanInterval time.Duration `mapstructure:"scan_interval"`
}

// Store interface for baseline hash storage
type Store interface {
	Get(path string) (string, error)
	Set(path, hash string) error
	Delete(path string) error
	Close() error
}

// FileInfo represents a monitored file's information
type FileInfo struct {
	Path    string
	Hash    string
	ModTime time.Time
}

// Processor implements file monitoring and baseline comparison
type Processor struct {
	logger    *zap.Logger
	config    Config
	store     Store
	tracer    trace.Tracer
	mu        sync.RWMutex
	stopChan  chan struct{}
	fileInfos map[string]FileInfo
}

// NewProcessor creates a new baseline processor
func NewProcessor(config Config, logger *zap.Logger, tracer trace.Tracer) (*Processor, error) {
	store, err := newBoltStore(config.BaselineStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create baseline store: %w", err)
	}

	return &Processor{
		logger:    logger,
		config:    config,
		store:     store,
		tracer:    tracer,
		stopChan:  make(chan struct{}),
		fileInfos: make(map[string]FileInfo),
	}, nil
}

// Start begins monitoring files
func (p *Processor) Start(ctx context.Context) error {
	// Initial scan
	if err := p.scanFiles(ctx); err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}

	// Start monitoring goroutine
	go p.monitor(ctx)

	return nil
}

// Stop stops monitoring files
func (p *Processor) Stop(ctx context.Context) error {
	close(p.stopChan)
	return p.store.Close()
}

// monitor periodically scans for file changes
func (p *Processor) monitor(ctx context.Context) {
	ticker := time.NewTicker(p.config.ScanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			if err := p.scanFiles(ctx); err != nil {
				p.logger.Error("file scan failed", zap.Error(err))
			}
		}
	}
}

// scanFiles checks all monitored paths for changes
func (p *Processor) scanFiles(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "baseline.scanFiles")
	defer span.End()

	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pattern := range p.config.MonitoredPaths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("glob pattern error: %w", err)
		}

		for _, path := range matches {
			if err := p.checkFile(ctx, path); err != nil {
				p.logger.Error("file check failed",
					zap.String("path", path),
					zap.Error(err))
				continue
			}
		}
	}

	return nil
}

// checkFile verifies if a single file has changed
func (p *Processor) checkFile(ctx context.Context, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}

	// Skip if file hasn't been modified
	if existing, ok := p.fileInfos[path]; ok {
		if !info.ModTime().After(existing.ModTime) {
			return nil
		}
	}

	// Calculate new hash
	hash, err := p.hashFile(path)
	if err != nil {
		return fmt.Errorf("hash calculation failed: %w", err)
	}

	// Check if hash has changed
	oldHash, err := p.store.Get(path)
	if err == nil && oldHash == hash {
		// Update mod time but hash hasn't changed
		p.fileInfos[path] = FileInfo{
			Path:    path,
			Hash:    hash,
			ModTime: info.ModTime(),
		}
		return nil
	}

	// Hash has changed or is new, emit change event
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("file.path", path),
		attribute.String("file.hash.old", oldHash),
		attribute.String("file.hash.new", hash),
		attribute.String("file.event", "changed"),
	)

	// Update store and cache
	if err := p.store.Set(path, hash); err != nil {
		return fmt.Errorf("failed to update store: %w", err)
	}

	p.fileInfos[path] = FileInfo{
		Path:    path,
		Hash:    hash,
		ModTime: info.ModTime(),
	}

	return nil
}

// hashFile calculates SHA-256 hash of a file
func (p *Processor) hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash calculation failed: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
