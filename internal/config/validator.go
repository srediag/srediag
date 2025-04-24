package config

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
)

// BaseValidator provides a base implementation of ConfigValidator
type BaseValidator struct {
	version ConfigVersion
	rules   []ValidationRule
}

// ValidationRule defines a single validation rule
type ValidationRule struct {
	Name        string
	Description string
	Validate    func(*confmap.Conf) error
}

// NewBaseValidator creates a new base validator
func NewBaseValidator(version ConfigVersion) *BaseValidator {
	return &BaseValidator{
		version: version,
		rules:   make([]ValidationRule, 0),
	}
}

// Version implements ConfigValidator
func (v *BaseValidator) Version() ConfigVersion {
	return v.version
}

// AddRule adds a validation rule
func (v *BaseValidator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// Validate implements ConfigValidator
func (v *BaseValidator) Validate(conf *confmap.Conf) error {
	if conf == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Check configuration version
	version, err := v.getConfigVersion(conf)
	if err != nil {
		return fmt.Errorf("invalid configuration version: %w", err)
	}

	if version.Major != v.version.Major {
		return fmt.Errorf("incompatible configuration version: expected %d.x.x, got %d.%d.%d",
			v.version.Major, version.Major, version.Minor, version.Patch)
	}

	// Apply all validation rules
	for _, rule := range v.rules {
		if err := rule.Validate(conf); err != nil {
			return fmt.Errorf("validation rule '%s' failed: %w", rule.Name, err)
		}
	}

	return nil
}

// getConfigVersion extracts version information from the configuration
func (v *BaseValidator) getConfigVersion(conf *confmap.Conf) (ConfigVersion, error) {
	raw := conf.ToStringMap()
	versionVal, ok := raw["version"]
	if !ok {
		return ConfigVersion{}, fmt.Errorf("version field is missing")
	}

	versionMap, ok := versionVal.(map[string]interface{})
	if !ok {
		return ConfigVersion{}, fmt.Errorf("version field is not a map")
	}

	major, ok := versionMap["major"].(float64)
	if !ok {
		return ConfigVersion{}, fmt.Errorf("major version is missing or invalid")
	}

	minor, ok := versionMap["minor"].(float64)
	if !ok {
		return ConfigVersion{}, fmt.Errorf("minor version is missing or invalid")
	}

	patch, ok := versionMap["patch"].(float64)
	if !ok {
		return ConfigVersion{}, fmt.Errorf("patch version is missing or invalid")
	}

	return ConfigVersion{
		Major: int(major),
		Minor: int(minor),
		Patch: int(patch),
	}, nil
}
