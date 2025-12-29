-- +goose Up
CREATE TABLE IF NOT EXISTS stats_daily (
    id BIGSERIAL PRIMARY KEY,
    day DATE NOT NULL,
    tx_count BIGINT DEFAULT 0,
    block_count BIGINT DEFAULT 0,
    avg_gas_price NUMERIC,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(day)
);

CREATE TABLE IF NOT EXISTS stats_hourly (
    id BIGSERIAL PRIMARY KEY,
    hour TIMESTAMPTZ NOT NULL,
    tx_count BIGINT DEFAULT 0,
    block_count BIGINT DEFAULT 0,
    avg_gas_price NUMERIC,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(hour)
);

-- +goose Down
DROP TABLE IF EXISTS stats_hourly;
DROP TABLE IF EXISTS stats_daily;
