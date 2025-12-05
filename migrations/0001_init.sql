-- Blocks partitioned by range of block number, each 1M blocks per partition.
CREATE TABLE IF NOT EXISTS blocks (
    number BIGINT NOT NULL,
    hash TEXT NOT NULL,
    parent_hash TEXT NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    gas_used BIGINT,
    gas_limit BIGINT,
    miner TEXT,
    nonce TEXT,
    difficulty TEXT,
    extra_data TEXT,
    logs_bloom TEXT,
    mix_hash TEXT,
    receipts_root TEXT,
    sha3_uncles TEXT,
    size_bytes BIGINT,
    state_root TEXT,
    tx_root TEXT,
    tx_count BIGINT DEFAULT 0,
    uncles TEXT[] DEFAULT '{}'::text[],
    tx_hashes TEXT[] DEFAULT '{}'::text[],
    CONSTRAINT blocks_number_positive CHECK (number >= 0),
    CONSTRAINT pk_blocks PRIMARY KEY (number, hash)
) PARTITION BY RANGE (number);

CREATE TABLE IF NOT EXISTS blocks_p0 PARTITION OF blocks
    FOR VALUES FROM (0) TO (1000000);

CREATE TABLE IF NOT EXISTS blocks_p1 PARTITION OF blocks
    FOR VALUES FROM (1000000) TO (2000000);
-- TODO: automate future partitions via cronjob/migration.

CREATE INDEX IF NOT EXISTS idx_blocks_number ON blocks USING btree (number);
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks USING btree (hash);

CREATE TABLE IF NOT EXISTS transactions (
    hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    "from" TEXT NOT NULL,
    "to" TEXT,
    value NUMERIC(78,0),
    status TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    CONSTRAINT pk_transactions PRIMARY KEY (block_number, hash)
) PARTITION BY RANGE (block_number);

CREATE TABLE IF NOT EXISTS transactions_p0 PARTITION OF transactions
    FOR VALUES FROM (0) TO (1000000);

CREATE TABLE IF NOT EXISTS transactions_p1 PARTITION OF transactions
    FOR VALUES FROM (1000000) TO (2000000);

CREATE INDEX IF NOT EXISTS idx_txs_block_number ON transactions USING btree (block_number);
CREATE INDEX IF NOT EXISTS idx_txs_from ON transactions USING btree ("from");
CREATE INDEX IF NOT EXISTS idx_txs_to ON transactions USING btree ("to");
CREATE INDEX IF NOT EXISTS idx_txs_hash ON transactions USING btree (hash);

CREATE TABLE IF NOT EXISTS logs (
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    address TEXT NOT NULL,
    topic0 TEXT,
    topic1 TEXT,
    topic2 TEXT,
    topic3 TEXT,
    data BYTEA,
    log_index INT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    CONSTRAINT pk_logs PRIMARY KEY (block_number, tx_hash, log_index)
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

-- DAG blocks stored separately for DAG chain consensus view.
CREATE TABLE IF NOT EXISTS dag_blocks (
    number BIGINT NOT NULL,
    hash TEXT NOT NULL,
    parent_hash TEXT,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT dag_blocks_number_positive CHECK (number >= 0),
    CONSTRAINT pk_dag_blocks PRIMARY KEY (number, hash)
) PARTITION BY RANGE (number);

CREATE TABLE IF NOT EXISTS dag_blocks_p0 PARTITION OF dag_blocks
    FOR VALUES FROM (0) TO (1000000);

CREATE TABLE IF NOT EXISTS dag_blocks_p1 PARTITION OF dag_blocks
    FOR VALUES FROM (1000000) TO (2000000);
-- TODO: automate further dag block partitions.

CREATE INDEX IF NOT EXISTS idx_dag_blocks_number ON dag_blocks USING btree (number);
CREATE INDEX IF NOT EXISTS idx_dag_blocks_hash ON dag_blocks USING btree (hash);
