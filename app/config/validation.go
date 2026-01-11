package config

import (
	"iac/utils/errors"
	"iac/utils/fs"
	"path/filepath"
)

func validateUserConfig(c *UserConfig) error {
	var err error

	c.Repository.DirectoryPath, err = fs.AbsPath(c.Repository.DirectoryPath)
	if err != nil {
		return err
	}

	return nil
}

func validateRepoConfig(c *RepoConfig, repoPath string) error {
	var err error

	repoPath, err = fs.AbsPath(repoPath)
	if err != nil {
		return err
	}

	c.Key.DirectoryPath, err = fs.MkdirIfNotExists(filepath.Join(repoPath, c.Key.DirectoryPath))
	if err != nil {
		return err
	}

	c.Stacks.DirectoryPath, err = fs.MkdirIfNotExists(filepath.Join(repoPath, c.Stacks.DirectoryPath))
	if err != nil {
		return err
	}

	c.Backups.DirectoryPath, err = fs.MkdirIfNotExists(filepath.Join(repoPath, c.Backups.DirectoryPath))
	if err != nil {
		return err
	}

	c.Environment.DirectoryPath, err = fs.MkdirIfNotExists(filepath.Join(repoPath, c.Environment.DirectoryPath))
	if err != nil {
		return err
	}

	if c.Key.DirectoryPath == c.Stacks.DirectoryPath || c.Key.DirectoryPath == c.Backups.DirectoryPath || c.Key.DirectoryPath == c.Environment.DirectoryPath || c.Stacks.DirectoryPath == c.Backups.DirectoryPath || c.Stacks.DirectoryPath == c.Environment.DirectoryPath || c.Backups.DirectoryPath == c.Environment.DirectoryPath {
		return errors.New("Key, Stacks, Backups and Environment directory paths must be different")
	}

	return nil
}
