Listing Service (Сервер): Хранит данные о жилье. Слушает порт :50051.

Booking Service (Клиент): Запрашивает данные о цене и доступности у Listing Service для оформления бронирования.

Proto: Общий контракт взаимодействия, описанный в файле listing.proto.

1. Генерация кода из Proto (Опционально)
Если вы вносили изменения в .proto файл, выполните команду в папке listingService:

Bash
protoc --go_out=. --go-grpc_out=. proto/listing.proto
2. Запуск Listing Service (Сервер)
Откройте первый терминал и перейдите в папку сервиса:

Bash
cd listingService
go mod tidy
go run .
Вы должны увидеть сообщение: 🚀 Listing Service успешно запущен на порту :50051

3. Запуск Booking Service (Клиент)
Откройте второй терминал и перейдите в его папку:

Bash
cd bookingService
go mod tidy
go run .
Сервис отправит запрос соседнему микросервису и выведет полученные данные в консоль.

