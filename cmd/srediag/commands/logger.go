package commands

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// getLogger returns a configured logger instance
func getLogger() *zap.Logger {
	// If we have settings with a logger, use that
	if settings.Logger != nil {
		return settings.Logger
	}

	// Otherwise create a new development logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(ExitGeneralError)
	}
	return logger.Named("srediag")
}
