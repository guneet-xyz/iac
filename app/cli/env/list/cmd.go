package list

import (
	"fmt"
	"iac/utils/env"
	"log/slog"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ListEnvEntries()
		if err != nil {
			slog.Error("Failed to list env entries", "error", err)
			return err
		}
		if len(entries) == 0 {
			fmt.Println("No entries found. Add one using 'iac env add' command.")
			return nil
		}
		if len(entries) == 1 {
			fmt.Println("Found 1 env entry:")
		} else {
			fmt.Printf("Found %d env entry:\n", len(entries))
		}

		for _, s := range entries {
			fmt.Println("- " + s)
		}
		return nil
	},
}
