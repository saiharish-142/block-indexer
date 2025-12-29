package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/block-indexer/pkg/models"
)

type ContractsRepo struct {
	pool *pgxpool.Pool
}

func NewContractsRepo(pool *pgxpool.Pool) *ContractsRepo {
	return &ContractsRepo{pool: pool}
}

func (r *ContractsRepo) Get(ctx context.Context, address string) (models.Contract, error) {
	var c models.Contract
	err := r.pool.QueryRow(ctx, `SELECT address, creator_tx_hash, created_in_block, name, symbol, decimals, verified, abi, bytecode, runtime_bytecode FROM contracts WHERE address = $1`, address).
		Scan(&c.Address, &c.CreatorTxHash, &c.CreatedInBlock, &c.Name, &c.Symbol, &c.Decimals, &c.Verified, &c.ABI, &c.Bytecode, &c.RuntimeBytecode)
	return c, err
}

func (r *ContractsRepo) Upsert(ctx context.Context, c models.Contract) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO contracts (address, creator_tx_hash, created_in_block, name, symbol, decimals, verified, abi, bytecode, runtime_bytecode)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (address) DO UPDATE SET verified = EXCLUDED.verified, abi = EXCLUDED.abi`, c.Address, c.CreatorTxHash, c.CreatedInBlock, c.Name, c.Symbol, c.Decimals, c.Verified, c.ABI, c.Bytecode, c.RuntimeBytecode)
	return err
}
