package manager

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFile represents a configuration file
type ConfigFile struct {
	Name     string
	Path     string
	Required bool
}

// DefaultConfigFiles returns the default configuration files
func DefaultConfigFiles() []ConfigFile {
	return []ConfigFile{
		{
			Name:     "srediag.yaml",
			Path:     "/etc/srediag/config/srediag.yaml",
			Required: true,
		},
		{
			Name:     "otel-config.yaml",
			Path:     "/etc/srediag/config/otel-config.yaml",
			Required: true,
		},
	}
}

// FindConfigFile finds a configuration file in the search paths
func FindConfigFile(name string) (string, error) {
	searchPaths := []string{
		".",
		"configs",
		"/etc/srediag/config",
		filepath.Join(os.Getenv("HOME"), ".srediag"),
	}

	for _, path := range searchPaths {
		filePath := filepath.Join(path, name)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}
	}

	return "", fmt.Errorf("configuration file %q not found in search paths", name)
}

// ValidateConfigFile validates a configuration file
func ValidateConfigFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file %q: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", path)
	}

	if info.Mode().Perm()&0444 == 0 {
		return fmt.Errorf("%q is not readable", path)
	}

	return nil
}
