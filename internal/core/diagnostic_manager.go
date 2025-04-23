package core

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// DefaultDiagnosticManager provides a default implementation of DiagnosticManager
type DefaultDiagnosticManager struct {
	logger      *zap.Logger
	diagnostics map[string]Diagnostic
	mu          sync.RWMutex
	healthy     bool
}

// NewDiagnosticManager creates a new diagnostic manager instance
func NewDiagnosticManager(logger *zap.Logger) *DefaultDiagnosticManager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &DefaultDiagnosticManager{
		logger:      logger,
		diagnostics: make(map[string]Diagnostic),
		healthy:     true,
	}
}

// Start implements Component
func (m *DefaultDiagnosticManager) Start(ctx context.Context) error {
	m.logger.Info("starting diagnostic manager")
	return nil
}

// Stop implements Component
func (m *DefaultDiagnosticManager) Stop(ctx context.Context) error {
	m.logger.Info("stopping diagnostic manager")
	return nil
}

// IsHealthy implements Component
func (m *DefaultDiagnosticManager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthy
}

// RegisterDiagnostic implements DiagnosticManager
func (m *DefaultDiagnosticManager) RegisterDiagnostic(name string, diagnostic Diagnostic) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.diagnostics[name]; exists {
		return fmt.Errorf("diagnostic %s already registered", name)
	}

	m.diagnostics[name] = diagnostic
	return nil
}

// UnregisterDiagnostic implements DiagnosticManager
func (m *DefaultDiagnosticManager) UnregisterDiagnostic(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.diagnostics[name]; !exists {
		return fmt.Errorf("diagnostic %s not found", name)
	}

	delete(m.diagnostics, name)
	return nil
}

// GetDiagnostic implements DiagnosticManager
func (m *DefaultDiagnosticManager) GetDiagnostic(name string) (Diagnostic, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	diagnostic, exists := m.diagnostics[name]
	if !exists {
		return nil, fmt.Errorf("diagnostic %s not found", name)
	}

	return diagnostic, nil
}

// ListDiagnostics implements DiagnosticManager
func (m *DefaultDiagnosticManager) ListDiagnostics() []Diagnostic {
	m.mu.RLock()
	defer m.mu.RUnlock()

	diagnostics := make([]Diagnostic, 0, len(m.diagnostics))
	for _, diagnostic := range m.diagnostics {
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}
