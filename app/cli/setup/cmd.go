package setup

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "setup",
	Short: "first time setup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Edit config file at ~/.iac/config.yaml to set up your environment`)
	},
}
