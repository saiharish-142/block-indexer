package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type DagBlocksRepo struct {
	pool *pgxpool.Pool
}

func (r *DagBlocksRepo) Insert(ctx context.Context, block models.DagBlock) error {
	query := `
		INSERT INTO dag_blocks (
			hash, block_order, height, weight, tx_root, state_root, parent_root,
			confirmations, txs_valid, difficulty, bits, pow_name, pow_type, pow_nonce,
			timestamp, reward, tx_fee, is_blue, main_chain, inserted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (hash) DO UPDATE SET
			block_order = EXCLUDED.block_order,
			height = EXCLUDED.height,
			is_blue = EXCLUDED.is_blue,
			main_chain = EXCLUDED.main_chain;
	`
	_, err := r.pool.Exec(ctx, query,
		block.Hash, block.Order, block.Height, block.Weight, block.TxRoot, block.StateRoot, block.ParentRoot,
		block.Confirmations, block.TxsValid, block.Difficulty, block.Bits, block.PowName, block.PowType, block.PowNonce,
		block.Timestamp, block.Reward, block.TxFee, block.IsBlue, block.MainChain, block.InsertedAt,
	)
	return err
}
