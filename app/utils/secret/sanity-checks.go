package secret

import (
	"iac/utils/errors"
	"iac/utils/exitcodes"
)

func checkIfMasterKeyExists() (int, error) {
	masterKey, err := GetMasterKey()
	if err != nil {
		return exitcodes.MasterKeyDoesNotExist, err
	}
	if masterKey == "" {
		return exitcodes.MasterKeyDoesNotExist, errors.New("master key does not exist in keyring")
	}
	return 0, nil
}

func SanityChecks() (int, error) {
	if exitCode, err := checkIfMasterKeyExists(); err != nil {
		return exitCode, err
	}
	return 0, nil
}
