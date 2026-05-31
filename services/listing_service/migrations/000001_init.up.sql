CREATE TABLE listings (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE
);

-- Таблица для Inbox (Идемпотентность, защита от дублей из RabbitMQ)
CREATE TABLE inbox (
    event_id UUID PRIMARY KEY,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);