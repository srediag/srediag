package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerCreation(t *testing.T) {
	manager := NewManager(nil, "/tmp/plugins")
	assert.NotNil(t, manager)
}

func TestManagerList(t *testing.T) {
	manager := NewManager(nil, "/tmp/plugins")
	list := manager.List()
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))
}

func TestLoaderCreation(t *testing.T) {
	loader := NewLoader(nil, nil)
	assert.NotNil(t, loader)
}
