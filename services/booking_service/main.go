package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	proto "github.com/Ternuraa/DistributedMicroservice/listing_service/proto"
	"github.com/Ternuraa/DistributedMicroservice/services/booking_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Структура для парсинга JSON из тела запроса фронтенда
type BookRequest struct {
	ListingID string `json:"listing_id"`
	UserID    string `json:"user_id"`
}

func main() {
	// 1. Инициализируем конфиг через наш пакет
	cfg, err := config.Init("../../configs/booking_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// 2. gRPC соединение с сервисом объявлений
	conn, err := grpc.Dial(cfg.ListingService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC connection error: %v", err)
	}
	defer conn.Close()
	client := proto.NewListingServiceClient(conn)

	// 3. Роуты
	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		// НАСТРОЙКА CORS: разрешаем запросы с любых источников в браузере
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обработка предварительного запроса (preflight) от браузера
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем, что метод запроса именно POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем JSON из тела запроса
		var reqBook BookRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBook); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Передаем ID, полученный из фронтенда, в gRPC клиент
		resp, err := client.GetListingInfo(ctx, &proto.ListingRequest{Id: reqBook.ListingID})
		if err != nil {
			log.Printf("❌ Ошибка gRPC вызова: %v", err)
			http.Error(w, "Listing not found", http.StatusNotFound)
			return
		}

		// Возвращаем успешный ответ обратно на фронтенд
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   resp,
		})
	})

	log.Printf("🚀 Booking Service запущен на порту %s (CORS включен)", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(cfg.Service.Port, nil))
}
