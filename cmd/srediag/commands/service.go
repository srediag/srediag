// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/service"
)

// newStartCmd creates a new command to start the SREDIAG service
func newStartCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the SREDIAG service",
		Long: `Start the SREDIAG service with the configured components and settings.
The service will run until interrupted.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get component factories
			exporters := ctx.ComponentManager.GetFactories(string(core.TypeExporter))
			extensions := ctx.ComponentManager.GetFactories(string(core.TypeExtension))
			processors := ctx.ComponentManager.GetFactories(string(core.TypeProcessor))
			receivers := ctx.ComponentManager.GetFactories(string(core.TypeReceiver))

			// Create service
			svc := service.NewService(
				ctx.GetLogger().UnderlyingZap(),
				receivers,
				processors,
				exporters,
				extensions,
			)

			// Start service
			if err := svc.Start(cmd.Context()); err != nil {
				ctx.GetLogger().Error(fmt.Sprintf("Failed to start service: %v", err))
				return fmt.Errorf("failed to start service: %w", err)
			}

			// Wait for interrupt
			<-cmd.Context().Done()

			// Stop service
			if err := svc.Stop(cmd.Context()); err != nil {
				ctx.GetLogger().Error(fmt.Sprintf("Failed to stop service: %v", err))
				return fmt.Errorf("failed to stop service: %w", err)
			}

			return nil
		},
	}

	return cmd
}
