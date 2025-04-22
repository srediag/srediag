package plugins

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	Name        string
	Version     string
	Type        string
	Description string
	Author      string
}

// Registry manages plugin registration and discovery
type registry struct {
	plugins map[string]interface{}
}

var Registry = &registry{plugins: make(map[string]interface{})}

func NewPluginRegistry() *registry {
	return &registry{plugins: make(map[string]interface{})}
}

func (r *registry) List() []string {
	keys := make([]string, 0, len(Registry.plugins))
	for k := range Registry.plugins {
		keys = append(keys, k)
	}
	return keys
}

func (r *registry) Get(name string) interface{} {
	return Registry.plugins[name]
}

func (r *registry) Register(name string, p interface{}) {
	r.plugins[name] = p
}

func (r *registry) Exists(name string) bool {
	_, exists := r.plugins[name]
	return exists
}
