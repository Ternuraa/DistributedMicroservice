package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
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
	cfg, err := LoadConfig("../../configs/listing_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.Service.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// Тут должна быть регистрация твоего сервера, например:
	// proto.RegisterListingServiceServer(s, &server{})

	fmt.Printf("Listing Service (gRPC) запущен на %s\n", cfg.Service.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
