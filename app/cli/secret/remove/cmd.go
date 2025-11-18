package remove

import (
	"iac/utils/secret"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:     "remove",
	Short:   "remove a secret",
	Aliases: []string{"delete", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := secret.DeleteSecret(name)
		if err != nil {
			slog.Error("Failed to remove secret", "error", err, "name", name)
			return err
		}
		slog.Info("Secret removed", "name", name)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
}
