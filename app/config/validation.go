package config

import (
	"iac/utils/errors"
	"iac/utils/fs"
)

func validateConfig(c *Config) error {
	var err error

	c.ServicesDirPath, err = fs.MkdirIfNotExists(c.ServicesDirPath)
	if err != nil {
		return err
	}

	c.BackupsDirPath, err = fs.MkdirIfNotExists(c.BackupsDirPath)
	if err != nil {
		return err
	}

	c.SecretsDirPath, err = fs.MkdirIfNotExists(c.SecretsDirPath)
	if err != nil {
		return err
	}

	if c.ServicesDirPath == c.BackupsDirPath || c.ServicesDirPath == c.SecretsDirPath || c.BackupsDirPath == c.SecretsDirPath {
		return errors.New("ServicesDirPath, BackupsDirPath and SecretsDirPath must be different paths")
	}

	return nil
}
