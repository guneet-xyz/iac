package config

import (
	"iac/utils/errors"
	"iac/utils/fs"
)

func validateConfig(c *Config) error {
	var err error

	c.Stacks.DirectoryPath, err = fs.MkdirIfNotExists(c.Stacks.DirectoryPath)
	if err != nil {
		return err
	}

	c.Backups.DirectoryPath, err = fs.MkdirIfNotExists(c.Backups.DirectoryPath)
	if err != nil {
		return err
	}

	c.Environment.DirectoryPath, err = fs.MkdirIfNotExists(c.Environment.DirectoryPath)
	if err != nil {
		return err
	}

	if c.Stacks.DirectoryPath == c.Backups.DirectoryPath || c.Stacks.DirectoryPath == c.Environment.DirectoryPath || c.Backups.DirectoryPath == c.Environment.DirectoryPath {
		return errors.New("Stacks, Backups and Environment directory paths must be different")
	}

	c.MasterKey.KcvPath, err = fs.AbsPath(c.MasterKey.KcvPath)
	if err != nil {
		return err
	}

	c.MasterKey.MasterKeyPath, err = fs.AbsPath(c.MasterKey.MasterKeyPath)
	if err != nil {
		return err
	}

	return nil
}
