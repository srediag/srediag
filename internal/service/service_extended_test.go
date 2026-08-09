package service

import (
	"context"
	"testing"

	"github.com/srediag/srediag/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestServiceCreationWithFactories(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)
	assert.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

func TestServiceStartStopExtended(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)

	ctx := context.Background()

	// Start should not error with empty factories
	err = service.Start(ctx)
	assert.NoError(t, err)

	// Stop should not error
	err = service.Stop(ctx)
	assert.NoError(t, err)
}

func TestServiceMultipleStarts(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)

	ctx := context.Background()

	// Multiple starts should not error
	err1 := service.Start(ctx)
	err2 := service.Start(ctx)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Clean up
	err = service.Stop(ctx)
	assert.NoError(t, err)
}

func TestServiceStopWithoutStart(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)

	ctx := context.Background()

	// Stop without start should not error
	err = service.Stop(ctx)
	assert.NoError(t, err)
}
