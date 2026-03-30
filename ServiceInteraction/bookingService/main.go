package main

import (
	"context"
	"log"
	"time"

	// Импортируем proto из ПАПКИ СОСЕДНЕГО СЕРВИСА (Listing)
	proto "github.com/Ternuraa/DistributedMicroservice/listingService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Устанавливаем соединение с Listing Service (порт 50051)
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось соединиться: %v", err)
	}
	defer conn.Close()

	// 2. Создаем клиента
	client := proto.NewListingServiceClient(conn)

	// 3. Делаем тестовый запрос (например, хотим забронировать жилье ID=123)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Отправляем gRPC запрос в Listing Service...")
	resp, err := client.GetListingInfo(ctx, &proto.ListingRequest{Id: "123"})

	if err != nil {
		log.Fatalf("Ошибка при вызове: %v", err)
	}

	// 4. Выводим результат
	log.Printf("Успех! Получены данные от соседа:")
	log.Printf("Жилье ID: %s", resp.Id)
	log.Printf("Цена: %.2f", resp.Price)
	log.Printf("Доступно: %v", resp.IsAvailable)
}
