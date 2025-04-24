package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

const defaultRefreshInterval = 30 * time.Second

// ServiceInfo represents information about a discovered service
type ServiceInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Endpoints []string          `json:"endpoints"`
	Metadata  map[string]string `json:"metadata"`
}

// DiscoveryProvider defines the interface for service discovery providers
type DiscoveryProvider interface {
	// Start starts the discovery provider
	Start(ctx context.Context) error
	// Stop stops the discovery provider
	Stop(ctx context.Context) error
	// Discover performs service discovery
	Discover(ctx context.Context) ([]ServiceInfo, error)
}

// Manager manages service discovery
type Manager struct {
	mu              sync.RWMutex
	logger          *zap.Logger
	providers       map[string]DiscoveryProvider
	services        map[string]ServiceInfo
	watchChan       chan ServiceInfo
	refreshInterval time.Duration
	watching        bool
	callbacks       []func(ServiceInfo)
}

// NewManager creates a new discovery manager
func NewManager(logger *zap.Logger, refreshInterval time.Duration) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}

	if refreshInterval <= 0 {
		refreshInterval = defaultRefreshInterval
	}

	return &Manager{
		logger:          logger,
		providers:       make(map[string]DiscoveryProvider),
		services:        make(map[string]ServiceInfo),
		watchChan:       make(chan ServiceInfo),
		refreshInterval: refreshInterval,
		callbacks:       make([]func(ServiceInfo), 0),
	}
}

// RegisterProvider registers a discovery provider
func (m *Manager) RegisterProvider(name string, provider DiscoveryProvider) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	m.providers[name] = provider
	m.logger.Info("Registered discovery provider", zap.String("name", name))
	return nil
}

// UnregisterProvider removes a discovery provider
func (m *Manager) UnregisterProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; !exists {
		return fmt.Errorf("provider %s not registered", name)
	}

	delete(m.providers, name)
	m.logger.Info("Unregistered discovery provider", zap.String("name", name))
	return nil
}

// Start starts the discovery manager
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.watching {
		return fmt.Errorf("discovery manager already started")
	}

	for name, provider := range m.providers {
		if err := provider.Start(ctx); err != nil {
			m.logger.Error("Failed to start discovery provider",
				zap.String("name", name),
				zap.Error(err))
			continue
		}
	}

	m.watching = true
	go m.watch(ctx)

	m.logger.Info("Started discovery manager",
		zap.Duration("interval", m.refreshInterval))
	return nil
}

// Stop stops the discovery manager
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.watching {
		return nil
	}

	m.watching = false
	close(m.watchChan)

	for name, provider := range m.providers {
		if err := provider.Stop(ctx); err != nil {
			m.logger.Error("Failed to stop discovery provider",
				zap.String("name", name),
				zap.Error(err))
		}
	}

	m.logger.Info("Stopped discovery manager")
	return nil
}

// GetServices returns all discovered services
func (m *Manager) GetServices() []ServiceInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	services := make([]ServiceInfo, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	return services
}

// Watch returns a channel that receives service updates
func (m *Manager) Watch() <-chan ServiceInfo {
	return m.watchChan
}

// AddCallback registers a callback function for service updates
func (m *Manager) AddCallback(callback func(ServiceInfo)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks = append(m.callbacks, callback)
}

// watch periodically discovers services
func (m *Manager) watch(ctx context.Context) {
	ticker := time.NewTicker(m.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.discover(ctx)
		}
	}
}

// discover performs service discovery using all providers
func (m *Manager) discover(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.watching {
		return
	}

	for name, provider := range m.providers {
		services, err := provider.Discover(ctx)
		if err != nil {
			m.logger.Error("Failed to discover services",
				zap.String("provider", name),
				zap.Error(err))
			continue
		}

		for _, service := range services {
			existing, exists := m.services[service.ID]
			if !exists || m.serviceChanged(existing, service) {
				m.services[service.ID] = service
				m.notifySubscribers(service)
			}
		}
	}
}

// serviceChanged checks if a service has changed
func (m *Manager) serviceChanged(old, new ServiceInfo) bool {
	if old.Name != new.Name || old.Type != new.Type || old.Version != new.Version {
		return true
	}

	if len(old.Endpoints) != len(new.Endpoints) {
		return true
	}

	// Compare endpoints
	for i, endpoint := range old.Endpoints {
		if endpoint != new.Endpoints[i] {
			return true
		}
	}

	// Compare metadata
	if len(old.Metadata) != len(new.Metadata) {
		return true
	}

	for key, oldValue := range old.Metadata {
		if newValue, exists := new.Metadata[key]; !exists || oldValue != newValue {
			return true
		}
	}

	return false
}

// notifySubscribers notifies all subscribers of a service update
func (m *Manager) notifySubscribers(service ServiceInfo) {
	select {
	case m.watchChan <- service:
	default:
		m.logger.Warn("Watch channel full, skipping notification",
			zap.String("service_id", service.ID))
	}

	for _, callback := range m.callbacks {
		go callback(service)
	}
}
