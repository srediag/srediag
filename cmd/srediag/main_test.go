package main

import (
	"errors"
	"testing"

	"github.com/srediag/srediag/internal/core"
)

// mockExit is used to capture os.Exit calls.
var exitCode int

func mockExit(code int) {
	exitCode = code
}

// TestMain_ExecuteError tests that main calls os.Exit(1) if commands.Execute returns an error.
func TestMain_ExecuteError(t *testing.T) {
	origExit := osExit
	defer func() { osExit = origExit }()
	exitCode = 0

	// Use a wrapper for executeFunc to allow mocking
	origExecute := executeFunc
	executeFunc = func(ctx *core.AppContext) error {
		return errors.New("mock error")
	}
	defer func() { executeFunc = origExecute }()

	osExit = mockExit

	main()

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

// TestMain_ExecuteSuccess tests that main does not call os.Exit if commands.Execute returns nil.
func TestMain_ExecuteSuccess(t *testing.T) {
	origExit := osExit
	defer func() { osExit = origExit }()
	exitCode = 0

	// Mock executeFunc to return nil
	origExecute := executeFunc
	defer func() { executeFunc = origExecute }()
	executeFunc = func(ctx *core.AppContext) error {
		return nil
	}

	osExit = mockExit

	main()

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func init() {
	// No redeclaration; osExit is now declared in main.go
}
