# Distributed Microservices System

Проект по реализации межсервисного взаимодействия в рамках курса по распределенным системам. Включает в себя синхронное взаимодействие по gRPC, внешнее REST API и маршрутизацию через API Gateway.

---

## Архитектура системы

Проект состоит из трех независимых микросервисов и единого входного шлюза:

1.  **Listing Service** (Port `50051`): gRPC-сервер, хранящий данные об объектах недвижимости.
2.  **User Service** (Port `8081`): REST-сервис для управления профилями пользователей.
3.  **Booking Service** (Port `8082`): "Умный" клиент, принимающий REST-запросы и запрашивающий данные у Listing по gRPC.
4.  **API Gateway (Nginx)** (Port `80`): Единая точка входа, распределяющая трафик.



---

## Инструкция по запуску

### 1. Подготовка инфраструктуры (Docker)
Убедитесь, что **Docker Desktop** запущен и работает. Затем соберите и запустите шлюз:

```bash
# Перейдите в папку шлюза
cd ServiceInteraction/api-gateway

# Сборка образа и запуск контейнера
docker build -t my-gateway .
docker run --name gateway -p 80:80 -d my-gateway

---

## Инструкция терминалы

###2. Запуск Backend-сервисов
Откройте три новых терминала и запустите сервисы в указанном порядке:

 Терминал 1: Listing Service (gRPC Сервер)

```bash
cd ServiceInteraction/listingService
go run .

Терминал 2: User Service (REST Сервер)

```bash
cd ServiceInteraction/userService
go run .
Терминал 3: Booking Service (Клиент gRPC + REST Сервер)

```bash
cd ServiceInteraction/bookingService
go run .