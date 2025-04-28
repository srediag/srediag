package main

import (
	"os"

	"github.com/srediag/srediag/cmd/srediag/commands"
	"github.com/srediag/srediag/internal/core"
)

// Allow injection of Execute and os.Exit for testing
var executeFunc = commands.Execute
var osExit = os.Exit

// main is the entry point for the srediag CLI.
func main() {
	// All CLI flag/env/config logic is handled in core.go/root command.
	// Here we only create the AppContext and execute the root command.
	ctx := &core.AppContext{}
	if err := executeFunc(ctx); err != nil {
		osExit(1)
	}
}
