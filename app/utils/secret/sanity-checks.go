package secret

import (
	"iac/utils/errors"
	"iac/utils/exitcodes"
)

func checkIfMasterKeyExists() (int, error) {
	exists, err := DoesMasterKeyExist()
	if err != nil {
		return exitcodes.MasterKeyStatError, err
	}
	if !exists {
		return exitcodes.MasterKeyNotFound, errors.New("Master key not found. Please run 'iac setup' to create one.")
	}
	return 0, nil
}

func SanityChecks() (int, error) {
	if exitCode, err := checkIfMasterKeyExists(); err != nil {
		return exitCode, err
	}
	return 0, nil
}
