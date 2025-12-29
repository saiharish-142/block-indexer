-- +goose Up
-- Blocks Table
CREATE TABLE IF NOT EXISTS dag_blocks (
    hash TEXT PRIMARY KEY,
    block_order BIGINT, -- 'order' is a keyword
    height BIGINT,
    weight TEXT,
    tx_root TEXT,
    state_root TEXT,
    parent_root TEXT,
    confirmations BIGINT,
    txs_valid BOOLEAN,
    difficulty BIGINT,
    bits TEXT,
    pow_name TEXT,
    pow_type INT,
    pow_nonce BIGINT,
    timestamp TIMESTAMPTZ,
    reward BIGINT,
    tx_fee BIGINT,
    is_blue SMALLINT,
    main_chain BOOLEAN,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

-- Edges (Parents) Table
CREATE TABLE IF NOT EXISTS dag_edges (
    block_hash TEXT NOT NULL,
    parent_hash TEXT NOT NULL,
    PRIMARY KEY (block_hash, parent_hash)
);

-- Tips Table
CREATE TABLE IF NOT EXISTS dag_tips (
    hash TEXT PRIMARY KEY,
    height BIGINT,
    status TEXT
);

-- Transactions Table (for DAG-specific view or raw storage)
CREATE TABLE IF NOT EXISTS dag_transactions (
    tx_hash TEXT PRIMARY KEY,
    block_hash TEXT NOT NULL,
    block_order BIGINT,
    tx_index INT,
    size INT,
    confirmations BIGINT,
    timestamp TIMESTAMPTZ,
    txs_valid BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_dag_blocks_order ON dag_blocks(block_order);
CREATE INDEX IF NOT EXISTS idx_dag_blocks_height ON dag_blocks(height);
CREATE INDEX IF NOT EXISTS idx_dag_edges_parent ON dag_edges(parent_hash);
CREATE INDEX IF NOT EXISTS idx_dag_tx_block ON dag_transactions(block_hash);

-- +goose Down
DROP TABLE IF EXISTS dag_transactions;
DROP TABLE IF EXISTS dag_tips;
DROP TABLE IF EXISTS dag_edges;
DROP TABLE IF EXISTS dag_blocks;
