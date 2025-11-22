package cli

import (
	"iac/cli/config"
	"iac/cli/info"
	"iac/cli/secret"
	"iac/cli/setup"
	"iac/cli/status"
	"iac/cli/up"
	"iac/cli/version"
	"iac/utils"
	"iac/utils/logger"
	"log/slog"
	"slices"

	"github.com/spf13/cobra"
)

var (
	verbose  bool
	showTime bool
)

var exceptions = []string{
	"setup",
	"version",
}

var Cmd = &cobra.Command{
	Use:   "iac",
	Short: "A custom, and perhaps over-engineered, IaC solution",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		logger.SetupLogger(logger.LoggerOptions{
			Verbose: verbose,
			Time:    showTime,
		})
		slog.Debug("Logger setup complete", "verbose", verbose)
		cmdName := cmd.Name()
		if !slices.Contains(exceptions, cmdName) {
			err := utils.SanityChecks()
			if err != nil {
				slog.Error("Sanity checks failed", "error", err)
				panic(err)
			}
		}
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
	Cmd.PersistentFlags().BoolVarP(&showTime, "show-time", "t", false, "Show timestamps in logs")
}
