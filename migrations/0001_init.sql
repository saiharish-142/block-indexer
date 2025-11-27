-- Blocks partitioned by range of block number, each 1M blocks per partition.
CREATE TABLE IF NOT EXISTS blocks (
    number BIGINT NOT NULL,
    hash TEXT PRIMARY KEY,
    parent_hash TEXT NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT blocks_number_positive CHECK (number >= 0)
) PARTITION BY RANGE (number);

CREATE TABLE IF NOT EXISTS blocks_p0 PARTITION OF blocks
    FOR VALUES FROM (0) TO (1000000);

CREATE TABLE IF NOT EXISTS blocks_p1 PARTITION OF blocks
    FOR VALUES FROM (1000000) TO (2000000);
-- TODO: automate future partitions via cronjob/migration.

CREATE INDEX IF NOT EXISTS idx_blocks_number ON blocks USING btree (number);

CREATE TABLE IF NOT EXISTS transactions (
    hash TEXT PRIMARY KEY,
    block_number BIGINT NOT NULL,
    "from" TEXT NOT NULL,
    "to" TEXT,
    value NUMERIC(78,0),
    status TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
) PARTITION BY RANGE (block_number);

CREATE TABLE IF NOT EXISTS transactions_p0 PARTITION OF transactions
    FOR VALUES FROM (0) TO (1000000);

CREATE TABLE IF NOT EXISTS transactions_p1 PARTITION OF transactions
    FOR VALUES FROM (1000000) TO (2000000);

CREATE INDEX IF NOT EXISTS idx_txs_block_number ON transactions USING btree (block_number);
CREATE INDEX IF NOT EXISTS idx_txs_from ON transactions USING btree ("from");
CREATE INDEX IF NOT EXISTS idx_txs_to ON transactions USING btree ("to");

CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    address TEXT NOT NULL,
    topic0 TEXT,
    topic1 TEXT,
    topic2 TEXT,
    topic3 TEXT,
    data BYTEA,
    log_index INT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
) PARTITION BY RANGE (block_number);

CREATE TABLE IF NOT EXISTS logs_p0 PARTITION OF logs
    FOR VALUES FROM (0) TO (1000000);

CREATE INDEX IF NOT EXISTS idx_logs_tx_hash ON logs(tx_hash);
CREATE INDEX IF NOT EXISTS idx_logs_address ON logs(address);
CREATE INDEX IF NOT EXISTS idx_logs_topic0 ON logs(topic0);

-- Simple addresses materialized table (can be derived in ETL later).
CREATE TABLE IF NOT EXISTS addresses (
    address TEXT PRIMARY KEY,
    first_seen_block BIGINT,
    last_seen_block BIGINT,
    tx_count BIGINT DEFAULT 0
);
