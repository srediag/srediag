package main

import (
	"os"

	"github.com/srediag/srediag/cmd/srediag/commands"
	"github.com/srediag/srediag/internal/core"
)

// main is the entry point for the srediag CLI.
func main() {
	// All CLI flag/env/config logic is handled in core.go/root command.
	// Here we only create the AppContext and execute the root command.
	ctx := &core.AppContext{}
	if err := commands.Execute(ctx); err != nil {
		os.Exit(1)
	}
}
