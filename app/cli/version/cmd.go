package version

import (
	"fmt"
	"iac/utils/meta"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "display application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(meta.Version)
	},
}
