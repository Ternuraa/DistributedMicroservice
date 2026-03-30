package main

import (
	"context"
	"log"

	// Импортируем сгенерированный код из папки proto
	"github.com/your-username/harbor/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListingServer — это наш "мозг" сервиса
type ListingServer struct {
	proto.UnimplementedListingServiceServer
}

// GetListingInfo — та самая функция, которую мы описали в .proto
func (s *ListingServer) GetListingInfo(ctx context.Context, req *proto.ListingRequest) (*proto.ListingResponse, error) {
	log.Printf("Получен запрос на поиск жилья с ID: %s", req.Id)

	// Имитируем поиск в базе данных (пока просто хардкод для теста)
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID не может быть пустым")
	}

	if req.Id == "404" {
		return nil, status.Errorf(codes.NotFound, "Жильё с таким ID не найдено")
	}

	// Если всё ок, возвращаем данные
	return &proto.ListingResponse{
		Id:          req.Id,
		Price:       1500.50,
		IsAvailable: true,
	}, nil
}
