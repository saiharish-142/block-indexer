-- +goose Up
CREATE TABLE IF NOT EXISTS search_index (
    id BIGSERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    key TEXT NOT NULL,
    display_name TEXT,
    metadata JSONB,
    inserted_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_search_key ON search_index(key);
CREATE INDEX IF NOT EXISTS idx_search_type ON search_index(type);

-- +goose Down
DROP TABLE IF EXISTS search_index;
