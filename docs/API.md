# SREDIAG API Reference

This document provides comprehensive API documentation for SREDIAG components.

## Table of Contents

- [Core API](#core-api)
- [Plugin API](#plugin-api)
- [Diagnostics API](#diagnostics-api)
- [Service API](#service-api)
- [Configuration API](#configuration-api)
- [HTTP Endpoints](#http-endpoints)

## Core API

### AppContext

The main application context that holds configuration, logger, and component manager.

```go
type AppContext struct {
    Config           *Config
    Logger           *Logger
    BuildInfo        *BuildInfo
    ComponentManager *ComponentManager
}
```

### Configuration

#### Config

Main configuration structure for SREDIAG.

```go
type Config struct {
    Service     ServiceConfig     `yaml:"service"`
    Build       BuildConfig       `yaml:"build"`
    Diagnose    DiagnoseConfig    `yaml:"diagnose"`
    Plugin      PluginConfig      `yaml:"plugin"`
    Log         LogConfig         `yaml:"log"`
}
```

#### ServiceConfig

```go
type ServiceConfig struct {
    Enabled     bool   `yaml:"enabled"`
    Host        string `yaml:"host"`
    Port        int    `yaml:"port"`
    TLSEnabled  bool   `yaml:"tls_enabled"`
    CertFile    string `yaml:"cert_file"`
    KeyFile     string `yaml:"key_file"`
}
```

#### BuildConfig

```go
type BuildConfig struct {
    OutputDir   string            `yaml:"output_dir"`
    Targets     []BuildTarget     `yaml:"targets"`
    Components  ComponentsConfig  `yaml:"components"`
}

type BuildTarget struct {
    OS   string `yaml:"os"`
    Arch string `yaml:"arch"`
}
```

### Logger

Structured logging interface based on Zap.

```go
type Logger struct {
    Level            string   `yaml:"level"`
    Format           string   `yaml:"format"`
    OutputPaths      []string `yaml:"output_paths"`
    ErrorOutputPaths []string `yaml:"error_output_paths"`
}

// Methods
func NewLogger(cfg *Logger) (*Logger, error)
func (l *Logger) Info(msg string, fields ...Field)
func (l *Logger) Error(msg string, fields ...Field)
func (l *Logger) Warn(msg string, fields ...Field)
func (l *Logger) Debug(msg string, fields ...Field)
```

## Plugin API

### PluginManager

Main interface for managing plugins.

```go
type PluginManager interface {
    Load(ctx context.Context, pluginType ComponentType, name string) error
    Get(name string) (IPluginInstance, bool)
    List() []PluginMetadata
    CheckHealth(ctx context.Context) map[string]*PluginHealth
    GetFactory(typ component.Type) (component.Factory, error)
    Unload(ctx context.Context, name string) error
}
```

#### Methods

##### Load

Load and initialize a plugin.

```go
func (m *PluginManager) Load(ctx context.Context, pluginType ComponentType, name string) error
```

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `pluginType`: Type of plugin (receiver, processor, exporter, extension)
- `name`: Name/identifier of the plugin

**Returns:**
- `error`: Error if loading fails

**Example:**
```go
ctx := context.Background()
err := pluginManager.Load(ctx, core.ReceiverType, "custom-receiver")
if err != nil {
    log.Fatalf("Failed to load plugin: %v", err)
}
```

##### Get

Retrieve a loaded plugin instance.

```go
func (m *PluginManager) Get(name string) (IPluginInstance, bool)
```

**Parameters:**
- `name`: Name of the plugin

**Returns:**
- `IPluginInstance`: Plugin instance interface
- `bool`: True if plugin exists

**Example:**
```go
plugin, exists := pluginManager.Get("custom-receiver")
if !exists {
    log.Printf("Plugin not found")
    return
}
```

##### List

Get metadata for all loaded plugins.

```go
func (m *PluginManager) List() []PluginMetadata
```

**Returns:**
- `[]PluginMetadata`: Slice of plugin metadata

**Example:**
```go
plugins := pluginManager.List()
for _, plugin := range plugins {
    fmt.Printf("Plugin: %s, Type: %s, Version: %s\n", 
        plugin.Name, plugin.Type, plugin.Version)
}
```

### PluginMetadata

Plugin information structure.

```go
type PluginMetadata struct {
    Name         string        `json:"name"`
    Version      string        `json:"version"`
    Type         ComponentType `json:"type"`
    Description  string        `json:"description"`
    Capabilities []string      `json:"capabilities"`
    SHA256       string        `json:"sha256"`
    Signature    string        `json:"signature,omitempty"`
}
```

### IPluginInstance

Interface for plugin instances.

```go
type IPluginInstance interface {
    Initialize(ctx context.Context, metadata PluginMetadata) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    HealthCheck(ctx context.Context) (*PluginHealth, error)
    Factory() (component.Factory, error)
}
```

## Diagnostics API

### DiagnoseManager

Main diagnostics orchestrator.

```go
type DiagnoseManager struct {
    logger *core.Logger
}

func NewDiagnoseManager(logger *core.Logger) *DiagnoseManager
func (m *DiagnoseManager) RunSystem() error
func (m *DiagnoseManager) RunPerformance() error  
func (m *DiagnoseManager) RunSecurity() error
```

#### Methods

##### RunSystem

Execute system diagnostics.

```go
func (m *DiagnoseManager) RunSystem() error
```

**Returns:**
- `error`: Error if diagnostics fail

**Example:**
```go
manager := diagnose.NewDiagnoseManager(logger)
if err := manager.RunSystem(); err != nil {
    log.Fatalf("System diagnostics failed: %v", err)
}
```

##### RunPerformance

Execute performance diagnostics.

```go
func (m *DiagnoseManager) RunPerformance() error
```

**Example:**
```go
if err := manager.RunPerformance(); err != nil {
    log.Fatalf("Performance diagnostics failed: %v", err)
}
```

##### RunSecurity

Execute security diagnostics.

```go
func (m *DiagnoseManager) RunSecurity() error
```

**Example:**
```go
if err := manager.RunSecurity(); err != nil {
    log.Fatalf("Security diagnostics failed: %v", err)
}
```

### Diagnostics Types

#### SystemDiagnostics

```go
type SystemDiagnostics struct {
    logger *core.Logger
}

func NewSystemDiagnostics(logger *core.Logger) *SystemDiagnostics
func (d *SystemDiagnostics) Run(ctx context.Context) error
```

#### PerformanceDiagnostics

```go
type PerformanceDiagnostics struct {
    logger *core.Logger
}

func NewPerformanceDiagnostics(logger *core.Logger) *PerformanceDiagnostics
func (d *PerformanceDiagnostics) Run(ctx context.Context) error
```

#### SecurityDiagnostics

```go
type SecurityDiagnostics struct {
    logger *core.Logger
}

func NewSecurityDiagnostics(logger *core.Logger) *SecurityDiagnostics
func (d *SecurityDiagnostics) Run(ctx context.Context) error
```

## Service API

### Service

Main service orchestrator for the OpenTelemetry Collector.

```go
type Service struct {
    logger     *core.Logger
    receivers  map[component.Type]component.Factory
    processors map[component.Type]component.Factory
    exporters  map[component.Type]component.Factory
    extensions map[component.Type]component.Factory
}
```

#### Constructor

```go
func NewService(
    logger *core.Logger,
    receivers map[component.Type]component.Factory,
    processors map[component.Type]component.Factory,
    exporters map[component.Type]component.Factory,
    extensions map[component.Type]component.Factory,
) *Service
```

#### Methods

##### Start

Start the service and all components.

```go
func (s *Service) Start(ctx context.Context) error
```

**Parameters:**
- `ctx`: Context for cancellation and timeouts

**Returns:**
- `error`: Error if startup fails

**Example:**
```go
ctx := context.Background()
if err := service.Start(ctx); err != nil {
    log.Fatalf("Failed to start service: %v", err)
}
```

##### Stop

Stop the service and all components.

```go
func (s *Service) Stop(ctx context.Context) error
```

**Parameters:**
- `ctx`: Context for cancellation and timeouts

**Returns:**
- `error`: Error if shutdown fails

**Example:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := service.Stop(ctx); err != nil {
    log.Printf("Failed to stop service gracefully: %v", err)
}
```

## Configuration API

### Configuration Loading

#### LoadConfigWithOverlay

Load configuration with environment variable and CLI flag overlays.

```go
func LoadConfigWithOverlay(
    spec interface{}, 
    cliFlags map[string]string, 
    opts ...ConfigOption
) error
```

**Parameters:**
- `spec`: Configuration struct to populate
- `cliFlags`: Map of CLI flags to overlay
- `opts`: Additional configuration options

**Returns:**
- `error`: Error if loading fails

**Example:**
```go
var config core.Config
flags := map[string]string{
    "log-level": "debug",
    "service.port": "8080",
}

if err := core.LoadConfigWithOverlay(&config, flags); err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

#### ValidateConfig

Validate configuration structure.

```go
func ValidateConfig(cfg *Config) error
```

**Parameters:**
- `cfg`: Configuration to validate

**Returns:**
- `error`: Validation errors

## HTTP Endpoints

SREDIAG exposes several HTTP endpoints for monitoring and management.

### Health Check

**Endpoint:** `GET /`  
**Port:** 13133 (default)

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

### Metrics

**Endpoint:** `GET /metrics`  
**Port:** 8889 (default)  
**Format:** Prometheus format

**Example metrics:**
```
# HELP otelcol_receiver_accepted_spans_total Number of spans successfully pushed into the pipeline.
# TYPE otelcol_receiver_accepted_spans_total counter
otelcol_receiver_accepted_spans_total{receiver="otlp",transport="grpc"} 1000

# HELP otelcol_processor_batch_batch_send_size_total Number of spans in a batch
# TYPE otelcol_processor_batch_batch_send_size_total histogram
otelcol_processor_batch_batch_send_size_total_bucket{processor="batch",le="10"} 5
```

### Plugin Management

#### List Plugins

**Endpoint:** `GET /api/v1/plugins`

**Response:**
```json
{
  "plugins": [
    {
      "name": "custom-receiver",
      "version": "1.0.0",
      "type": "receiver",
      "status": "running",
      "health": "healthy"
    }
  ]
}
```

#### Load Plugin

**Endpoint:** `POST /api/v1/plugins`

**Request:**
```json
{
  "name": "custom-receiver",
  "type": "receiver",
  "config": {
    "endpoint": "0.0.0.0:9999"
  }
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Plugin loaded successfully"
}
```

#### Plugin Health

**Endpoint:** `GET /api/v1/plugins/{name}/health`

**Response:**
```json
{
  "name": "custom-receiver",
  "status": "healthy",
  "last_check": "2025-01-01T12:00:00Z",
  "details": {
    "uptime": "24h30m",
    "memory_usage": "45MB"
  }
}
```

### Diagnostics

#### Run Diagnostics

**Endpoint:** `POST /api/v1/diagnose/{type}`

**Parameters:**
- `type`: One of `system`, `performance`, `security`

**Request:**
```json
{
  "output_format": "json",
  "include_details": true
}
```

**Response:**
```json
{
  "status": "completed",
  "duration": "5.2s",
  "results": {
    "checks_passed": 15,
    "checks_failed": 2,
    "warnings": 3
  },
  "details": [
    {
      "check": "cpu_usage",
      "status": "pass",
      "value": "15%"
    }
  ]
}
```

## Error Handling

### Error Types

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### Common Error Codes

- `400`: Bad Request - Invalid input parameters
- `404`: Not Found - Resource not found
- `409`: Conflict - Resource already exists
- `500`: Internal Server Error - Server error
- `503`: Service Unavailable - Service temporarily unavailable

### Error Response Format

```json
{
  "error": {
    "code": 400,
    "message": "Invalid plugin type",
    "details": "Plugin type must be one of: receiver, processor, exporter, extension"
  }
}
```

This API reference provides comprehensive documentation for integrating with and extending SREDIAG.