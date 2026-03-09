# План: Сервис управления перемещением аппаратов Goldex

## Контекст

Goldex — сеть аппаратов по скупке/продаже изделий из драгоценных металлов. Аппараты перемещаются между локациями (ТЦ, банки и т.д.). Сейчас данные хранятся в MySQL в плоской таблице `bots` с текстовыми полями локации. Нужно построить полноценный сервис с нормализованной БД, REST/gRPC API и фронтендом. Результат — `docker compose up` и всё работает.

**Исходные данные:** MySQL-дамп с 25 аппаратами, 23 уникальными локациями (2 аппарата на "ТРЦ Мозаика").

**Выбранный стек:**
- Backend: Go, Clean Architecture, gRPC + gRPC-Gateway
- Frontend: Vue 3 + TypeScript + Vite + Leaflet/OSM
- БД: PostgreSQL 16
- Тесты: unit + интеграционные (базовые)

---

## Часть 1: Инфраструктура и схема БД

**Коммит после:** `docker compose up` поднимает PostgreSQL с данными.

### Файлы:
- `docker-compose.yml` (PostgreSQL + migrate контейнер)
- `backend/migrations/001_init_schema.up.sql` / `.down.sql`
- `backend/migrations/002_seed_data.up.sql` / `.down.sql`

### Схема PostgreSQL:

```sql
CREATE TABLE locations (
    id          BIGSERIAL PRIMARY KEY,
    address     VARCHAR(240) NOT NULL,
    place_name  VARCHAR(240) NOT NULL DEFAULT '',
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_locations_address_place ON locations(address, place_name);

CREATE TABLE machines (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    serial_number   VARCHAR(20) NOT NULL UNIQUE,
    enabled         BOOLEAN NOT NULL DEFAULT TRUE,
    location_id     BIGINT REFERENCES locations(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_machines_location_id ON machines(location_id);

CREATE TABLE movement_history (
    id               BIGSERIAL PRIMARY KEY,
    machine_id       BIGINT NOT NULL REFERENCES machines(id),
    from_location_id BIGINT REFERENCES locations(id),  -- NULL для первой установки
    to_location_id   BIGINT NOT NULL REFERENCES locations(id),
    moved_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_movement_history_machine_id ON movement_history(machine_id);
CREATE INDEX idx_movement_history_moved_at ON movement_history(moved_at);
CREATE INDEX idx_movement_history_machine_moved ON movement_history(machine_id, moved_at);
```

### Миграция данных (002_seed_data.up.sql):
1. Вставить уникальные локации (дедупликация по address+place_name, парсинг координат через `split_part`)
2. Вставить аппараты с оригинальными ID, ссылаясь на location_id
3. Создать начальные записи movement_history (from=NULL, to=текущая, moved_at=created_at)

---

## Часть 2: Backend — каркас Go-проекта, proto-файлы, кодогенерация

**Коммит после:** Go-проект компилируется, proto сгенерированы.

### Файлы:
- `backend/go.mod`
- `backend/cmd/server/main.go` (заглушка)
- `backend/internal/domain/` — machine.go, location.go, movement.go
- `backend/proto/goldex/v1/` — machine.proto, location.proto, movement.proto, statistics.proto, service.proto
- `backend/buf.yaml`, `backend/buf.gen.yaml`
- `backend/pkg/gen/goldex/v1/` — сгенерированный код (коммитится)

### Proto API дизайн (service.proto):
- **Machines CRUD:** List/Get/Create/Update/Delete → `/api/v1/machines`
- **Locations CRUD:** List/Get/Create/Update/Delete → `/api/v1/locations`
- **Movements:** MoveMachine → `POST /api/v1/machines/{id}/move`, GetHistory → `GET /api/v1/machines/{id}/movements`
- **Statistics:** GetLocationDurations, GetMovementsCount, GetMachineTimeline → `/api/v1/statistics/...`

