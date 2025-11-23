package verifymasterkey

import (
	"fmt"
	"iac/utils/secret"
	"log/slog"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var Cmd = &cobra.Command{
	Use:   "verify-master-key",
	Short: "Verify the master key for encrypted secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter master phrase to generate master key :")
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			slog.Error("Error while trying to read master key from terminal", "error", err)
			return
		}
		passphrase := string(bytes)
		slog.Debug("Setting master key")
		valid, err := secret.VerifyPassphraseWithKcv(passphrase)
		if err != nil {
			slog.Error("Error while trying to verify master key", "error", err)
			return
		}
		if !valid {
			fmt.Println("Master key verification failed")
			return
		}
		fmt.Println("Master key verification succeeded")
	},
}
