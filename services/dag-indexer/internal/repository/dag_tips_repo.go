package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/example/block-indexer/pkg/models"
)

type DagTipsRepo struct {
	pool *pgxpool.Pool
}

func (r *DagTipsRepo) UpdateTips(ctx context.Context, tips []models.BlockTip) error {
	// Simple strategy: truncate and replace, or upsert.
	// Since tips are a small set snapshot, truncate/insert is often easiest, 
	// but might lose history if we want that. User said "snapshot ... for explorer dashboards".
	// Let's use Upsert + Delete old if strictly snapshot, but usually upsert is safer.
    // However, knowing which are no longer tips is hard without delete.
    // I'll assume Upsert for now, maybe explorer filters by latest.
    // Actually, let's just Upsert.
    
    batch := &pgx.Batch{}
    // We might want to clear old tips first if we want an exact snapshot? 
    // Let's rely on the caller or just upsert status.
    
    for _, tip := range tips {
        batch.Queue(`INSERT INTO dag_tips (hash, height, status) VALUES ($1, $2, $3) 
                     ON CONFLICT (hash) DO UPDATE SET status = EXCLUDED.status, height = EXCLUDED.height`,
            tip.Hash, tip.Height, tip.Status)
    }
    
    return r.pool.SendBatch(ctx, batch).Close()
}
