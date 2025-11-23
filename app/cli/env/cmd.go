package env

import (
	"iac/cli/env/add"
	"iac/cli/env/get"
	"iac/cli/env/list"
	"iac/cli/env/remove"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "environment",
	Short:   "manage environment",
	Aliases: []string{"env"},
}

func init() {
	Cmd.AddCommand(add.Cmd)
	Cmd.AddCommand(get.Cmd)
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(remove.Cmd)
}
