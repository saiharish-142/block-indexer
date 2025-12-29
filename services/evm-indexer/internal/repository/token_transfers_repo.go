package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/example/block-indexer/pkg/models"
)

type TokenTransfersRepo struct {
	pool *pgxpool.Pool
}

func (r *TokenTransfersRepo) Insert(ctx context.Context, transfer models.TokenTransfer) error {
	query := `
		INSERT INTO token_transfers (
			tx_hash, block_hash, block_number, log_index, token_address, from_address, to_address,
			value, token_id, type, is_canonical, inserted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query,
		transfer.TxHash, transfer.BlockHash, transfer.BlockNumber, transfer.LogIndex, transfer.TokenAddress,
		transfer.From, transfer.To, transfer.Value, transfer.TokenID, transfer.Type,
		transfer.IsCanonical, transfer.InsertedAt,
	)
	return err
}
