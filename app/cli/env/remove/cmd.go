package remove

import (
	"iac/utils/env"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:     "remove",
	Short:   "remove an env entry",
	Aliases: []string{"delete", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := env.DeleteEnv(name)
		if err != nil {
			slog.Error("Failed to remove env entry", "error", err, "name", name)
			return err
		}
		slog.Info("Env entry removed", "name", name)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
}
