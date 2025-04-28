package build

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// Package build provides the build orchestration layer for SREDIAG.
//
// This file provides utilities for updating the builder YAML configuration, including version synchronization with go.mod and plugin directory cleanup.
//
// Usage:
//   - Use UpdateBuilder to update builder YAML config in-place, filling in missing versions for known Otel core components and cleaning up unknowns.
//   - All config loading should use LoadBuildConfig for schema compliance and validation.
//
// Best Practices:
//   - Always check for errors from update methods.
//   - Use printSummary for reporting changes to stdout.
//
// TODO:
//   - Add context.Context to all methods for cancellation and timeouts.
//   - Add more granular error reporting for YAML and file operations.

// UpdateBuilder updates the builder YAML config in-place, filling in missing versions for known Otel core components using go.mod, commenting out and cleaning up unknowns, and writing a Makefile summary target.
//
// Parameters:
//   - yamlPath: Path to the builder YAML configuration file.
//   - goModPath: Path to the go.mod file for authoritative Otel versions.
//   - pluginGenDir: Directory containing generated plugin code.
//
// Returns:
//   - error: If parsing, reading, or writing fails, returns a detailed error.
//
// Side Effects:
//   - Modifies the builder YAML file in-place.
//   - Removes plugin directories for unknown or missing versions.
//   - Prints a summary of changes to stdout.
func UpdateBuilder(yamlPath, goModPath, pluginGenDir string) error {
	// 1. Parse go.mod for authoritative Otel versions
	coreVersions, err := parseOtelCoreVersions(goModPath)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// 2. Read YAML
	orig, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML: %w", err)
	}
	var root yaml.Node
	if err := yaml.Unmarshal(orig, &root); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// 3. Track changes for summary
	var changed, commented, removed []string

	// 4. Walk YAML and update versions
	updateComponentVersions(&root, coreVersions, pluginGenDir, &changed, &commented, &removed)

	// 5. Write YAML back in-place
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&root); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}
	if err := os.WriteFile(yamlPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}

	// 6. Print summary to stdout
	printSummary(changed, commented, removed)

	return nil
}

// parseOtelCoreVersions parses the go.mod file for authoritative Otel core component versions.
//
// Parameters:
//   - goModPath: Path to the go.mod file.
//
// Returns:
//   - map[string]string: Map of Otel core module paths to their versions.
//   - error: If reading or parsing fails, returns a detailed error.
func parseOtelCoreVersions(goModPath string) (map[string]string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	core := map[string]string{}
	re := regexp.MustCompile(`(go\.opentelemetry\.io/collector/[^ ]+) v([0-9.\-a-zA-Z]+)`) // e.g. go.opentelemetry.io/collector/receiver v1.30.0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			core[m[1]] = m[2]
		}
	}
	return core, scanner.Err()
}

// updateComponentVersions walks the YAML node tree and updates component versions, comments out unknowns, and removes plugin directories as needed.
//
// Parameters:
//   - root: Root YAML node.
//   - coreVersions: Map of Otel core module paths to their versions.
//   - pluginGenDir: Directory containing generated plugin code.
//   - changed: Pointer to slice tracking changed components.
//   - commented: Pointer to slice tracking commented components.
//   - removed: Pointer to slice tracking removed plugin directories.
func updateComponentVersions(root *yaml.Node, coreVersions map[string]string, pluginGenDir string, changed, commented, removed *[]string) {
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return
	}
	m := root.Content[0]
	for i := 0; i < len(m.Content); i += 2 {
		key := m.Content[i]
		val := m.Content[i+1]
		if key.Kind == yaml.ScalarNode && isComponentSection(key.Value) && val.Kind == yaml.SequenceNode {
			newContent := make([]*yaml.Node, 0, len(val.Content))
			for _, item := range val.Content {
				if item.Kind == yaml.MappingNode {
					var goModIdx = -1
					var modPath, modVer string
					for j := 0; j < len(item.Content); j += 2 {
						k := item.Content[j]
						v := item.Content[j+1]
						if k.Value == "gomod" {
							goModIdx = j + 1
							parts := strings.Fields(v.Value)
							if len(parts) == 2 {
								modPath, modVer = parts[0], parts[1]
							} else if len(parts) == 1 {
								modPath = parts[0]
							}
						}
					}
					if goModIdx != -1 && modVer == "" && modPath != "" {
						if v, ok := coreVersions[modPath]; ok {
							item.Content[goModIdx].Value = modPath + " " + v
							*changed = append(*changed, modPath)
							newContent = append(newContent, item)
						} else {
							*commented = append(*commented, modPath)
							// Remove plugin/generated dir
							pluginDir := filepath.Join(pluginGenDir, key.Value, filepath.Base(modPath))
							if err := os.RemoveAll(pluginDir); err == nil {
								*removed = append(*removed, pluginDir)
							}
							// Optionally, add a comment to the sequence for traceability
							commentNode := &yaml.Node{
								Kind:  yaml.ScalarNode,
								Value: "# REMOVED: " + modPath + " (missing version, removed by update utility)",
							}
							newContent = append(newContent, commentNode)
						}
					} else {
						newContent = append(newContent, item)
					}
				} else {
					newContent = append(newContent, item)
				}
			}
			val.Content = newContent
		}
	}
}

// isComponentSection returns true if the given string is a recognized component section in the builder YAML.
//
// Parameters:
//   - s: Section name string.
//
// Returns:
//   - bool: True if the section is a recognized component section.
func isComponentSection(s string) bool {
	switch s {
	case "receivers", "exporters", "processors", "extensions", "connectors":
		return true
	}
	return false
}

// printSummary prints a summary of updated, commented, and removed components to stdout.
//
// Parameters:
//   - changed: Slice of updated component names.
//   - commented: Slice of commented component names.
//   - removed: Slice of removed plugin directory paths.
func printSummary(changed, commented, removed []string) {
	if len(changed) > 0 {
		fmt.Printf("Updated versions for: %s\n", strings.Join(changed, ", "))
	}
	if len(commented) > 0 {
		fmt.Printf("Removed (missing version): %s\n", strings.Join(commented, ", "))
	}
	if len(removed) > 0 {
		fmt.Printf("Removed plugin/generated dirs: %s\n", strings.Join(removed, ", "))
	}
	if len(changed) == 0 && len(commented) == 0 && len(removed) == 0 {
		fmt.Println("No changes made by update utility.")
	}
}
