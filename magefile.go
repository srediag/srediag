//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target when running `mage`
var Default = All

func All() {
	mg.SerialDeps(Deps, Format, Lint, Test, Build)
}

func Deps() error {
	fmt.Println(">> Downloading modules")
	return sh.RunV("go", "mod", "download")
}

func Format() error {
	fmt.Println(">> Formatting code")
	if err := sh.RunV("go", "fmt", "./cmd/...", "./internal/..."); err != nil {
		return err
	}
	return sh.RunV("goimports", "-w", "./cmd", "./internal")
}

func Lint() error {
	fmt.Println(">> Running linter")
	return sh.RunV("golangci-lint", "run", "--config", ".golangci.yml")
}

func Test() error {
	fmt.Println(">> Running tests")
	return sh.RunV("go", "test", "./...", "-coverprofile=coverage.out", "-timeout", "30s")
}

func Build() error {
	fmt.Println(">> Building binary")
	ver, _ := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	return sh.RunV("go", "build", "-ldflags=-X main.Version="+ver, "-o", "bin/srediag", "./cmd/srediag")
}
