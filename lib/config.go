package lib

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains information needed to create a deployment
type Config struct {
	DeployServer struct {
		HostPort     string `yaml:"port"`
		RegistryUser string `yaml:"registry_user"`
		RegistryPass string `yaml:"registry_pass"`
	} `yaml:"deploy_server"`

	Container struct {
		ImageName     string   `yaml:"image_name"`
		ContainerPort string   `yaml:"port"`
		ContainerEnv  []string `yaml:"env"`
		RestartPolicy string   `yaml:"restart_policy"`
	} `yaml:"container"`
}

// LoadConfigFile load configuration from a config.yml file
func LoadConfigFile() (Config, error) {
	var cfg Config

	f, err := os.Open("config.yml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
