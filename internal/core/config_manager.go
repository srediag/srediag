package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// DefaultConfigManager is the default implementation of ConfigManager
type DefaultConfigManager struct {
	logger   *zap.Logger
	config   interface{}
	watcher  *fsnotify.Watcher
	configCh chan interface{}
	mu       sync.RWMutex
	healthy  bool
	running  bool
	stopChan chan struct{}
}

// NewConfigManager creates a new instance of DefaultConfigManager
func NewConfigManager(logger *zap.Logger) (*DefaultConfigManager, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &DefaultConfigManager{
		logger:   logger,
		watcher:  watcher,
		configCh: make(chan interface{}),
		healthy:  true,
		stopChan: make(chan struct{}),
	}, nil
}

// Start initializes the config manager
func (cm *DefaultConfigManager) Start(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.running {
		return fmt.Errorf("config manager is already running")
	}

	cm.logger.Info("starting config manager")
	cm.running = true

	// Start watching for config changes
	go cm.watchConfig()

	return nil
}

// Stop stops the config manager
func (cm *DefaultConfigManager) Stop(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.running {
		return fmt.Errorf("config manager is not running")
	}

	cm.logger.Info("stopping config manager")
	close(cm.stopChan)
	cm.running = false

	if err := cm.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close file watcher: %w", err)
	}

	return nil
}

// IsHealthy returns the health status
func (cm *DefaultConfigManager) IsHealthy() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.healthy
}

// LoadConfig loads configuration from the given path
func (cm *DefaultConfigManager) LoadConfig(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.running {
		return fmt.Errorf("config manager is not running")
	}

	// Remove existing watch
	if err := cm.watcher.Remove(path); err != nil && !isNotExist(err) {
		cm.logger.Warn("failed to remove existing watch", zap.Error(err))
	}

	// Add new watch
	if err := cm.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	// TODO: Implement actual config loading based on file type (yaml, json, etc)
	// For now, just set a placeholder
	cm.config = map[string]interface{}{
		"loaded_at": time.Now(),
		"path":      path,
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *DefaultConfigManager) GetConfig() interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// ValidateConfig validates the given configuration
func (cm *DefaultConfigManager) ValidateConfig(cfg interface{}) error {
	// TODO: Implement config validation based on schema
	return nil
}

// WatchConfig watches for configuration changes
func (cm *DefaultConfigManager) WatchConfig(ctx context.Context) (<-chan interface{}, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.running {
		return nil, fmt.Errorf("config manager is not running")
	}

	return cm.configCh, nil
}

// watchConfig watches for file system events on the config file
func (cm *DefaultConfigManager) watchConfig() {
	for {
		select {
		case <-cm.stopChan:
			return
		case event, ok := <-cm.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				cm.logger.Info("config file modified", zap.String("path", event.Name))
				// TODO: Implement config reload
				select {
				case cm.configCh <- cm.config:
				default:
					cm.logger.Warn("config channel is blocked")
				}
			}
		case err, ok := <-cm.watcher.Errors:
			if !ok {
				return
			}
			cm.logger.Error("config watcher error", zap.Error(err))
			cm.mu.Lock()
			cm.healthy = false
			cm.mu.Unlock()
		}
	}
}

// isNotExist returns true if the error is a "not exists" error
func isNotExist(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "no such file or directory" || err.Error() == "file does not exist"
}

// SaveConfig implements ConfigManager
func (cm *DefaultConfigManager) SaveConfig(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal the configuration to YAML
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write the configuration to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	return nil
}
