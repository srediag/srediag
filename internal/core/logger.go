// Package logging provides a unified logging system for SREDIAG that integrates
// with OpenTelemetry Collector's logging infrastructure.
package core

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/featuregate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide additional functionality
type Logger struct {
	logger            *zap.Logger
	gates             *featuregate.Registry
	Level             string                 // Level is the minimum enabled logging level
	Format            string                 // Format specifies the output format (json, console)
	OutputPaths       []string               // OutputPaths is a list of URLs or file paths to write logging output to
	ErrorOutputPaths  []string               // ErrorOutputPaths is a list of URLs to write internal logger errors to
	InitialFields     map[string]interface{} // InitialFields are fields to be included in every log entry
	Development       bool                   // Development puts the logger in development mode
	DisableCaller     bool                   // DisableCaller stops annotating logs with the calling function's file name and line number
	DisableStacktrace bool                   // DisableStacktrace disables automatic stacktrace capturing on error level and above
	Sampling          *zap.SamplingConfig    // Sampling sets a sampling strategy for the logger
}

// defaultConfig provides the default logging configuration
func defaultConfig() *Logger {
	return &Logger{
		Level:            "info",
		Format:           "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    make(map[string]interface{}),
		Development:      false,
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(cfg *Logger) (*Logger, error) {
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
		logger:            logger,
		gates:             featuregate.NewRegistry(),
		Level:             cfg.Level,
		Format:            cfg.Format,
		OutputPaths:       cfg.OutputPaths,
		ErrorOutputPaths:  cfg.ErrorOutputPaths,
		InitialFields:     cfg.InitialFields,
		Development:       cfg.Development,
		DisableCaller:     cfg.DisableCaller,
		DisableStacktrace: cfg.DisableStacktrace,
		Sampling:          cfg.Sampling,
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
	return l.Level
}

// SetLogLevel changes the logging level
func (l *Logger) SetLogLevel(level string) {
	l.Level = level
}

// WithComponent returns a logger with the component field set
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		logger:            l.logger.With(zap.String("component", component)),
		gates:             l.gates,
		Level:             l.Level,
		Format:            l.Format,
		OutputPaths:       l.OutputPaths,
		ErrorOutputPaths:  l.ErrorOutputPaths,
		InitialFields:     l.InitialFields,
		Development:       l.Development,
		DisableCaller:     l.DisableCaller,
		DisableStacktrace: l.DisableStacktrace,
		Sampling:          l.Sampling,
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
	return l.logger.Sync()
}

// Info logs a message at InfoLevel
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs a message at ErrorLevel
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Debug logs a message at DebugLevel
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Warn logs a message at WarnLevel
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Add a public method to expose the underlying zap.Logger
func (l *Logger) UnderlyingZap() *zap.Logger {
	return l.logger
}

// ZapError returns a zap.Field for an error
func ZapError(err error) zap.Field {
	return zap.Error(err)
}

// ZapString returns a zap.Field for a string key/value
func ZapString(key, val string) zap.Field {
	return zap.String(key, val)
}

// ZapInt returns a zap.Field for an int key/value
func ZapInt(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// ZapReflect returns a zap.Field for a reflect value
func ZapReflect(key string, val interface{}) zap.Field {
	return zap.Reflect(key, val)
}
