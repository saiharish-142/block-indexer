package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DagEdgesRepo struct {
	pool *pgxpool.Pool
}

func (r *DagEdgesRepo) Insert(ctx context.Context, blockHash, parentHash string) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO dag_edges (block_hash, parent_hash) VALUES ($1, $2) ON CONFLICT DO NOTHING`, blockHash, parentHash)
	return err
}
