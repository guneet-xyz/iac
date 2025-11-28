package config

type RepositoryConfig struct {
	DirectoryPath string `toml:"directory_path"`
	OriginURL     string `toml:"origin_url"`
}

type MasterKeyConfig struct {
	KcvPath       string `toml:"kcv_path"`
	MasterKeyPath string `toml:"master_key_path"`
}

type EnvironmentConfig struct {
	DirectoryPath string `toml:"directory_path"`
}

type StacksConfig struct {
	DirectoryPath string `toml:"directory_path"`
}

type BackupsConfig struct {
	DirectoryPath string `toml:"directory_path"`
}

type Config struct {
	Repository  RepositoryConfig  `toml:"repository"`
	MasterKey   MasterKeyConfig   `toml:"master_key"`
	Environment EnvironmentConfig `toml:"environment"`
	Stacks      StacksConfig      `toml:"stacks"`
	Backups     BackupsConfig     `toml:"backups"`
}
