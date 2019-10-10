package lib

import "os"

// Config contains information needed to create a deployment
type Config struct {
	ImageName     string
	RegistryUser  string
	RegistryPass  string
	HostPort      string
	ContainerPort string
}

// NewEnvConfig create a new config based on enviroment variables
func NewEnvConfig() (Config, error) {
	return Config{
		ImageName:     os.Getenv("OHPLOY_IMAGE_NAME"),
		RegistryUser:  os.Getenv("OHPLOY_REGISTRY_USER"),
		RegistryPass:  os.Getenv("OHPLOY_REGISTRY_PASSWORD"),
		HostPort:      os.Getenv("OHPLOY_HOST_PORT"),
		ContainerPort: os.Getenv("OHPLOY_CONTAINER_PORT"),
	}, nil
}
