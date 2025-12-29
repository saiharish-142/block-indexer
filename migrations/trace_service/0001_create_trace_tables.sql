-- +goose Up
CREATE TABLE IF NOT EXISTS tx_trace_queue (
    tx_hash BYTEA PRIMARY KEY,
    block_number BIGINT,
    status SMALLINT DEFAULT 0,
    last_error TEXT,
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tx_traces (
    id BIGSERIAL PRIMARY KEY,
    tx_hash BYTEA NOT NULL,
    trace JSONB,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS state_diffs (
    id BIGSERIAL PRIMARY KEY,
    tx_hash BYTEA NOT NULL,
    address BYTEA,
    slot BYTEA,
    before_value BYTEA,
    after_value BYTEA,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS state_diffs;
DROP TABLE IF EXISTS tx_traces;
DROP TABLE IF EXISTS tx_trace_queue;
