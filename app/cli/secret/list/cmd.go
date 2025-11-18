package list

import (
	"fmt"
	"iac/utils/secret"
	"log/slog"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		secrets, err := secret.GetSecrets()
		if err != nil {
			slog.Error("Failed to list secrets", "error", err)
			return err
		}
		if len(secrets) == 0 {
			fmt.Println("No secrets found. Add one using 'iac secret add' command.")
			return nil
		}
		if len(secrets) == 1 {
			fmt.Println("Found 1 secret:")
		} else {
			fmt.Printf("Found %d secrets:\n", len(secrets))
		}

		for _, s := range secrets {
			fmt.Println("- " + s)
		}
		return nil
	},
}
