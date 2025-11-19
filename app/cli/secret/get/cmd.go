package get

import (
	"fmt"
	"iac/utils/secret"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "get a secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		value, err := secret.GetSecret(name)
		if err != nil {
			slog.Error("Failed to get secret", "error", err, "name", name)
			return err
		}
		fmt.Println(value)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
}
