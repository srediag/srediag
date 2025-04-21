package plugins

// Plugin interface definitions for srediag
// Each plugin must implement Init and Run

type Plugin interface {
	// Init initializes the plugin with raw configuration
	Init(config map[string]interface{}) error
	// Run executes the plugin's primary function (collect, process, act)
	Run() error
}

// Register and discover plugins here
type registry struct {
	plugins map[string]Plugin
}

var Registry = &registry{plugins: make(map[string]Plugin)}

func NewPluginRegistry() *registry {
	return &registry{plugins: make(map[string]Plugin)}
}

func (r *registry) List() []string {
	keys := make([]string, 0, len(Registry.plugins))
	for k := range Registry.plugins {
		keys = append(keys, k)
	}
	return keys
}

func (r *registry) Get(name string) Plugin {
	return Registry.plugins[name]
}

func (r *registry) Register(name string, p interface{}) Plugin {
	var plugin Plugin
	switch v := p.(type) {
	case func() Plugin:
		plugin = v()
	case Plugin:
		plugin = v
	default:
		panic("Invalid plugin type passed to Register; must be Plugin or func() Plugin")
	}
	r.plugins[name] = plugin
	return plugin
}

func (r *registry) Exists(name string) bool {
	_, exists := r.plugins[name]
	return exists
}
