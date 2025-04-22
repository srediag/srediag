package diagnostic

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config/diagnostic"
)

// Component represents a diagnostic component
type Component interface {
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
	// IsHealthy returns the health status of the component
	IsHealthy() bool
	// GetType returns the diagnostic type
	GetType() Type
	// GetName returns the diagnostic name
	GetName() string
	// GetVersion returns the diagnostic version
	GetVersion() string
	// Configure configures the diagnostic with the given configuration
	Configure(cfg interface{}) error
}

// Type represents the type of a diagnostic component
type Type string

const (
	// TypeSystem represents a system diagnostic component
	TypeSystem Type = "system"
	// TypeKubernetes represents a Kubernetes diagnostic component
	TypeKubernetes Type = "kubernetes"
	// TypeCloud represents a cloud diagnostic component
	TypeCloud Type = "cloud"
	// TypeSecurity represents a security diagnostic component
	TypeSecurity Type = "security"
)

// Base provides a base implementation of Component
type Base struct {
	logger  *zap.Logger
	meter   metric.Meter
	typ     Type
	name    string
	version string
	healthy bool
	running bool

	// Metrics
	healthMetric metric.Float64ObservableGauge
	upMetric     metric.Float64ObservableGauge
}

// NewBase creates a new base diagnostic component
func NewBase(logger *zap.Logger, meter metric.Meter, typ Type, name, version string) *Base {
	if logger == nil {
		logger = zap.NewNop()
	}

	d := &Base{
		logger:  logger,
		meter:   meter,
		typ:     typ,
		name:    name,
		version: version,
		healthy: true,
	}

	// Initialize metrics
	var err error
	d.healthMetric, err = meter.Float64ObservableGauge(
		"srediag.diagnostic.health",
		metric.WithDescription("Health status of diagnostic component"),
		metric.WithUnit("1"),
	)
	if err != nil {
		logger.Error("failed to create health metric", zap.Error(err))
	}

	d.upMetric, err = meter.Float64ObservableGauge(
		"srediag.diagnostic.up",
		metric.WithDescription("Up status of diagnostic component"),
		metric.WithUnit("1"),
	)
	if err != nil {
		logger.Error("failed to create up metric", zap.Error(err))
	}

	// Register callbacks for metrics
	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		if d.healthMetric != nil {
			o.ObserveFloat64(d.healthMetric, boolToFloat64(d.healthy))
		}
		if d.upMetric != nil {
			o.ObserveFloat64(d.upMetric, boolToFloat64(d.running))
		}
		return nil
	}, d.healthMetric, d.upMetric)
	if err != nil {
		logger.Error("failed to register metric callbacks", zap.Error(err))
	}

	return d
}

// Start implements Component
func (d *Base) Start(ctx context.Context) error {
	d.running = true
	return nil
}

// Stop implements Component
func (d *Base) Stop(ctx context.Context) error {
	d.running = false
	return nil
}

// IsHealthy implements Component
func (d *Base) IsHealthy() bool {
	return d.healthy
}

// GetType implements Component
func (d *Base) GetType() Type {
	return d.typ
}

// GetName implements Component
func (d *Base) GetName() string {
	return d.name
}

// GetVersion implements Component
func (d *Base) GetVersion() string {
	return d.version
}

// Configure implements Component
func (d *Base) Configure(cfg interface{}) error {
	return nil
}

// boolToFloat64 converts a boolean to a float64
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// System represents a system diagnostic component
type System struct {
	*Base
	config *diagnostic.SystemConfig
}

// NewSystem creates a new system diagnostic component
func NewSystem(logger *zap.Logger, meter metric.Meter, cfg *diagnostic.SystemConfig) *System {
	return &System{
		Base: NewBase(
			logger,
			meter,
			TypeSystem,
			"system",
			"1.0.0",
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *System) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	config, ok := cfg.(*diagnostic.SystemConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Kubernetes represents a Kubernetes diagnostic component
type Kubernetes struct {
	*Base
	config *diagnostic.KubernetesConfig
}

// NewKubernetes creates a new Kubernetes diagnostic component
func NewKubernetes(logger *zap.Logger, meter metric.Meter, cfg *diagnostic.KubernetesConfig) *Kubernetes {
	return &Kubernetes{
		Base: NewBase(
			logger,
			meter,
			TypeKubernetes,
			"kubernetes",
			"1.0.0",
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Kubernetes) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	config, ok := cfg.(*diagnostic.KubernetesConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Cloud represents a cloud diagnostic component
type Cloud struct {
	*Base
	config *diagnostic.CloudConfig
}

// NewCloud creates a new cloud diagnostic component
func NewCloud(logger *zap.Logger, meter metric.Meter, cfg *diagnostic.CloudConfig) *Cloud {
	return &Cloud{
		Base: NewBase(
			logger,
			meter,
			TypeCloud,
			"cloud",
			"1.0.0",
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Cloud) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	config, ok := cfg.(*diagnostic.CloudConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}

// Security represents a security diagnostic component
type Security struct {
	*Base
	config *diagnostic.SecurityConfig
}

// NewSecurity creates a new security diagnostic component
func NewSecurity(logger *zap.Logger, meter metric.Meter, cfg *diagnostic.SecurityConfig) *Security {
	return &Security{
		Base: NewBase(
			logger,
			meter,
			TypeSecurity,
			"security",
			"1.0.0",
		),
		config: cfg,
	}
}

// Configure implements Component
func (d *Security) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	config, ok := cfg.(*diagnostic.SecurityConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}
