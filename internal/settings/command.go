package settings

import (
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/components"
	"github.com/srediag/srediag/internal/plugin"
)

// CommandSettings holds settings for command execution
type CommandSettings struct {
	ComponentManager *components.Manager
	PluginManager    *plugin.Manager
	Logger           *zap.Logger
}

// GetLogger returns the logger from command settings
func (s *CommandSettings) GetLogger() *zap.Logger {
	if s.Logger == nil {
		return zap.NewNop()
	}
	return s.Logger
}
