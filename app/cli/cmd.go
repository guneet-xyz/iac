package cli

import (
	"iac/cli/config"
	"iac/cli/env"
	"iac/cli/setup"
	"iac/cli/stack"
	"iac/cli/status"
	verifymasterkey "iac/cli/verify-master-key"
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
	Cmd.AddCommand(env.Cmd)
	Cmd.AddCommand(setup.Cmd)
	Cmd.AddCommand(stack.Cmd)
	Cmd.AddCommand(status.Cmd)
	Cmd.AddCommand(verifymasterkey.Cmd)
	Cmd.AddCommand(version.Cmd)
	Cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	Cmd.PersistentFlags().BoolVarP(&showTime, "show-time", "t", false, "Show timestamps in logs")
}
