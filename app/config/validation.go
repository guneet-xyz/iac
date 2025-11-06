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

	if c.ServicesDirPath == c.BackupsDirPath {
		return errors.New("ServicesDirPath and BackupsDirPath cannot be the same")
	}

	return nil
}
