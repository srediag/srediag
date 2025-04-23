package main

import (
	"fmt"
	"os"

	"github.com/srediag/srediag/cmd/srediag/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
