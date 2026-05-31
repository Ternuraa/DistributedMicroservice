package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ternuraa/DistributedMicroservice/services/user_service/internal/config"
)

func main() {
	cfg, err := config.Init("../../configs/user_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User Profile: Elizaveta Savelieva")
	})

	log.Printf("🚀 User Service запущен на %s", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(cfg.Service.Port, nil))
}
