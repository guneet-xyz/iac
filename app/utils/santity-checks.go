package utils

import (
	"iac/utils/docker"
	"iac/utils/env/secret"
	"iac/utils/git"
	"os"
)

func SanityChecks() error {
	if exitCode, err := docker.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	if exitCode, err := git.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	if exitCode, err := secret.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	return nil
}
