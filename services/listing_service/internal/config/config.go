package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Port string `yaml:"port"`
	} `yaml:"service"`
}

func Init(path string) (*Config, error) {
	cfg := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
