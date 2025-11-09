package info

import (
	"iac/cli/info/services"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "info",
	Short: "info about stuff",
}

func init() {
	Cmd.AddCommand(services.Cmd)
}
