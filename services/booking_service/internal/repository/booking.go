package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type BookingRepository struct {
	DB *sql.DB
}

type BookingEvent struct {
	BookingID string `json:"booking_id"`
	ListingID string `json:"listing_id"`
}

func (r *BookingRepository) CreateBooking(ctx context.Context, bID, uID, lID uuid.UUID) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Создаем бронь
	_, err = tx.ExecContext(ctx,
		"INSERT INTO bookings (id, user_id, listing_id, status) VALUES ($1, $2, $3, 'Pending')",
		bID, uID, lID)
	if err != nil {
		return err
	}

	// 2. Формируем событие для RabbitMQ
	event := BookingEvent{BookingID: bID.String(), ListingID: lID.String()}
	payload, _ := json.Marshal(event)

	// 3. Записываем в Outbox атомарно
	_, err = tx.ExecContext(ctx,
		"INSERT INTO outbox (id, event_type, payload) VALUES ($1, 'BookingCreated', $2)",
		uuid.New(), payload)
	if err != nil {
		return err
	}

	return tx.Commit()
}
