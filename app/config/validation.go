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

	c.EnvDirPath, err = fs.MkdirIfNotExists(c.EnvDirPath)
	if err != nil {
		return err
	}

	if c.ServicesDirPath == c.BackupsDirPath || c.ServicesDirPath == c.EnvDirPath || c.BackupsDirPath == c.EnvDirPath {
		return errors.New("ServicesDirPath, BackupsDirPath and EnvDirPath must be different paths")
	}

	c.KcvPath, err = fs.AbsPath(c.KcvPath)
	if err != nil {
		return err
	}

	c.MasterKeyPath, err = fs.AbsPath(c.MasterKeyPath)
	if err != nil {
		return err
	}

	return nil
}
