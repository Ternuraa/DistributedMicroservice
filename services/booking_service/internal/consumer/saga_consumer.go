package consumer

import (
	"context"
	"database/sql"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SagaEvent struct {
	BookingID string `json:"booking_id"`
}

func HandleSagaResponse(ctx context.Context, db *sql.DB, msg amqp.Delivery) {
	var ev SagaEvent
	json.Unmarshal(msg.Body, &ev)

	var targetStatus string
	if msg.RoutingKey == "BookingAccepted" {
		targetStatus = "Confirmed"
	} else if msg.RoutingKey == "BookingRejected" {
		targetStatus = "Cancelled" // Компенсация
	} else {
		msg.Ack(false)
		return
	}

	_, err := db.ExecContext(ctx, "UPDATE bookings SET status = $1 WHERE id = $2", targetStatus, ev.BookingID)
	if err == nil {
		msg.Ack(false) // Успешно обновили статус
	} else {
		msg.Nack(false, true) // Ошибка БД, пробуем позже
	}
}
