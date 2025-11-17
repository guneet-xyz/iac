package list

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation for adding a secret goes here
		return nil
	},
}
