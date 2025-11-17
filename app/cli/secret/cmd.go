package secret

import (
	"iac/cli/secret/add"
	"iac/cli/secret/get"
	"iac/cli/secret/list"
	"iac/cli/secret/remove"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "secret",
	Short:   "manage secrets",
	Aliases: []string{"secrets"},
}

func init() {
	Cmd.AddCommand(add.Cmd)
	Cmd.AddCommand(get.Cmd)
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(remove.Cmd)
}
