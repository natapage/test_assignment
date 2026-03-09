# Goldex Machine Management Service

## Описание проекта

Сервис управления перемещением аппаратов Goldex между локациями (ТЦ, банки и т.д.). Нормализованная БД (PostgreSQL), REST/gRPC API (Go), фронтенд (Vue 3 + Leaflet).

## Команды

### Запуск всего проекта
```bash
docker compose up --build
```

### Запуск тестов (backend)
```bash
cd backend
go test ./...
```

### Пересборка отдельного сервиса
```bash
docker compose up --build -d backend
docker compose up --build -d frontend
```

### Proto-кодогенерация (через Docker)
```bash
cd backend
docker run --rm -v "$(pwd):/workspace" -w /workspace bufbuild/buf:latest generate
```

### Применение миграций вручную
Миграции применяются автоматически при старте backend. Для ручного запуска:
```bash
docker compose run --rm migrate
```

## Архитектура

### Clean Architecture (backend)
- `internal/domain/` — доменные модели (Machine, Location, Movement)
- `internal/repository/` — интерфейсы репозиториев
- `internal/repository/postgres/` — PostgreSQL реализация (pgx/v5)
- `internal/usecase/` — бизнес-логика
- `internal/delivery/grpc/` — gRPC хендлеры + converter
- `cmd/server/main.go` — точка входа, DI, graceful shutdown

### Proto-first API
- Proto-файлы: `backend/proto/goldex/v1/`
- Сгенерированный код: `backend/pkg/gen/goldex/v1/`
- buf.build для кодогенерации (protoc-gen-go, grpc, gateway)

### БД: PostgreSQL 16
- Таблицы: `locations`, `machines`, `movement_history`
- Миграции: `backend/migrations/` (golang-migrate, embedded в бинарник)
- Данные из MySQL-дампа мигрированы в seed-миграции

### Frontend: Vue 3 + TypeScript + Vite
- Leaflet/OpenStreetMap для карты
- fetch-обёртка над REST API (`src/api/client.ts`)
- nginx проксирует `/api/` на backend

## Ключевые решения

- **Транзакции**: MoveMachine выполняется в транзакции (insert movement + update machine.location_id)
- **Контекст транзакций**: через `context.WithValue` — репозитории извлекают `pgx.Tx` из контекста
- **ID аппаратов**: сохранены оригинальные ID из MySQL-дампа
- **Миграции**: embedded в Go-бинарник через `embed.FS`
- **CORS**: middleware в main.go для фронтенда

## Порты

| Сервис | Порт |
|--------|------|
| Frontend | 3000 |
| Backend HTTP (gRPC-Gateway) | 8080 |
| Backend gRPC | 50051 |
| PostgreSQL | 5432 |

## Данные

- 25 аппаратов, 24 локации (2 аппарата на "ТРЦ Мозаика")
- Координаты из MySQL-дампа (широта/долгота)
- Начальная history: from=NULL, to=текущая локация
