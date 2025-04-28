// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/service"
)

// NewServiceCmd creates the root 'service' command with all subcommands.
// The subcommands include:
// - start: Start the SREDIAG service
// - stop: Stop the SREDIAG service
// - restart: Restart the SREDIAG service
// - reload: Hot-reload YAML configuration
// - detach: Fork to background (daemonize)
// - status: Show service health and resource usage
// - health: Exit 0 if /healthz is ready
// - profile: Gather CPU+heap profile bundle
// - tail-logs: Stream live service logs
// - validate: Dry-run parse YAML + plugin refs
// - install-unit: Create & enable systemd unit
// - uninstall-unit: Remove systemd unit
// - gc: Purge stale PID/socket/logs
// Only CLI wiring is present here; all business logic is delegated to internal/service CLI_* functions.
func NewServiceCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Operate the SREDIAG agent (start, stop, reload, status, etc.)",
		Long:  "Manage the long-lived SREDIAG daemon and embedded OpenTelemetry Collector.",
	}

	cmd.AddCommand(
		newServiceStartCmd(ctx),
		newServiceStopCmd(ctx),
		newServiceRestartCmd(ctx),
		newServiceReloadCmd(ctx),
		newServiceDetachCmd(ctx),
		newServiceStatusCmd(ctx),
		newServiceHealthCmd(ctx),
		newServiceProfileCmd(ctx),
		newServiceTailLogsCmd(ctx),
		newServiceValidateCmd(ctx),
		newServiceInstallUnitCmd(ctx),
		newServiceUninstallUnitCmd(ctx),
		newServiceGcCmd(ctx),
	)
	return cmd
}

// newServiceStartCmd wires the 'start' subcommand to service.CLI_Start.
func newServiceStartCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the SREDIAG service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Start(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceStopCmd wires the 'stop' subcommand to service.CLI_Stop.
func newServiceStopCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the SREDIAG service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Stop(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceRestartCmd wires the 'restart' subcommand to service.CLI_Restart.
func newServiceRestartCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the SREDIAG service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Restart(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceReloadCmd wires the 'reload' subcommand to service.CLI_Reload.
func newServiceReloadCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Hot-reload YAML configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Reload(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceDetachCmd wires the 'detach' subcommand to service.CLI_Detach.
func newServiceDetachCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Fork to background (daemonize)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Detach(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceStatusCmd wires the 'status' subcommand to service.CLI_Status.
func newServiceStatusCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show service health and resource usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Status(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceHealthCmd wires the 'health' subcommand to service.CLI_Health.
func newServiceHealthCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Exit 0 if /healthz is ready",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Health(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceProfileCmd wires the 'profile' subcommand to service.CLI_Profile.
func newServiceProfileCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Gather CPU+heap profile bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Profile(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceTailLogsCmd wires the 'tail-logs' subcommand to service.CLI_TailLogs.
func newServiceTailLogsCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tail-logs",
		Short: "Stream live service logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_TailLogs(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceValidateCmd wires the 'validate' subcommand to service.CLI_Validate.
func newServiceValidateCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Dry-run parse YAML + plugin refs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Validate(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceInstallUnitCmd wires the 'install-unit' subcommand to service.CLI_InstallUnit.
func newServiceInstallUnitCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-unit",
		Short: "Create & enable systemd unit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_InstallUnit(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceUninstallUnitCmd wires the 'uninstall-unit' subcommand to service.CLI_UninstallUnit.
func newServiceUninstallUnitCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall-unit",
		Short: "Remove systemd unit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_UninstallUnit(ctx, cmd, args)
		},
	}
	return cmd
}

// newServiceGcCmd wires the 'gc' subcommand to service.CLI_Gc.
func newServiceGcCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Purge stale PID/socket/logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.CLI_Gc(ctx, cmd, args)
		},
	}
	return cmd
}
