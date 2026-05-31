package main

import (
	"log"
	"net"

	"github.com/Ternuraa/DistributedMicroservice/services/listing_service/internal/config"
	proto "github.com/Ternuraa/DistributedMicroservice/services/listing_service/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Init("../../configs/listing_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1"+cfg.Service.Port) // Жестко привязываем к 127.0.0.1 для Windows
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// РАСКОММЕНТИРОВАЛИ И ПОДКЛЮЧИЛИ: Регистрируем твой сервер из соседнего файла
	proto.RegisterListingServiceServer(s, &ListingServer{})

	log.Printf("🚀 Listing Service (gRPC) запущен на 127.0.0.1%s", cfg.Service.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
