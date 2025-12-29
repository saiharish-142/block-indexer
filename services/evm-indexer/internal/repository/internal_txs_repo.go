package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/example/block-indexer/pkg/models"
)

type InternalTxsRepo struct {
	pool *pgxpool.Pool
}

func (r *InternalTxsRepo) Insert(ctx context.Context, tx models.InternalTransaction) error {
	query := `
		INSERT INTO internal_transactions (
			tx_hash, block_hash, block_number, from_address, to_address, value,
			gas_limit, gas_used, call_type, trace_address, error, is_canonical, inserted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query,
		tx.TxHash, tx.BlockHash, tx.BlockNumber, tx.From, tx.To, tx.Value,
		tx.GasLimit, tx.GasUsed, tx.CallType, tx.TraceAddress, tx.Error, tx.IsCanonical, tx.InsertedAt,
	)
	return err
}
