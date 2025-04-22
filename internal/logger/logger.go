// internal/logger/logger.go
package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// Init initializes the global logger instance
func Init(debug bool) error {
	var err error
	once.Do(func() {
		var cfg zap.Config
		if debug {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			cfg = zap.NewProductionConfig()
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		logger, err = cfg.Build(zap.AddCallerSkip(1))
	})
	return err
}

// Get returns the global logger instance
func Get() (*zap.Logger, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger not initialized")
	}
	return logger, nil
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}
	return logger.With(fields...)
}

// Debug logs a message at debug level
func Debug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Info logs a message at info level
func Info(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Warn logs a message at warn level
func Warn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Error logs a message at error level
func Error(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

// Fatal logs a message at fatal level
func Fatal(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}
