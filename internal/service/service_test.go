package service

import (
	"context"
	"testing"

	"github.com/srediag/srediag/internal/core"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
)

func TestServiceCreation(t *testing.T) {
	service := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, service)
}

func TestServiceStartStop(t *testing.T) {
	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	// Test with proper logger
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	assert.NoError(t, err)

	service := NewService(logger, receivers, processors, exporters, extensions)

	ctx := context.Background()
	err = service.Start(ctx)
	assert.NoError(t, err)

	err = service.Stop(ctx)
	assert.NoError(t, err)
}
