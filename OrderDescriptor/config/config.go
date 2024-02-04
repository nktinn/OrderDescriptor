package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	HTTP     HTTPConfig     `yaml:"http"`
	Nats     Nats           `yaml:"nats"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslMode"`
}

type HTTPConfig struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
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
