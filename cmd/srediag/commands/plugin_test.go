package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPluginCmd(t *testing.T) {
	settings := setupTestSettings(t)

	// Test with nil options
	t.Run("nil options", func(t *testing.T) {
		cmd := NewPluginCmd(nil)
		assert.NotNil(t, cmd)
		assert.Equal(t, "plugin", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
	})

	// Test with custom options
	t.Run("custom options", func(t *testing.T) {
		opts := &Options{
			Settings: settings,
			LogConfig: LogConfig{
				Level:  "debug",
				Format: "json",
			},
		}
		cmd := NewPluginCmd(opts)
		assert.NotNil(t, cmd)
		assert.Equal(t, "plugin", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
	})
}
