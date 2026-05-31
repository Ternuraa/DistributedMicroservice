package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Port string `yaml:"port"`
	} `yaml:"service"`
	ListingService struct {
		Address string `yaml:"address"`
	} `yaml:"listing_service"`
}

// Init загружает конфиг. В продакшене часто путь передают через флаг или ENV
func Init(path string) (*Config, error) {
	config := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, err
	}
	return config, nil
}
