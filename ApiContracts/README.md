gRPC API Contracts

Все внутренние взаимодействия в Harbor типизированы с помощью gRPC

Файлы контрактов:
 `user.proto`, `listing.proto`, `booking.proto`, `notification.proto`, `review.proto`

 
Сводка endpoints:
 [API-Contracts-Endpoints.md](API-Contracts-Endpoints.md).



Почему это важно для проекта:



Type Safety: Ошибки в структуре данных отлавливаются на этапе компиляции, а не в рантайме.



Performance: Использование Protobuf (бинарный формат) вместо JSON снижает нагрузку на сеть на 30-40%.



Contract-First Design: Сначала проектируется интерфейс взаимодействия, затем пишется бизнес-логика.



Ключевые сценарии:



BookingService -> ListingService.GetListingDetails (Получение цены).



BookingService -> UserService.GetUserStatus (Проверка прав гостя).



AnyService -> NotificationService.SendSystemNotification (Асинхронные уведомления).

