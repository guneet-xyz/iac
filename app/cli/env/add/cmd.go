package add

import (
	"iac/utils/secret"
	"log/slog"

	"github.com/spf13/cobra"
)

type EnvType string

var EnvTypeSecret EnvType = "secret"
var EnvTypeVariable EnvType = "variable"

var (
	envType string
	name    string
	value   string
	file    string
)

var Cmd = &cobra.Command{
	Use:     "add",
	Short:   "add a new environment secret / variable",
	Aliases: []string{"create", "new", "set"},
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(name) == 0 {
			slog.Error("name option is required")
			return nil
		}
		// TODO: name validations
		// TODO: name transformations

		if len(file) > 0 {
			if len(value) > 0 {
				slog.Error("cannot use both value and file options")
				return nil
			}

			// TODO: implement reading from file
			slog.Error("not implemented")
			return nil
		} else if len(value) > 0 {
			// TODO: value validations. (is this possible?)
		} else {
			slog.Error("either value or file option must be provided")
			return nil
		}

		switch EnvType(envType) {
		case EnvTypeSecret:
			slog.Info("adding new secret", "name", name)
			return secret.SetSecret(name, value)
		case EnvTypeVariable:
			slog.Info("adding new variable", "name", name)
			slog.Error("not implemented")
			return nil
		default:
			slog.Error("unknown environment type", "type", envType)
			return nil
		}
	},
}

func init() {
	Cmd.Flags().StringVar(&envType, "type", "", "type of the environment entry (secret or variable)")
	Cmd.Flags().StringVar(&name, "name", "", "name of the secret")
	Cmd.Flags().StringVar(&value, "value", "", "value of the secret")
	Cmd.Flags().StringVar(&file, "file", "", "file containing the secret value")
}
