package config

type RepositoryConfig struct {
	DirectoryPath string `mapstructure:"directory_path"`
	OriginURL     string `mapstructure:"origin_url"`
}

type MasterKeyConfig struct {
	KcvPath       string `mapstructure:"kcv_path"`
	MasterKeyPath string `mapstructure:"master_key_path"`
}

type EnvironmentConfig struct {
	DirectoryPath string `mapstructure:"directory_path"`
}

type StacksConfig struct {
	DirectoryPath string `mapstructure:"directory_path"`
}

type BackupsConfig struct {
	DirectoryPath string `mapstructure:"directory_path"`
}

type Config struct {
	Repository  RepositoryConfig  `mapstructure:"repository"`
	MasterKey   MasterKeyConfig   `mapstructure:"master_key"`
	Environment EnvironmentConfig `mapstructure:"environment"`
	Stacks      StacksConfig      `mapstructure:"stacks"`
	Backups     BackupsConfig     `mapstructure:"backups"`
}

type UserConfig struct {
	Repository RepositoryConfig `mapstructure:"repository"`
}

type RepoConfig struct {
	Key         KeyConfig         `mapstructure:"key"`
	Environment EnvironmentConfig `mapstructure:"environment"`
	Stacks      StacksConfig      `mapstructure:"stacks"`
	Backups     BackupsConfig     `mapstructure:"backups"`
}

type KeyConfig struct {
	DirectoryPath string `mapstructure:"directory_path"`
}
