package plugins

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/lifecycle"
)

// Plugin representa a interface de um plugin SREDIAG
type Plugin interface {
	// Initialize inicializa o plugin com a configuração
	Init(config map[string]interface{}) error

	// Start inicia o plugin
	Start(ctx context.Context) error

	// Stop para o plugin
	Stop(ctx context.Context) error

	// Info retorna os metadados do plugin
	Info() Info
}

// Info contém metadados sobre um plugin
type Info struct {
	Name        string
	Version     string
	Type        string
	Description string
	Author      string
}

// Manager gerencia o ciclo de vida dos plugins
type Manager struct {
	*lifecycle.BaseManager
	config  config.PluginsConfig
	logger  *zap.Logger
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager cria uma nova instância do gerenciador de plugins
func NewManager(cfg config.PluginsConfig, logger *zap.Logger) *Manager {
	return &Manager{
		BaseManager: lifecycle.NewBaseManager(),
		config:      cfg,
		logger:      logger,
		plugins:     make(map[string]Plugin),
	}
}

// Start inicializa e inicia todos os plugins habilitados
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(false); err != nil {
		return err
	}

	// Carrega plugins do diretório configurado
	if err := m.loadPlugins(); err != nil {
		return fmt.Errorf("falha ao inicializar: %v", err)
	}

	// Inicializa e inicia cada plugin habilitado
	for name, p := range m.plugins {
		if !m.isEnabled(name) {
			m.logger.Info("pulando plugin desabilitado", zap.String("name", name))
			continue
		}

		if err := p.Init(m.config.Settings[name]); err != nil {
			return fmt.Errorf("falha ao inicializar plugin %s: %v", name, err)
		}

		if err := p.Start(ctx); err != nil {
			return fmt.Errorf("falha ao iniciar plugin %s: %v", name, err)
		}

		m.logger.Info("plugin iniciado", zap.String("name", name))
	}

	m.SetRunning(true)
	m.logger.Info("gerenciador de plugins iniciado com sucesso")

	return nil
}

// Stop para todos os plugins de forma graciosa
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(true); err != nil {
		return err
	}

	var errs []error

	// Para cada plugin habilitado
	for name, p := range m.plugins {
		if !m.isEnabled(name) {
			continue
		}

		if err := p.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("falha ao parar plugin %s: %v", name, err))
		}
		m.logger.Info("plugin parado", zap.String("name", name))
	}

	m.SetRunning(false)

	if len(errs) > 0 {
		return fmt.Errorf("falha ao parar plugins: %v", errs)
	}

	m.logger.Info("gerenciador de plugins parado com sucesso")
	return nil
}

// LoadPlugin carrega um plugin do sistema de arquivos
func (m *Manager) LoadPlugin(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Abre o plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("falha ao abrir plugin %s: %w", path, err)
	}

	// Procura o símbolo do plugin
	sym, err := plug.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin %s não exporta o símbolo 'New': %w", path, err)
	}

	// Verifica a interface do plugin
	p, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s não implementa a interface Plugin", path)
	}

	// Obtém informações do plugin
	info := p.Info()

	// Registra o plugin
	m.plugins[info.Name] = p
	m.logger.Info("plugin carregado com sucesso",
		zap.String("name", info.Name),
		zap.String("version", info.Version),
		zap.String("type", info.Type),
	)

	return nil
}

// loadPlugins carrega todos os plugins do diretório configurado
func (m *Manager) loadPlugins() error {
	pattern := filepath.Join(m.config.Directory, "*.so")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("falha ao listar plugins: %w", err)
	}

	for _, path := range matches {
		if err := m.LoadPlugin(path); err != nil {
			m.logger.Error("falha ao carregar plugin",
				zap.String("path", path),
				zap.Error(err),
			)
		}
	}

	return nil
}

// GetPlugin retorna um plugin pelo nome
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.plugins[name]
	return p, ok
}

// ListPlugins retorna uma lista de todos os plugins carregados
func (m *Manager) ListPlugins() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Info, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p.Info())
	}
	return plugins
}

// isEnabled verifica se um plugin está habilitado na configuração
func (m *Manager) isEnabled(name string) bool {
	if len(m.config.Enabled) == 0 {
		return true // se nenhum plugin estiver explicitamente habilitado, todos estão
	}
	for _, enabled := range m.config.Enabled {
		if enabled == name {
			return true
		}
	}
	return false
}

// GetPluginConfig retorna a configuração de um plugin específico
func (m *Manager) GetPluginConfig(name string) map[string]interface{} {
	return m.config.Settings[name]
}
