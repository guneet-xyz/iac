package setup

import (
	"fmt"
	"iac/utils/secret"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Cmd = &cobra.Command{
	Use:   "setup",
	Short: "first time setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(`Edit config file at ~/.iac/config.yaml to set up your environment`)
		masterKey, err := secret.GetMasterKey()
		if err == nil && masterKey != "" {
			fmt.Println("Master key is already set up.")
			return nil
		}

		fmt.Println("Enter master phrase to generate master key :")
		var phrase string
		_, err = fmt.Scanln(&phrase)
		if err != nil {
			return err
		}
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return err
		}
		plaintextKey := string(bytes)
		err = secret.SetMasterKey(plaintextKey)
		if err != nil {
			return err
		}
		return nil
	},
}
