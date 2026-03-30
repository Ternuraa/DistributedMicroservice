package main

import (
	"context"
	"log"

	// Импортируем сгенерированный код.
	// Путь строится как: Название модуля + путь к папке proto
	proto "github.com/Ternuraa/DistributedMicroservice/listingService/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListingServer — структура нашего сервера
type ListingServer struct {
	proto.UnimplementedListingServiceServer
}

// GetListingInfo — реализация метода из .proto файла
func (s *ListingServer) GetListingInfo(ctx context.Context, req *proto.ListingRequest) (*proto.ListingResponse, error) {
	log.Printf("📥 Получен запрос на жильё с ID: %s", req.Id)

	// Валидация: ID не должен быть пустым
	if req.Id == "" {
		log.Printf("⚠️ Ошибка: пустой ID")
		return nil, status.Errorf(codes.InvalidArgument, "ID не может быть пустым")
	}

	// Имитация ошибки "Не найдено" (специально для теста)
	if req.Id == "404" {
		log.Printf("⚠️ Ошибка: жильё %s не найдено", req.Id)
		return nil, status.Errorf(codes.NotFound, "Жильё с таким ID не существует")
	}

	// Успешный ответ (заглушка данных)
	log.Printf("✅ Отправляем данные для ID: %s", req.Id)
	return &proto.ListingResponse{
		Id:          req.Id,
		Price:       1500.50,
		IsAvailable: true,
	}, nil
}
