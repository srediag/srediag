package diagnostic

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Component represents a diagnostic component
type Component interface {
	types.IDiagnostic
}

// Base provides a base implementation of Component
type Base struct {
	logger  *zap.Logger
	meter   metric.Meter
	typ     types.ComponentType
	name    string
	version string
	healthy bool
	running bool

	// Metrics
	healthMetric metric.Float64ObservableGauge
	upMetric     metric.Float64ObservableGauge
}

// NewBase creates a new base diagnostic component
func NewBase(logger *zap.Logger, meter metric.Meter, typ types.ComponentType, name, version string) *Base {
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
func (d *Base) GetType() types.ComponentType {
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

// GetInterval implements Component
func (d *Base) GetInterval() string {
	return "30s" // Default interval
}

// GetThresholds implements Component
func (d *Base) GetThresholds() map[string]float64 {
	return make(map[string]float64) // Default empty thresholds
}

// Collect implements Component
func (d *Base) Collect(ctx context.Context) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil // Default empty collection
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
			logger,
			meter,
			types.ComponentTypeCore,
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
	config, ok := cfg.(*types.SystemConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
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
			logger,
			meter,
			types.ComponentTypeCore,
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
	config, ok := cfg.(*types.KubernetesConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
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
			logger,
			meter,
			types.ComponentTypeCore,
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
	config, ok := cfg.(*types.CloudConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
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
			logger,
			meter,
			types.ComponentTypeCore,
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
	config, ok := cfg.(*types.SecurityConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}
	d.config = config
	return nil
}
