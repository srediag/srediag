package telemetry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/lifecycle"
)

const (
	defaultMetricInterval = 10 // segundos
)

// Manager gerencia a configuração e coleta de telemetria
type Manager struct {
	*lifecycle.BaseManager
	logger     *zap.Logger
	config     config.TelemetryConfig
	version    string
	mu         sync.RWMutex
	tracerProv *sdktrace.TracerProvider
	meterProv  *sdkmetric.MeterProvider
}

// NewManager cria uma nova instância do gerenciador de telemetria
func NewManager(cfg config.TelemetryConfig, version string, logger *zap.Logger) (*Manager, error) {
	if !cfg.Enabled {
		logger.Info("telemetria desabilitada pela configuração")
		return &Manager{
			BaseManager: lifecycle.NewBaseManager(),
			logger:      logger,
			config:      cfg,
			version:     version,
		}, nil
	}

	return &Manager{
		BaseManager: lifecycle.NewBaseManager(),
		logger:      logger,
		config:      cfg,
		version:     version,
	}, nil
}

// Start inicializa e inicia a coleta de telemetria
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(false); err != nil {
		return err
	}

	// Se a telemetria estiver desabilitada, retorna sem erro
	if !m.config.Enabled {
		m.logger.Info("telemetria desabilitada, não será iniciada")
		return nil
	}

	// Cria recurso com informações do serviço
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(m.config.ServiceName),
			semconv.ServiceVersion(m.version),
			semconv.DeploymentEnvironment(m.config.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("falha ao criar recurso: %v", err)
	}

	// Configura conexão gRPC
	//nolint:staticcheck // TODO: Atualizar para NewClient em uma versão futura
	conn, err := grpc.Dial(m.config.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("falha ao criar conexão gRPC: %v", err)
	}

	// Configura exportador de traces
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return fmt.Errorf("falha ao criar exportador de traces: %v", err)
	}

	// Configura amostragem baseada na configuração
	var sampler sdktrace.Sampler
	switch m.config.Sampling.Type {
	case "always_on":
		sampler = sdktrace.AlwaysSample()
	case "always_off":
		sampler = sdktrace.NeverSample()
	case "probabilistic":
		sampler = sdktrace.TraceIDRatioBased(m.config.Sampling.Rate)
	default:
		sampler = sdktrace.AlwaysSample()
	}

	// Configura provedor de traces
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sampler),
	)
	m.tracerProv = tracerProvider

	// Configura exportador de métricas
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return fmt.Errorf("falha ao criar exportador de métricas: %v", err)
	}

	// Configura provedor de métricas
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(time.Duration(defaultMetricInterval)*time.Second),
			),
		),
	)
	m.meterProv = meterProvider

	// Configura provedores globais
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	m.SetRunning(true)
	m.logger.Info("gerenciador de telemetria iniciado com sucesso")

	return nil
}

// Stop desliga a coleta de telemetria de forma graciosa
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(true); err != nil {
		return err
	}

	// Se a telemetria estiver desabilitada, retorna sem erro
	if !m.config.Enabled {
		return nil
	}

	var errs []error

	// Desliga provedor de traces
	if m.tracerProv != nil {
		if err := m.tracerProv.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("falha ao desligar provedor de traces: %v", err))
		}
	}

	// Desliga provedor de métricas
	if m.meterProv != nil {
		if err := m.meterProv.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("falha ao desligar provedor de métricas: %v", err))
		}
	}

	m.SetRunning(false)

	if len(errs) > 0 {
		return fmt.Errorf("falha ao desligar telemetria: %v", errs)
	}

	m.logger.Info("gerenciador de telemetria parado com sucesso")
	return nil
}

// Tracer retorna um tracer nomeado
func (m *Manager) Tracer(name string) trace.Tracer {
	if !m.config.Enabled || m.tracerProv == nil {
		return noop.NewTracerProvider().Tracer(name)
	}
	return m.tracerProv.Tracer(name)
}

// Meter retorna um meter nomeado
func (m *Manager) Meter(name string) metric.Meter {
	if !m.config.Enabled || m.meterProv == nil {
		return sdkmetric.NewMeterProvider().Meter(name)
	}
	return m.meterProv.Meter(name)
}
