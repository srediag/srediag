package deduplication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// Config holds the configuration for the deduplication processor
type Config struct {
	// MaxAge is the maximum age of cache entries
	MaxAge time.Duration `mapstructure:"max_age"`

	// RepeatAttr is the attribute name for repeat count
	RepeatAttr string `mapstructure:"repeat_attr"`
}

// Entry represents a cache entry
type Entry struct {
	Hash        string
	FirstSeen   time.Time
	LastSeen    time.Time
	RepeatCount int64
}

// Processor implements the OpenTelemetry SpanProcessor interface
type Processor struct {
	logger *zap.Logger
	config Config
	cache  map[string]Entry
	mu     sync.RWMutex
}

// NewProcessor creates a new deduplication processor
func NewProcessor(config Config, logger *zap.Logger) (*Processor, error) {
	return &Processor{
		logger: logger,
		config: config,
		cache:  make(map[string]Entry),
	}, nil
}

// OnStart is called when a span starts
func (p *Processor) OnStart(_ context.Context, _ sdktrace.ReadWriteSpan) {
	// No-op: deduplication happens on end
}

// OnEnd is called when a span ends
func (p *Processor) OnEnd(s sdktrace.ReadOnlySpan) {
	hash := p.computeHash(s)

	p.mu.Lock()
	defer p.mu.Unlock()

	if entry, exists := p.cache[hash]; exists {
		// Update existing entry
		entry.LastSeen = time.Now()
		entry.RepeatCount++
		p.cache[hash] = entry

		// Add repeat count attribute
		if rw, ok := s.(sdktrace.ReadWriteSpan); ok {
			rw.SetAttributes(
				attribute.Int64(p.config.RepeatAttr, entry.RepeatCount),
			)
		}
	} else {
		// Create new entry
		p.cache[hash] = Entry{
			Hash:        hash,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			RepeatCount: 1,
		}
	}
}

// Shutdown stops the processor
func (p *Processor) Shutdown(ctx context.Context) error {
	return nil
}

// computeHash generates a deterministic hash for a span
func (p *Processor) computeHash(s sdktrace.ReadOnlySpan) string {
	h := sha256.New()

	// Hash span name
	h.Write([]byte(s.Name()))

	// Hash attributes (sorted to ensure deterministic order)
	attrs := s.Attributes()
	sortedAttrs := make([]attribute.KeyValue, len(attrs))
	copy(sortedAttrs, attrs)

	// Sort attributes by key
	sort.Slice(sortedAttrs, func(i, j int) bool {
		return sortedAttrs[i].Key < sortedAttrs[j].Key
	})

	for _, attr := range sortedAttrs {
		h.Write([]byte(attr.Key))
		h.Write([]byte(attr.Value.AsString()))
	}

	return hex.EncodeToString(h.Sum(nil))
}

// ForceCleanup triggers cache cleanup
func (p *Processor) ForceCleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for key, entry := range p.cache {
		if now.Sub(entry.LastSeen) > p.config.MaxAge {
			delete(p.cache, key)
		}
	}
}
