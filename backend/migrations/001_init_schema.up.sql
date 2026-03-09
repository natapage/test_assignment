-- 001_init_schema.up.sql
-- Normalized schema for Goldex machine management service

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
    from_location_id BIGINT REFERENCES locations(id),
    to_location_id   BIGINT NOT NULL REFERENCES locations(id),
    moved_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_movement_history_machine_id ON movement_history(machine_id);
CREATE INDEX idx_movement_history_moved_at ON movement_history(moved_at);
CREATE INDEX idx_movement_history_machine_moved ON movement_history(machine_id, moved_at);
