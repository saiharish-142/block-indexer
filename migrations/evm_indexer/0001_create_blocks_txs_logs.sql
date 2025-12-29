-- +goose Up
CREATE TABLE IF NOT EXISTS blocks (
    hash TEXT PRIMARY KEY,
    parent_hash TEXT,
    number BIGINT NOT NULL,
    timestamp TIMESTAMPTZ,
    miner TEXT,
    gas_used NUMERIC,
    gas_limit NUMERIC,
    base_fee_per_gas NUMERIC,
    extra_data TEXT,
    tx_count INT,
    is_canonical BOOLEAN DEFAULT true,
    raw JSONB, -- For any extra fields
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    hash TEXT PRIMARY KEY,
    block_hash TEXT,
    block_number BIGINT,
    tx_index INT,
    from_address TEXT,
    to_address TEXT,
    nonce BIGINT,
    value NUMERIC,
    gas_price NUMERIC,
    gas_used NUMERIC,
    gas_limit NUMERIC,
    status INT, -- 1 success, 0 failure
    method_sig TEXT,
    input TEXT,
    contract_address TEXT,
    type INT,
    is_canonical BOOLEAN DEFAULT true,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    tx_hash TEXT,
    block_hash TEXT,
    block_number BIGINT,
    tx_index INT,
    log_index INT,
    address TEXT,
    topics TEXT[],
    data TEXT,
    is_canonical BOOLEAN DEFAULT true,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS internal_transactions (
    id BIGSERIAL PRIMARY KEY,
    tx_hash TEXT,
    block_hash TEXT,
    block_number BIGINT,
    from_address TEXT,
    to_address TEXT,
    value NUMERIC,
    gas_limit NUMERIC,
    gas_used NUMERIC,
    call_type TEXT, -- CALL, DELEGATECALL, etc.
    trace_address TEXT, -- "0,1,2"
    error TEXT,
    is_canonical BOOLEAN DEFAULT true,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS token_transfers (
    id BIGSERIAL PRIMARY KEY,
    tx_hash TEXT,
    block_hash TEXT,
    block_number BIGINT,
    log_index INT,
    token_address TEXT,
    from_address TEXT,
    to_address TEXT,
    value NUMERIC, -- For ERC20
    token_id NUMERIC, -- For ERC721/1155
    type TEXT, -- ERC20, ERC721, ERC1155
    is_canonical BOOLEAN DEFAULT true,
    inserted_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_blocks_number ON blocks(number);
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks(hash);
CREATE INDEX IF NOT EXISTS idx_txs_block_number ON transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_logs_block_number ON logs(block_number);
CREATE INDEX IF NOT EXISTS idx_logs_topic0 ON logs((topics[1]));
CREATE INDEX IF NOT EXISTS idx_int_txs_block ON internal_transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_transfers_token ON token_transfers(token_address);
CREATE INDEX IF NOT EXISTS idx_transfers_from ON token_transfers(from_address);
CREATE INDEX IF NOT EXISTS idx_transfers_to ON token_transfers(to_address);

-- +goose Down
DROP TABLE IF EXISTS token_transfers;
DROP TABLE IF EXISTS internal_transactions;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS blocks;
