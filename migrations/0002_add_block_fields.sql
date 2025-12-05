-- Add extended EVM block fields to support full header storage.
ALTER TABLE IF EXISTS blocks
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
    ADD COLUMN IF NOT EXISTS tx_count BIGINT,
    ADD COLUMN IF NOT EXISTS uncles TEXT[],
    ADD COLUMN IF NOT EXISTS tx_hashes TEXT[];

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'blocks' AND column_name = 'tx_count') THEN
        ALTER TABLE blocks ALTER COLUMN tx_count SET DEFAULT 0;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'blocks' AND column_name = 'uncles') THEN
        ALTER TABLE blocks ALTER COLUMN uncles SET DEFAULT '{}'::text[];
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'blocks' AND column_name = 'tx_hashes') THEN
        ALTER TABLE blocks ALTER COLUMN tx_hashes SET DEFAULT '{}'::text[];
    END IF;
END $$;
