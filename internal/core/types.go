package core

// ComponentType represents a type of component in the system.
type ComponentType string

const (
	TypeCore      ComponentType = "core"      // TypeCore represents a core component
	TypePlugin    ComponentType = "plugin"    // TypePlugin represents a plugin component
	TypeConnector ComponentType = "connector" // TypeConnector represents a connector component
	TypeExporter  ComponentType = "exporter"  // TypeExporter represents an exporter component
	TypeExtension ComponentType = "extension" // TypeExtension represents an extension component
	TypeProcessor ComponentType = "processor" // TypeProcessor represents a processor component
	TypeReceiver  ComponentType = "receiver"  // TypeReceiver represents a receiver component
)
