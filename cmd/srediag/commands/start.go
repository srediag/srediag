package commands

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/service"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the SREDIAG service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get component factories
			receivers := cmdSettings.ComponentManager.GetReceivers()
			processors := cmdSettings.ComponentManager.GetProcessors()
			exporters := cmdSettings.ComponentManager.GetExporters()
			extensions := cmdSettings.ComponentManager.GetExtensions()

			// Create service
			svc := service.NewService(
				cmdSettings.Logger,
				receivers,
				processors,
				exporters,
				extensions,
			)

			// Start service
			if err := svc.Start(context.Background()); err != nil {
				cmdSettings.Logger.Error("Failed to start service", zap.Error(err))
				return err
			}

			// Wait for interrupt
			<-cmd.Context().Done()

			// Stop service
			if err := svc.Stop(context.Background()); err != nil {
				cmdSettings.Logger.Error("Failed to stop service", zap.Error(err))
				return err
			}

			return nil
		},
	}

	return cmd
}
