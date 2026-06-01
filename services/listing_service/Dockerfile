# Шаг 1: Сборка приложения
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходники
COPY . .

# Собираем бинарник (обрати внимание на имя listing-service)
RUN CGO_ENABLED=0 GOOS=linux go build -o listing-service ./main.go

# Шаг 2: Финальный минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /app/listing-service .

# Копируем миграции, если они есть (если нет, эту строчку можно удалить)
COPY --from=builder /app/migrations ./migrations

# Открываем порты
EXPOSE 8080 50051

# Команда для запуска
CMD ["./listing-service"]