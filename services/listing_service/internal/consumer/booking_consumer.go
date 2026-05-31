package consumer

import (
	"context"
	"database/sql"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BookingEvent struct {
	BookingID string `json:"booking_id"`
	ListingID string `json:"listing_id"`
}

func HandleBookingCreated(ctx context.Context, db *sql.DB, ch *amqp.Channel, msg amqp.Delivery) {
	var event BookingEvent
	json.Unmarshal(msg.Body, &event)
	eventID := msg.MessageId

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		msg.Nack(false, true) // Возвращаем в очередь при сбое БД
		return
	}
	defer tx.Rollback()

	// 1. Паттерн Inbox: проверка на дубликаты
	_, err = tx.ExecContext(ctx, "INSERT INTO inbox (event_id) VALUES ($1)", eventID)
	if err != nil {
		// Событие уже обработано
		msg.Ack(false)
		return
	}

	// 2. Бизнес-логика: проверка доступности
	var isAvailable bool
	err = tx.QueryRowContext(ctx, "SELECT is_available FROM listings WHERE id = $1 FOR UPDATE", event.ListingID).Scan(&isAvailable)

	var replyEventType string

	if err == nil && isAvailable {
		// Успех: Бронируем
		tx.ExecContext(ctx, "UPDATE listings SET is_available = FALSE WHERE id = $1", event.ListingID)
		replyEventType = "BookingAccepted"
	} else {
		// Провал: Жилье занято или не найдено (Компенсирующий сценарий)
		replyEventType = "BookingRejected"
	}

	tx.Commit()

	// 3. Отправляем ответ в RabbitMQ (Saga)
	ch.PublishWithContext(ctx, "saga_events", replyEventType, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg.Body,
	})

	msg.Ack(false) // Подтверждаем успешную обработку исходного сообщения
}
