package db

import (
	"context"

	"github.com/example/block-indexer/internal/pb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CopyBlocks ingests a slice of blocks into Postgres using CopyFrom for throughput.
func CopyBlocks(ctx context.Context, pool *pgxpool.Pool, blocks []pb.BlockSummary) error {
	rows := make([][]any, 0, len(blocks))
	for _, b := range blocks {
		rows = append(rows, []any{b.Number, b.Hash, b.ParentHash, b.Timestamp})
	}

	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"blocks"},
		[]string{"number", "hash", "parent_hash", "timestamp"},
		pgx.CopyFromRows(rows),
	)
	return err
}
