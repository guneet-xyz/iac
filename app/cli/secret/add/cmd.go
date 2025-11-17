package add

import "github.com/spf13/cobra"

var (
	name  string
	value string
	file  string
)

var Cmd = &cobra.Command{
	Use:     "add",
	Short:   "add a new secret",
	Aliases: []string{"create", "new"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation for adding a secret goes here
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
	Cmd.Flags().StringVarP(&value, "value", "v", "", "value of the secret")
	Cmd.Flags().StringVarP(&file, "file", "f", "", "file containing the secret value")
}
