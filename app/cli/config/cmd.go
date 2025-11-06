package config

import (
	"fmt"
	"iac/cli/config/show"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	Cmd.AddCommand(show.Cmd)
}
