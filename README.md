### BookStore

Необходимо реализовать 2 микросервиса, которые будут общаться между собой по GRPC или HTTP. Один из которых сохраняет или получает данные из Postgress или опционально redis.А другой предоставляет Api для работы с системой и опционально шлет сообщения в Kafka

Все действия внутри сервиса необходимо легировать на 3 уровнях debug, info, error.

Все сервисы и зависимости должны подниматься в докер контейнерах через команду docker-compose

Работа над проектом должна происходить в git репозитория
