package get

import "github.com/spf13/cobra"

var (
	name string
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "get a secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation for adding a secret goes here
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
}
