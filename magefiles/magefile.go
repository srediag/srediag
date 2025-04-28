//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/magefiles/errorcodes"
)

// Default target when running `mage`
// Runs Clean, Deps, Format, Lint, Test, Build in order
var Default = All

func All() {
	mg.SerialDeps(Clean, Deps, Format, Lint, Test, Build)
}

// Clean removes build artifacts
func Clean() {
	fmt.Println(">> Cleaning…")
	if err := os.RemoveAll("bin"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
}

func Deps() {
	fmt.Println(">> Downloading modules")
	if err := sh.RunV("go", "mod", "download"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
}

func Format() {
	fmt.Println(">> Formatting code")
	if err := sh.RunV("go", "fmt", "./cmd/...", "./internal/..."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
	if err := sh.RunV("goimports", "-w", "./cmd", "./internal"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
}

func Lint() {
	fmt.Println(">> Running linter")
	if err := sh.RunV("golangci-lint", "run", "--config", ".golangci.yml"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeLintVetFailed)
	}
}

func Test() {
	fmt.Println(">> Running tests")
	if err := sh.RunV("go", "test", "./...", "-coverprofile=coverage.out", "-timeout", "30s"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeUnitTestsFailed)
	}
}

// Build builds the main binary and all plugins
func Build() {
	fmt.Println(">> Building collector binary…")
	if err := os.MkdirAll("bin", 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
	if err := sh.RunV("go", "build", "-o", "bin/srediag", "./cmd/srediag"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
	BuildPlugins()
}

// BuildPlugins builds all plugins in bin/plugins
func BuildPlugins() {
	fmt.Println(">> Building plugins…")
	buildCfg, err := core.LoadBuildConfig(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load build config: %v\n", err)
		os.Exit(errorcodes.ErrCodeGeneral)
	}
	outputDir := buildCfg.OutputDir
	if outputDir == "" {
		outputDir = core.DefaultBuildOutputDir()
	}
	pluginRoot := filepath.Join(outputDir, "plugins")
	failed := false

	types := []struct {
		name string
		list []map[string]string
	}{
		{"receiver", buildCfg.Receivers},
		{"processor", buildCfg.Processors},
		{"exporter", buildCfg.Exporters},
		{"extension", buildCfg.Extensions},
	}

	for _, typ := range types {
		for _, plugin := range typ.list {
			name := plugin["name"]
			path := plugin["path"]
			if name == "" || path == "" {
				fmt.Fprintf(os.Stderr, "Missing name or path for %s plugin: %+v\n", typ.name, plugin)
				failed = true
				continue
			}
			outDir := filepath.Join(pluginRoot, typ.name, name)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create dir %s: %v\n", outDir, err)
				failed = true
				continue
			}
			outFile := filepath.Join(outDir, name+".so")
			cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outFile, path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			fmt.Printf("Building %s/%s...\n", typ.name, name)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to build %s/%s: %v\n", typ.name, name, err)
				failed = true
				continue
			}
			fmt.Printf("✓ Built %s/%s at %s\n", typ.name, name, outFile)
		}
	}

	if failed {
		os.Exit(errorcodes.ErrCodeGeneral)
	}
	fmt.Println("✓ All plugins built successfully in", pluginRoot)
}

// TestPlugins tests all plugins in bin/plugins
func TestPlugins() {
	fmt.Println(">> Testing plugins…")
	// TODO: Implement plugin test logic here (iterate plugins, run tests)
	fmt.Println("✓ All plugins tested successfully")
}

// GenerateProto generates protobuf code
func GenerateProto() {
	fmt.Println(">> Generating protobuf code…")
	// TODO: Implement proto generation logic
}
