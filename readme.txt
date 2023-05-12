cmd
    messenger
        main.go
internal
    controllers // API
        user.go
    services // business logic
        user.go
    repository // SQL
        user.go
    entities // business entities
        user.go

entities:
    - channel: чат/диалог
    - user: пользователи
    - message: сообщения

DB (у каждого микросервиса своя БД):
    - DB PostgreSQL "user"
    - DB PostgreSQL "channel"
    - DB PostgreSQL "message"