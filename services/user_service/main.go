package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Port string `yaml:"port"`
	} `yaml:"service"`
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	err = yaml.Unmarshal(buf, &c)
	return &c, err
}

func main() {
	cfg, err := LoadConfig("../../configs/user_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User Profile Data")
	})

	fmt.Printf("User Service (REST) запущен на %s\n", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(cfg.Service.Port, nil))
}
