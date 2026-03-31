package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Обработчик для получения профиля
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		// Имитация базы данных: возвращаем данные
		user := User{
			ID:    id,
			Name:  "Savelieva Elizaveta",
			Email: "savel@itmo.ru",
		}

		log.Printf("📥 Запрос профиля для ID: %s", id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	log.Println("👤 User Service (REST) успешно запущен на порту :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
