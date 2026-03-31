package main

import (
	"log"
	"net"

	// Исправленный импорт с псевдонимом proto
	proto "github.com/Ternuraa/DistributedMicroservice/listingService/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Создаем TCP-слушателя на порту 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить TCP-слушателя: %v", err)
	}

	// 2. Создаем экземпляр gRPC сервера
	grpcServer := grpc.NewServer()

	// 3. Регистрируем нашу логику (из server.go) в gRPC сервере
	// Убедись, что ListingServer описан в файле server.go в этой же папке!
	proto.RegisterListingServiceServer(grpcServer, &ListingServer{})

	// 4. Включаем рефлексию
	reflection.Register(grpcServer)

	log.Println("Listing Service успешно запущен на порту :50051")
	log.Println("Ожидание входящих gRPC вызовов...")

	// 5. Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
