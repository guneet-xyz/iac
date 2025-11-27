package stack

import (
	"iac/cli/stack/info"
	"iac/cli/stack/status"
	"iac/cli/stack/up"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "stack",
	Short: "manage stacks",
}

func init() {
	Cmd.AddCommand(info.Cmd)
	Cmd.AddCommand(status.Cmd)
	Cmd.AddCommand(up.Cmd)
}
