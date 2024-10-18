# Task Scheduler (TODO-list)

## Описание проекта

Этот проект представляет собой веб-сервер, реализующий функциональность простейшего планировщика задач (аналог TODO-листа). Планировщик позволяет добавлять задачи с датой дедлайна, заголовком и комментарием, редактировать и удалять их. Поддерживается функциональность повторяющихся задач, которые переносятся на следующую дату по заданному правилу при выполнении.

## Основные функции 
- Добавление задачи.
- Получение списка всех задач.
- Удаление задачи.
- Получение информации о конкретной задаче.
- Изменение параметров задачи.
- Отметка задачи как выполненной.

## Выполненные задания со звёздочкой
- Возможность определять порт при запуске сервера через переменную окружения `TODO_PORT`.
- Возможность определять путь к файлу базы данных через переменную окружения `TODO_DBFILE`.
- Реализована аутентификация.
- Создание Docker-образа.

## Запуск проекта 

### Пример файла .env
- `TODO_PORT = :7540`   Порт, на котором будет работать сервер
- `TODO_DBFILE = scheduler.db`  Путь к файлу базы данных SQLite
- `TODO_PASSWORD = password`   Пароль для аутентификации

### Запуск тестов
`go test ./tests`

### Сборка Docker-образа: 
`docker build -t task-scheduler .`

### Запуск контейнера:
`docker run -p 7540:7540 task-scheduler`
