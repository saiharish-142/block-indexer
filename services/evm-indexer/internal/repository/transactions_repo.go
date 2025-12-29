package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type TransactionsRepo struct {
	pool *pgxpool.Pool
}

func (r *TransactionsRepo) Insert(ctx context.Context, tx models.Transaction) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO transactions 
		(hash, block_hash, block_number, tx_index, from_address, to_address, nonce, value, gas_price, gas_used, gas_limit, status, method_sig, input, contract_address, type, is_canonical) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (hash) DO NOTHING`,
		tx.Hash, tx.BlockHash, tx.BlockNumber, tx.Index, tx.From, tx.To, tx.Nonce, tx.Value, tx.GasPrice, tx.GasUsed, tx.GasLimit, tx.Status, tx.MethodSignature, tx.Input, tx.ContractAddress, tx.Type, tx.IsCanonical)
	return err
}