### Кодогенерация: buf.build
- protoc-gen-go, protoc-gen-go-grpc, grpc-gateway, openapiv2

---

## Часть 3: Backend — Repository и UseCase слои

**Коммит после:** бизнес-логика реализована, unit-тесты проходят.

### Файлы:
- `backend/internal/repository/interfaces.go`
- `backend/internal/repository/postgres/` — machine_repo.go, location_repo.go, movement_repo.go, statistics_repo.go
- `backend/internal/usecase/` — machine_usecase.go, location_usecase.go, movement_usecase.go, statistics_usecase.go

### Ключевая логика MoveMachine:
1. Проверить существование аппарата и целевой локации
2. Проверить что целевая локация != текущая
3. В транзакции: записать в movement_history + обновить machine.location_id

### Статистика SQL:
- Дни на локации: `LEAD()` window function по movement_history
- Перемещения за период: `COUNT(*) WHERE moved_at BETWEEN ... GROUP BY machine_id`
- Timeline: отсортированная история с JOIN на locations

### Библиотека БД: pgx/v5

---

## Часть 4: Backend — Delivery (gRPC), main.go, Dockerfile

**Коммит после:** `docker compose up` — работающий backend API.

### Файлы:
- `backend/internal/delivery/grpc/` — machine_handler.go, location_handler.go, movement_handler.go, statistics_handler.go, converter.go
- `backend/cmd/server/main.go` (полная реализация)
- `backend/Dockerfile` (multi-stage: golang:1.23-alpine → alpine:3.19)
- Обновить `docker-compose.yml` (добавить backend сервис)

### main.go:
1. Подключение к PostgreSQL (pgxpool)
2. Применение миграций (golang-migrate, embed)
3. Создание репозиториев → usecase-ов → хендлеров
4. Запуск gRPC сервера (:50051)
5. Запуск gRPC-Gateway HTTP proxy (:8080) + Swagger
6. Graceful shutdown

---

## Часть 5: Frontend — Vue 3 приложение

**Коммит после:** полностью рабочий фронтенд.

### Файлы:
- `frontend/` — package.json, vite.config.ts, tsconfig.json, index.html
- `frontend/src/` — main.ts, App.vue, router/index.ts, types/index.ts
- `frontend/src/api/client.ts` — fetch-обёртка над REST API
- `frontend/src/views/` — MachinesView, LocationsView, MapView, StatisticsView
- `frontend/src/components/` — MachineList, LocationList, MachineMap, MoveMachineDialog, MovementHistory, StatisticsPage
- `frontend/Dockerfile` (node:20-alpine → nginx:alpine)
- `frontend/nginx.conf` — проксирование /api/ на backend

### Страницы:
- `/machines` — таблица аппаратов с кнопкой "Переместить"
- `/locations` — таблица локаций
- `/map` — Leaflet-карта с маркерами аппаратов (popup с информацией)
- `/statistics` — выбор аппарата, таблица дней на локации, перемещения за период, timeline

---

## Часть 6: Тесты, документация, финализация

**Коммит после:** готовый продукт.

### Файлы:
- `backend/internal/usecase/*_test.go` — unit-тесты с mock-репозиториями
- `backend/tests/integration_test.go` — интеграционные тесты (testcontainers-go)
- `README.md` — инструкция по запуску
- `CLAUDE.md` — инструкции для AI-ассистента
- Финальная версия `docker-compose.yml`

---

## Верификация

1. `docker compose up --build` — все сервисы стартуют без ошибок
2. `curl http://localhost:8080/api/v1/machines` — возвращает 25 аппаратов с локациями
3. `curl http://localhost:8080/api/v1/locations` — возвращает 23 локации
4. POST перемещение аппарата → проверить обновление location_id и запись в history
5. Фронтенд на http://localhost:3000 — карта с маркерами, списки, статистика
6. `cd backend && go test ./...` — тесты проходят
7. Swagger UI доступен по http://localhost:8080/swagger/
