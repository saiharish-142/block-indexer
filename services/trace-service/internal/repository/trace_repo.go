package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type TraceRepo struct {
	pool *pgxpool.Pool
}

func NewTraceRepo(pool *pgxpool.Pool) *TraceRepo {
	return &TraceRepo{pool: pool}
}

func (r *TraceRepo) NextQueued(ctx context.Context, limit int) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT tx_hash FROM tx_trace_queue WHERE status = 0 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var hashes []string
	for rows.Next() {
		var h string
		if err := rows.Scan(&h); err == nil {
			hashes = append(hashes, h)
		}
	}
	return hashes, nil
}

func (r *TraceRepo) MarkProcessing(ctx context.Context, hash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tx_trace_queue SET status = 1, updated_at = now() WHERE tx_hash = $1`, hash)
	return err
}

func (r *TraceRepo) StoreTrace(ctx context.Context, trace models.TxTrace) error {
	return pgx.BeginTxFunc(ctx, r.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO tx_traces (tx_hash, trace) VALUES ($1, $2) ON CONFLICT (tx_hash) DO UPDATE SET trace = EXCLUDED.trace`, trace.TxHash, trace.Trace); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE tx_trace_queue SET status = 2, updated_at = now() WHERE tx_hash = $1`, trace.TxHash)
		return err
	})
}

func (r *TraceRepo) GetTrace(ctx context.Context, hash string) (models.TxTrace, error) {
	var trace models.TxTrace
	err := r.pool.QueryRow(ctx, `SELECT tx_hash, trace FROM tx_traces WHERE tx_hash = $1`, hash).Scan(&trace.TxHash, &trace.Trace)
	return trace, err
}
