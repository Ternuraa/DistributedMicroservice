CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    listing_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL -- 'Pending', 'Confirmed', 'Cancelled'
);

-- Таблица для Outbox (Гарантированная доставка в RabbitMQ)
CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);