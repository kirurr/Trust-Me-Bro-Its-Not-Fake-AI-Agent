# Trust Me Bro - AI Agent

## Запуск

```bash
# Запуск всех сервисов (web, backend, rabbitmq, postgres)
docker compose up -d

# Просмотр логов
docker compose logs -f

# Остановка
docker compose down
```

## TUI

```bash
# Запуск терминального интерфейса
docker compose --profile tui run --rm --build tui
```

## Сервисы

- **Web** — http://localhost:5173
- **Backend API** — http://localhost:8080
- **RabbitMQ** — http://localhost:15672 (management UI, guest/guest)
- **PostgreSQL** — localhost:5432

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `POSTGRES_USER` | Пользователь БД | admin |
| `POSTGRES_PASSWORD` | ��ароль БД | secret |
| `POSTGRES_DB` | Имя базы данных | mydb |