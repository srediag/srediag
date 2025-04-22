# SREDIAG Plugin Types

SREDIAG supports three main types of plugins, each designed for specific functionality:

## 1. Diagnostic Plugins

Diagnostic plugins focus on system and application analysis.

### Core Features of Diagnostic Plugins

- System metrics collection
- Performance analysis
- Health checking
- Resource monitoring
- Event correlation

### Example Implementation of Diagnostic Plugins

```go
type DiagnosticPlugin struct {
    // OpenTelemetry integration
    receiver otlpreceiver.Receiver
    
    // Diagnostic specific fields
    collectors []Collector
    analyzers  []Analyzer
    reporters  []Reporter
}

type Collector interface {
    Collect(ctx context.Context) (*DiagnosticData, error)
    GetMetadata() CollectorMetadata
}

type Analyzer interface {
    Analyze(ctx context.Context, data *DiagnosticData) (*Analysis, error)
    GetCapabilities() []AnalyzerCapability
}

type Reporter interface {
    Report(ctx context.Context, analysis *Analysis) error
    GetFormats() []ReportFormat
}
```

### Configuration Example of Diagnostic Plugins

```yaml
plugins:
  diagnostic:
    system_analyzer:
      enabled: true
      version: "1.0.0"
      collectors:
        - name: "cpu"
          interval: "30s"
          thresholds:
            usage: 80
        - name: "memory"
          interval: "30s"
          thresholds:
            usage: 90
      analyzers:
        - name: "performance"
          correlation: true
        - name: "health"
          checks: ["all"]
      reporters:
        - name: "json"
        - name: "prometheus"
```

## 2. Analysis Plugins

Analysis plugins focus on data processing and insight generation.

### Core Features of Analysis Plugins

- Data processing
- Pattern recognition
- Anomaly detection
- Trend analysis
- Prediction models

### Example Implementation of Analysis Plugins

```go
type AnalysisPlugin struct {
    // OpenTelemetry integration
    processor processor.Processor
    
    // Analysis specific fields
    algorithms []Algorithm
    models     []Model
    results    []ResultHandler
}

type Algorithm interface {
    Process(ctx context.Context, data *AnalysisData) (*Result, error)
    GetParameters() []Parameter
}

type Model interface {
    Train(ctx context.Context, data *TrainingData) error
    Predict(ctx context.Context, input *InputData) (*Prediction, error)
}

type ResultHandler interface {
    Handle(ctx context.Context, result *Result) error
    GetOutputTypes() []OutputType
}
```

### Configuration Example of Analysis Plugins

```yaml
plugins:
  analysis:
    performance_analyzer:
      enabled: true
      version: "2.0.0"
      algorithms:
        - name: "anomaly_detection"
          sensitivity: "high"
          window: "1h"
        - name: "trend_analysis"
          period: "24h"
      models:
        - name: "resource_prediction"
          type: "lstm"
          parameters:
            layers: 3
            units: 64
      output:
        format: "json"
        destination: "file:///var/log/analysis"
```

## 3. Integration Plugins

Integration plugins handle connectivity with external systems.

### Core Features of Integration Plugins

- Protocol adaptation
- Data transformation
- Authentication
- Rate limiting
- Error handling

### Example Implementation of Integration Plugins

```go
type IntegrationPlugin struct {
    // OpenTelemetry integration
    exporter otlpexporter.Exporter
    
    // Integration specific fields
    connectors []Connector
    adapters   []Adapter
    transforms []Transform
}

type Connector interface {
    Connect(ctx context.Context, config *ConnectionConfig) error
    GetStatus() ConnectionStatus
}

type Adapter interface {
    Convert(ctx context.Context, data interface{}) (interface{}, error)
    GetSupportedTypes() []DataType
}

type Transform interface {
    Transform(ctx context.Context, data interface{}) (interface{}, error)
    GetCapabilities() []TransformCapability
}
```

### Configuration Example of Integration Plugins

```yaml
plugins:
  integration:
    oci_integration:
      enabled: true
      version: "1.0.0"
      connection:
        type: "oci"
        region: "us-ashburn-1"
        credentials:
          type: "instance_principal"
      services:
        - name: "compute"
          resources: ["instance", "volume"]
        - name: "monitoring"
          metrics: ["cpu", "memory"]
      transform:
        format: "otlp"
        mapping:
          cpu_utilization: "system.cpu.usage"
          memory_utilization: "system.memory.usage"
```

## Plugin Lifecycle

All plugins follow a standard lifecycle:

1. **Registration**

   ```go
   func (p *Plugin) Register(ctx context.Context) error
   ```

2. **Initialization**

   ```go
   func (p *Plugin) Init(ctx context.Context, config *Config) error
   ```

3. **Start**

   ```go
   func (p *Plugin) Start(ctx context.Context) error
   ```

4. **Stop**

   ```go
   func (p *Plugin) Stop(ctx context.Context) error
   ```

5. **Cleanup**

   ```go
   func (p *Plugin) Cleanup(ctx context.Context) error
   ```

## Plugin Development Guidelines

1. **OpenTelemetry Integration**
   - Use OpenTelemetry interfaces when possible
   - Follow OpenTelemetry data model
   - Support standard protocols

2. **Error Handling**
   - Provide detailed error information
   - Implement graceful degradation
   - Support retry mechanisms

3. **Configuration**
   - Use YAML for configuration
   - Support environment variables
   - Validate configuration

4. **Testing**
   - Write unit tests
   - Include integration tests
   - Provide example configurations

## Further Reading

- [Plugin Development Guide](development.md)
- [Plugin Best Practices](best-practices.md)
- [OpenTelemetry Integration](../architecture/opentelemetry.md)
