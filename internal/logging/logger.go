// Package logging provides a unified logging system for SREDIAG that integrates
// with OpenTelemetry Collector's logging infrastructure.
package logging

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/featuregate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	// Level is the minimum enabled logging level
	Level string
	// Format specifies the output format (json, console)
	Format string
	// OutputPaths is a list of URLs or file paths to write logging output to
	OutputPaths []string
	// ErrorOutputPaths is a list of URLs to write internal logger errors to
	ErrorOutputPaths []string
	// InitialFields are fields to be included in every log entry
	InitialFields map[string]interface{}
	// Development puts the logger in development mode
	Development bool
	// DisableCaller stops annotating logs with the calling function's file name and line number
	DisableCaller bool
	// DisableStacktrace disables automatic stacktrace capturing on error level and above
	DisableStacktrace bool
	// Sampling sets a sampling strategy for the logger
	Sampling *zap.SamplingConfig
}

// Logger wraps zap.Logger to provide additional functionality
type Logger struct {
	*zap.Logger
	config *LoggerConfig
	gates  *featuregate.Registry
}

// defaultConfig provides the default logging configuration
func defaultConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:            "info",
		Format:           "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    make(map[string]interface{}),
		Development:      false,
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(cfg *LoggerConfig) (*Logger, error) {
	if cfg == nil {
		cfg = defaultConfig()
	}

	// Create basic encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure for development if needed
	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	// Create zap config
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(parseLevel(cfg.Level)),
		Development:       cfg.Development,
		DisableCaller:     cfg.DisableCaller,
		DisableStacktrace: cfg.DisableStacktrace,
		Sampling:          cfg.Sampling,
		Encoding:          cfg.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       cfg.OutputPaths,
		ErrorOutputPaths:  cfg.ErrorOutputPaths,
		InitialFields:     cfg.InitialFields,
	}

	// Build the logger
	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return &otelCore{
				Core: core,
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Logger{
		Logger: logger,
		config: cfg,
		gates:  featuregate.NewRegistry(),
	}, nil
}

// parseLevel converts a level string to zapcore.Level
func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// otelCore wraps zapcore.Core to integrate with OpenTelemetry
type otelCore struct {
	zapcore.Core
}

// With adds structured context to the Core.
func (c *otelCore) With(fields []zapcore.Field) zapcore.Core {
	return &otelCore{
		Core: c.Core.With(fields),
	}
}

// Check determines whether the supplied Entry should be logged.
func (c *otelCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write serializes the Entry and any Fields supplied at the log site and writes them to their destination.
func (c *otelCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// Here we could add OpenTelemetry specific handling if needed
	return c.Core.Write(ent, fields)
}

// Sync flushes any buffered log entries
func (c *otelCore) Sync() error {
	return c.Core.Sync()
}

// GetLogLevel returns the current log level
func (l *Logger) GetLogLevel() string {
	return l.config.Level
}

// SetLogLevel changes the logging level
func (l *Logger) SetLogLevel(level string) {
	l.config.Level = level
}

// WithComponent returns a logger with the component field set
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.With(zap.String("component", component)),
		config: l.config,
	}
}

// WithFeatureGates adds OpenTelemetry feature gates to the logger
func (l *Logger) WithFeatureGates(gates *featuregate.Registry) *Logger {
	if gates == nil {
		gates = featuregate.NewRegistry()
	}
	l.gates = gates
	return l
}

// Shutdown flushes any buffered log entries
func (l *Logger) Shutdown() error {
	return l.Sync()
}
