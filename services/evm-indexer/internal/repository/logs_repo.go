package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type LogsRepo struct {
	pool *pgxpool.Pool
}

func (r *LogsRepo) Insert(ctx context.Context, log models.Log) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO logs
		(tx_hash, block_hash, block_number, tx_index, log_index, address, topics, data, is_canonical)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, log.TxHash, log.BlockHash, log.BlockNumber, log.TxIndex, log.Index, log.Address, log.Topics, log.Data, log.IsCanonical)
	return err
}
