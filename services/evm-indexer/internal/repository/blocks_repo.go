package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type BlocksRepo struct {
	pool *pgxpool.Pool
}

func (r *BlocksRepo) Insert(ctx context.Context, block models.Block) error {
	raw, _ := json.Marshal(block.Raw)
	_, err := r.pool.Exec(ctx, `INSERT INTO blocks (hash, parent_hash, number, timestamp, tx_count, is_canonical, raw) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hash) DO UPDATE SET is_canonical = EXCLUDED.is_canonical, raw = EXCLUDED.raw`, block.Hash, block.ParentHash, block.Number, toUnix(block.Timestamp), block.TxCount, block.IsCanonical, raw)
	return err
}

func toUnix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}
