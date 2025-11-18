package add

import (
	"iac/utils/secret"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	name  string
	value string
	file  string
)

var Cmd = &cobra.Command{
	Use:     "add",
	Short:   "add a new secret",
	Aliases: []string{"create", "new", "set"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := secret.SetSecret(name, value)
		if err != nil {
			slog.Error("Failed to add secret", "error", err, "name", name)
			return err
		}
		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&name, "name", "", "name of the secret")
	Cmd.Flags().StringVar(&value, "value", "", "value of the secret")
	Cmd.Flags().StringVar(&file, "file", "", "file containing the secret value")
}
