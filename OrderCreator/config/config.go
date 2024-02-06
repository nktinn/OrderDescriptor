package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nats Nats `yaml:"nats"`
}

type Nats struct {
	ClusterID string `yaml:"clusterID"`
	ClientID  string `yaml:"clientID"`
	URL       string `yaml:"url"`
	Subject   string `yaml:"subject"`
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}
	return cfg, nil
}
