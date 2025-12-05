package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

// ensureSQL creates required partitioned tables if migrations have not been applied.
const ensureSQL = `DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'blocks' AND n.nspname = current_schema()) THEN
        CREATE TABLE blocks (
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
        CREATE TABLE blocks_p0 PARTITION OF blocks FOR VALUES FROM (0) TO (1000000);
        CREATE TABLE blocks_p1 PARTITION OF blocks FOR VALUES FROM (1000000) TO (2000000);
        CREATE INDEX idx_blocks_number ON blocks USING btree (number);
        CREATE INDEX idx_blocks_hash ON blocks USING btree (hash);
    END IF;

    -- ensure extended block columns exist even if table predated this binary
    ALTER TABLE blocks
        ADD COLUMN IF NOT EXISTS gas_used BIGINT,
        ADD COLUMN IF NOT EXISTS gas_limit BIGINT,
        ADD COLUMN IF NOT EXISTS miner TEXT,
        ADD COLUMN IF NOT EXISTS nonce TEXT,
        ADD COLUMN IF NOT EXISTS difficulty TEXT,
        ADD COLUMN IF NOT EXISTS extra_data TEXT,
        ADD COLUMN IF NOT EXISTS logs_bloom TEXT,
        ADD COLUMN IF NOT EXISTS mix_hash TEXT,
        ADD COLUMN IF NOT EXISTS receipts_root TEXT,
        ADD COLUMN IF NOT EXISTS sha3_uncles TEXT,
        ADD COLUMN IF NOT EXISTS size_bytes BIGINT,
        ADD COLUMN IF NOT EXISTS state_root TEXT,
        ADD COLUMN IF NOT EXISTS tx_root TEXT,
        ADD COLUMN IF NOT EXISTS tx_count BIGINT DEFAULT 0,
        ADD COLUMN IF NOT EXISTS uncles TEXT[] DEFAULT '{}'::text[],
        ADD COLUMN IF NOT EXISTS tx_hashes TEXT[] DEFAULT '{}'::text[];
    ALTER TABLE blocks ALTER COLUMN tx_count SET DEFAULT 0;
    ALTER TABLE blocks ALTER COLUMN uncles SET DEFAULT '{}'::text[];
    ALTER TABLE blocks ALTER COLUMN tx_hashes SET DEFAULT '{}'::text[];

    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'transactions' AND n.nspname = current_schema()) THEN
        CREATE TABLE transactions (
            hash TEXT NOT NULL,
            block_number BIGINT NOT NULL,
            "from" TEXT NOT NULL,
            "to" TEXT,
            value NUMERIC(78,0),
            status TEXT,
            created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
            CONSTRAINT pk_transactions PRIMARY KEY (block_number, hash)
        ) PARTITION BY RANGE (block_number);
        CREATE TABLE transactions_p0 PARTITION OF transactions FOR VALUES FROM (0) TO (1000000);
        CREATE TABLE transactions_p1 PARTITION OF transactions FOR VALUES FROM (1000000) TO (2000000);
        CREATE INDEX idx_txs_block_number ON transactions USING btree (block_number);
        CREATE INDEX idx_txs_from ON transactions USING btree ("from");
        CREATE INDEX idx_txs_to ON transactions USING btree ("to");
        CREATE INDEX idx_txs_hash ON transactions USING btree (hash);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'logs' AND n.nspname = current_schema()) THEN
        CREATE TABLE logs (
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
        CREATE TABLE logs_p0 PARTITION OF logs FOR VALUES FROM (0) TO (1000000);
        CREATE INDEX idx_logs_tx_hash ON logs(tx_hash);
        CREATE INDEX idx_logs_address ON logs(address);
        CREATE INDEX idx_logs_topic0 ON logs(topic0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'addresses' AND n.nspname = current_schema()) THEN
        CREATE TABLE addresses (
            address TEXT PRIMARY KEY,
            first_seen_block BIGINT,
            last_seen_block BIGINT,
            tx_count BIGINT DEFAULT 0
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'dag_blocks' AND n.nspname = current_schema()) THEN
        CREATE TABLE dag_blocks (
            number BIGINT NOT NULL,
            hash TEXT NOT NULL,
            parent_hash TEXT,
            timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
            CONSTRAINT dag_blocks_number_positive CHECK (number >= 0),
            CONSTRAINT pk_dag_blocks PRIMARY KEY (number, hash)
        ) PARTITION BY RANGE (number);
        CREATE TABLE dag_blocks_p0 PARTITION OF dag_blocks FOR VALUES FROM (0) TO (1000000);
        CREATE TABLE dag_blocks_p1 PARTITION OF dag_blocks FOR VALUES FROM (1000000) TO (2000000);
        CREATE INDEX idx_dag_blocks_number ON dag_blocks USING btree (number);
        CREATE INDEX idx_dag_blocks_hash ON dag_blocks USING btree (hash);
    END IF;
END $$;`

// EnsureSchema bootstraps the required tables if migrations have not run.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    _, err := pool.Exec(ctx, ensureSQL)
    return err
}
