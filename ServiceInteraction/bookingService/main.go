package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	proto "github.com/Ternuraa/DistributedMicroservice/listingService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Настраиваем подключение к Listing Service (gRPC)
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := proto.NewListingServiceClient(conn)

	// Создаем HTTP обработчик для бронирования
	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		listingID := r.URL.Query().Get("id")

		// Синхронный вызов соседа по gRPC
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.GetListingInfo(ctx, &proto.ListingRequest{Id: listingID})
		if err != nil {
			http.Error(w, "Listing not found", http.StatusNotFound) // Обработка 404
			return
		}

		// Если нашли жилье, "бронируем"
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"price":  resp.Price,
			"msg":    "Бронирование подтверждено!",
		})
	})

	log.Println("Booking Service (REST + gRPC Client) запущен на :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
