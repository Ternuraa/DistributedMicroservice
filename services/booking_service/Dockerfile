# Шаг 1: Сборка приложения
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходники
COPY . .

# Собираем бинарник под Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o booking-service ./main.go

# Шаг 2: Финальный минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник из шага сборки
COPY --from=builder /app/booking-service .

# Если сервису нужны миграции базы данных при старте, копируем и их
COPY --from=builder /app/migrations ./migrations

# Открываем порты (например, 8080 для HTTP или 50051 для gRPC — подставь свои, если отличаются)
EXPOSE 8080 50051

# Команда для запуска
CMD ["./booking-service"]