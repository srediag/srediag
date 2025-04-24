package diagnostic

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Component represents a diagnostic component
type Component interface {
	// GetName returns the component name
	GetName() string

	// GetVersion returns the component version
	GetVersion() string

	// GetType returns the component type
	GetType() types.ComponentType

	// GetStatus returns the component status
	GetStatus() types.ComponentStatus

	// IsHealthy returns the component health status
	IsHealthy() bool

	// Configure configures the component with settings
	Configure(settings types.ComponentSettings) error

	// Start starts the component
	Start(ctx context.Context) error

	// Stop stops the component
	Stop(ctx context.Context) error

	// Collect implements IDiagnostic
	Collect(ctx context.Context) (map[string]interface{}, error)

	// GetInterval implements IDiagnostic
	GetInterval() string

	// GetThresholds implements IDiagnostic
	GetThresholds() map[string]float64
}

// Base provides a base implementation of Component
type Base struct {
	mu      sync.RWMutex
	logger  *zap.Logger
	tracer  trace.Tracer
	meter   metric.Meter
	healthy bool
	running bool
	name    string
	version string
	ctype   types.ComponentType

	// Metrics
	healthMetric metric.Int64UpDownCounter
	upMetric     metric.Int64UpDownCounter
}

// NewBase creates a new base component
func NewBase(name string, version string, ctype types.ComponentType, logger *zap.Logger, meter metric.Meter, tracer trace.Tracer) *Base {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Base{
		name:    name,
		version: version,
		ctype:   ctype,
		logger:  logger,
		meter:   meter,
		tracer:  tracer,
		healthy: true,
	}
}

// GetName returns the component name
func (b *Base) GetName() string {
	return b.name
}

// GetVersion returns the component version
func (b *Base) GetVersion() string {
	return b.version
}

// GetType returns the component type
func (b *Base) GetType() types.ComponentType {
	return b.ctype
}

// GetStatus returns the component status
func (b *Base) GetStatus() types.ComponentStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.running {
		return types.ComponentStatusStopped
	}
	if !b.healthy {
		return types.ComponentStatusError
	}
	return types.ComponentStatusRunning
}

// IsHealthy returns the component health status
func (b *Base) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthy
}

// Configure configures the component with settings
func (b *Base) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	telemetrySettings := settings.GetTelemetrySettings()
	if telemetrySettings == nil {
		return fmt.Errorf("telemetry settings cannot be nil")
	}

	b.logger = telemetrySettings.Logger
	b.tracer = telemetrySettings.Tracer
	b.meter = telemetrySettings.Meter

	// Initialize metrics
	var err error
	b.healthMetric, err = b.meter.Int64UpDownCounter(
		fmt.Sprintf("%s_health", b.name),
		metric.WithDescription("Health status of the component"),
	)
	if err != nil {
		return fmt.Errorf("failed to create health metric: %w", err)
	}

	b.upMetric, err = b.meter.Int64UpDownCounter(
		fmt.Sprintf("%s_up", b.name),
		metric.WithDescription("Running status of the component"),
	)
	if err != nil {
		return fmt.Errorf("failed to create up metric: %w", err)
	}

	return nil
}

// Start starts the component
func (b *Base) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return nil
	}

	b.running = true
	b.upMetric.Add(ctx, 1)
	b.logger.Info("Started component", zap.String("name", b.name))
	return nil
}

// Stop stops the component
func (b *Base) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	b.running = false
	b.upMetric.Add(ctx, -1)
	b.logger.Info("Stopped component", zap.String("name", b.name))
	return nil
}

// System represents a system diagnostic component
type System struct {
	*Base
	config *types.SystemConfig
}

// NewSystem creates a new system diagnostic component
func NewSystem(logger *zap.Logger, meter metric.Meter, cfg *types.SystemConfig) *System {
	return &System{
		Base: NewBase(
			"system",
			"1.0.0",
			types.ComponentTypeCore,
			logger,
			meter,
			nil,
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *System) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	config, ok := settings.GetInterface("config").(*types.SystemConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Collect implements IDiagnostic
func (d *System) Collect(ctx context.Context) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

// GetInterval implements IDiagnostic
func (d *System) GetInterval() string {
	return d.config.Interval
}

// GetThresholds implements IDiagnostic
func (d *System) GetThresholds() map[string]float64 {
	return map[string]float64{
		"cpu":    d.config.CPULimit,
		"memory": d.config.MemoryLimit,
		"disk":   d.config.DiskLimit,
	}
}

// Kubernetes represents a Kubernetes diagnostic component
type Kubernetes struct {
	*Base
	config *types.KubernetesConfig
}

// NewKubernetes creates a new Kubernetes diagnostic component
func NewKubernetes(logger *zap.Logger, meter metric.Meter, cfg *types.KubernetesConfig) *Kubernetes {
	return &Kubernetes{
		Base: NewBase(
			"kubernetes",
			"1.0.0",
			types.ComponentTypeCore,
			logger,
			meter,
			nil,
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Kubernetes) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	config, ok := settings.GetInterface("config").(*types.KubernetesConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Collect implements IDiagnostic
func (d *Kubernetes) Collect(ctx context.Context) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

// GetInterval implements IDiagnostic
func (d *Kubernetes) GetInterval() string {
	return "30s" // Default interval for Kubernetes diagnostics
}

// GetThresholds implements IDiagnostic
func (d *Kubernetes) GetThresholds() map[string]float64 {
	return make(map[string]float64) // No default thresholds for Kubernetes
}

// Cloud represents a cloud diagnostic component
type Cloud struct {
	*Base
	config *types.CloudConfig
}

// NewCloud creates a new cloud diagnostic component
func NewCloud(logger *zap.Logger, meter metric.Meter, cfg *types.CloudConfig) *Cloud {
	return &Cloud{
		Base: NewBase(
			"cloud",
			"1.0.0",
			types.ComponentTypeCore,
			logger,
			meter,
			nil,
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Cloud) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	config, ok := settings.GetInterface("config").(*types.CloudConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Collect implements IDiagnostic
func (d *Cloud) Collect(ctx context.Context) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

// GetInterval implements IDiagnostic
func (d *Cloud) GetInterval() string {
	return "60s" // Default interval for Cloud diagnostics
}

// GetThresholds implements IDiagnostic
func (d *Cloud) GetThresholds() map[string]float64 {
	return make(map[string]float64) // No default thresholds for Cloud
}

// Security represents a security diagnostic component
type Security struct {
	*Base
	config *types.SecurityConfig
}

// NewSecurity creates a new security diagnostic component
func NewSecurity(logger *zap.Logger, meter metric.Meter, cfg *types.SecurityConfig) *Security {
	return &Security{
		Base: NewBase(
			"security",
			"1.0.0",
			types.ComponentTypeCore,
			logger,
			meter,
			nil,
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Security) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	config, ok := cfg.(*types.SecurityConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}
