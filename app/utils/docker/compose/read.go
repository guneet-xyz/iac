package compose

import "iac/utils/yaml"

type ConfigService struct {
	ContainerName string `yaml:"container_name"`
}

type Config struct {
	Services map[string]ConfigService `yaml:"services"`
}

func ReadConfig(path string) (Config, error) {
	var config Config

	err := yaml.FromFile(path, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
