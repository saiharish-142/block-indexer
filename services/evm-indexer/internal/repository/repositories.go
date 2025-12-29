package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type Repositories struct {
	Blocks           *BlocksRepo
	Transactions     *TransactionsRepo
	Logs             *LogsRepo
	InternalTxs      *InternalTxsRepo
	TokenTransfers   *TokenTransfersRepo
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Blocks:           &BlocksRepo{pool: pool},
		Transactions:     &TransactionsRepo{pool: pool},
		Logs:             &LogsRepo{pool: pool},
		InternalTxs:      &InternalTxsRepo{pool: pool},
		TokenTransfers:   &TokenTransfersRepo{pool: pool},
	}
}

func (r *Repositories) InsertEnvelope(ctx context.Context, env models.BlockEnvelope) error {
	if err := r.Blocks.Insert(ctx, env.Block); err != nil {
		return err
	}
	for _, tx := range env.Transactions {
		if err := r.Transactions.Insert(ctx, tx); err != nil {
			return err
		}
	}
	for _, log := range env.Logs {
		if err := r.Logs.Insert(ctx, log); err != nil {
			return err
		}
	}
	for _, itx := range env.InternalTxs {
		if err := r.InternalTxs.Insert(ctx, itx); err != nil {
			return err
		}
	}
	for _, tt := range env.TokenTransfers {
		if err := r.TokenTransfers.Insert(ctx, tt); err != nil {
			return err
		}
	}
	return nil
}
