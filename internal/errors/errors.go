package errors

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
)

// ComponentError represents a component error
type ComponentError struct {
	Type    component.Type
	Name    string
	Message string
	Err     error
}

// Error implements error.Error
func (e *ComponentError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s/%s: %s: %v", e.Type, e.Name, e.Message, e.Err)
	}
	return fmt.Sprintf("%s/%s: %s", e.Type, e.Name, e.Message)
}

// Unwrap implements errors.Unwrap
func (e *ComponentError) Unwrap() error {
	return e.Err
}

// NewComponentError creates a new component error
func NewComponentError(typ component.Type, name string, message string, err error) error {
	return &ComponentError{
		Type:    typ,
		Name:    name,
		Message: message,
		Err:     err,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
	Err     error
}

// Error implements error.Error
func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("config error: %s: %s: %v", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("config error: %s: %s", e.Field, e.Message)
}

// Unwrap implements errors.Unwrap
func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error
func NewConfigError(field string, message string, err error) error {
	return &ConfigError{
		Field:   field,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Error implements error.Error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %v: %s", e.Field, e.Value, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) error {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
