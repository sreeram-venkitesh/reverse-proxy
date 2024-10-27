package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Definitions for the config file
// which defines routers and services
type Config struct {
	Port     int       `yaml:"port"`
	Routers  []Router  `yaml:"routers"`
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Router struct {
	Host    string `yaml:"host"`
	Service string `yaml:"service"`
}

func (c *Config) GetServiceUrl(name string) (string, error) {
	for _, svc := range c.Services {
		if svc.Name == name {
			return svc.URL, nil
		}
	}
	return "", fmt.Errorf("Service not found\n")
}

func LoadConfig(configPath string) (Config, error) {
	var runtimeConfig Config

	// Read the config file
	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		return runtimeConfig, fmt.Errorf("Error reading config.yaml: %v", err)
	}

	err = yaml.Unmarshal(yamlConfig, &runtimeConfig)
	if err != nil {
		return runtimeConfig, fmt.Errorf("Error parsing config.yaml: %v", err)
	}

	return runtimeConfig, nil
}
