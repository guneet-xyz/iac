package cli

import (
	"iac/cli/config"
	"iac/cli/info"
	"iac/cli/secret"
	"iac/cli/setup"
	"iac/cli/status"
	"iac/cli/up"
	"iac/cli/version"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var Cmd = &cobra.Command{
	Use:   "iac",
	Short: "A custom, and perhaps over-engineered, IaC solution",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		slog.Debug("Setting up logger", "verbose", verbose)
		handlerOptions := tint.Options{
			Level: slog.LevelInfo,
		}
		if verbose {
			handlerOptions.Level = slog.LevelDebug
		}
		handler := tint.NewHandler(os.Stdout, &handlerOptions)
		logger := slog.New(handler)
		slog.SetDefault(logger)
		slog.Debug("Logger setup complete", "verbose", verbose)
	},
}

func init() {
	Cmd.AddCommand(config.Cmd)
	Cmd.AddCommand(info.Cmd)
	Cmd.AddCommand(secret.Cmd)
	Cmd.AddCommand(setup.Cmd)
	Cmd.AddCommand(status.Cmd)
	Cmd.AddCommand(up.Cmd)
	Cmd.AddCommand(version.Cmd)
	Cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}
