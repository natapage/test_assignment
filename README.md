# Goldex — Сервис управления перемещением аппаратов

Сервис для управления перемещением аппаратов Goldex между локациями. Включает REST/gRPC API, фронтенд с картой и статистикой.

## Стек

- **Backend:** Go 1.23, gRPC + gRPC-Gateway, PostgreSQL 16, pgx/v5
- **Frontend:** Vue 3, TypeScript, Vite, Leaflet/OpenStreetMap
- **Инфраструктура:** Docker Compose

## Запуск

```bash
docker compose up --build
```

После запуска:

| Сервис | URL |
|--------|-----|
| Фронтенд | http://localhost:3000 |
| REST API | http://localhost:8080/api/v1/ |
| gRPC | localhost:50051 |

## API эндпоинты

### Аппараты
- `GET /api/v1/machines` — список аппаратов
- `GET /api/v1/machines/{id}` — аппарат по ID
- `POST /api/v1/machines` — создать аппарат
- `PUT /api/v1/machines/{id}` — обновить аппарат
- `DELETE /api/v1/machines/{id}` — удалить аппарат

### Локации
- `GET /api/v1/locations` — список локаций
- `GET /api/v1/locations/{id}` — локация по ID
- `POST /api/v1/locations` — создать локацию
- `PUT /api/v1/locations/{id}` — обновить локацию
- `DELETE /api/v1/locations/{id}` — удалить локацию

### Перемещения
- `POST /api/v1/machines/{machine_id}/move` — переместить аппарат
- `GET /api/v1/machines/{machine_id}/movements` — история перемещений

### Статистика
- `GET /api/v1/statistics/durations/{machine_id}` — дни на локациях
- `GET /api/v1/statistics/movements-count?from=...&to=...` — количество перемещений за период
- `GET /api/v1/statistics/timeline/{machine_id}` — timeline перемещений

## Тесты

```bash
cd backend
go test ./...
```

## Структура проекта

```
├── docker-compose.yml
├── backend/
│   ├── cmd/server/main.go          # Точка входа
│   ├── internal/
│   │   ├── domain/                  # Доменные модели
│   │   ├── repository/              # Интерфейсы репозиториев
│   │   │   └── postgres/            # PostgreSQL реализация
│   │   ├── usecase/                 # Бизнес-логика
│   │   └── delivery/grpc/           # gRPC хендлеры
│   ├── migrations/                  # SQL миграции
│   ├── proto/goldex/v1/             # Proto-файлы
│   └── pkg/gen/goldex/v1/           # Сгенерированный код
└── frontend/
    ├── src/
    │   ├── views/                   # Страницы (Vue)
    │   ├── api/client.ts            # API клиент
    │   └── router/                  # Маршрутизация
    └── nginx.conf                   # Проксирование API
```
