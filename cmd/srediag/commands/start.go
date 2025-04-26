// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/service"
)

// newStartCmd creates a new command to start the SREDIAG service
func newStartCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the SREDIAG service",
		Long: `Start the SREDIAG service with the configured components and settings.
The service will run until interrupted.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get component factories
			receivers := opts.Settings.ComponentManager.GetReceivers()
			processors := opts.Settings.ComponentManager.GetProcessors()
			exporters := opts.Settings.ComponentManager.GetExporters()
			extensions := opts.Settings.ComponentManager.GetExtensions()

			// Create service
			svc := service.NewService(
				opts.Settings.Logger,
				receivers,
				processors,
				exporters,
				extensions,
			)

			// Start service
			if err := svc.Start(cmd.Context()); err != nil {
				opts.Settings.Logger.Error("Failed to start service", zap.Error(err))
				return fmt.Errorf("failed to start service: %w", err)
			}

			// Wait for interrupt
			<-cmd.Context().Done()

			// Stop service
			if err := svc.Stop(cmd.Context()); err != nil {
				opts.Settings.Logger.Error("Failed to stop service", zap.Error(err))
				return fmt.Errorf("failed to stop service: %w", err)
			}

			return nil
		},
	}

	return cmd
}
