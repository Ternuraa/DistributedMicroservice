package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	proto "github.com/Ternuraa/DistributedMicroservice/listing_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// Структура для чтения конфига
type Config struct {
	Service struct {
		Port string `yaml:"port"`
	} `yaml:"service"`
	ListingService struct {
		Address string `yaml:"address"`
	} `yaml:"listing_service"`
}

// Функция загрузки конфига
func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	// 1. Загружаем конфигурацию
	// Путь "../../configs/booking_config.yaml", так как мы находимся в services/booking_service/
	cfg, err := LoadConfig("../../configs/booking_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// 2. Настраиваем gRPC подключение, используя адрес из конфига
	conn, err := grpc.Dial(cfg.ListingService.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к Listing Service: %v", err)
	}
	defer conn.Close()
	client := proto.NewListingServiceClient(conn)

	// 3. Обработчик HTTP
	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		listingID := r.URL.Query().Get("id")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.GetListingInfo(ctx, &proto.ListingRequest{Id: listingID})
		if err != nil {
			http.Error(w, "Listing not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"price":  resp.Price,
			"msg":    "Бронирование подтверждено через конфиг!",
		})
	})

	// 4. Запуск сервера на порту из конфига
	fmt.Printf("Booking Service запущен на %s\n", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(cfg.Service.Port, nil))
}
