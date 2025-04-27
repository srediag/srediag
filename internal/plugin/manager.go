// Package plugin provides plugin management functionality
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/cloudwego/shmipc-go"
	"go.opentelemetry.io/collector/component"

	"github.com/srediag/srediag/internal/core"
)

// defaultPluginDir is the default directory for plugin socket files
const defaultPluginDir = "/tmp/srediag/plugins"

// manager implements the PluginManager interface
type PluginManager struct {
	logger    *core.Logger
	pluginDir string
	plugins   map[string]*pluginInstance
	mu        sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager(logger *core.Logger, pluginDir string) *PluginManager {
	if pluginDir == "" {
		pluginDir = defaultPluginDir
	}

	return &PluginManager{
		logger:    logger,
		pluginDir: pluginDir,
		plugins:   make(map[string]*pluginInstance),
	}
}

// Load initializes a plugin of the specified type
func (m *PluginManager) Load(ctx context.Context, pluginType core.ComponentType, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin already loaded")
	}

	pluginPath := filepath.Join(m.pluginDir, string(pluginType)+"s", name)
	if _, err := os.Stat(pluginPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin not found")
		}
		return fmt.Errorf("failed to check plugin: %w", err)
	}

	shmPath := fmt.Sprintf("/tmp/srediag-%s-%s.ipc", pluginType, name)
	conf := shmipc.DefaultSessionManagerConfig()
	if runtime.GOOS == "darwin" {
		conf.ShareMemoryPathPrefix = "/tmp/srediag-plugin-ipc"
		conf.QueuePath = "/tmp/srediag-plugin-ipc_queue"
	} else {
		conf.ShareMemoryPathPrefix = "/dev/shm/srediag-plugin-ipc"
	}
	conf.Network = "unix"
	conf.Address = shmPath

	sessionManager, err := shmipc.NewSessionManager(conf)
	if err != nil {
		return fmt.Errorf("failed to create session manager: %w", err)
	}

	// Start the plugin process, passing the ipc address as argumento
	cmd := exec.Command(pluginPath, "--ipc", shmPath)
	if err := cmd.Start(); err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to start plugin: %w", err)
	}

	// Obtain a stream from the session manager for comunicação.
	stream, err := sessionManager.GetStream()
	if err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to get stream: %w", err)
	}
	// Após uso, o stream será devolvido automaticamente pelo sessionManager, portanto não chamamos PutBack aqui.

	// Prepare and send the initialization request
	initReq := IPCRequest{Method: "Initialize", Params: json.RawMessage(`{}`)}
	reqData, err := json.Marshal(initReq)
	if err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to marshal initialization request: %w", err)
	}

	writer := stream.BufferWriter()
	if err := writer.WriteString(string(reqData)); err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to write initialization request: %w", err)
	}

	// Flush the buffer to send the data to the plugin.
	if err := stream.Flush(true); err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to flush stream: %w", err)
	}

	// Read the response from the plugin.
	reader := stream.BufferReader()
	// Note: ajuste o tamanho conforme o protocolo definido com o plugin.
	respData, err := reader.ReadBytes(512)
	if err != nil {
		sessionManager.Close()
		return fmt.Errorf("failed to read initialization response: %w", err)
	}

	var resp IPCResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		sessionManager.Close()
		return fmt.Errorf("bad response: %w", err)
	}
	if resp.Error != "" {
		sessionManager.Close()
		return fmt.Errorf("plugin initialization error: %s", resp.Error)
	}

	// Armazena a instância do plugin junto com o session manager ativo.
	m.plugins[name] = &pluginInstance{
		metadata: PluginMetadata{Name: name, Type: pluginType},
		ch:       sessionManager,
		cmd:      cmd,
	}

	return nil
}

// Get returns a loaded plugin instance
func (m *PluginManager) Get(name string) (IPluginInstance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, false
	}

	return &clientInstance{
		metadata: plugin.metadata,
		ch:       plugin.ch,
	}, true
}

// List returns metadata for all loaded plugins
func (m *PluginManager) List() []PluginMetadata {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]PluginMetadata, 0, len(m.plugins))
	for _, p := range m.plugins {
		list = append(list, p.metadata)
	}
	return list
}

// CheckHealth performs health checks on all plugins
func (m *PluginManager) CheckHealth(ctx context.Context) map[string]*PluginHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]*PluginHealth)
	for name := range m.plugins {
		// Implement health check logic using shmipc-go
		results[name] = &PluginHealth{
			Status:    "unknown",
			LastCheck: time.Now(),
			Error:     "",
		}
	}

	return results
}

// clientInstance implements the Instance interface for a remote plugin
type clientInstance struct {
	metadata PluginMetadata
	ch       *shmipc.SessionManager
}

func (i *clientInstance) Initialize(ctx context.Context, metadata PluginMetadata) error {
	if i.ch == nil {
		return fmt.Errorf("plugin session manager not initialized")
	}
	stream, err := i.ch.GetStream()
	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}
	defer i.ch.PutBack(stream)

	params, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	initReq := IPCRequest{Method: "Initialize", Params: params}
	reqData, err := json.Marshal(initReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	writer := stream.BufferWriter()
	if err := writer.WriteString(string(reqData)); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	if err := stream.Flush(true); err != nil {
		return fmt.Errorf("failed to flush stream: %w", err)
	}
	reader := stream.BufferReader()
	respData, err := reader.ReadBytes(512)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	var resp IPCResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("bad response: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("plugin initialization error: %s", resp.Error)
	}
	return nil
}

func (i *clientInstance) Start(ctx context.Context) error {
	if i.ch == nil {
		return fmt.Errorf("plugin session manager not initialized")
	}
	stream, err := i.ch.GetStream()
	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}
	defer i.ch.PutBack(stream)

	initReq := IPCRequest{Method: "Start", Params: json.RawMessage(`{}`)}
	reqData, err := json.Marshal(initReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	writer := stream.BufferWriter()
	if err := writer.WriteString(string(reqData)); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	if err := stream.Flush(true); err != nil {
		return fmt.Errorf("failed to flush stream: %w", err)
	}
	reader := stream.BufferReader()
	respData, err := reader.ReadBytes(512)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	var resp IPCResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return fmt.Errorf("bad response: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("plugin start error: %s", resp.Error)
	}
	return nil
}

func (i *clientInstance) Stop(ctx context.Context) error {
	// Implement stop logic using shmipc-go
	return nil
}

func (i *clientInstance) HealthCheck(ctx context.Context) (*PluginHealth, error) {
	// Implement health check logic using shmipc-go
	return nil, nil
}

func (i *clientInstance) Factory() (component.Factory, error) {
	// Implement factory logic using shmipc-go
	return nil, nil
}

// GetFactory returns a factory for the given component type
func (m *PluginManager) GetFactory(typ component.Type) (component.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Find plugin with matching type
	for _, plugin := range m.plugins {
		if string(plugin.metadata.Type) == typ.String() {
			instance := &clientInstance{
				metadata: plugin.metadata,
			}
			return instance.Factory()
		}
	}

	return nil, fmt.Errorf("no factory found for type %s", typ)
}

func (m *PluginManager) Unload(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	plugin.ch.Close()
	if plugin.cmd != nil && plugin.cmd.Process != nil {
		err := plugin.cmd.Process.Kill() // Best effort
		if err != nil {
			m.logger.Warn("Failed to kill plugin process", core.ZapString("name", name), core.ZapError(err))
		}
	}

	delete(m.plugins, name)

	return nil
}
