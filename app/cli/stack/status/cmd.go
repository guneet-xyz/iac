package status

import (
	"fmt"
	"iac/utils/stack"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "status",
	Run: func(cmd *cobra.Command, args []string) {
		stacks, err := stack.GetStackNames()
		if err != nil {
			fmt.Println("Error retrieving stacks:", err)
			return
		}

		fmt.Printf("Found %d stacks:\n", len(stacks))
		for _, stack := range stacks {
			fmt.Printf("- %s\n", stack)
		}
	},
}

// check config status
// check how many stacks are found
// check how many are deployed
