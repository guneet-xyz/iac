package setup

import (
	"fmt"
	"iac/utils/env/secret"
	"log/slog"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Cmd = &cobra.Command{
	Use:   "setup",
	Short: "first time setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(`Edit config file at ~/.iac/config.yaml to set up your environment`)

		exists, err := secret.DoesMasterKeyExist()
		if err != nil {
			slog.Error("Error checking if master key exists", "error", err)
			return err
		}
		if exists {
			fmt.Println("Master key already exists. Setup is already complete.")
			return nil
		}

		fmt.Println("Enter master phrase to generate master key :")
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return err
		}
		plaintextKey := string(bytes)
		slog.Debug("Setting master key")
		err = secret.SetMasterKey(plaintextKey)
		if err != nil {
			return err
		}
		return nil
	},
}
