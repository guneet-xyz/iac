package secret

import (
	"iac/utils/errors"
	"log/slog"

	"github.com/zalando/go-keyring"
)

func keyringGet(name string) (string, error) {
	slog.Debug("Retrieving secret from keyring", "name", name)
	secret, err := keyring.Get("iac", name)
	if err != nil {
		return "", errors.New("failed to retrieve secret from keyring", "name", name, "error", err)
	}
	slog.Debug("Secret retrieved from keyring", "name", name)
	return secret, nil
}

func keyringSet(name string, value string) error {
	slog.Debug("Storing secret in keyring", "name", name)
	err := keyring.Set("iac", name, value)
	if err != nil {
		return errors.New("failed to store secret in keyring", "name", name, "error", err)
	}
	slog.Debug("Secret stored in keyring", "name", name)
	return nil
}
