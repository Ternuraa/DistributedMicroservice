package worker

import (
	"database/sql"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RunOutboxWorker запускает бесконечный цикл обработки событий
func RunOutboxWorker(db *sql.DB, ch *amqp.Channel) {
	ticker := time.NewTicker(1 * time.Second) // Проверяем новые события каждую секунду
	defer ticker.Stop()

	for range ticker.C {
		// 1. Выбираем события для отправки
		rows, err := db.Query("SELECT id, event_type, payload FROM outbox LIMIT 10")
		if err != nil {
			log.Printf("Ошибка при чтении из outbox: %v", err)
			continue
		}

		for rows.Next() {
			var id, eType, payload string
			if err := rows.Scan(&id, &eType, &payload); err != nil {
				continue
			}

			// 2. Публикуем в RabbitMQ
			err = ch.Publish(
				"",               // exchange
				"booking_events", // routing key
				false,            // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(payload),
				})

			if err != nil {
				log.Printf("Не удалось отправить сообщение %s: %v", id, err)
				continue
			}

			// 3. Удаляем из outbox только после успешной отправки
			_, err = db.Exec("DELETE FROM outbox WHERE id = $1", id)
			if err != nil {
				log.Printf("Ошибка удаления из outbox: %v", err)
			} else {
				log.Printf("Событие %s успешно отправлено в RabbitMQ", id)
			}
		}
		rows.Close()
	}
}
