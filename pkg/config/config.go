package config

import (
	"fmt"
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