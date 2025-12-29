-- +goose Up
CREATE TABLE IF NOT EXISTS contracts (
    address BYTEA PRIMARY KEY,
    creator_tx_hash BYTEA,
    created_in_block BIGINT,
    name TEXT,
    symbol TEXT,
    decimals INT,
    verified BOOLEAN DEFAULT false,
    abi JSONB,
    bytecode BYTEA,
    runtime_bytecode BYTEA,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS contracts;
